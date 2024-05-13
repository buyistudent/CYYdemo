package shared

import (
	"regexp"
	"time"
)

const TIME_LAYOUT = "2006-01-02 15:04:01"
const DATA_LAYOUT = "2006-01-02"
const MOYTH_LAYOUT = "2006-01"

// 字符串时间格式化
func Format(dateStr string) time.Time {
	parse, _ := time.Parse(TIME_LAYOUT, dateStr)
	return parse
}

// time格式转字符串时间
func FormatTimeToStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TIME_LAYOUT)
}

// time格式转字符串时间
func FormatToStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DATA_LAYOUT)
}

func FormatToDateTimeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TIME_LAYOUT)
}

// 判断两个时间大小
// -1 ：t1小于t2  0:t1等于t2  1：t1大于t2
func CompareData(time1 time.Time, time2 time.Time) int {
	format1 := time1.Format(DATA_LAYOUT)
	format2 := time2.Format(DATA_LAYOUT)
	t1, _ := time.Parse(DATA_LAYOUT, format1)
	t2, _ := time.Parse(DATA_LAYOUT, format2)
	var result int
	if t1.Equal(t2) {
		result = 0
	} else if t1.Before(t2) {
		result = -1
	} else if t1.After(t2) {
		result = 1
	}
	return result
}

// 获取指定时间加一天
func FindDateAddDay(t time.Time) time.Time {
	nowT := t.Format(DATA_LAYOUT)
	parse, _ := time.Parse(DATA_LAYOUT, nowT)
	d, _ := time.ParseDuration("24h")
	times := parse.Add(d)
	return times
}

// 获取指定时间加一个月 减一天
func FindDateAddmonths(t time.Time) time.Time {

	times := t.AddDate(0, 1, -1)

	return times
}

// 获取指定时间减一天
func FindDateLessDay(t time.Time) time.Time {
	nowT := t.Format(DATA_LAYOUT)
	parse, _ := time.Parse(DATA_LAYOUT, nowT)
	d, _ := time.ParseDuration("-24h")
	times := parse.Add(d)
	return times
}

func IsNumber(str string) bool {
	pattern := "^[0-9]+$"
	match, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return match
}
