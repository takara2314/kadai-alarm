// 参考: https://qiita.com/dahiyu/items/7ffd6ee0b2afa0ea46bd
// 参考: https://lancers.work/golang/golang-tweet/
// 参考: https://developers.timetreeapp.com/ja/docs/api
// 参考: https://qiita.com/hnakamur/items/ce87adfe04e932dab2aa
// 参考: https://maku77.github.io/hugo/go/cast.html
// 感謝: https://mholt.github.io/json-to-go/
// 感謝: https://jsonformatter.curiousconcept.com/

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// 本人に伝える時間
var alarmTimes map[string][]interface{} = make(map[string][]interface{})

func main() {
	// 提出期限の指定した時間前になったら本人に通知
	go doOnScheduleTime([]int{0, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}, []int{0, 15, 30, 45})
	// 定期的にTimeTreeのスケジュールを取得
	go getScheduleRegularly([]int{0, 7, 9, 11, 13, 15, 17, 19, 21, 23})
	// 毎時指定した時間にpingを送って返させて、dynoを停止させないようにする
	go pingRegularly([]int{0, 15, 30, 45})
	// アクセスしたときにTimeTreeのスケジュールを取得
	http.HandleFunc("/", serverMainFunc)
	http.HandleFunc("/get", serverGetFunc)
	http.HandleFunc("/now", serverNowFunc)
	http.HandleFunc("/ping", serverPingFunc)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

// serverMainFunc はアクセスされたときに処理を行う関数
func serverMainFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("スケジュールを手動取得する場合は\"/get\"ページ、現在のスケジュール一覧を見る場合は\"/now\"ページにアクセスしてください。"))
}

// serverGetFunc はアクセスされたときに処理を行う関数
func serverGetFunc(w http.ResponseWriter, r *http.Request) {
	getSchedule()

	fmt.Println("スケジュールが手動更新されました。(ページアクセス)")
	w.Write([]byte("スケジュールをTimeTreeから取得しました。詳細はコンソールからご確認ください。"))
}

// serverNowFunc はアクセスされたときに処理を行う関数
func serverNowFunc(w http.ResponseWriter, r *http.Request) {
	nowAlarms()
	w.Write([]byte("現在のスケジュールをコンソールに表示しました。"))
}

// serverPingFunc はアクセスされたときに処理を行う関数
func serverPingFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// getScheduleRegularly は毎日定めた時間にTimeTreeのスケジュールを取得する関数
func getScheduleRegularly(getHours []int) {
	for {
		var nowTime time.Time = time.Now()
		// 現在時刻がh時0分なら条件分岐開始
		if nowTime.Minute() == 0 {
			// 対象の時刻(h時)になったらスケジュール取得
			if sliceContain(getHours, nowTime.Hour()) {
				getSchedule()
				fmt.Println("スケジュールが自動更新されました。(定期更新)")
			}
		}
	}
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
