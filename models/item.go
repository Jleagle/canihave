package models

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	amaz "github.com/Jleagle/canihave/amazon"
	"github.com/Jleagle/canihave/helpers"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/metal3d/go-slugify"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

const (
	TYPE_MANUAL  string = "manual"
	TYPE_SCRAPE  string = "scrape"
	TYPE_SIMILAR string = "similar"
	TYPE_NODE    string = "node"
	TYPE_SEARCH  string = "search"
)

type Item struct {
	ID          string
	DateCreated int64
	DateUpdated int64
	DateScanned int64
	Name        string
	Link        string
	Source      string
	SalesRank   int
	Photo       string
	Node        string
	NodeName    string
	Price       int
	Region      string
	Hits        int
	Type        string
	CompanyName string
}

func (i *Item) GetAmazonLink() string {
	if i.Region == location.US {
		return strings.Replace(i.Link, "www", "smile", 1)
	}
	return i.Link
}

func (i *Item) GetSlug() string {
	return slugify.Marshal(i.Name, true)
}

func (i *Item) GetPath() string {
	return "/" + i.ID + "/" + i.GetSlug()
}

func (i *Item) GetLink() string {
	return "https://canihave.one" + i.GetPath()
}

func (i *Item) GetPrice() float32 {
	return float32(i.Price) / 100
}

func (i *Item) GetCurrency() string {
	return location.GetCurrency(i.Region)
}

func (i *Item) GetFlag() string {
	return "/assets/flags/" + i.Region + ".gif"
}

func (i *Item) NeedsScanning() bool {

	return i.DateScanned < time.Now().AddDate(0, 0, -7).Unix() && helpers.InArray(i.Type, []string{TYPE_SCRAPE, TYPE_MANUAL})
}

func GetWithExtras(id string, region string, itemType string, source string) (item Item, err error) {

	item, err = Get(id, region, itemType, source)
	if err != nil {
		logger.Err("Can't get item from anywhere", err)
		return item, err
	}

	if item.NeedsScanning() {
		saveSimilarItems(id, region, itemType)
		saveNodeitems(item.Node, item.Region)
		saveReviews()
		updateDateScanned(id)
	}

	return item, err
}

func Get(id string, region string, itemType string, source string) (item Item, err error) {

	// Get from cache
	item, err = getFromMemcache(id)
	if err == nil {
		logger.Info("Retrieving " + id + " from cache")
		return item, err
	}

	// Get from MySQL
	item, err = getFromMysql(id)
	if err == nil {
		logger.Info("Retrieving " + id + " from SQL")
		saveToMemcache(item)
		return item, err
	}

	// Get from Amazon
	item, err = getFromAmazon(id, region, itemType)
	if err == nil {
		logger.Info("Retrieving " + id + " from Amazon")
		saveToMysql(item)
		saveToMemcache(item)
		return item, err
	}

	return item, err
}

func GetMulti(ids []string, region string, itemType string) (items []Item) {

	// Memcache
	mcItems, err := store.GetMemcacheMulti(ids)
	if err != nil {
		logger.Err("Can't get from memcache", err)
	}

	for _, v := range mcItems {
		item := decodeItem(v.Value)
		items = append(items, item)
		ids = helpers.RemFromArray(ids, store.MEMCACHE_APP_KEY+item.ID)
	}

	// MySQL
	builder := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": ids})
	rows := store.Query(builder)
	defer rows.Close()

	for rows.Next() {
		i := Item{}
		err := rows.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Type, &i.CompanyName)
		if err != nil && err.Error() == "sql: no rows in result set" {
			// No problem
		} else if err != nil {
			logger.Err("Can't scan item", err)
		}

		items = append(items, i)
		ids = helpers.RemFromArray(ids, i.ID)
	}

	// Amazon
	for _, v := range ids {

		response, err := amaz.GetItemDetails([]string{v}, region)

		if err == nil && len(response.Items.Item) > 0 {

			item := amazonItemToItem(response.Items.Item[0], itemType, region)
			items = append(items, item)
		} else {

			logger.Err("Can't get item from amazon", err)
		}
	}

	return items
}

func updateDateScanned(id string) (err error) {
	builder := squirrel.Update("items").Set("DateScanned", time.Now().Unix()).Where("id = ?", id)
	err = store.Update(builder)
	if err != nil {
		logger.Info("Can't update DateScanned", err)
	}

	delErr := store.DeleteMemcacheItem(id)
	if delErr != nil {
		logger.Info("Can't delete memcache object", err)
	}

	return err
}

func amazonItemToItem(amazonItem amazon.Item, itemType string, region string) (item Item) {

	var err error
	if amazonItem.ItemAttributes.ListPrice.Amount == "" {
		item.Price = 0
	} else {
		item.Price, err = strconv.Atoi(amazonItem.ItemAttributes.ListPrice.Amount)
		if err != nil {
			log.Fatal("Error converting string to int")
		}
	}

	item.ID = amazonItem.ASIN
	item.Name = amazonItem.ItemAttributes.Title
	item.Link = amazonItem.DetailPageURL
	item.SalesRank = amazonItem.SalesRank
	item.Photo = amazonItem.LargeImage.URL
	item.Node = "0" //todo
	item.NodeName = amazonItem.ItemAttributes.ProductGroup
	item.CompanyName = amazonItem.ItemAttributes.Manufacturer
	item.Type = itemType
	item.Region = region

	return item
}

func getFromMemcache(id string) (item Item, err error) {

	mcItem, err := store.GetMemcacheItem(id)
	if err == nil {
		return decodeItem(mcItem.Value), err
	}

	return item, err
}

func saveToMemcache(item Item) (success bool, err error) {

	err = store.SetMemcacheItem(item.ID, encodeItem(item))
	if err == nil {
		return true, err
	}
	logger.Err("Failed to save to memcache", err)
	return false, err
}

func getFromMysql(id string) (i Item, err error) {

	// Make the query
	builder := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": id}).Limit(1)
	row := store.QueryRow(builder)
	err = row.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Type, &i.CompanyName)

	if err != nil && err.Error() == "sql: no rows in result set" {
		return i, err
	} else if err != nil {
		logger.Err("Can't retrieve from MySQL", err)
		return i, err
	} else {
		return i, err
	}
}

func saveToMysql(i Item) (success bool, err error) {

	if i.DateCreated == 0 {
		i.DateCreated = time.Now().Unix()
	}
	if i.DateUpdated == 0 {
		i.DateUpdated = time.Now().Unix()
	}

	if i.Region == "" {
		logger.Err("Item has no region")
		return false, errors.New("can't save item into mysql with no region")
	}

	builder := squirrel.Insert("items")
	builder = builder.Columns("id", "dateCreated", "dateUpdated", "dateScanned", "name", "link", "source", "salesRank", "photo", "node", "nodeName", "price", "region", "hits", "type", "companyName")
	builder = builder.Values(i.ID, i.DateCreated, i.DateUpdated, i.DateScanned, i.Name, i.Link, i.Source, i.SalesRank, i.Photo, i.Node, i.NodeName, i.Price, i.Region, i.Hits, i.Type, i.CompanyName)

	err = store.Insert(builder)
	if sqlerr, ok := err.(*mysql.MySQLError); ok {
		if sqlerr.Number == mysqlerr.ER_DUP_ENTRY {
			//logger.Info("Trying to insert dupe entry", err)
			return true, nil
		}
	}

	if err == nil {
		return true, nil
	} else {
		logger.Err("Trying to add item to Mysql", err)
		return false, err
	}
}

func getFromAmazon(id string, region string, itemType string) (item Item, err error) {

	response, err := amaz.GetItemDetails([]string{id}, region)

	if err == nil && len(response.Items.Item) > 0 {

		return amazonItemToItem(response.Items.Item[0], itemType, region), err
	} else {

		return item, err
	}
}

func decodeItem(raw []byte) (item Item) {
	err := json.Unmarshal(raw, &item)
	if err != nil {
		logger.Err("Error decoding item to JSON", err)
	}
	return item
}

func encodeItem(item Item) []byte {
	enc, err := json.Marshal(item)
	if err != nil {
		logger.Err("Error encoding item to JSON", err)
	}
	return enc
}
