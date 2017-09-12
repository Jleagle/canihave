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

type Item struct {
	ID          string
	DateCreated int
	DateUpdated int
	DateScanned int
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

func (i *Item) GetDetailsLink() string {
	slug := slugify.Marshal(i.Name, true)
	return "/" + i.ID + "/" + slug
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

func (i *Item) IncrementHits() (item Item) {

	conn := store.GetMysqlConnection()

	_, err := conn.Exec("UPDATE items SET hits = hits + 1 WHERE id = ?", i.ID)
	if err != nil {
		fmt.Println(err)
	}

	i.Hits++

	return item
}

func (i *Item) GetWithExtras() {

	i.Get()

	if i.Status == "" && i.Region != "" && i.DateScanned == 0 {
		findSimilar(i.ID, i.Region)
		//saveNodeItems(i.Node, i.Region)
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
		panic(err)
	}
}

func (i *Item) getFromMysql() (found bool) {

	// Make the query
	builder := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": i.ID}).Limit(1)

	row := store.QueryRow(builder)
	err := row.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Status, &i.Type, &i.CompanyName)
	if err != nil {
		//fmt.Printf("%v", err.Error())
		return false
	}

	return true
}

func (i *Item) saveToMysql() {

	// todo, check this works
	date := int(time.Now().Unix())
	if i.DateCreated == 0 {
		i.DateCreated = date
	}
	if i.DateUpdated == 0 {
		i.DateUpdated = date
	}

	if i.Region == "" {
		panic("no region")
	}

	builder := squirrel.Insert("items")
	builder = builder.Columns("id", "dateCreated", "dateUpdated", "dateScanned", "name", "link", "source", "salesRank", "photo", "node", "nodeName", "price", "region", "hits", "status", "type", "companyName")
	builder = builder.Values(i.ID, i.DateCreated, i.DateUpdated, i.DateScanned, i.Name, i.Link, i.Source, i.SalesRank, i.Photo, i.Node, i.NodeName, i.Price, i.Region, i.Hits, i.Status, i.Type, i.CompanyName)

	_, err := store.Insert(builder)

	if sqlerr, ok := err.(*mysql.MySQLError); ok {
		if sqlerr.Number == 1062 { // Duplicate entry
			return
		}
		if sqlerr.Number == 1040 { // Too many connections
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
	price, _ := strconv.Atoi(amazonItem.ItemAttributes.ListPrice.Amount)

	i.Name = amazonItem.ItemAttributes.Title
	i.Link = amazonItem.DetailPageURL
	i.SalesRank = amazonItem.SalesRank
	i.Photo = amazonItem.LargeImage.URL
	i.Node = "0" //todo
	i.NodeName = amazonItem.ItemAttributes.ProductGroup
	i.Price = price
	i.CompanyName = "" //todo

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
