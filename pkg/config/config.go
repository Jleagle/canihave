package config

import (
	"os"
)

const (
	EnvLocal = "local"
	EnvProd  = "production"
)

var (
	AlgoliaSearch   = os.Getenv("CANIHAVE_ALGOLIA_API_KEY_SEARCH")
	Environment     = os.Getenv("CANIHAVE_ENV")
	Port            = os.Getenv("CANIHAVE_PORT")
	AmazonAPIKey    = os.Getenv("CANIHAVE_AWS_ACCESS_KEY_ID")
	AmazonAPISecret = os.Getenv("CANIHAVE_AWS_SECRET_ACCESS_KEY")
)

func IsLocal() bool {
	return Environment != EnvLocal
}

func IsProd() bool {
	return Environment == EnvProd
}
