package location

import (
	"github.com/oschwald/geoip2-golang"
	"net"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	BR string = "BR"
	CA string = "CA"
	CN string = "CN"
	DE string = "DE"
	ES string = "ES"
	FR string = "FR"
	IN string = "IN"
	IT string = "IT"
	JP string = "JP"
	MX string = "MX"
	UK string = "UK"
	US string = "US"
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
	ip := net.ParseIP(header)

	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Geo Lite Lookup: %v - %v", ip, record.Country.IsoCode)
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
			ret = UK
		default:
			ret = US
		}

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

func SetAmazonEnviromentVars(region string) {

	os.Setenv("AWS_PRODUCT_REGION", region)

	switch region {
	case "BR":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "CA":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "CN":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "DE":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "ES":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "FR":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "IN":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "IT":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "JP":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "MX":
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case "UK":
		os.Setenv("AWS_ASSOCIATE_TAG", "canihaveone00-21")
	case "US":
		os.Setenv("AWS_ASSOCIATE_TAG", "canihaveone-20")
	}
}

func TLDToRegion(tld string) string {
	switch tld {
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
	case "uk":
		return strings.ToUpper(tld)
	}

	return US
}

func GetCurrency(region string) string {

	switch region {
	case "BR":
		return "R$"
	case "CA":
	case "US":
	case "MX":
		return "$"
	case "CN":
		return "¥"
	case "DE":
	case "ES":
	case "FR":
	case "IT":
		return "€"
	case "IN":
		return "₹"
	case "JP":
		return "¥"
	case "uk":
		return "£"
	}

	return ""
}
