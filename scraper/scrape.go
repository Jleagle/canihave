package scraper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/Jleagle/canihave/helpers"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/models"
)

const (
	MANUAL string = "manual"

	ShitYouCanAfford string = "shityoucanafford"
	WannaSpend       string = "wannaspend"
	DatTwenty        string = "dattwenty"
	Canopy           string = "canopy"
	FiveStar         string = "fivestar"
	BoughtItOnce     string = "boughtitonce"
)

func ScrapeHandler(social bool) {

	// todo, get all items and pass them in so not to add them

	getSingle(social, ShitYouCanAfford, "http://shityoucanafford.com/")
	getSingle(social, DatTwenty, "http://dattwenty.com/pages/home")
	getSingle(social, WannaSpend, "http://www.wannaspend.com/")
	getSingle(social, Canopy, "https://canopy.co/ajax/merged_feed_products?limit=100")
	getSingle(social, FiveStar, "https://fivestar.io/index-3b0dc4e7b4c5c55e5da4.js")
	getSingle(social, BoughtItOnce, "http://boughtitonce.com/")
}

func getSingle(social bool, source string, url string) {

	body, code := doCurl(url)
	if code == 200 {

		r := regexp.MustCompile(`http(.*?)amazon.([a-z]{2,3})/(.*?)/([A-Z0-9]{10})`)
		links := r.FindAllString(body, -1)

		links = helpers.RemoveDuplicatesUnordered(links)
		links = helpers.ArrayReverse(links)

		item := models.Item{}
		for _, link := range links {

			m := r.FindStringSubmatch(link)
			fmt.Printf("%v", m)

			item.Region = location.TLDToRegion(m[2])
			item.ID = m[4]
			item.Source = source
			item.Type = models.TYPE_SCRAPE
			item.GetWithExtras()

			if item.Region == "" {
				logger.Err("Item has no region when being scraped")
			}
		}

	} else {
		fmt.Printf("%# v", source+" seems to be down")
	}
}

func doCurl(url string) (body string, code int) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Googlebot")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(bytes), resp.StatusCode
}
