package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	dbDriver  string = os.Getenv("DRIVER_NAME")
	dbURL     string = os.Getenv("DATABASE_URL")
	tableName string = "homeworks"
	db        *sql.DB
)

func init() {
	var err error
	db, err = sql.Open(dbDriver, dbURL)
	if err != nil {
		panic(err)
	}

	// 現在の課題ID、省略された課題の教科名、アラーム時刻を取得
	var sqlStatement string = fmt.Sprintf("SELECT id, omitted, alarmtime FROM %s", tableName)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(nil)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name string
		var alarmTime time.Time

		rows.Scan(&id, &name, &alarmTime)
		// アラーム一覧になければ課題アラームを追加
		if _, exist := alarmTimes[id]; !exist {
			alarmTimes[id] = []interface{}{alarmTime, name}
		}
	}

	// 現在のアラーム一覧
	nowAlarms()
}
