package store

import "github.com/bradfitz/gomemcache/memcache"

var memcacheConnection *memcache.Client

const (
	MEMCACHE_APP_KEY string = "canihave-"
)

func getMemcacheConnection() *memcache.Client {

	if memcacheConnection == nil {
		memcacheConnection = memcache.New("127.0.0.1:11211")
	}
	return memcacheConnection
}

func GetMemcacheItem(key string) (item *memcache.Item, err error) {
	return getMemcacheConnection().Get(MEMCACHE_APP_KEY + key)
}

func GetMemcacheMulti(keys []string) (items map[string]*memcache.Item, err error) {

	for k, v := range keys {
		keys[k] = MEMCACHE_APP_KEY + v
	}
	return getMemcacheConnection().GetMulti(keys)
}

func SetMemcacheItem(key string, value []byte) (err error) {
	return getMemcacheConnection().Set(&memcache.Item{Key: MEMCACHE_APP_KEY + key, Value: value})
}

func DeleteMemcacheItem(key string) (err error) {
	return getMemcacheConnection().Delete(MEMCACHE_APP_KEY + key)
}
