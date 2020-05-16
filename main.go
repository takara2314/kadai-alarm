package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/callback", callbackFunc)
	http.ListenAndServe(":8080", nil)
}

// callbackFunc はアクセスされたときに処理を行う関数
func callbackFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GETされちゃいました！")
	fmt.Println("code:", r.FormValue("code"))
	fmt.Println("state", r.FormValue("state"))
}
