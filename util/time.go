package util

import (
	"strconv"
	"time"
)

func MicroTime() int64 {
	return time.Now().UnixNano() / 1000
}

func Timestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func ToInt64(v interface{}, defaultValue int64) int64 {
	if v == nil {
		return defaultValue
	}

	switch v := v.(type) {
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaultValue
		}
		return i
	case int8:
		return int64(v)
	case uint8:
		return int64(v)
	case uint64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case uint32:
		return int64(v)
	case float64:
		return int64(v)
	case []byte:
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return defaultValue
		}
		return i
	case nil:
		return 0
	}
	return defaultValue
}
