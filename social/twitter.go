package social

import (
	"github.com/Jleagle/canihave/models"
	"github.com/ChimeraCoder/anaconda"
	"os"
)

func postToTwitter(item models.Item) {

	anaconda.SetConsumerKey(os.Getenv("SQL_PW"))
	anaconda.SetConsumerSecret(os.Getenv("SQL_PW"))
	_ := anaconda.NewTwitterApi(os.Getenv("SQL_PW"), os.Getenv("SQL_PW"))

	//api.tweet

}
