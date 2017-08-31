package models

import (
	"fmt"
	"log"
	"time"

	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
	"github.com/patrickmn/go-cache"
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
	Currency     string

	Status string
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
		i.Currency = item.Currency
	}
	return found
}

func (i *Item) saveToMemcache() {

	x := store.GetGoCache()
	x.Set(i.ID, i, cache.DefaultExpiration)
}

func (i *Item) getFromMysql() (found bool) {

	// Make the query
	query := squirrel.Select("*").From("items").Where(squirrel.Eq{"id": i.ID}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	db := store.GetMysqlConnection()
	err = db.QueryRow(sql, args...).Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.ProductGroup, &i.Price, &i.Currency)
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

	// Make query
	//sql, args, err := sq.Insert("items").Columns("name", "age").Values("moe", 13).ToSql()

	conn := store.GetMysqlConnection()

	// todo, switch to query builder
	// Prepare statement for inserting data
	insert, err := conn.Prepare("INSERT INTO items (id, dateCreated, dateUpdated, `name`, link, source, salesRank, photo, productGroup, price, currency) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	// run query
	_, err = insert.Exec(i.ID, i.DateCreated, i.DateUpdated, i.Name, i.Link, i.Source, i.SalesRank, i.Photo, i.ProductGroup, i.Price, i.Currency)
	if err != nil {
		panic(err.Error())
	}
}

func (i *Item) getFromAmazon() (found bool) {

	// Setup Amazon
	client, err := amazon.NewFromEnvionment()
	if err != nil {
		log.Fatal(err)
	}

	// Make API call
	res, err := client.ItemLookup(amazon.ItemLookupParameters{
		ResponseGroups: []amazon.ItemLookupResponseGroup{
			amazon.ItemLookupResponseGroupLarge,
		},
		IDType:  amazon.IDTypeASIN,
		ItemIDs: []string{i.ID},
	}).Do()

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
	i.Currency = amazonItem.ItemAttributes.ListPrice.CurrencyCode

	return true
}
