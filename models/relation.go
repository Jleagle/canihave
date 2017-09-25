package models

import (
	"log"
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

func saveSimilarItems(id string, region string, itemType string) {

	similar, err := amaz.GetSimilarItems(id, region)

	if err != nil && strings.Contains(err.Error(), "AWS.ECommerceService.NoSimilarities") {
		return
	} else if err != nil {
		log.Fatal(err)
	}

	for _, amazonItem := range similar.Items.Item {

		item := amazonItemToItem(amazonItem, itemType, region)

		saveToMysql(item)
		saveToMemcache(item)

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
			logger.Err("Can't insert related item", err)
		}
	}
}

func saveNodeitems(node string, region string) {

	//nodeItems, err := amaz.GetNodeDetails(node, region)
}

func (i Item) GetRelated(itemType string) (items []Item) {

	builder := squirrel.Select("relatedId").From("relations").Where("id = ? AND type = ?", i.ID, itemType)
	rows := store.Query(builder)
	defer rows.Close()

	ids := []string{}
	for rows.Next() {

		var id string
		err := rows.Scan(&id)
		if err != nil {
			logger.Err("Can't scan related item", err)
		}
		ids = append(ids, id)
	}

	return GetMulti(ids, i.Region, itemType)
}
