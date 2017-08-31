package location

import (
	"github.com/oschwald/geoip2-golang"
	"net"
	"fmt"
	"log"
	"net/http"
)

func getISO(r *http.Request) string {

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

func GetAmazonRegion(w http.ResponseWriter, r *http.Request) string {

	var ret string

	var cookie, err = r.Cookie("flag")
	if err == nil {
		ret = cookie.Value
	} else {

		iso := getISO(r)

		switch iso {
		case "BR":
		case "CA":
		case "CN":
		case "DE":
		case "ES":
		case "FR":
		case "IN":
		case "IT":
		case "JP":
		case "MX":
			ret = iso
		case "UK":
		case "GB":
			ret = "UK"
		}
		ret = "US"

		// Set cookie
		SetCookie(w, ret)
	}

	return ret
}

func SetCookie(w http.ResponseWriter, flag string) {
	cookie := &http.Cookie{
		Name:     "flag",
		Value:    flag,
		HttpOnly: false,
		MaxAge:   0,
	}
	http.SetCookie(w, cookie)
}
