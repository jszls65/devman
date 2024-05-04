// @Title
// @Author  zls  2023/10/15 23:30
package utils

import (
	"time"
)

// 时间戳转日期字符串
func TimestampToStr(ts int) string {
	date := time.Unix(int64(ts), 0)
	return date.Format(time.DateTime)
}

// 字符串转时间戳
func StrToTimestamp(dateStr string) (int, error) {
	location, err := time.ParseInLocation(time.DateTime, dateStr, time.Local)
	if err != nil {
		return 0, err
	}
	return int(location.Unix()), nil
}

// 获取当前时间戳
func GetNowTimestamp() int64 {
	return time.Now().Unix()
}

// 获取当前日期 字符串 "2006-01-02 15:04:05"
func GetNowStr() string {
	return time.Now().Format(time.DateTime)
}
