package amazon

import (
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
	"log"
	"time"
	"fmt"
	"os"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/models"
)

var RateLimit <-chan time.Time

func getAmazonClient(region string) (client *amazon.Client) {

	client, err := amazon.New(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		location.GetAmazonTag(region),
		amazon.Region(region),
	)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func GetItemDetails(item models.Item) (*amazon.ItemLookupResponse, error) {

	client := getAmazonClient(item.Region)

	<-RateLimit

	resp, err := client.ItemLookup(amazon.ItemLookupParameters{
		ResponseGroups: []amazon.ItemLookupResponseGroup{
			amazon.ItemLookupResponseGroupMedium,
		},
		IDType:  amazon.IDTypeASIN,
		ItemIDs: []string{item.ID},
	}).Do()

	return resp, err
}

func GetSimilarItems(item models.Item) (*amazon.SimilarityLookupResponse) {

	client := getAmazonClient(item.Region)

	<-RateLimit

	resp, err := client.SimilarityLookup(amazon.SimilarityLookupParameters{
		ResponseGroups: []amazon.SimilarityLookupResponseGroup{
			amazon.SimilarityLookupResponseGroupLarge,
		},
		ItemIDs: []string{item.ID},
	}).Do()

	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func GetReviews(item models.Item) {

}

func TopInNode(node int) {

}

func GetNodeDetails(node int) {

	client, err := amazon.NewFromEnvionment()
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.BrowseNodeLookup(amazon.BrowseNodeLookupParameters{
		ResponseGroups: []amazon.BrowseNodeLookupResponseGroup{
			amazon.BrowseNodeLookupResponseGroupBrowseNodeInfo,
			amazon.BrowseNodeLookupResponseGroupNewReleases,
			amazon.BrowseNodeLookupResponseGroupMostGifted,
			amazon.BrowseNodeLookupResponseGroupTopSellers,
			amazon.BrowseNodeLookupResponseGroupMostWishedFor,
		},
		BrowseNodeID: "492352",
	}).Do()
	if err != nil {
		log.Fatal(err)
	}
	browseNode := res.BrowseNodes()[0]
	fmt.Printf("%v: %v\n", browseNode.ID, browseNode.Name)
}

func Search(search string) {

	client, err := amazon.NewFromEnvionment()
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.ItemSearch(amazon.ItemSearchParameters{
		SearchIndex:    amazon.SearchIndexMusic,
		ResponseGroups: []amazon.ItemSearchResponseGroup{amazon.ItemSearchResponseGroupLarge},
		Keywords:       "Pat Metheny",
	}).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d results found\n\n", res.Items.TotalResults)
	for _, item := range res.Items.Item {
		fmt.Printf(`-------------------------------
[Title] %v
[URL]   %v
`, item.ItemAttributes.Title, item.DetailPageURL)
	}
}
