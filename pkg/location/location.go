package location

import (
	"errors"
	"net"
	"strings"
	"sync"

	"github.com/Jleagle/canihave/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
	"github.com/oschwald/maxminddb-golang"
)

const (
	cookieRegion = "region"
)

var regionNames = map[amazon.Region]string{
	amazon.RegionBrazil:  "Brazil",
	amazon.RegionCanada:  "Canada",
	amazon.RegionChina:   "China",
	amazon.RegionGermany: "Deutschland",
	amazon.RegionSpain:   "España",
	amazon.RegionFrance:  "France",
	amazon.RegionIndia:   "India",
	amazon.RegionItaly:   "Italia",
	amazon.RegionJapan:   "Japan",
	amazon.RegionMexico:  "Mexico",
	amazon.RegionUK:      "United Kingdom",
	amazon.RegionUS:      "United States",
}

var regionTags = map[amazon.Region]string{
	amazon.RegionBrazil:  "",
	amazon.RegionCanada:  "",
	amazon.RegionChina:   "",
	amazon.RegionGermany: "",
	amazon.RegionSpain:   "",
	amazon.RegionFrance:  "",
	amazon.RegionIndia:   "",
	amazon.RegionItaly:   "",
	amazon.RegionJapan:   "",
	amazon.RegionMexico:  "",
	amazon.RegionUK:      "canihaveone00-21",
	amazon.RegionUS:      "canihaveone-20",
}

func GetAmazonTag(region amazon.Region) string {
	return regionTags[region]
}

func GetRegions() []amazon.Region {
	return []amazon.Region{
		amazon.RegionBrazil,
		amazon.RegionCanada,
		amazon.RegionChina,
		amazon.RegionGermany,
		amazon.RegionSpain,
		amazon.RegionFrance,
		amazon.RegionIndia,
		amazon.RegionItaly,
		amazon.RegionJapan,
		amazon.RegionMexico,
		amazon.RegionUK,
		amazon.RegionUS,
	}
}

var (
	ErrInvalidIP = errors.New("invalid ip")

	maxMindLock sync.Mutex
	maxMindFile *maxminddb.Reader
)

func GetRegion(c *fiber.Ctx) (region amazon.Region) {

	if region := amazon.Region(c.Query("region")); region.IsValid() {
		return region
	}

	if region := amazon.Region(c.Cookies(cookieRegion)); region.IsValid() {
		return region
	}

	region = func() amazon.Region {

		// Local
		if config.IsLocal() {
			return amazon.RegionUK
		}

		// Get from cloudflare
		cf := strings.ToUpper(c.Get("cf-ipcountry"))
		if cf == "GB" {
			cf = "UK"
		}
		if region := amazon.Region(cf); region.IsValid() {
			return region
		}

		// Get from MaxMind
		if region, err := getRegionFromIP(c); err == nil {
			return region
		}

		return amazon.RegionUS
	}()

	c.Cookie(&fiber.Cookie{
		Name:     cookieRegion,
		Value:    string(region),
		MaxAge:   0,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "strict",
		// Path:     "",
		// Domain:   "",
		// Expires:  time.Time{},
	})

	return region
}

func getRegionFromIP(c *fiber.Ctx) (region amazon.Region, err error) {

	if c.IP() == "" {
		return "", ErrInvalidIP
	}

	maxMindLock.Lock()
	defer maxMindLock.Unlock()

	if maxMindFile == nil {
		maxMindFile, err = maxminddb.Open("location/GeoLite2-Country.mmdb")
		if err != nil {
			return "", err
		}
	}

	ip := net.ParseIP(c.IP())
	if ip == nil {
		return "", ErrInvalidIP
	}

	record := &Record{}
	err = maxMindFile.Lookup(ip, record)
	if err != nil {
		return "", err
	}
	if record.Country.ISOCode == "GB" {
		record.Country.ISOCode = "UK"
	}

	return amazon.Region(record.Country.ISOCode), nil
}

// Record More fields available @ https://github.com/oschwald/geoip2-golang/blob/master/reader.go#L85
// Only using what we need is faster
type Record struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

func TLDToRegion(tld string) amazon.Region {
	switch tld {
	case "br", "ca", "cn", "de", "es", "fr", "in", "it", "jp", "mx", "uk":
		return amazon.Region(strings.ToUpper(tld))
	}

	return amazon.RegionUS
}

func GetCurrencySign(region amazon.Region) string {

	switch region {
	case amazon.RegionBrazil:
		return "R$"
	case amazon.RegionCanada, amazon.RegionUS, amazon.RegionMexico:
		return "$"
	case amazon.RegionChina:
		return "¥"
	case amazon.RegionGermany, amazon.RegionSpain, amazon.RegionFrance, amazon.RegionItaly:
		return "€"
	case amazon.RegionIndia:
		return "₹"
	case amazon.RegionJapan:
		return "¥"
	case amazon.RegionUK:
		return "£"
	}

	return ""
}
