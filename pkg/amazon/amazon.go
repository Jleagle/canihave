package amazon

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Jleagle/canihave/pkg/config"
	"github.com/Jleagle/canihave/pkg/helpers"
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

var RateLimit = time.Tick(time.Millisecond * 1000)

var amazonConnection *amazon.Client

func getAmazonClient(region amazon.Region) (*amazon.Client, error) {

	<-RateLimit

	var err error

	if amazonConnection == nil {

		amazonConnection, err = amazon.New(
			config.AmazonAPIKey,
			config.AmazonAPISecret,
			location.GetAmazonTag(region),
			region,
		)
		if err != nil {
			logger.Err("Can't create Amazon client", err)
		}
	}

	return amazonConnection, err
}

func GetItemDetails(ids []string, region amazon.Region) (resp *amazon.ItemLookupResponse, err error) {

	ids = helpers.RemoveDuplicatesUnordered(ids)

	if len(ids) > 10 {
		return nil, errors.New("you can only query 10 items at a time")
	}

	client, err := getAmazonClient(region)
	if err != nil {
		return nil, err
	}

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
func GetItemDetailsBulk(ids []string, region amazon.Region) (ret []amazon.Item) {

	// Chunk the IDs into 10s
	var items [][]string
	var subItems []string

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

func GetSimilarItems(id string, region amazon.Region) (resp *amazon.SimilarityLookupResponse, err error) {

	client, err := getAmazonClient(region)
	if err != nil {
		return nil, err
	}

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

func GetNodeDetails(node string, region amazon.Region) (nil, err error) {

	client, err := getAmazonClient(region)
	if err != nil {
		return nil, err
	}

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
	return nil, err
}

func Search(search string, region amazon.Region) (nil, err error) {

	client, err := getAmazonClient(region)
	if err != nil {
		return nil, err
	}

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
	return nil, err
}
