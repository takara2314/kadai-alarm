package main

import (
	"fmt"
	"time"
)

// doOnScheduleTime は定期的にアラーム時刻になっていないかをチェックし、なっている場合は報告を行う関数
func doOnScheduleTime(getHours []int, getMinutes []int) {
	for {
		var nowTime time.Time = time.Now()
		var nowHour, nowMinute int = nowTime.Hour(), nowTime.Minute()

		// 現在時刻がh時m分なら条件分岐開始
		if sliceContain(getHours, nowHour) && sliceContain(getMinutes, nowMinute) {
			// マップの各キー名のIDを取り出す
			var homeworkIDs []string = make([]string, 0)
			for key := range alarmTimes {
				homeworkIDs = append(homeworkIDs, key)
			}

			var alarmHomeworks []string = make([]string, 0)
			// アラーム時刻に近くなっていないかをそれぞれチェック
			for i := 0; i < len(alarmTimes); i++ {
				// interface{}型からtime.Time型にキャスト (型アサーション)
				var alarmTime time.Time = alarmTimes[homeworkIDs[i]][0].(time.Time)

				// fmt.Println(alarmTime)
				// fmt.Println(nowTime)
				// fmt.Println(alarmTime.Sub(nowTime))
				// 現在時刻の1時間以内にアラーム予定時刻が来るとき
				if int(alarmTime.Sub(nowTime).Hours()) <= 1 {
					fmt.Printf("%vの課題が未提出かもしれません！", alarmTimes[homeworkIDs[i]][1])
					alarmHomeworks = append(alarmHomeworks, fmt.Sprintf("%v", alarmTimes[homeworkIDs[i]][1]))
				}
			}

			// 報告する課題があったとき
			if len(alarmHomeworks) == 1 {
				// 報告するので、アラームリストから削除
				delete(alarmTimes, alarmHomeworks[0])
				// 報告する課題の数が1つだけなら、教科名をそのままにして報告
				updateProfile(alarmHomeworks[0])

			} else if len(alarmHomeworks) > 1 {
				// 報告する課題の数が1より多いなら、それぞれの教科名の頭文字をとって報告
				var subjectShowed string = ""
				for _, str := range alarmHomeworks {
					// 課題数が5より大きい場合は、「複数」と表示させる
					if len(alarmHomeworks) <= 5 {
						subjectShowed += subjectNameOneCharConv(str)
					} else {
						subjectShowed = "複数"
					}
					// 報告するので、アラームリストから削除
					delete(alarmTimes, str)
				}
				updateProfile(subjectShowed)
			}

			time.Sleep(1 * time.Minute)
		}
	}
}

// subjectNameOneCharConv は省略された教科名をさらに1文字にしたものを返す関数
func subjectNameOneCharConv(name string) string {
	for i, oneChar := range subjectsOneChar {
		if subjectsOmitted[i] == name {
			return oneChar
		}
	}
	return "学"
}
