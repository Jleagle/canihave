package store

import "github.com/bradfitz/gomemcache/memcache"

var memcacheConnection *memcache.Client

func getMemcacheConnection() *memcache.Client {

	if memcacheConnection == nil {
		memcacheConnection = memcache.New("127.0.0.1:11211")
	}
	return memcacheConnection
}

func GetMemcacheItem(key string) (item *memcache.Item, err error) {
	return getMemcacheConnection().Get("canihave-" + key)
}

func SetMemcacheItem(key string, value []byte) (err error) {
	return getMemcacheConnection().Set(&memcache.Item{Key: "canihave-" + key, Value: value})
}

func GetMemcacheMulti(keys []string) (items map[string]*memcache.Item, err error) {

	for k, v := range keys {
		keys[k] = "canihave-" + v
	}
	return getMemcacheConnection().GetMulti(keys)
}
