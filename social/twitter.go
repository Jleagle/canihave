package social

import (
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Jleagle/canihave/models"
)

func postToTwitter(item models.Item) {

	anaconda.SetConsumerKey(os.Getenv("CANIHAVE_TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CANIHAVE_TWITTER_CONSUMER_SECRET"))
	_ := anaconda.NewTwitterApi(os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN"), os.Getenv("CANIHAVE_TWITTER_ACCESS_TOKEN_SECRET"))

	//api.tweet

}
