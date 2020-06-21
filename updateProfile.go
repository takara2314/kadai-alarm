package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func updateProfile(subject string, toNormal bool) {
	api := getTwitterAPI()

	var postUserName string
	var userName string = "ふぉくしーど"

	if !toNormal {
		postUserName = userName + "@" + subject + "の課題未提出かもよ"
	} else {
		postUserName = userName
	}

	changeThings := url.Values{}
	changeThings.Set("name", postUserName)
	api.PostAccountUpdateProfile(changeThings)
	fmt.Println("ツイッターの名前を変更しました…")
}

func getTwitterAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))

	return anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}
