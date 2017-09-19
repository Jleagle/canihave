package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	amaz "github.com/Jleagle/canihave/amazon"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/VividCortex/mysqlerr"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-sql-driver/mysql"
	"github.com/metal3d/go-slugify"
)

const (
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

func (i *Item) IncrementHits() {

	i.Hits++

	builder := squirrel.Update("items").Set("hits", squirrel.Expr("hits + 1")).Where("id = ?", i.ID)
	err := store.Update(builder)

	if err != nil {
		logger.Err("Cant increment hits query: " + err.Error())
	}

	// Clear cache
	store.GetMemcacheConnection().Delete(i.ID)
}

func (i *Item) GetWithExtras() {

	i.Get()

	lastWeek := time.Now().AddDate(0, 0, -7)

	if i.Status == "" && i.Region != "" && i.DateScanned < lastWeek.Unix() {
		findSimilar(i.ID, i.Region)
		findNodeitems(i.Node, i.Region)
		findReviews()

		// Update DateScanned
		store.Update(squirrel.Update("items").Set("DateScanned", time.Now().Unix()).Where("id = ?"))
	}

}

func (i *Item) Get() {

	if i.Status != "" {
		return
	}

	if i.ID == "" {
		log.Fatal("Item needs an id")
	}

	// Get from cache
	if i.getFromMemcache() {
		logger.Info("Retrieving " + i.ID + " from cache")
		return
	}

	// Get from MySQL
	if i.getFromMysql() {
		logger.Info("Retrieving " + i.ID + " from SQL")
		i.saveToMemcache()
		return
	}

	// Get from Amazon
	if i.getFromAmazon() {
		fmt.Println("Retrieving " + i.ID + " from Amazon")
		i.saveToMysql()
		i.saveToMemcache()
		return
	}

	// Clear the data so it doesnt remember any items from before
	// todo, do we need this?
	i.DateCreated = 0
	i.DateUpdated = 0
	i.DateScanned = 0
	i.Name = ""
	i.Link = ""
	i.Source = ""
	i.SalesRank = 0
	i.Photo = ""
	i.Node = ""
	i.NodeName = ""
	i.Price = 0
	i.CompanyName = ""

	// Save invalid IDs so we dont query Amazon for them again
	if strings.Contains(i.Status, "AWS.InvalidParameterValue") {

		i.saveToMysql()
		i.saveToMemcache()
		return
	}

	// Try again
	if strings.Contains(i.Status, "RequestThrottled") {

		i.Get()
		return
	}
}

func (i *Item) getFromMemcache() (found bool) {

	byteArray, err := store.GetMemcacheConnection().Get(i.ID)

	if err == memcache.ErrCacheMiss {
		return false
	} else if err != nil {
		fmt.Println("Error fetching from memcache", err)
		return false
	}

	item := DecodeItem(byteArray.Value)
	item = interfaceToItem(item)

	i.DateCreated = item.DateCreated
	i.DateUpdated = item.DateUpdated
	i.DateScanned = item.DateScanned
	i.Name = item.Name
	i.Link = item.Link
	i.Source = item.Source
	i.SalesRank = item.SalesRank
	i.Photo = item.Photo
	i.Node = item.Node
	i.NodeName = item.NodeName
	i.Price = item.Price
	i.Region = item.Region
	i.Hits = item.Hits
	i.Status = item.Status
	i.Type = item.Type
	i.CompanyName = item.CompanyName

	return true
}

func (i *Item) saveToMemcache() {

	err := store.GetMemcacheConnection().Set(&memcache.Item{Key: i.ID, Value: EncodeItem(*i)})
	if err != nil {
		logger.Err("Failed to save to memcache: " + err.Error())
	}
}

func (i *Item) getFromMysql() (found bool) {

	// Make the query
	builder := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": i.ID}).Limit(1)

	row := store.QueryRow(builder)
	err := row.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Status, &i.Type, &i.CompanyName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// No problem
		} else {
			logger.Err("Can't retrieve from MySQL: " + err.Error())
		}
		return false
	}

	return true
}

func (i *Item) saveToMysql() {

	// todo, check this works
	date := time.Now().Unix()
	if i.DateCreated == 0 {
		i.DateCreated = date
	}
	if i.DateUpdated == 0 {
		i.DateUpdated = date
	}

	if i.Region == "" {
		logger.Err("Item has no region")
		// todo, return error
	}

	builder := squirrel.Insert("items")
	builder = builder.Columns("id", "dateCreated", "dateUpdated", "dateScanned", "name", "link", "source", "salesRank", "photo", "node", "nodeName", "price", "region", "hits", "status", "type", "companyName")
	builder = builder.Values(i.ID, i.DateCreated, i.DateUpdated, i.DateScanned, i.Name, i.Link, i.Source, i.SalesRank, i.Photo, i.Node, i.NodeName, i.Price, i.Region, i.Hits, i.Status, i.Type, i.CompanyName)

	err := store.Insert(builder)

	if sqlerr, ok := err.(*mysql.MySQLError); ok {
		if sqlerr.Number == mysqlerr.ER_DUP_ENTRY { // Duplicate entry
			logger.Info("Trying to insert dupe entry: " + err.Error())
			return
		}
	}

	if err != nil {
		logger.Err("Trying to add item to Mysql: " + err.Error())
	}
}

func (i *Item) getFromAmazon() (found bool) {

	response, err := amaz.GetItemDetails(i.ID, i.Region)

	if len(response.Items.Item) > 0 && err == nil {

		amazonItem := response.Items.Item[0]

		// Price
		var price int = 0
		if amazonItem.ItemAttributes.ListPrice.Amount != "" {
			price, err = strconv.Atoi(amazonItem.ItemAttributes.ListPrice.Amount)
			if err != nil {
				log.Fatal("Error converting string to int")
			}
		}

		i.Name = amazonItem.ItemAttributes.Title
		i.Link = amazonItem.DetailPageURL
		i.SalesRank = amazonItem.SalesRank
		i.Photo = amazonItem.LargeImage.URL
		i.Node = "0" //todo
		i.NodeName = amazonItem.ItemAttributes.ProductGroup
		i.Price = price
		i.CompanyName = amazonItem.ItemAttributes.Manufacturer

		return true
	} else {

		i.Status = err.Error()
		return false
	}
}

func DecodeItem(raw []byte) (item Item) {
	err := json.Unmarshal(raw, &item)
	if err != nil {
		logger.Err("Error decoding item to JSON: " + err.Error())
	}
	return item
}

func EncodeItem(item Item) []byte {
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
