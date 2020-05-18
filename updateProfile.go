package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func updateProfile(subject string, toNormal bool) {
	api := getTwitterAPI()

	// _, err := api.PostTweet("さっき言ったやつ、結構順調に進んでる！", nil)
	// if err != nil {
	// 	panic(err)
	// }

	var postUserName string
	if !toNormal {
		var userName string = os.Getenv("USER_NAME")
		postUserName = userName + "@" + subject + "の課題未提出かもよ"
	} else {
		postUserName = os.Getenv("USER_NAME")
	}

	changeThings := url.Values{}
	changeThings.Set("name", postUserName)
	api.PostAccountUpdateProfile(changeThings)
}

func getTwitterAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))

	fmt.Println("ツイッターの名前を変更しました…")
	return anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}
