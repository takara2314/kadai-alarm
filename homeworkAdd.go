package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// insertData はデータベースに入れる内容を格納する構造体
type insertData struct {
	HomeworkID     string
	Subject        string
	SubjectOmitted string
	Title          string
	DueTime        time.Time
	AlarmTime      time.Time
}

func homeworkAdd(homeworkID string, subject string, subjectOmitted string, homework string, dueTime time.Time) {
	homeworkInfo := insertData{
		HomeworkID:     homeworkID,
		Subject:        subject,
		SubjectOmitted: subjectOmitted,
		Title:          homework,
		DueTime:        dueTime,
		AlarmTime:      alarmTimeAdd(homeworkID, dueTime, subjectOmitted),
	}
	dbAdd(&homeworkInfo)
}

// dbAdd はDBテーブルに同じデータがなければ追加する関数
func dbAdd(homeworkInfo *insertData) {
	var sqlStatement string = fmt.Sprintf("SELECT id FROM %s", tableName)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(nil)
	}
	defer rows.Close()

	var dbSubjectIDs map[string]int = make(map[string]int, 0)
	for rows.Next() {
		var temp string
		rows.Scan(&temp)
		dbSubjectIDs[temp] = 1
	}

	// 同じものがなければ書き込み
	if _, exist := dbSubjectIDs[homeworkInfo.HomeworkID]; !exist {
		dbWrite(homeworkInfo)
	} else {
		fmt.Println("同じものがありました…")
	}
}

// dbWrite はdbテーブルに課題情報を書き込む関数
func dbWrite(homeworkInfo *insertData) {
	var sqlStatement string = fmt.Sprintf("INSERT INTO %s (id, subject, omitted, title, duetime, alarmtime) VALUES ($1, $2, $3, $4, $5, $6)", tableName)

	_, err := db.Exec(
		sqlStatement,
		homeworkInfo.HomeworkID,
		homeworkInfo.Subject,
		homeworkInfo.SubjectOmitted,
		homeworkInfo.Title,
		homeworkInfo.DueTime,
		homeworkInfo.AlarmTime,
	)
	if err != nil {
		panic(err)
	}
}

// dbDelete はテーブル内のデータをすべて削除する関数
func dbDelete() {
	var sqlStatement string = "DELETE FROM " + tableName

	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}

// alarmTimeAdd は課題アラームを追加し、アラーム時刻を返す
func alarmTimeAdd(id string, dueTime time.Time, name string) time.Time {
	// 提出期限の何時間前かを取得し、その時刻を求める
	hourBeforeDep, _ := strconv.Atoi(os.Getenv("HOUR_BEFORE_DEPARTURE"))
	var alarmTime time.Time = dueTime.Add(-time.Duration(hourBeforeDep) * time.Hour)

	// アラーム一覧になければ課題アラームを追加
	if _, exist := alarmTimes[id]; !exist {
		alarmTimes[id] = []interface{}{alarmTime, name}
	}

	return alarmTime
}
