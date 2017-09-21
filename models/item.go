package models

import (
	"encoding/json"
	"errors"
	"fmt"
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
	TYPE_MANUAL    string = "manual"
	TYPE_SCRAPE    string = "scrape"
	TYPE_SIMILAR   string = "similar"
	TYPE_NODE      string = "node"
	TYPE_SEARCH    string = "search"
	TYPE_INCORRECT string = "incorrect"
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
	Status      string
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

func Get(id string, region string) (item Item, err error) {

	// Get from cache
	item, err = getFromMemcache(id)
	if err != nil {
		logger.Info("Retrieving " + id + " from cache")
		return item, err
	}

	// Get from MySQL
	item, err = getFromMysql(id)
	if err != nil {
		logger.Info("Retrieving " + id + " from SQL")
		saveToMemcache(item)
		return item, err
	}

	// Get from Amazon
	item, err = getFromAmazon(id, region)
	if err != nil {
		fmt.Println("Retrieving " + id + " from Amazon")
		saveToMysql(item)
		saveToMemcache(item)
		return item, err
	}

	return item, err
}

func GetWithExtras(i Item) {

	if i.Region == "" {
		i.Region = location.US
	}

	if i.Type == "" {
		i.Type = TYPE_MANUAL
	}

	item, err := Get(i.ID, i.Region)
	if err != nil {
		logger.Err("Something went wrong") // todo
	}

	lastWeek := time.Now().AddDate(0, 0, -7)

	if item.Status == "" && item.Region != "" && item.DateScanned < lastWeek.Unix() {
		findSimilar(item.ID, item.Region)
		findNodeitems(item.Node, item.Region)
		findReviews()

		// Update DateScanned
		err := store.Update(squirrel.Update("items").Set("DateScanned", time.Now().Unix()).Where("id = ?", i.ID))
		if err != nil {
			logger.Info("Can't update DateScanned: " + err.Error())
		}
	}
}

func GetMulti(ids []string, region string) (items []Item) {

	// Memcache
	mcItems, err := store.GetMemcacheMulti(ids)
	if err != nil {
		logger.Err("Can't get from memcache: " + err.Error())
	}

	for _, v := range mcItems {
		item := decodeItem(v.Value)
		items = append(items, item)
		ids = helpers.RemFromArray(ids, item.ID)
	}

	// MySQL
	builder := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": ids})
	rows := store.Query(builder)
	defer rows.Close()

	for rows.Next() {
		i := Item{}
		err := rows.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Status, &i.Type, &i.CompanyName)
		if err.Error() == "sql: no rows in result set" {
			// No problem
		} else if err != nil {
			logger.Err("Can't scan item: " + err.Error())
		}

		items = append(items, i)
		ids = helpers.RemFromArray(ids, i.ID)
	}

	// Amazon
	for _, v := range ids {

		response, err := amaz.GetItemDetails([]string{v}, region)

		if len(response.Items.Item) > 0 && err == nil {

			item := amazonItemToItem(response.Items.Item[0])
			items = append(items, item)
		} else {

			logger.Err("Can't get item from amazon: " + err.Error())
		}
	}

	return items
}

func IncrementHits(id string) {

	builder := squirrel.Update("items").Set("hits", squirrel.Expr("hits + 1")).Where("id = ?", id)
	err := store.Update(builder)
	if err != nil {
		logger.Err("Cant increment hits query: " + err.Error())
	}
}

func amazonItemToItem(amazonItem amazon.Item) (item Item) {

	var err error
	if amazonItem.ItemAttributes.ListPrice.Amount == "" {
		item.Price = 0
	} else {
		item.Price, err = strconv.Atoi(amazonItem.ItemAttributes.ListPrice.Amount)
		if err != nil {
			log.Fatal("Error converting string to int")
		}
	}

	item.Name = amazonItem.ItemAttributes.Title
	item.Link = amazonItem.DetailPageURL
	item.SalesRank = amazonItem.SalesRank
	item.Photo = amazonItem.LargeImage.URL
	item.Node = "0" //todo
	item.NodeName = amazonItem.ItemAttributes.ProductGroup
	item.CompanyName = amazonItem.ItemAttributes.Manufacturer

	return item
}

func getFromMemcache(id string) (item Item, err error) {

	mcItem, err := store.GetMemcacheItem(id)
	if err == nil {
		return decodeItem(mcItem.Value), err
	}
	logger.Err("Failed to get form memcache: " + err.Error())
	return item, err
}

func saveToMemcache(item Item) (success bool, err error) {

	err = store.SetMemcacheItem(item.ID, encodeItem(item))
	if err == nil {
		return true, err
	}
	logger.Err("Failed to save to memcache: " + err.Error())
	return false, err
}

func getFromMysql(id string) (i Item, err error) {

	// Make the query
	builder := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": i.ID}).Limit(1)
	row := store.QueryRow(builder)
	err = row.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Status, &i.Type, &i.CompanyName)

	if err.Error() == "sql: no rows in result set" {
		return i, err
	} else if err != nil {
		logger.Err("Can't retrieve from MySQL: " + err.Error())
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
	builder = builder.Columns("id", "dateCreated", "dateUpdated", "dateScanned", "name", "link", "source", "salesRank", "photo", "node", "nodeName", "price", "region", "hits", "status", "type", "companyName")
	builder = builder.Values(i.ID, i.DateCreated, i.DateUpdated, i.DateScanned, i.Name, i.Link, i.Source, i.SalesRank, i.Photo, i.Node, i.NodeName, i.Price, i.Region, i.Hits, i.Status, i.Type, i.CompanyName)

	err = store.Insert(builder)
	if sqlerr, ok := err.(*mysql.MySQLError); ok {
		if sqlerr.Number == mysqlerr.ER_DUP_ENTRY { // Duplicate entry
			logger.Info("Trying to insert dupe entry: " + err.Error())
			return true, nil
		}
	}

	if err == nil {
		return true, nil
	} else {
		logger.Err("Trying to add item to Mysql: " + err.Error())
		return false, err
	}
}

func getFromAmazon(id string, region string) (item Item, err error) {

	response, err := amaz.GetItemDetails([]string{id}, region)

	if len(response.Items.Item) > 0 && err == nil {

		return amazonItemToItem(response.Items.Item[0]), err
	} else {

		item.Status = err.Error()
		return item, err
	}
}

func decodeItem(raw []byte) (item Item) {
	err := json.Unmarshal(raw, item)
	if err != nil {
		logger.Err("Error decoding item to JSON: " + err.Error())
	}
	return item
}

func encodeItem(item Item) []byte {
	enc, err := json.Marshal(item)
	if err != nil {
		logger.Err("Error encoding item to JSON: " + err.Error())
	}
	return enc
}

func interfaceToItem(tst interface{}) (ret Item) {
	ret, _ = tst.(Item)
	return ret
}
