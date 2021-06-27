package social

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Jleagle/canihave/pkg/config"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/Jleagle/canihave/pkg/mysql"
	"github.com/huandu/facebook"
)

func Post(item mysql.Item) {
	postToReddit(item)
	postToInstagram(item)
	postToTwitter(item)
	postToFacebook(item)
}

func postToReddit(item mysql.Item) {

}

func postToInstagram(item mysql.Item) {

}

func postToTwitter(item mysql.Item) {

	if config.IsProd() {

		anaconda.SetConsumerKey(os.Getenv("CANIHAVE_TWITTER_CONSUMER_KEY"))
		anaconda.SetConsumerSecret(os.Getenv("CANIHAVE_TWITTER_CONSUMER_SECRET"))
		api := anaconda.NewTwitterApi(os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN"), os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN_SECRET"))

		message := item.Name + " - " + item.GetLink() + " - " + item.GetAmazonLink()

		logger.Info("Posting to Twitter: " + message)
		go func() {

			_, err := api.PostTweet(message, url.Values{})
			if err != nil {
				logger.Err("Can't tweet", err)
			}
		}()
	}
}

func postToFacebook(item mysql.Item) {

	if config.IsProd() {

		params := facebook.Params{
			"access_token": os.Getenv("CANIHAVE_FACEBOOK_TOKEN"),
			"message":      "Item Name: " + item.Name + "\n\nCanihave.one: " + item.GetLink() + "\n\nBuy on Amazon: " + item.GetAmazonLink(),
			"link":         item.GetLink(),
			"published":    true,
		}

		go func() {

			_, err := facebook.Post("/pleasecanihaveone/feed", params)
			if err != nil {
				logger.Err("Can't facebook", err)
			}
		}()
	}
}
