package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ktnyt/go-moji"
)

// ScheduleStruct はスケジュール一覧を格納する構造体
type ScheduleStruct struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Title         string      `json:"title"`
			AllDay        bool        `json:"all_day"`
			StartAt       time.Time   `json:"start_at"`
			StartTimezone string      `json:"start_timezone"`
			EndAt         time.Time   `json:"end_at"`
			EndTimezone   string      `json:"end_timezone"`
			Location      string      `json:"location"`
			URL           interface{} `json:"url"`
			UpdatedAt     time.Time   `json:"updated_at"`
			CreatedAt     time.Time   `json:"created_at"`
			Category      string      `json:"category"`
			Description   interface{} `json:"description"`
			Recurrence    interface{} `json:"recurrence"`
			RecurringUUID interface{} `json:"recurring_uuid"`
		} `json:"attributes"`
		Relationships struct {
			Label struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"label"`
			Creator struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"creator"`
			Attendees struct {
				Data []struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"attendees"`
		} `json:"relationships"`
	} `json:"data"`
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

// getSchedule はTimeTreeにアクセスして予定を取得する関数
func getSchedule() {
	var baseURL string = "https://timetreeapis.com/calendars/"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	reqURL.Path = path.Join(reqURL.Path, os.Getenv("CALENDAR_ID"))
	reqURL.Path = path.Join(reqURL.Path, "upcoming_events")
	reqURLvar, _ := url.ParseQuery(reqURL.RawQuery)
	reqURLvar.Add("days", "7")
	reqURL.RawQuery = reqURLvar.Encode()

	// GETするURLを確認
	// fmt.Println(reqURL.String())

	req, _ := http.NewRequest("GET", reqURL.String(), nil)

	req.Header.Add("Accept", "application/vnd.timetree.v1+json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TIMETREE_TOKEN"))

	response, _ := http.DefaultClient.Do(req)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	// APIレスポンスを確認
	// fmt.Println(string(body))
	// fmt.Println("取得回数上限:", response.Header.Values("X-Ratelimit-Limit"))
	// fmt.Println("残機:", response.Header.Values("X-Ratelimit-Remaining"))
	// fmt.Println("リセットまであと(UNIX):", response.Header.Values("X-Ratelimit-Reset"))

	var jsonData ScheduleStruct
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		panic(err)
	}

	// 現在のデータベースのhomeworksコレクションを全初期化
	dbDelete()
	for i := 0; i < len(jsonData.Data); i++ {
		// スケジュールの時刻をUTCからJSTに変換
		jsonData.Data[i].Attributes.StartAt = timeDiffConv(jsonData.Data[i].Attributes.StartAt)
		jsonData.Data[i].Attributes.EndAt = timeDiffConv(jsonData.Data[i].Attributes.EndAt)
		jsonData.Data[i].Attributes.UpdatedAt = timeDiffConv(jsonData.Data[i].Attributes.UpdatedAt)
		jsonData.Data[i].Attributes.CreatedAt = timeDiffConv(jsonData.Data[i].Attributes.CreatedAt)

		// 課題ラベルが貼られたスケジュールをデータベースに追加
		// そのラベルが貼られたもののIDは32文字と決まっているので、それ以外は無視
		if strings.Split(jsonData.Data[i].Relationships.Label.Data.ID, ",")[0] == os.Getenv("LABEL_ID") && idIsFinite(jsonData.Data[i].ID) {
			var homeworkID string = jsonData.Data[i].ID
			var subject, subjectOmitted, title string = subjectSeparate(jsonData.Data[i].Attributes.Title)
			var dueAt time.Time = jsonData.Data[i].Attributes.EndAt
			var alarmAt time.Time = alarmTimeCalc(dueAt)

			// 現在のアラームに指定する課題を追加
			dbAdd(dbClient, map[string]interface{}{
				"homeworkID": homeworkID,
				"subject":    subject,
				"omitted":    subjectOmitted,
				"title":      title,
				"dueAt":      dueAt,
				"alarmAt":    alarmAt,
			})

			alarmAdd(homeworkID, dueTime, alarmAt)
		}
	}
}

// timeDiffConv は時差変換をして返す関数
func timeDiffConv(tTime time.Time) (rTime time.Time) {
	// 純度が高いUTCにする
	rTime = tTime.UTC()

	// UTC → JST
	var jst *time.Location = time.FixedZone("Asia/Tokyo", 9*60*60)
	rTime = rTime.In(jst)

	return
}

// subjectSeparate はタイトルから教科と課題名を抜き出す関数
func subjectSeparate(title string) (string, string, string) {
	for i, str := range subjectsTimeTree {
		if strings.Contains(homeworkFormat(title), str) {
			subject = str
			subjectOmitted = subjectsOmitted[i]
			homeworkTitle = strings.TrimLeft(homeworkFormat(title), subject)

			return subject, subjectOmitted, homeworkTitle
		}
	}

	// 教科が見つからなかった場合はどうしようもないので、
	// 「○○の課題」に合いそうな感じの教科名を入れる
	return "学校", "学校", title
}

// homeworkFormat は与えられたタイトルの数字やアルファベットを半角にし、空白を埋める関数
func homeworkFormat(subject string) string {
	// 全角英数を半角英数に
	subject = moji.Convert(subject, moji.ZE, moji.HE)
	// 全角スペースを半角スペースに
	subject = moji.Convert(subject, moji.ZS, moji.HS)

	return strings.Join(strings.Split(subject, " "), "")
}

// idIsFinite はIDが32文字であるかどうかをチェックする関数
func idIsFinite(tStr string) bool {
	if len(tStr) == 32 {
		return true
	}
	return false
}
