package models

import (
	"fmt"
	"log"
	"strings"

	"time"

	amaz "github.com/Jleagle/canihave/amazon"
	"github.com/Jleagle/canihave/store"
)

func getSimilar(id string, region string) {

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

		fmt.Println("Saving " + item.ID + " as a similar item")
		item.saveAsNewMysqlRow()
		item.saveToMemcache()

		// Save similar-link
		date := time.Now().Format("2006-01-02 15:04:05")

		conn := store.GetMysqlConnection()
		_, err := conn.Query("INSERT INTO relation (id, related_id, date_created, type) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE type=type", id, item.ID, date, RELATION_TYPE_SIMILAR)
		if err != nil {
			panic(err.Error())
		}
	}
}
