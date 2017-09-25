package location

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Jleagle/canihave/logger"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
	"github.com/oschwald/geoip2-golang"
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

var regions map[string]string

func SetRegions() {

	regions = map[string]string{
		US: "United States",
		UK: "United Kingdom",
		//location.DE: "Deutschland",
		//location.FR: "France",
		//location.JP: "Japan",
		//location.CA: "Canada",
		//location.CN: "China",
		//location.IT: "Italia",
		//location.ES: "España",
		//location.IN: "India",
		//location.BR: "Brazil",
		//location.MX: "Mexico",
	}
}

func GetRegions() map[string]string {
	return regions
}

func IsValidRegion(region string) (result bool) {
	reg := amazon.Region(region)
	return reg.IsValid()
}

func getISO(r *http.Request) string {

	// Connect to database
	file, err := geoip2.Open("location/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// If you are using strings that may be invalid, check that header is not nil
	header := r.Header.Get("x-forwarded-for")
	if header == "" {
		return ""
	}
	ip := net.ParseIP(header)

	record, err := file.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	logger.Err("Geo Lite Lookup: "+string(ip)+" - "+record.Country.IsoCode, err)
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

		setCookie(w, value)
	}

	return value
}

func setCookie(w http.ResponseWriter, region string) {

	if IsValidRegion(region) {
		cookie := &http.Cookie{
			Name:     "region",
			Value:    region,
			HttpOnly: false,
			MaxAge:   0,
		}
		http.SetCookie(w, cookie)
	}
}

func GetAmazonTag(region string) (tag string) {

	switch region {
	case BR:
		return ""
	case CA:
		return ""
	case CN:
		return ""
	case DE:
		return ""
	case ES:
		return ""
	case FR:
		return ""
	case IN:
		return ""
	case IT:
		return ""
	case JP:
		return ""
	case MX:
		return ""
	case UK:
		return "canihaveone00-21"
	}

	return "canihaveone-20"
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
		setCookie(w, region)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}
}
