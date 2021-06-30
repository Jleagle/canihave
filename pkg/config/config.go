package config

import (
	"os"
)

const (
	EnvLocal = "local"
	EnvProd  = "production"
)

var (
	AlgoliaAppID            = os.Getenv("CANIHAVE_ALGOLIA_API_ID")
	AlgoliaSearch           = os.Getenv("CANIHAVE_ALGOLIA_API_KEY_SEARCH")
	AmazonAPIKey            = os.Getenv("CANIHAVE_AWS_ACCESS_KEY_ID")
	AmazonAPISecret         = os.Getenv("CANIHAVE_AWS_SECRET_ACCESS_KEY")
	Environment             = os.Getenv("CANIHAVE_ENV")
	FacebookToken           = os.Getenv("CANIHAVE_FACEBOOK_TOKEN")
	Memcache                = os.Getenv("CANIHAVE_MEMCACHE_DNS")
	MySQLDB                 = os.Getenv("CANIHAVE_MYSQL_DB")
	MySQLIP                 = os.Getenv("CANIHAVE_MYSQL_IP")
	MySQLPass               = os.Getenv("CANIHAVE_MYSQL_PASS")
	MySQLUser               = os.Getenv("CANIHAVE_MYSQL_USER")
	Port                    = os.Getenv("CANIHAVE_PORT")
	TwitterAccessToken      = os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN")
	TwitterAcessTokenSecret = os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN_SECRET")
	TwitterConsumerKey      = os.Getenv("CANIHAVE_TWITTER_CONSUMER_KEY")
	TwitterConsumerSecret   = os.Getenv("CANIHAVE_TWITTER_CONSUMER_SECRET")
)

func IsLocal() bool {
	return Environment != EnvLocal
}

func IsProd() bool {
	return Environment == EnvProd
}
