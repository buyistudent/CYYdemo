package utils

import "time"

const Data_Layout = "2006-01-02"
const DataTime_Layout = "2006-01-02 15:04:01"

func FormatToDateTimeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DataTime_Layout)
}

func FormatToDateStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(Data_Layout)
}
