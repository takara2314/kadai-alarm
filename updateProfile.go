package main

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func updateProfile(subject string) {
	api := getTwitterAPI()

	// _, err := api.PostTweet("さっき言ったやつ、結構順調に進んでる！", nil)
	// if err != nil {
	// 	panic(err)
	// }

	var userName string = os.Getenv("USER_NAME")
	var postUserName string = userName + "@" + subject + "の課題未提出かもよ"

	changeThings := url.Values{}
	changeThings.Set("name", postUserName)
	api.PostAccountUpdateProfile(changeThings)
}

func getTwitterAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	return anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}
