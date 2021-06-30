package amazon

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Jleagle/canihave/pkg/config"
	"github.com/Jleagle/canihave/pkg/helpers"
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
	"go.uber.org/zap"
)

var (
	clients    = map[amazon.Region]*amazon.Client{}
	clientRate = time.Tick(time.Millisecond * 1000)
	clientLock sync.Mutex
)

func getAmazonClient(region amazon.Region) (*amazon.Client, error) {

	clientLock.Lock()
	defer clientLock.Unlock()

	if val, ok := clients[region]; ok {
		return val, nil
	}

	client, err := amazon.New(
		config.AmazonAPIKey,
		config.AmazonAPISecret,
		location.GetAmazonTag(region),
		region,
	)
	if err != nil {
		return nil, err
	}

	clients[region] = client

	return client, err
}

func GetItemDetails(ids []string, region amazon.Region) (resp *amazon.ItemLookupResponse, err error) {

	ids = helpers.RemoveDuplicates(ids)

	if len(ids) > 10 {
		return nil, errors.New("you can only query 10 items at a time")
	}

	<-clientRate

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

		logger.Logger.Info("Retrying amazon API call because RequestThrottled")
		return GetItemDetails(ids, region)
	}

	return resp, err
}

func GetItems(ids []string, region amazon.Region) (resp *amazon.ItemLookupResponse, err error) {

	ids = helpers.RemoveDuplicates(ids)

	if len(ids) > 10 {
		return nil, errors.New("you can only query 10 items at a time")
	}

	<-clientRate

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

		logger.Logger.Info("Retrying amazon API call because RequestThrottled")
		return GetItems(ids, region)
	}

	return resp, err
}

func GetSimilarItems(id string, region amazon.Region) (resp *amazon.SimilarityLookupResponse, err error) {

	<-clientRate

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

func GetNode(node string, region amazon.Region) (nil, err error) {

	<-clientRate

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
		logger.Logger.Fatal("", zap.Error(err))
	}

	browseNode := res.BrowseNodes()[0]
	fmt.Printf("%v: %v\n", browseNode.ID, browseNode.Name)
	return nil, err
}

func SearchItems(search string, region amazon.Region) (nil, err error) {

	<-clientRate

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
		logger.Logger.Fatal("", zap.Error(err))
	}

	fmt.Printf("%d results found\n\n", res.Items.TotalResults)
	return nil, err
}
