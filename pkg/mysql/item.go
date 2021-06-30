package mysql

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	amazonHelper "github.com/Jleagle/canihave/pkg/amazon"
	"github.com/Jleagle/canihave/pkg/helpers"
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/Masterminds/squirrel"
	"github.com/VividCortex/mysqlerr"
	"github.com/gosimple/slug"
	"github.com/memcachier/mc/v3"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
	"go.uber.org/zap"
)

const (
	TypeManual  string = "manual"
	TypeScraper string = "scrape"
	TypeSimilar string = "similar"
	TypeNode    string = "node"
	TypeSearch  string = "search"
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
	Region      amazon.Region
	Hits        int
	Type        string
	CompanyName string
}

func (i *Item) GetAmazonLink() string {
	if i.Region == amazon.RegionUS {
		return strings.Replace(i.Link, "www", "smile", 1)
	}
	return i.Link
}

func (i *Item) GetSlug() string {
	return slug.Make(i.Name)
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
	return location.GetCurrencySign(i.Region)
}

func (i *Item) GetFlag() string {
	return "/assets/flags/" + string(i.Region) + ".gif"
}

func (i *Item) NeedsScanning() bool {

	return i.DateScanned < time.Now().AddDate(0, 0, -7).Unix() && helpers.InSlice(i.Type, []string{TypeScraper, TypeManual})
}

func GetItem(id string, region amazon.Region) (item Item, err error) {

	// Memcache
	err = memcacheClient.Get("item-"+id, &item)
	if err == nil {
		return item, err
	} else if err != mc.ErrNotFound {
		return item, err
	}

	// MySQL
	sql, args, err := squirrel.Select("*").From("items").Where("id = ?", id).ToSql()
	if err != nil {
		return item, err
	}

	c, err := getClient()
	if err != nil {
		return item, err
	}

	err = c.Select(&item, sql, args...)
	if err == nil {
		return item, err
	} else if err != mc.ErrNotFound { // todo, fix error type
		return item, err
	}

	// Amazon
	response, err := amazonHelper.GetItemDetails([]string{id}, region)
	if err == nil && len(response.Items.Item) > 0 {
		return amazonItemToItem(response.Items.Item[0], itemType, region), err
	} else {
		return item, err
	}

	return item, nil
}

func GetItems(ids []int) (items []Item, err error) {

	builder := squirrel.Select("*")
	builder = builder.From("items")
	builder = builder.Where("id IN ?", ids)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	c, err := getClient()
	if err != nil {
		return nil, err
	}

	err = c.Select(&items, sql, args...)
	return items, err
}

func GetWithExtras(id string, region amazon.Region, itemType string, source string) (item Item, err error) {

	item, err = Get(id, region, itemType, source)
	if err != nil {
		logger.Logger.Error("Can't get item from anywhere", zap.Error(err))
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

func Get(id string, region amazon.Region, itemType string, source string) (item Item, err error) {

	// Get from cache
	item, err = getFromMemcache(id)
	if err == nil {
		logger.Logger.Info("Retrieving " + id + " from cache")
		return item, err
	}

	// Get from MySQL
	item, err = getFromMysql(id)
	if err == nil {
		logger.Logger.Info("Retrieving " + id + " from SQL")
		saveToMemcache(item)
		return item, err
	}

	// Get from Amazon
	item, err = getFromAmazon(id, region, itemType)
	if err == nil {
		logger.Logger.Info("Retrieving " + id + " from Amazon")
		saveToMysql(item)
		saveToMemcache(item)
		return item, err
	}

	return item, err
}

func GetMulti(ids []string, region amazon.Region, itemType string) (items []Item) {

	// Memcache
	mcItems, err := memcache.GetMemcacheMulti(ids)
	if err != nil {
		logger.Logger.Error("Can't get from memcache", zap.Error(err))
	}

	for _, v := range mcItems {
		item := decodeItem(v.Value)
		items = append(items, item)
		ids = helpers.RemoveFromSlice(ids, item.ID)
	}

	// MySQL
	builder := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": ids})
	rows := Query(builder)
	defer rows.Close()

	for rows.Next() {
		i := Item{}
		err := rows.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Type, &i.CompanyName)
		if err != nil && err.Error() == "sql: no rows in result set" {
			// No problem
		} else if err != nil {
			logger.Logger.Error("Can't scan item", zap.Error(err))
		}

		items = append(items, i)
		ids = helpers.RemoveFromSlice(ids, i.ID)
	}

	// Amazon
	for _, v := range ids {

		response, err := amazonHelper.GetItemDetails([]string{v}, region)

		if err == nil && len(response.Items.Item) > 0 {

			item := amazonItemToItem(response.Items.Item[0], itemType, region)
			items = append(items, item)
		} else {

			logger.Logger.Error("Can't get item from amazon", zap.Error(err))
		}
	}

	return items
}

func updateDateScanned(id string) (err error) {

	builder := squirrel.Update("items").Set("DateScanned", time.Now().Unix()).Where("id = ?", id)
	err = Update(builder)
	if err != nil {
		logger.Logger.Info("Can't update DateScanned", zap.Error(err))
	}

	delErr := memcache.DeleteMemcacheItem(id)
	if delErr != nil {
		logger.Logger.Info("Can't delete memcache object", zap.Error(err))
	}

	return err
}

func amazonItemToItem(amazonItem amazon.Item, itemType string, region amazon.Region) (item Item) {

	var err error
	if amazonItem.ItemAttributes.ListPrice.Amount == "" {
		item.Price = 0
	} else {
		item.Price, err = strconv.Atoi(amazonItem.ItemAttributes.ListPrice.Amount)
		if err != nil {
			logger.Logger.Fatal("Error converting string to int", zap.Error(err))
		}
	}

	item.ID = amazonItem.ASIN
	item.Name = amazonItem.ItemAttributes.Title
	item.Link = amazonItem.DetailPageURL
	item.SalesRank = amazonItem.SalesRank
	item.Photo = amazonItem.LargeImage.URL
	item.Node = "0" // todo
	item.NodeName = amazonItem.ItemAttributes.ProductGroup
	item.CompanyName = amazonItem.ItemAttributes.Manufacturer
	item.Type = itemType
	item.Region = region

	return item
}

func getFromMemcache(id string) (item Item, err error) {

	mcItem, err := memcache.GetMemcacheItem(id)
	if err == nil {
		return decodeItem(mcItem.Value), err
	}

	return item, err
}

func saveToMemcache(item Item) (success bool, err error) {

	err = memcache.SetMemcacheItem(item.ID, encodeItem(item))
	if err == nil {
		return true, err
	}
	logger.Logger.Error("Failed to save to memcache", zap.Error(err))
	return false, err
}

func saveToMysql(i Item) (success bool, err error) {

	if i.DateCreated == 0 {
		i.DateCreated = time.Now().Unix()
	}
	if i.DateUpdated == 0 {
		i.DateUpdated = time.Now().Unix()
	}

	if i.Region == "" {
		logger.Logger.Error("Item has no region")
		return false, errors.New("can't save item into mysql with no region")
	}

	builder := squirrel.Insert("items")
	builder = builder.Columns("id", "dateCreated", "dateUpdated", "dateScanned", "name", "link", "source", "salesRank", "photo", "node", "nodeName", "price", "region", "hits", "type", "companyName")
	builder = builder.Values(i.ID, i.DateCreated, i.DateUpdated, i.DateScanned, i.Name, i.Link, i.Source, i.SalesRank, i.Photo, i.Node, i.NodeName, i.Price, i.Region, i.Hits, i.Type, i.CompanyName)

	err = Insert(builder)
	if sqlerr, ok := err.(*MySQLError); ok {
		if sqlerr.Number == mysqlerr.ER_DUP_ENTRY {
			// logger.Info("Trying to insert dupe entry", zap.Error(err))
			return true, nil
		}
	}

	if err == nil {
		return true, nil
	} else {
		logger.Logger.Error("Trying to add item to Mysql", zap.Error(err))
		return false, err
	}
}

func decodeItem(raw []byte) (item Item) {
	err := json.Unmarshal(raw, &item)
	if err != nil {
		logger.Logger.Error("Error decoding item to JSON", zap.Error(err))
	}
	return item
}

func encodeItem(item Item) []byte {
	enc, err := json.Marshal(item)
	if err != nil {
		logger.Logger.Error("Error encoding item to JSON", zap.Error(err))
	}
	return enc
}
