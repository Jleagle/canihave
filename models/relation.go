package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	amaz "github.com/Jleagle/canihave/amazon"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
)

type Relation struct {
	ID          string
	RelatedID   string
	DateCreated string
	Type        string
}

func findSimilar(id string, region string) {

	similar, err := amaz.GetSimilarItems(id, region)

	if err != nil && strings.Contains(err.Error(), "AWS.ECommerceService.NoSimilarities") {
		return
	} else if err != nil {
		log.Fatal(err)
	}

	for _, amazonItem := range similar.Items.Item {

		// Save item
		var price int = 0
		if amazonItem.ItemAttributes.ListPrice.Amount != "" {
			price, err = strconv.Atoi(amazonItem.ItemAttributes.ListPrice.Amount)
			if err != nil {
				log.Fatal("Error converting string to int")
			}
		}

		item := Item{}
		item.ID = amazonItem.ASIN
		item.Name = amazonItem.ItemAttributes.Title
		item.Link = amazonItem.DetailPageURL
		item.SalesRank = amazonItem.SalesRank
		item.Photo = amazonItem.LargeImage.URL
		item.Node = "0" //todo
		item.NodeName = amazonItem.ItemAttributes.ProductGroup
		item.Price = price
		item.CompanyName = amazonItem.ItemAttributes.Manufacturer
		item.Type = TYPE_SIMILAR
		item.Region = region

		item.saveToMysql()
		item.saveToMemcache()

		// Save the relation
		builder := squirrel.Insert("relations")
		builder = builder.Columns("id", "relatedId", "dateCreated", "type")
		builder = builder.Values(id, item.ID, time.Now().Unix(), TYPE_SIMILAR)

		err := store.Insert(builder)

		if sqlerr, ok := err.(*mysql.MySQLError); ok {
			if sqlerr.Number == 1062 { // Duplicate entry
				continue
			}
		}

		if err != nil {
			logger.Err("Can't insert related item: " + err.Error())
		}
	}
}

func findNodeitems(node string, region string) {

	//nodeItems, err := amaz.GetNodeDetails(node, region)
}

func (i Item) GetSimilar() (items []Item) {

	// Get relations
	builder := squirrel.Select("relatedId").From("relations").Where("id = ? AND type = ?", i.ID, TYPE_SIMILAR)
	rows := store.Query(builder)
	defer rows.Close()

	relation := Relation{}
	relations := []Relation{}

	for rows.Next() {

		err := rows.Scan(&relation.RelatedID)
		if err != nil {
			fmt.Println(err)
		}

		relations = append(relations, relation)
	}

	if len(relations) < 1 {
		return []Item{}
	}

	for _, v := range relations {

		item := Item{}
		item.ID = v.RelatedID
		item.Get()

		items = append(items, item)
	}

	return items
}
