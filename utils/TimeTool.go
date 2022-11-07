package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type TimeTool struct {
}

// 将传入的“秒”解析为3种时间单位
func (t *TimeTool) ResolveTime(seconds float32) (day string, hour string, minute string, second string) {

	tempDay := int(seconds) / (24 * 3600)
	tempHour := (int(seconds) - tempDay*3600*24) / 3600
	tempMinute := (int(seconds) - tempDay*24*3600 - tempHour*3600) / 60
	tempSecond2 := float32(int(seconds) - tempDay*24*3600 - tempHour*3600 - tempMinute*60)

	arr := strings.Split(fmt.Sprintf("%v", seconds), ".")
	if len(arr) > 1 {
		temp := fmt.Sprintf("%v.%v", tempSecond2, arr[1])
		val, _ := strconv.ParseFloat(temp, 32)
		tempSecond2 = float32(val)
	}

	day = t.timeToStr(float32(tempDay))
	hour = t.timeToStr(float32(tempHour))
	minute = t.timeToStr(float32(tempMinute))
	second = t.timeToStr(tempSecond2)
	return
}

func (t *TimeTool) timeToStr(timeVal float32) (timeStr string) {
	if timeVal < 10 {
		timeStr = fmt.Sprintf("0%v", timeVal)
	} else {
		timeStr = fmt.Sprintf("%v", timeVal)
	}
	return
}
