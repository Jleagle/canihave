package amazon

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Jleagle/canihave/helpers"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

var RateLimit <-chan time.Time

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
		logger.Err("Can't create Amazon client", err)
	}

	return client
}

func GetItemDetails(ids []string, region string) (resp *amazon.ItemLookupResponse, err error) {

	ids = helpers.RemoveDuplicatesUnordered(ids)

	if len(ids) > 10 {
		return nil, errors.New("you can only query 10 items at a time")
	}

	client := getAmazonClient(region)

	resp, err = client.ItemLookup(amazon.ItemLookupParameters{
		ResponseGroups: []amazon.ItemLookupResponseGroup{
			amazon.ItemLookupResponseGroupMedium,
			amazon.ItemLookupResponseGroupReviews,
			amazon.ItemLookupResponseGroupEditorialReview,
		},
		IDType:  amazon.IDTypeASIN,
		ItemIDs: ids,
	}).Do()

	if err != nil && strings.Contains(err.Error(), "RequestThrottled") {

		logger.Info("Retrying amazon API call because RequestThrottled")
		return GetItemDetails(ids, region)
	}

	return resp, err
}

// Will return an error if any of the IDs fail
func GetItemDetailsBulk(ids []string, region string) (ret []amazon.Item) {

	// Chunk the IDs into 10s
	items := [][]string{}
	subItems := []string{}

	for _, v := range ids {

		subItems = append(subItems, v)

		if len(subItems) == 10 {
			items = append(items, subItems)
			subItems = []string{}
		}
	}

	if len(subItems) > 0 {
		items = append(items, subItems)
		subItems = []string{}
	}

	count := 0
	for _, subItem := range items {

		count += 10
		fmt.Print(count)

		amazonItems, err := GetItemDetails(subItem, region)
		if err != nil && strings.Contains(err.Error(), "AWS.InvalidParameterValue") {

			// One of the bunch failed
			r := regexp.MustCompile(`Value: ([A-Z0-9]{10}) is not`)
			links := r.FindAllString(err.Error(), 1)

			logger.Err("One item in a batch failed: "+links[1], err)

		} else if err != nil {

			// Fail
			logger.Err("Can't get amazon items", err)
		} else {

			// Success
			for _, v := range amazonItems.Items.Item {
				ret = append(ret, v)
			}
		}
	}

	return ret
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
