package amazon

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

var RateLimit <-chan time.Time

var retries int = 0

func SetRateLimit() {

	RateLimit = time.Tick(time.Millisecond * 1000) // Amazon has 1 request per second limit
}

func getAmazonClient(region string) (client *amazon.Client) {

	<-RateLimit

	client, err := amazon.New(
		os.Getenv("CANIHAVE_AWS_ACCESS_KEY_ID"),
		os.Getenv("CANIHAVE_AWS_SECRET_ACCESS_KEY"),
		location.GetAmazonTag(region),
		amazon.Region(region),
	)
	if err != nil {
		logger.Err("Can't create Amazon client: " + err.Error())
	}

	return client
}

func GetItemDetails(id string, region string) (resp *amazon.ItemLookupResponse, err error) {

	client := getAmazonClient(region)

	resp, err = client.ItemLookup(amazon.ItemLookupParameters{
		ResponseGroups: []amazon.ItemLookupResponseGroup{
			amazon.ItemLookupResponseGroupMedium,
			amazon.ItemLookupResponseGroupReviews,
			amazon.ItemLookupResponseGroupEditorialReview,
		},
		IDType:  amazon.IDTypeASIN,
		ItemIDs: []string{id},
	}).Do()

	if err != nil && strings.Contains(err.Error(), "RequestThrottled") {

		return GetItemDetails(id, region)
	}

	return resp, err
}

func GetSimilarItems(id string, region string) (resp *amazon.SimilarityLookupResponse, err error) {

	client := getAmazonClient(region)

	resp, err = client.SimilarityLookup(amazon.SimilarityLookupParameters{
		ResponseGroups: []amazon.SimilarityLookupResponseGroup{
			amazon.SimilarityLookupResponseGroupLarge,
		},
		ItemIDs: []string{id},
	}).Do()

	if err != nil && strings.Contains(err.Error(), "RequestThrottled") {

		return GetSimilarItems(id, region)
	}

	return resp, err
}

func GetNodeDetails(node string, region string) {

	client := getAmazonClient(region)

	res, err := client.BrowseNodeLookup(amazon.BrowseNodeLookupParameters{
		ResponseGroups: []amazon.BrowseNodeLookupResponseGroup{
			amazon.BrowseNodeLookupResponseGroupBrowseNodeInfo,
			amazon.BrowseNodeLookupResponseGroupNewReleases,
			amazon.BrowseNodeLookupResponseGroupMostGifted,
			amazon.BrowseNodeLookupResponseGroupTopSellers,
			amazon.BrowseNodeLookupResponseGroupMostWishedFor,
		},
		BrowseNodeID: node,
	}).Do()

	if err != nil {
		log.Fatal(err) // Remove
	}

	browseNode := res.BrowseNodes()[0]
	fmt.Printf("%v: %v\n", browseNode.ID, browseNode.Name)
}

func Search(search string, region string) {

	client := getAmazonClient(region)

	res, err := client.ItemSearch(amazon.ItemSearchParameters{
		SearchIndex: amazon.SearchIndexAll,
		ResponseGroups: []amazon.ItemSearchResponseGroup{
			amazon.ItemSearchResponseGroupLarge,
		},
		Keywords: search,
	}).Do()

	if err != nil {
		log.Fatal(err) // Remove
	}

	fmt.Printf("%d results found\n\n", res.Items.TotalResults)
}
