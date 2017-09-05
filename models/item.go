package models

import (
	"fmt"
	"log"
	"time"

	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
	"github.com/patrickmn/go-cache"
	"github.com/Jleagle/canihave/location"
	"strings"
	"strconv"
)

// item is the database row
type Item struct {
	ID           string
	DateCreated  string
	DateUpdated  string
	Name         string
	Desc         string
	Link         string
	Source       string
	SalesRank    int
	Photo        string
	ProductGroup string
	Price        string
	Region       string

	Status string
}

func (i *Item) GetLink() string {
	if i.Region == location.US {
		return strings.Replace(i.Link, "www", "smile", 1)
	}
	return i.Link
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

func (i *Item) Get() {

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
}

func (i *Item) getFromMemcache() (found bool) {

	foo, found := store.GetGoCache().Get(i.ID)
	if found {
		item, _ := foo.(Item) // Cast it back to item

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
	}
	return found
}

func (i *Item) saveToMemcache() {

	c := store.GetGoCache()
	c.Set(i.ID, i, cache.DefaultExpiration)
}

func (i *Item) getFromMysql() (found bool) {

	// Make the query
	query := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": i.ID}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	conn := store.GetMysqlConnection()
	err = conn.QueryRow(sql, args...).Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.ProductGroup, &i.Price, &i.Region)
	if err != nil {
		//fmt.Printf("%v", err.Error())
		return false
	}

	return true
}

func (i *Item) saveToMysql() {

	if i.Price == "" {
		i.Price = "0"
	}

	// run query
	_, err := store.GetInsertPrep().Exec(i.ID, i.DateCreated, i.DateUpdated, i.Name, i.Link, i.Source, i.SalesRank, i.Photo, i.ProductGroup, i.Price, i.Region)
	if err != nil {
		panic(err.Error())
	}
}

func (i *Item) getFromAmazon() (found bool) {

	location.SetAmazonEnviromentVars(i.Region)

	// Amazon rate limit
	time.Sleep(1100 * time.Millisecond)

	// Setup Amazon
	client, err := amazon.NewFromEnvionment()
	if err != nil {
		log.Fatal(err)
	}

	// Make API call
	res, err := client.ItemLookup(amazon.ItemLookupParameters{
		ResponseGroups: []amazon.ItemLookupResponseGroup{
			amazon.ItemLookupResponseGroupMedium,
		},
		IDType:  amazon.IDTypeASIN,
		ItemIDs: []string{i.ID},
	}).Do()

	//fmt.Printf("%# v", pretty.Formatter(res))

	if err != nil {
		i.Status = err.Error()
		fmt.Printf("%v", err.Error())
		return false
	}

	var amazonItem amazon.Item
	if len(res.Items.Item) > 0 {
		amazonItem = res.Items.Item[0]
	} else {
		i.Status = "Not found in Amazon"
		return false
	}

	// Make struct
	date := time.Now().Format("2006-01-02")
	i.DateCreated = date
	i.DateUpdated = date
	i.Name = amazonItem.ItemAttributes.Title
	i.Link = amazonItem.DetailPageURL
	i.SalesRank = amazonItem.SalesRank
	i.Photo = amazonItem.LargeImage.URL
	i.ProductGroup = amazonItem.ItemAttributes.ProductGroup
	i.Price = amazonItem.ItemAttributes.ListPrice.Amount

	return true
}
