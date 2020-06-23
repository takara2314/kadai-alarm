// 参考: https://qiita.com/dahiyu/items/7ffd6ee0b2afa0ea46bd
// 参考: https://lancers.work/golang/golang-tweet/
// 参考: https://developers.timetreeapp.com/ja/docs/api
// 参考: https://qiita.com/hnakamur/items/ce87adfe04e932dab2aa
// 参考: https://maku77.github.io/hugo/go/cast.html
// 参考: https://firebase.google.com/docs/firestore/
// 感謝: https://mholt.github.io/json-to-go/
// 感謝: https://jsonformatter.curiousconcept.com/

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// 定期的にTimeTreeのスケジュールを取得
	// 今回は7時から23時までの奇数時の0分に設定
	go getScheduleRegularly([]int{0, 7, 9, 11, 13, 15, 17, 19, 21, 23})
	// 提出期限の指定した時間になったら本人に通知
	// 今回は6時から24時の0分、15分、30分、45分に設定
	go doOnScheduleTime([]int{0, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}, []int{0, 15, 30, 45})

	// サイトアクセスされたときに案内を表示
	http.HandleFunc("/", serverMainFunc)
	// サイトアクセスされたときにTimeTreeのスケジュールを取得
	http.HandleFunc("/get", serverGetFunc)
	// サイトアクセスされたときに現在の通知スケジュールをコンソールで表示
	http.HandleFunc("/now", serverNowFunc)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

// serverMainFunc <= "/"
func serverMainFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("スケジュールを手動取得する場合は\"/get\"ページ、現在のスケジュール一覧を見る場合は\"/now\"ページにアクセスしてください。"))
}

// serverGetFunc <= "/get"
func serverGetFunc(w http.ResponseWriter, r *http.Request) {
	getSchedule()

	fmt.Println("スケジュールが手動更新されました。(ページアクセス)")
	w.Write([]byte("スケジュールをTimeTreeから取得しました。詳細はコンソールからご確認ください。"))
}

// serverNowFunc <= "/now"
func serverNowFunc(w http.ResponseWriter, r *http.Request) {
	nowAlarms()
	w.Write([]byte("現在のスケジュールをコンソールに表示しました。"))
}

// sliceContain はint型のスライスに該当する値があったらtrueを返す関数
func sliceContain(tSlice []int, tNum int) bool {
	for _, num := range tSlice {
		if num == tNum {
			return true
		}
	}
	return false
}

// nowAlarms は課題リスト(マップ)から課題一覧を標準出力する関数
func nowAlarms() {
	// マップの各キー名のIDを取り出す
	var homeworkIDs []string = make([]string, 0)
	for key := range alarmTimes {
		homeworkIDs = append(homeworkIDs, key)
	}

	var nowTime time.Time = time.Now()
	fmt.Printf("【現在の課題一覧】現在時刻> %s\n", nowTime.Format("1月2日 15時4分"))
	for i := 0; i < len(alarmTimes); i++ {
		// interface{}型からtime.Time型にキャスト (型アサーション)
		var alarmTime time.Time = alarmTimes[homeworkIDs[i]][0].(time.Time)

		fmt.Printf("%v(%s): アラーム時刻> %s\n",
			alarmTimes[homeworkIDs[i]][1],
			homeworkIDs[i],
			alarmTime.Format("1月2日15時4分"),
		)
	}
}
