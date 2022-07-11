package utils

import "time"

func TodayBegin() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func IsToday(t time.Time) bool {
	now := time.Now()
	return now.Year() == t.Year() &&
		now.Month() == t.Month() &&
		now.Day() == t.Day()
}
