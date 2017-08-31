package scraper

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/Jleagle/canihave/models"
	"github.com/go-chi/chi"
	"time"
)

func ScrapeHandler(w http.ResponseWriter, r *http.Request) {

	// todo, check env var to stop people hitting this url

	id := chi.URLParam(r, "id")
	if id != "" {
		item := models.Item{}
		item.ID = id
		item.Get()
	} else {
		shitYouCanAfford()
		importItems()
	}
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

		item.ID = value
		item.Source = "shityoucanafford"
		item.Get()

		time.Sleep(time.Second)
	}
}

func importItems() bool {

	items := []string{
		"0735216207",
		"1433805618",
		"1501175564",
		"B004RMK4BC",
		"B004S8F7QM",
		"B00E1EN92W",
		"B00EB4ADQW",
		"B00GAC1D2G",
		"B00IOY8XWQ",
		"B00JM5GW10",
		"B00NB86OYE",
		"B00O4OR4GQ",
		"B00OQVZDJM",
		"B00P77ZAN8",
		"B00REQKWGA",
		"B00REQL3AE",
		"B00U3FPN4U",
		"B00U3FPN4U",
		"B00UT823WQ",
		"B00X4WHP5E",
		"B00ZV9PXP2",
		"B00ZV9RDKK",
		"B00ZV9RDKK",
		"B00ZV9RDKK",
		"B0186JAEWK",
		"B01BH83OOM",
		"B01DFTCV90",
		"B01E3QM34W",
		"B01EQYX9NU",
		"B01GEW27DA",
		"B01GEW27DA",
		"B01I499BNA",
		"B01J24C0TI",
		"B01J90MSDS",
		"B01J90MSDS",
		"B01J94SBEY",
		"B01J94SWWU",
		"B01KMSKNGU",
		"B01M71IUZ7",
		"B01MRG7T0D",
		"B06XDC9RBJ",
		"B071JRMKBH",
	}

	i := models.Item{}
	for _, id := range items {

		i.ID = id
		i.Source = "import"
		i.Get()

		time.Sleep(time.Second)
	}

	return true
}
