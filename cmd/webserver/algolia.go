package main

import (
	"fmt"

	"github.com/Jleagle/canihave/pkg/config"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func searchx() {

	client := search.NewClient(config.AlgoliaAppID, config.AlgoliaSearch)

	index := client.InitIndex("products")

	fmt.Println(index)
}
