package main

import (
	"net/http"
	"time"
)

// pingRegularly は毎時指定した時間にpingを送って返させて、dynoを停止させないようにする関数
func pingRegularly(scheduledMinute []int) {
	for {
		if sliceContain(scheduledMinute, time.Now().Minute()) {
			var baseURL string = "https://simple-pinger.herokuapp.com/ping"

			req, _ := http.NewRequest("GET", baseURL, nil)

			_, _ = http.DefaultClient.Do(req)
			time.Sleep(1 * time.Minute)
		}

	}
}
