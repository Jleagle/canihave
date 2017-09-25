package social

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Jleagle/canihave/environment"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/models"
)

func PostToTwitter(item models.Item) {

	if environment.IsLive() {
		anaconda.SetConsumerKey(os.Getenv("CANIHAVE_TWITTER_CONSUMER_KEY"))
		anaconda.SetConsumerSecret(os.Getenv("CANIHAVE_TWITTER_CONSUMER_SECRET"))
		api := anaconda.NewTwitterApi(os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN"), os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN_SECRET"))

		postInNewRoutine(item, api)
	}
}

func postInNewRoutine(item models.Item, api *anaconda.TwitterApi) {

	tweet := item.Name + " - " + item.GetLink() + " - " + item.Link

	logger.Info("Posting to Twitter: " + tweet)
	go api.PostTweet(tweet, url.Values{})
}
