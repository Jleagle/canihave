package scraper

import (
	"net/http"
	"regexp"

	"github.com/Jleagle/canihave/models"
	"github.com/go-chi/chi"
	"fmt"
	"github.com/Jleagle/canihave/location"
	"io/ioutil"
	"log"
)

const (
	Import           string = "import"
	ShitYouCanAfford string = "shityoucanafford"
	WannaSpend       string = "wannaspend"
	DatTwenty        string = "dattwenty"
)

func ScrapeHandler(w http.ResponseWriter, r *http.Request) {

	region := location.GetAmazonRegion(w, r)
	location.SetAmazonEnviromentVars(region)
	// todo, check env var to stop people hitting this url

	id := chi.URLParam(r, "id")
	if id != "" {

		item := models.Item{}
		item.ID = id
		item.Get()

	} else {

		getSingle("http://shityoucanafford.com/", ShitYouCanAfford)
		getSingle("http://dattwenty.com/pages/home", DatTwenty)
		getSingle("http://www.wannaspend.com/", WannaSpend)
	}
}

func getSingle(url string, source string) {

	body, code := doCurl(url)
	if code == 200 {

		r := regexp.MustCompile(`http(.*?)amazon.([a-z]{2,3})/(.*?)/([A-Z0-9]{10})`)
		links := r.FindAllString(body, -1)

		item := models.Item{}
		for _, link := range links {

			m := r.FindStringSubmatch(link)
			//fmt.Printf("%v", m)

			item.Region = location.TLDToRegion(m[2])
			item.ID = m[4]
			item.Source = source
			item.Get()
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
