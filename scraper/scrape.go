package scraper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/Jleagle/canihave/models"
)

func ScrapeHandler(w http.ResponseWriter, r *http.Request) {

	// todo, check env var to stop people hitting this url
	shitYouCanAfford()
}

func shitYouCanAfford() {

	resp, err := http.Get("http://shityoucanafford.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	r := regexp.MustCompile("[A-Z0-9]{10}")
	matches := r.FindAllString(string(body), 10)

	item := models.Item{}
	for _, value := range matches {
		fmt.Printf("%v", "Adding "+item.ID)
		item.ID = value
		item.Source = "shityoucanafford"
		item.Get()
	}
}
