package social

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Jleagle/canihave/pkg/config"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/Jleagle/canihave/pkg/mysql"
	"github.com/huandu/facebook"
	"go.uber.org/zap"
)

func Post(item mysql.Item) {
	go postToReddit(item)
	go postToInstagram(item)
	go postToTwitter(item)
	go postToFacebook(item)
}

func postToReddit(item mysql.Item) {

}

func postToInstagram(item mysql.Item) {

}

func postToTwitter(item mysql.Item) {

	anaconda.SetConsumerKey(config.TwitterConsumerKey)
	anaconda.SetConsumerSecret(config.TwitterConsumerSecret)
	api := anaconda.NewTwitterApi(config.TwitterAccessToken, config.TwitterAcessTokenSecret)

	message := item.Name + " - " + item.GetLink() + " - " + item.GetAmazonLink()

	logger.Logger.Info("Posting to Twitter: " + message)
	go func() {

		_, err := api.PostTweet(message, url.Values{})
		if err != nil {
			logger.Logger.Error("Can't tweet", zap.Error(err))
		}
	}()
}

func postToFacebook(item mysql.Item) {

	params := facebook.Params{
		"access_token": config.FacebookToken,
		"message":      "Item Name: " + item.Name + "\n\nCanihave.one: " + item.GetLink() + "\n\nBuy on Amazon: " + item.GetAmazonLink(),
		"link":         item.GetLink(),
		"published":    true,
	}

	go func() {

		_, err := facebook.Post("/pleasecanihaveone/feed", params)
		if err != nil {
			logger.Logger.Error("Can't facebook", zap.Error(err))
		}
	}()
}
