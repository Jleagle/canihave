package memcache

import (
	"io"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

var memcacheConnection *memcache.Client

const (
	appKey string = "canihave-"
)

func GetMemcacheItem(key string) (item *memcache.Item, err error) {
	item, err = getMemcacheConnection().Get(appKey + key)
	if err == io.EOF {
		resetMemcacheConnection()
		return GetMemcacheItem(key)
	}
	return item, err
}

func GetMemcacheMulti(keys []string) (items map[string]*memcache.Item, err error) {

	var keys2 []string

	for _, v := range keys {
		keys2 = append(keys2, appKey+v)
	}
	return getMemcacheConnection().GetMulti(keys2)
}

func SetMemcacheItem(key string, value []byte) (err error) {
	return getMemcacheConnection().Set(&memcache.Item{Key: appKey + key, Value: value})
}

func DeleteMemcacheItem(key string) (err error) {
	return getMemcacheConnection().Delete(appKey + key)
}

func getMemcacheConnection() *memcache.Client {

	if memcacheConnection == nil {
		memcacheConnection = memcache.New(os.Getenv("CANIHAVE_MEMCACHE_DNS"))
	}
	return memcacheConnection
}

func resetMemcacheConnection() {
	memcacheConnection = nil
}
