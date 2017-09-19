package social

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Jleagle/canihave/models"
)

func PostToTwitter(item models.Item) {

	anaconda.SetConsumerKey(os.Getenv("CANIHAVE_TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CANIHAVE_TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN"), os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN_SECRET"))

	go api.PostTweet(item.Name+" - "+item.GetLink()+" - "+item.Link, url.Values{})
}
