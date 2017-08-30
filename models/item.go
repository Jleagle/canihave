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
	ID              string
	DateCreated     string
	DateUpdated     string
	Name            string
	Desc            string
	Link            string
	Source          string
	SalesRank       string
	Images          string
	ProductGroup    string
	ProductTypeName string
	Status          string
	Price           string
}

func (i Item) GetUKPixel() string {
	return "//ir-uk.amazon-adsystem.com/e/ir?t=canihaveone00-21&l=am2&o=2&a=" + i.ID
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
		i.Images = item.Images
		i.ProductGroup = item.ProductGroup
		i.ProductTypeName = item.ProductTypeName
	}
	return found
}

func (i *Item) getFromMysql() (found bool) {

	// Make the query
	query := squirrel.Select("*").From("items").Where("id = ?", i.ID).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	db := store.GetMysqlConnection()
	err = db.QueryRow(sql, args...).Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.Name, &i.Desc, &i.Link, &i.Source, &i.SalesRank, &i.Images, &i.ProductGroup, &i.ProductTypeName)
	if err != nil {
		return false
	}

	return true
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
		log.Fatal(err)
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
	i.Images = "[]"
	i.ProductGroup = amazonItem.ItemAttributes.ProductGroup
	i.ProductTypeName = amazonItem.ItemAttributes.ProductTypeName
	i.Price = amazonItem.ItemAttributes.ListPrice.FormattedPrice

	return true
}

func (i *Item) saveToMemcache() {

	x := store.GetGoCache()
	x.Set(i.ID, i, cache.DefaultExpiration)
}

func (i *Item) saveToMysql() {

	// Make query
	//sql, args, err := sq.Insert("items").Columns("name", "age").Values("moe", 13).Values("larry", sq.Expr("? + 5", 12)).ToSql()

	conn := store.GetMysqlConnection()

	// todo, switch to query builder
	// Prepare statement for inserting data
	insert, err := conn.Prepare("INSERT INTO items (id, dateCreated, dateUpdated, `name`, `desc`, link, source, salesRank, images, productGroup, productTypeName) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	// run query
	_, err = insert.Exec(i.ID, i.DateCreated, i.DateUpdated, i.Name, i.Desc, i.Link, i.Source, i.SalesRank, i.Images, i.ProductGroup, i.ProductTypeName)
	if err != nil {
		panic(err.Error())
	}
}

func ImportItems() bool {

	items := []string{
		"0735216207",
		"1433805618",
		"1501175564",
		"B004RMK4BC",
		"B004S8F7QM",
		"B00E1EN92W",
		"B00EB4ADQW",
		"B00GAC1D2G",
		"B00IOY8XWQ",
		"B00JM5GW10",
		"B00NB86OYE",
		"B00O4OR4GQ",
		"B00OQVZDJM",
		"B00P77ZAN8",
		"B00REQKWGA",
		"B00REQL3AE",
		"B00U3FPN4U",
		"B00U3FPN4U",
		"B00UT823WQ",
		"B00X4WHP5E",
		"B00ZV9PXP2",
		"B00ZV9RDKK",
		"B00ZV9RDKK",
		"B00ZV9RDKK",
		"B0186JAEWK",
		"B01BH83OOM",
		"B01DFTCV90",
		"B01E3QM34W",
		"B01EQYX9NU",
		"B01GEW27DA",
		"B01GEW27DA",
		"B01I499BNA",
		"B01J24C0TI",
		"B01J90MSDS",
		"B01J90MSDS",
		"B01J94SBEY",
		"B01J94SWWU",
		"B01KMSKNGU",
		"B01M71IUZ7",
		"B01MRG7T0D",
		"B06XDC9RBJ",
		"B071JRMKBH",
	}

	for _, id := range items {
		i := Item{}
		i.ID = id
		i.Source = "import"
		i.Get()
	}

	return true
}
