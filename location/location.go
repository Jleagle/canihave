package location

import (
	"github.com/oschwald/geoip2-golang"
	"net"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
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

func GetAmazonRegion(w http.ResponseWriter, r *http.Request) (region string) {

	var value string
	var cookie, err = r.Cookie("region")
	if err == nil && cookie.Value != "" {
		value = cookie.Value
	} else {

		iso := getISO(r)

		switch iso {
		case "BR", "CA", "CN", "DE", "ES", "FR", "IN", "IT", "JP", "MX":
			value = iso
		case "UK", "GB":
			value = UK
		default:
			value = US
		}

		setCookie(w, amazon.Region(value))
	}

	return value
}

func setCookie(w http.ResponseWriter, region amazon.Region) {

	if region.IsValid() {
		cookie := &http.Cookie{
			Name:     "region",
			Value:    string(region),
			HttpOnly: false,
			MaxAge:   0,
		}
		http.SetCookie(w, cookie)
	}
}

func SetAmazonEnviromentVars(region string) {

	os.Setenv("AWS_PRODUCT_REGION", region)

	switch region {
	case BR:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case CA:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case CN:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case DE:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case ES:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case FR:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case IN:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case IT:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case JP:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case MX:
		os.Setenv("AWS_ASSOCIATE_TAG", "")
	case UK:
		os.Setenv("AWS_ASSOCIATE_TAG", "canihaveone00-21")
	case US:
		os.Setenv("AWS_ASSOCIATE_TAG", "canihaveone-20")
	}
}

func TLDToRegion(tld string) string {
	switch tld {
	case "br", "ca", "cn", "de", "es", "fr", "in", "it", "jp", "mx", "uk":
		return strings.ToUpper(tld)
	}

	return US
}

func GetCurrency(region string) string {

	switch region {
	case BR:
		return "R$"
	case CA, US, MX:
		return "$"
	case CN:
		return "¥"
	case DE, ES, FR, IT:
		return "€"
	case IN:
		return "₹"
	case JP:
		return "¥"
	case UK:
		return "£"
	}

	return ""
}

func DetectLanguageChange(w http.ResponseWriter, r *http.Request) {

	region := r.URL.Query().Get("region")
	if region != "" {
		setCookie(w, amazon.Region(region))
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}
}
