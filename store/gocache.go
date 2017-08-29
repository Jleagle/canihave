package store

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var c *cache.Cache

func GetGoCache() *cache.Cache {

	if c == nil {

		c = cache.New(5*time.Minute, 10*time.Minute)
		return c
	}

	return c
}
