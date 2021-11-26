package util

import "time"

func MicroTime() int64 {
	return time.Now().UnixNano() / 1000
}
