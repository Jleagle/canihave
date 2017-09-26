package scraper

import (
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
	SOURCE_Manual           string = "manual"
	SOURCE_ShitYouCanAfford string = "shityoucanafford"
	SOURCE_WannaSpend       string = "wannaspend"
	SOURCE_DatTwenty        string = "dattwenty"
	SOURCE_Canopy           string = "canopy"
	SOURCE_FiveStar         string = "fivestar"
	SOURCE_BoughtItOnce     string = "boughtitonce"
)

func ScrapeHandler(social bool) {

	// todo, get all items and pass them in so not to add them

	getSingle(social, SOURCE_ShitYouCanAfford, "http://shityoucanafford.com/")
	getSingle(social, SOURCE_DatTwenty, "http://dattwenty.com/pages/home")
	getSingle(social, SOURCE_WannaSpend, "http://www.wannaspend.com/")
	getSingle(social, SOURCE_Canopy, "https://canopy.co/ajax/merged_feed_products?limit=100")
	getSingle(social, SOURCE_FiveStar, "https://fivestar.io/index-3b0dc4e7b4c5c55e5da4.js")
	getSingle(social, SOURCE_BoughtItOnce, "http://boughtitonce.com/")
}

func getSingle(postToSocial bool, source string, url string) {

	body, code := doCurl(url)
	if code == 200 {

		r := regexp.MustCompile(`http(.*?)amazon.([a-z]{2,3})/(.*?)/([A-Z0-9]{10})`)
		links := r.FindAllString(body, -1)

		links = helpers.RemoveDuplicatesUnordered(links)
		links = helpers.ArrayReverse(links)

		for _, link := range links {

			m := r.FindStringSubmatch(link)

			_, err := models.GetWithExtras(m[4], location.TLDToRegion(m[2]), models.TYPE_SCRAPE, source)
			if err != nil {
				logger.Err("Can't get with extras", err)
				continue
			}

			//social.PostToTwitter(item)
			//social.PostToFacebook(item)
		}

	} else {
		logger.Info(source + " seems to be down")
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
