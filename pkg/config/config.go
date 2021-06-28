package config

import (
	"os"
)

const (
	EnvLocal = "local"
	EnvProd  = "production"
)

var (
	Environment             = os.Getenv("CANIHAVE_ENV")
	Port                    = os.Getenv("CANIHAVE_PORT")
	Memcache                = os.Getenv("CANIHAVE_MEMCACHE_DNS")
	FacebookToken           = os.Getenv("CANIHAVE_FACEBOOK_TOKEN")
	AlgoliaSearch           = os.Getenv("CANIHAVE_ALGOLIA_API_KEY_SEARCH")
	AlgoliaAppID            = os.Getenv("CANIHAVE_ALGOLIA_API_ID")
	AmazonAPIKey            = os.Getenv("CANIHAVE_AWS_ACCESS_KEY_ID")
	AmazonAPISecret         = os.Getenv("CANIHAVE_AWS_SECRET_ACCESS_KEY")
	TwitterConsumerKey      = os.Getenv("CANIHAVE_TWITTER_CONSUMER_KEY")
	TwitterConsumerSecret   = os.Getenv("CANIHAVE_TWITTER_CONSUMER_SECRET")
	TwitterAccessToken      = os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN")
	TwitterAcessTokenSecret = os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN_SECRET")
	MySQLUser               = os.Getenv("CANIHAVE_MYSQL_USER")
	MySQLPass               = os.Getenv("CANIHAVE_MYSQL_PASS")
	MySQLIP                 = os.Getenv("CANIHAVE_MYSQL_IP")
	MySQLDB                 = os.Getenv("CANIHAVE_MYSQL_DB")
)

func IsLocal() bool {
	return Environment != EnvLocal
}

func IsProd() bool {
	return Environment == EnvProd
}
