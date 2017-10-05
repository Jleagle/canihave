package store

import (
	"io"

	"github.com/bradfitz/gomemcache/memcache"
)

var memcacheConnection *memcache.Client

const (
	MEMCACHE_APP_KEY string = "canihave-"
)

func GetMemcacheItem(key string) (item *memcache.Item, err error) {
	item, err = getMemcacheConnection().Get(MEMCACHE_APP_KEY + key)
	if err == io.EOF {
		resetMemcacheConnection()
		return GetMemcacheItem(key)
	}
	return item, err
}

func GetMemcacheMulti(keys []string) (items map[string]*memcache.Item, err error) {

	keys2 := []string{}

	for _, v := range keys {
		keys2 = append(keys2, MEMCACHE_APP_KEY+v)
	}
	return getMemcacheConnection().GetMulti(keys2)
}

func SetMemcacheItem(key string, value []byte) (err error) {
	return getMemcacheConnection().Set(&memcache.Item{Key: MEMCACHE_APP_KEY + key, Value: value})
}

func DeleteMemcacheItem(key string) (err error) {
	return getMemcacheConnection().Delete(MEMCACHE_APP_KEY + key)
}

func getMemcacheConnection() *memcache.Client {

	if memcacheConnection == nil {
		memcacheConnection = memcache.New("127.0.0.1:11211")
	}
	return memcacheConnection
}

func resetMemcacheConnection() {
	memcacheConnection = nil
}
