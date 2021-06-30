package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/Jleagle/canihave/pkg/config"
	"github.com/Jleagle/canihave/pkg/helpers"
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/Jleagle/canihave/pkg/mysql"
	"github.com/Jleagle/canihave/pkg/social"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"go.uber.org/zap"
)

const (
	siteShitYouCanAfford = "shityoucanafford"
	siteWannaSpend       = "wannaspend"
	siteDatTwenty        = "dattwenty"
	siteCanopy           = "canopy"
	siteFiveStar         = "fivestar"
	siteBoughtItOnce     = "boughtitonce"
)

func searchx() {

	client := search.NewClient(config.AlgoliaAppID, config.AlgoliaSearch)

	index := client.InitIndex("products")

	fmt.Println(index)
}

func ScrapeHandler(social bool) {

	// todo, get all items and pass them in so not to add them

	getSingle(social, siteShitYouCanAfford, "http://shityoucanafford.com/")
	getSingle(social, siteDatTwenty, "http://dattwenty.com/pages/home")
	getSingle(social, siteWannaSpend, "http://www.wannaspend.com/")
	getSingle(social, siteCanopy, "https://canopy.co/ajax/merged_feed_products?limit=100")
	getSingle(social, siteFiveStar, "https://fivestar.io/index-3b0dc4e7b4c5c55e5da4.js")
	getSingle(social, siteBoughtItOnce, "http://boughtitonce.com/")
}

func getSingle(postToSocial bool, source string, url string) {

	body, code := doCurl(url)
	if code == 200 {

		r := regexp.MustCompile(`http(.*?)amazon.([a-z]{2,3})/(.*?)/([A-Z0-9]{10})`)
		links := r.FindAllString(body, -1)

		links = helpers.RemoveDuplicates(links)
		links = helpers.Reverse(links)

		for _, link := range links {

			m := r.FindStringSubmatch(link)

			item, err := mysql.GetWithExtras(m[4], location.TLDToRegion(m[2]), mysql.TypeScraper, source)
			if err != nil {
				logger.Logger.Error("Can't get with extras", zap.Error(err))
				continue
			}

			logger.Logger.Info("Adding " + item.ID)
			if false {
				social.Post(item)
			}
		}

	} else {
		logger.Logger.Info(source + " seems to be down")
	}
}

func doCurl(url string) (body string, code int) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Logger.Fatal("", zap.Error(err))
	}

	req.Header.Set("User-Agent", "Googlebot")

	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Fatal("", zap.Error(err))
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Fatal("", zap.Error(err))
	}

	return string(bytes), resp.StatusCode
}
