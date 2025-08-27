package utils

import "time"

var TimeFormat = "2006-01-02T15:04:05-07:00"

func TimeToString(t time.Time) string {
	return t.Format(TimeFormat)
}
