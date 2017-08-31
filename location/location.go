package location

import (
	"github.com/oschwald/geoip2-golang"
	"net"
	"fmt"
	"log"
	"net/http"
)

func getISO(r *http.Request) string {

	// todo, wrap a cooking check around this
	// Connect to database
	db, err := geoip2.Open("location/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// If you are using strings that may be invalid, check that header is not nil
	header := r.Header.Get("x-forwarded-for")
	if header == "" {
		return ""
	}
	fmt.Printf("Looking up location for: %v", header)
	ip := net.ParseIP(header)
	fmt.Printf("Looking up location for: %v", ip)

	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	return record.Country.IsoCode
}

func GetAmazonRegion(r *http.Request) string {

	switch getISO(r) {
	case "UK":
	case "GB":
		return "UK"
	}
	return "US"
}
