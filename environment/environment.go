package environment

import "os"

var environment string

func SetEnv() {

	environment = os.Getenv("CANIHAVE_ENV")
}

func GetEnv() string {
	return environment
}

func IsLocal() bool {
	return environment != "local"
}

func IsLive() bool {
	return environment == "live"
}
