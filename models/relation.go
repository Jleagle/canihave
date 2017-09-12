package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	amaz "github.com/Jleagle/canihave/amazon"
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

const (
	RELATION_TYPE_SIMILAR   string = "similar"
	RELATION_TYPE_SAME_NODE string = "node"
)

func findSimilar(id string, region string) {

	similar, err := amaz.GetSimilarItems(id, region)

	if err != nil && strings.Contains(err.Error(), "AWS.ECommerceService.NoSimilarities") {
		return
	} else if err != nil {
		log.Fatal(err)
	}

	for _, amazonItem := range similar.Items.Item {

		// Save item
		price, _ := strconv.Atoi(amazonItem.ItemAttributes.ListPrice.Amount)

		item := Item{}

		item.ID = amazonItem.ASIN
		item.Name = amazonItem.ItemAttributes.Title
		item.Link = amazonItem.DetailPageURL
		item.SalesRank = amazonItem.SalesRank
		item.Photo = amazonItem.LargeImage.URL
		item.Node = "0" //todo
		item.NodeName = amazonItem.ItemAttributes.ProductGroup
		item.Price = price
		item.CompanyName = "" //todo
		item.Status = ""
		item.Type = TYPE_SIMILAR
		item.Region = region

		item.saveToMysql()
		item.saveToMemcache()

		// Save the relation
		builder := squirrel.Insert("relations")
		builder = builder.Columns("id", "relatedId", "dateCreated", "type")
		builder = builder.Values(id, item.ID, time.Now().Unix(), RELATION_TYPE_SIMILAR)

		_, err := store.Insert(builder)

		if sqlerr, ok := err.(*mysql.MySQLError); ok {
			if sqlerr.Number == 1062 { // Duplicate entry
				continue
			}
		}

		if err != nil {
			panic(err.Error())
		}
	}
}

func (i Item) GetSimilar() (items []Item) {

	// Get relations
	builder := squirrel.Select("relatedId").From("relations").Where("id = ? AND type = ?", i.ID, RELATION_TYPE_SIMILAR)
	rows := store.Query(builder)

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
