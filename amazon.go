package main

import (
	"fmt"
	"log"

	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

func getItems(items []string) {
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

func importItems() {

	items := []string{
		"B00FLYWNYQ",
		"B00AZBIZTW",
		"B01M2CTKH4",
		"B01M14ATO0",
		"B002N5MHLK",
		"B00LBK7OSY",
		"B005I33OVG",
		"B01DO7Y1AK",
		"B004UQ40IS",
		"B009MIK21S",
		"B00005MEGJ",
		"B0087T6CAI",
		"B000OZ9VLU",
		"B007VTP62U",
		"B0018DVYUS",
		"B00GBUPUOY",
		"B00FL43S3G",
		"B00O3HN4TU",
		"B00186098I",
		"B006MHEFWY",
		"B0002HE13I",
		"B00004S1DB",
		"B000K0FGE0",
		"B0024YTD08",
		"B000GGTYC8",
		"B0015SBILG",
		"B00BCEK2LA",
		"B0007PN9ZQ",
		"B001GBCXFW",
		"B00817YWPS",
		"B006P64GK8",
		"B0018LNXTU",
		"B00070E8LA",
		"B0000224VG",
		"B0076NOGPY",
		"B004WMFNRW",
		"B004XC7K6S",
		"B0058Y83Z2",
		"B0009JKG9M",
		"B001W2CJX6",
		"B003GXF9OA",
		"B002A8JO48",
		"B001NCDE84",
		"B0071OUJDQ",
		"B000CNY6UK",
		"B0039YY2QM",
		"B003I85GT6",
		"B001U52C9Q",
		"B00023RYS6",
		"B00JLDM98I",
		"B00EDRGLL8",
		"B001N444I6",
		"B00063RWUM",
		"B000EJPDOK",
		"B001MA0QY2",
	}

	// Connect to SQL
	db, _ := connectToSQL()
	defer db.Close()

	// Prepare statement for inserting data
	insert, err := db.Prepare("INSERT INTO items (id, date_created, date_updated, `name`, `desc`, source) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	for _, id := range items {

		_, err = insert.Exec(id, "2010-10-10", "2010-10-10", id, id, "source")
		if err != nil {
			panic(err.Error())
		}
	}
}
