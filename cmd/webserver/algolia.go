package main

import (
	"fmt"

	config2 "github.com/Jleagle/canihave/pkg/config"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func searchx() {

	client := search.NewClient("BYC8V5K9TT", config2.AlgoliaSearch)

	index := client.InitIndex("products")

	fmt.Println(index)

}
