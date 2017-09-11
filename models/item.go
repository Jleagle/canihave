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
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-sql-driver/mysql"
	"github.com/metal3d/go-slugify"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

const (
	TYPE_SCRAPE    string = "scrape"
	TYPE_SIMILAR   string = "similar"
	TYPE_NODE      string = "node"
	TYPE_SEARCH    string = "search"
	TYPE_INCORRECT string = "incorrect"
)

// item is the database row
type Item struct {
	ID           string
	DateCreated  string
	DateUpdated  string
	Name         string
	Link         string
	Source       string
	SalesRank    int
	Photo        string
	ProductGroup string
	Price        string
	Region       string
	Hits         int
	Status       string
	Type         string
}

func (i *Item) GetAmazonLink() string {
	if i.Region == location.US {
		return strings.Replace(i.Link, "www", "smile", 1)
	}
	return i.Link
}

func (i *Item) GetDetailsLink() string {
	slug := slugify.Marshal(i.Name, true)
	return "/" + i.ID + "/" + slug
}

func (i *Item) GetPrice() float32 {
	x, err := strconv.Atoi(i.Price)
	if err != nil {
		fmt.Println(err)
	}
	return float32(x) / 100
}

func (i *Item) GetCurrency() string {
	return location.GetCurrency(i.Region)
}

func (i *Item) GetFlag() string {
	return "/assets/flags/" + i.Region + ".gif"
}

func (i *Item) Reset() {

	//i.ID = ""
	i.DateCreated = ""
	i.DateUpdated = ""
	i.Name = ""
	i.Link = ""
	//i.Source = ""
	i.SalesRank = 0
	i.Photo = ""
	i.ProductGroup = ""
	i.Price = ""
	i.Region = ""
	i.Hits = 0
	//i.Status = ""
	i.Type = ""
}

func (i *Item) IncrementHits() (item Item) {

	conn := store.GetMysqlConnection()
	_, err := conn.Exec("UPDATE items SET hits = hits + 1 WHERE id = ?", i.ID)
	if err != nil {
		fmt.Println(err)
	}

	i.Hits++

	return item
}

func (i *Item) GetAll() {
	i.Get()
	saveSimilar(i.ID, i.Region)
	//saveNodeItems(i.Node, i.Region)
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
		fmt.Println("Retrieving " + i.ID + " from cache")
		return
	}

	// Get from MySQL
	if i.getFromMysql() {
		fmt.Println("Retrieving " + i.ID + " from SQL")
		i.saveToMemcache()
		return
	}

	// Get from Amazon
	if i.getFromAmazon() {
		fmt.Println("Retrieving " + i.ID + " from Amazon")
		i.saveAsNewMysqlRow()
		i.saveToMemcache()
		return
	}

	// Save errors into cache too
	if strings.Contains(i.Status, "AWS.InvalidParameterValue") {

		i.Reset()
		i.saveAsNewMysqlRow()
		i.saveToMemcache()

	} else if strings.Contains(i.Status, "RequestThrottled") {

		i.Get()
	}
}

//func (i *Item) Save() (result sql.Result) {
//
//	//builder := squirrel.Update("items").Limit(1)
//	//builder= builder.Set("id", "dateCreated", "dateUpdated", "name", "link", "source", "salesRank", "photo", "productGroup", "price", "region", "hits", "status", "type")
//
//	conn := store.GetMysqlConnection()
//	result, err := conn.Exec("UPDATE items SET id = ?, dateCreated = ?, dateUpdated = ?, name = ?, link = ?, source = ?, salesRank = ?, photo = ?, productGroup = ?, price = ?, region = ?, hits = ?, status = ? WHERE id = ?", i.ID, i.DateCreated, i.DateUpdated, i.Name, i.Link, i.Source, i.SalesRank, i.Photo, i.ProductGroup, i.Price, i.Region, i.Hits, i.Status, i.Type)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	return result
//}

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
	i.Name = item.Name
	i.Link = item.Link
	i.Source = item.Source
	i.SalesRank = item.SalesRank
	i.Photo = item.Photo
	i.ProductGroup = item.ProductGroup
	i.Price = item.Price
	i.Region = item.Region
	i.Hits = item.Hits
	i.Status = item.Status
	i.Type = item.Type

	return true
}

func (i *Item) saveToMemcache() {

	err := store.GetMemcacheConnection().Set(&memcache.Item{Key: i.ID, Value: EncodeItem(*i)})
	if err != nil {
		panic(err)
	}
}

func (i *Item) getFromMysql() (found bool) {

	// Make the query
	query := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": i.ID}).Limit(1)
	s, args, err := query.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	conn := store.GetMysqlConnection()
	err = conn.QueryRow(s, args...).Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.ProductGroup, &i.Price, &i.Region, &i.Hits, &i.Status, &i.Type)
	if err != nil {
		//fmt.Printf("%v", err.Error())
		return false
	}

	return true
}

func (i *Item) saveAsNewMysqlRow() {

	if i.Price == "" {
		i.Price = "0"
	}

	date := time.Now().Format("2006-01-02 15:04:05")
	if i.DateCreated == "" {
		i.DateCreated = date
	}
	if i.DateUpdated == "" {
		i.DateUpdated = date
	}

	// run query
	builder := squirrel.Insert("items")
	builder = builder.Columns("id", "dateCreated", "dateUpdated", "name", "link", "source", "salesRank", "photo", "productGroup", "price", "region", "hits", "status", "type")
	builder = builder.Values(i.ID, i.DateCreated, i.DateUpdated, i.Name, i.Link, i.Source, i.SalesRank, i.Photo, i.ProductGroup, i.Price, i.Region, i.Hits, i.Status, i.Type)

	_, err := store.Insert(builder)

	if sqlerr, ok := err.(*mysql.MySQLError); ok {
		if sqlerr.Number == 1062 { // Duplicate entry
			return
		}
	}

	if err != nil {
		panic(err.Error())
	}
}

func (i *Item) getFromAmazon() (found bool) {

	response, err := amaz.GetItemDetails(i.ID, i.Region)

	if err != nil {
		i.Status = err.Error()
		return false
	}

	var amazonItem amazon.Item
	if len(response.Items.Item) > 0 {
		amazonItem = response.Items.Item[0]
	} else {
		i.Status = "Not found in Amazon"
		return false
	}

	// Make struct
	amazonItemToItem(i, amazonItem)
	i.Type = TYPE_SCRAPE

	return true
}

func DecodeItem(raw []byte) (item Item) {
	err := json.Unmarshal(raw, &item)
	if err != nil {
		fmt.Println("Error decoding item to JSON", err)
	}
	return item
}

func EncodeItem(item Item) []byte {
	enc, err := json.Marshal(item)
	if err != nil {
		fmt.Println("Error encoding item to JSON", err)
	}
	return enc
}

func interfaceToItem(tst interface{}) (ret Item) {
	ret, _ = tst.(Item)
	return ret
}

func amazonItemToItem(item *Item, amazonItem amazon.Item) {

	item.Name = amazonItem.ItemAttributes.Title
	item.Link = amazonItem.DetailPageURL
	item.SalesRank = amazonItem.SalesRank
	item.Photo = amazonItem.LargeImage.URL
	item.ProductGroup = amazonItem.ItemAttributes.ProductGroup
	item.Price = amazonItem.ItemAttributes.ListPrice.Amount
	item.Hits = 0
	item.Status = ""
}
