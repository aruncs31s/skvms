package utils

import (
	"time"
)

// Asia/Kolkata is the timezone for India, which is UTC+5:30

func GetBeginningOfDay() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}
func GetEndOfDay() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	year, month, day := now.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, loc)
}
