package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// タイムゾーンを設定
const location = "Asia/Tokyo"

var (
	oldTwitterName string = "ふぉくしーど"
	ctx            context.Context
	app            *firebase.App
	collectionName string = "homeworks"
	alarmTimes            = make(map[string][]interface{})
)

func init() {
	var err error

	// Google App Engine はタイムゾーン指定できないので、Go側でタイムゾーンを指定する
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	// 利用するデータベースを宣言
	ctx = context.Background()
	sa := option.WithCredentialsFile("kadai-alarm-5365dc12423d.json")

	// データベースを開くこのアプリを初期化
	app, err = firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("データベースに接続成功しました。")
	}

	// 現在のデータベースのhomeworksコレクションを全初期化
	dbDelete()
	// TimeTreeのスケジュールを取得
	getSchedule()

	fmt.Println("起動完了しました。指定された時間になりましたら処理が開始されます。")
}
