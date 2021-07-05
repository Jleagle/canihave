package algolia

import (
	"github.com/Jleagle/canihave/pkg/config"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

const indexName = "products"

var (
	client = search.NewClient(config.AlgoliaAppID, config.AlgoliaSearch)
	index  = client.InitIndex(indexName)
)

type Product struct {
	ObjectID string `json:"objectID"`
	Name     string `json:"name"`
}

func (p Product) Save() error {

	_, err := index.SaveObject(p)
	return err
}

type Products []Product

func (p Products) Save() error {

	_, err := index.SaveObjects(p)
	return err
}
