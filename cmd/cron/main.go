package main

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/Jleagle/canihave/cmd/cron/sites"
	"github.com/Jleagle/canihave/pkg/helpers"
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/Jleagle/canihave/pkg/mysql"
	"github.com/Jleagle/canihave/pkg/social"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {

	c := cron.New(
		cron.WithLogger(logger.CronLogger{}),
		cron.WithSeconds(),
	)

	for _, task := range sites.Sites {
		// In a func here so `task` gets copied into a new memory location and can not be replaced at a later time
		func(task sites.Site) {

			_, err := c.AddFunc(task.GetID(), func() {

			})
			if err != nil {
				logger.Logger.Error("adding cron func")
			}
		}(task)
	}
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

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Fatal("", zap.Error(err))
	}

	return string(bytes), resp.StatusCode
}
