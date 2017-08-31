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

func CacheFunc(key string, fn CacheItem) (interface{}) {

	c := GetGoCache()

	val, exists := c.Get(key)
	if exists {
		return val
	}

	val = fn()
	c.Set(key, val, cache.DefaultExpiration)
	return val
}

type CacheItem func() []interface{}
