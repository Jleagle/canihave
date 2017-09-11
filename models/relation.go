package models

import (
	"fmt"
	"log"
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

func saveSimilar(id string, region string) {

	similar, err := amaz.GetSimilarItems(id, region)

	if err != nil && strings.Contains(err.Error(), "AWS.ECommerceService.NoSimilarities") {
		return
	} else if err != nil {
		log.Fatal(err)
	}

	for _, amazonItem := range similar.Items.Item {

		// Save item
		item := Item{}
		item.ID = amazonItem.ASIN
		amazonItemToItem(&item, amazonItem)
		item.Type = TYPE_SIMILAR

		item.saveAsNewMysqlRow()
		item.saveToMemcache()

		// Save the relation
		date := time.Now().Format("2006-01-02 15:04:05")

		builder := squirrel.Insert("relations")
		builder = builder.Columns("id", "related_id", "date_created", "type")
		builder = builder.Values(id, item.ID, date, RELATION_TYPE_SIMILAR)

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
	builder := squirrel.Select("related_id").From("relations").Where("id = ? AND type = ?", i.ID, RELATION_TYPE_SIMILAR)
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
