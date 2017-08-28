package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ngs/go-amazon-product-advertising-api/amazon"
	cache "github.com/patrickmn/go-cache"
)

// item is the database row
type item struct {
	ID              string
	DateCreated     string
	DateUpdated     string
	Name            string
	Desc            string
	Link            string
	Source          string
	SalesRank       int
	Images          string
	ProductGroup    string
	ProductTypeName string
}

func (i item) GetUKLink() string {
	return "https://www.amazon.co.uk/dp/" + i.ID + "?tag=canihaveone00-21"
}

func (i item) GetUKPixel() string {
	return "//ir-uk.amazon-adsystem.com/e/ir?t=canihaveone00-21&l=am2&o=2&a=" + i.ID
}

func (i item) get() item {
	c := cache.New(5*time.Minute, 10*time.Minute)

	// Get from cache
	itemInterface, found := c.Get(i.ID)
	if found {
		return itemInterface.(item)
	}

	// Get from MySQL
	db := connectToSQL()

	error := db.QueryRow("SELECT * FROM items WHERE id = ? LIMIT 1", i.ID).Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.Name, &i.Desc, &i.Link, &i.Source, &i.SalesRank, &i.Images, &i.ProductGroup, &i.ProductTypeName)
	if error == sql.ErrNoRows {

		// Get from amazon
		i = i.getFromAmazon()
		i = i.saveToMysql()
		i = i.saveToMemcache()

	} else if error != nil {
		fmt.Println(error)
	}

	i.saveToMemcache()

	return i
}

func (i item) getFromMemcache() item {
	c := cache.New(5*time.Minute, 10*time.Minute)
	foo, found := c.Get("foo")
	if found {
		return foo.(item)
	}
	return i.getFromMysql()
}

func (i item) getFromMysql() item {

	return i
}

func (i item) getFromAmazon() item {

	client, error := amazon.NewFromEnvionment()
	if error != nil {
		log.Fatal(error)
	}

	res, error := client.ItemLookup(amazon.ItemLookupParameters{
		ResponseGroups: []amazon.ItemLookupResponseGroup{
			amazon.ItemLookupResponseGroupLarge,
		},
		IDType:  amazon.IDTypeASIN,
		ItemIDs: []string{i.ID},
	}).Do()
	if error != nil {
		log.Fatal(error)
	}

	for _, item := range res.Items.Item {
		//fmt.Printf("%# v", pretty.Formatter(item))
		i.Name = item.ItemAttributes.Title
	}

	// Some presets
	date := time.Now().Format("2006-01-02")
	i.DateCreated = date
	i.DateUpdated = date
	i.Source = "1" //todo

	return i
}

func (i item) saveToMemcache() item {
	c := cache.New(5*time.Minute, 10*time.Minute)
	c.Set(i.ID, i, cache.DefaultExpiration)
	return i
}

func (i item) saveToMysql() item {

	db := connectToSQL()

	// todo, switch to query builder
	// Prepare statement for inserting data
	insert, error := db.Prepare("INSERT INTO items (id, dateCreated, dateUpdated, `name`, `desc`, link, source, salesRank, images, productGroup, productTypeName) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if error != nil {
		panic(error.Error())
	}
	defer insert.Close()

	// run query
	_, error = insert.Exec(i.ID, i.DateCreated, i.DateUpdated, i.Name, i.Desc, i.Link, i.Source, i.SalesRank, i.Images, i.ProductGroup, i.ProductTypeName)
	if error != nil {
		panic(error.Error())
	}

	return i
}

func (i item) importItems() bool {

	items := []string{
		"B00FLYWNYQ", "B00AZBIZTW", "B01M2CTKH4", "B01M14ATO0", "B002N5MHLK", "B00LBK7OSY", "B005I33OVG", "B01DO7Y1AK", "B004UQ40IS",
		"B009MIK21S", "B00005MEGJ", "B0087T6CAI", "B000OZ9VLU", "B007VTP62U", "B0018DVYUS", "B00GBUPUOY", "B00FL43S3G", "B00O3HN4TU",
		"B00186098I", "B006MHEFWY", "B0002HE13I", "B00004S1DB", "B000K0FGE0", "B0024YTD08", "B000GGTYC8", "B0015SBILG", "B00BCEK2LA",
		"B0007PN9ZQ", "B001GBCXFW", "B00817YWPS", "B006P64GK8", "B0018LNXTU", "B00070E8LA", "B0000224VG", "B0076NOGPY", "B004WMFNRW",
		"B004XC7K6S", "B0058Y83Z2", "B0009JKG9M", "B001W2CJX6", "B003GXF9OA", "B002A8JO48", "B001NCDE84", "B0071OUJDQ", "B000CNY6UK",
		"B0039YY2QM", "B003I85GT6", "B001U52C9Q", "B00023RYS6", "B00JLDM98I", "B00EDRGLL8", "B001N444I6", "B00063RWUM", "B000EJPDOK",
		"B001MA0QY2"}

	for _, id := range items {
		i.ID = id
		i.get()
	}

	return true
}
