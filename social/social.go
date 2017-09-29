package social

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Jleagle/canihave/environment"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/models"
	"github.com/huandu/facebook"
)

// todo, post to reddit user, instagram

func PostToTwitter(item models.Item) {

	if environment.IsLive() {
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

func PostToFacebook(item models.Item) {

	if environment.IsLive() {

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
