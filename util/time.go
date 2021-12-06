package util

import "time"

func MicroTime() int64 {
	return time.Now().UnixNano() / 1000
}

func Timestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
