package amazon

import (
	"fmt"
	"log"

	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

func GetItems(items []string) {
	client, err := amazon.NewFromEnvionment()
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.ItemLookup(amazon.ItemLookupParameters{
		ResponseGroups: []amazon.ItemLookupResponseGroup{
			amazon.ItemLookupResponseGroupLarge,
		},
		IDType:  amazon.IDTypeASIN,
		ItemIDs: items,
	}).Do()
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range res.Items.Item {
		fmt.Printf(`-------------------------------
[Title] %v
[URL]   %v
`, item.ItemAttributes.Title, item.DetailPageURL)
	}
}
