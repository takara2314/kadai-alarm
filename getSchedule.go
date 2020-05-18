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

// getSchedule はアクセスされたときに処理を行う関数
func getSchedule() {
	var baseURL string = "https://timetreeapis.com/calendars/"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	reqURL.Path = path.Join(reqURL.Path, os.Getenv("CALENDAR_ID"))
	reqURL.Path = path.Join(reqURL.Path, "upcoming_events")
	reqURLvar, _ := url.ParseQuery(reqURL.RawQuery)
	// タイムゾーン指定しても何故か要求通りに返ってこない
	// reqURLvar.Add("timezone", "Asia/Tokyo")
	reqURLvar.Add("days", "7")
	reqURL.RawQuery = reqURLvar.Encode()

	fmt.Println(reqURL.String())

	req, _ := http.NewRequest("GET", reqURL.String(), nil)

	req.Header.Add("Accept", "application/vnd.timetree.v1+json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TIMETREE_TOKEN"))

	response, _ := http.DefaultClient.Do(req)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	// fmt.Println(string(body))
	// fmt.Println("取得回数上限:", response.Header.Values("X-Ratelimit-Limit"))
	// fmt.Println("残機:", response.Header.Values("X-Ratelimit-Remaining"))
	// fmt.Println("リセットまであと(UNIX):", response.Header.Values("X-Ratelimit-Reset"))

	var jsonData ScheduleStruct
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(jsonData.Data); i++ {
		// UTC時刻をJST時刻に変換
		jsonData.Data[i].Attributes.StartAt = timeDiffConv(jsonData.Data[i].Attributes.StartAt)
		jsonData.Data[i].Attributes.EndAt = timeDiffConv(jsonData.Data[i].Attributes.EndAt)
		jsonData.Data[i].Attributes.UpdatedAt = timeDiffConv(jsonData.Data[i].Attributes.UpdatedAt)
		jsonData.Data[i].Attributes.CreatedAt = timeDiffConv(jsonData.Data[i].Attributes.CreatedAt)

		// 仮として追加
		jsonData.Data[i].Attributes.EndAt, _ = time.Parse("2006-01-02 15-04-05", "2020-05-18 04-00-00")
		jsonData.Data[i].Attributes.EndAt = timeDiffConv(jsonData.Data[i].Attributes.EndAt)

		// 課題ラベルが貼られたスケジュールをデータベースに追加
		// IDが32文字以外ならスルー
		if strings.Split(jsonData.Data[i].Relationships.Label.Data.ID, ",")[0] == os.Getenv("LABEL_ID") && idIsFinite(jsonData.Data[i].ID) {
			var homeworkID string = jsonData.Data[i].ID
			var subject, subjectOmitted, homework string = subjectSeparate(jsonData.Data[i].Attributes.Title)
			var dueTime time.Time = jsonData.Data[i].Attributes.EndAt

			// fmt.Println("追加します！", homeworkID, subject, subjectOmitted, homework, dueTime)

			homeworkAdd(homeworkID, subject, subjectOmitted, homework, dueTime)
		}
	}
}

// timeDiffConv は時差変換をして返す関数
func timeDiffConv(tTime time.Time) (rTime time.Time) {
	// よりUTCらしくする
	rTime = tTime.UTC()

	// UTC → JST
	var jst *time.Location = time.FixedZone("Asia/Tokyo", 9*60*60)
	rTime = rTime.In(jst)

	return
}

// subjectSeparate はタイトルから教科と課題名を抜き出す関数
func subjectSeparate(title string) (subject string, subjectOmitted string, homework string) {
	for i, str := range subjectsTimeTree {
		if strings.Contains(homeworkFormat(title), str) {
			subject = str
			subjectOmitted = subjectsOmitted[i]
			homework = strings.TrimLeft(homeworkFormat(title), subject)

			return subject, subjectOmitted, homework
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
