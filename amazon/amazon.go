package amazon

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Jleagle/canihave/location"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

var RateLimit <-chan time.Time

func getAmazonClient(region string) (client *amazon.Client) {

	<-RateLimit

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

func GetItemDetails(id string, region string) (*amazon.ItemLookupResponse, error) {

	client := getAmazonClient(region)

	resp, err := client.ItemLookup(amazon.ItemLookupParameters{
		ResponseGroups: []amazon.ItemLookupResponseGroup{
			amazon.ItemLookupResponseGroupMedium,
			amazon.ItemLookupResponseGroupReviews,
			amazon.ItemLookupResponseGroupEditorialReview,
		},
		IDType:  amazon.IDTypeASIN,
		ItemIDs: []string{id},
	}).Do()

	return resp, err
}

func GetSimilarItems(id string, region string) *amazon.SimilarityLookupResponse {

	client := getAmazonClient(region)

	resp, err := client.SimilarityLookup(amazon.SimilarityLookupParameters{
		ResponseGroups: []amazon.SimilarityLookupResponseGroup{
			amazon.SimilarityLookupResponseGroupLarge,
		},
		ItemIDs: []string{id},
	}).Do()

	if err != nil {
		log.Fatal(err)
	}

	return resp
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
		log.Fatal(err)
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
		log.Fatal(err)
	}

	fmt.Printf("%d results found\n\n", res.Items.TotalResults)
}
