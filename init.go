package main

import (
	"database/sql"
	"fmt"
	"os"

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
	} else {
		fmt.Println("データベースに接続成功しました。")
	}

	// 久しぶりの起動を予期して、最初に初期化してフレッシュな予定を取得する
	dbDelete()
	getSchedule()

	// 現在のアラーム一覧
	nowAlarms()
}
