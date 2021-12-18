package base

import (
	"math"
	"testing"
)

func TestGetRate(t *testing.T) {
	t.Logf("num: %v\n", GetRate(10, 100, true))
	t.Logf("pageSize: %v\n", GetPageSize(101, 10))
}

// GetRate 获取率
// molecule 分子, denominator 分母, reverse 是否取反
func GetRate(molecule, denominator int64, reverse bool) float64 {
	rate := float64(molecule) / float64(denominator)
	if reverse {
		rate = 1 - rate
	}
	return rate
}

// GetPageSize 获取总页数
func GetPageSize(total, num int64) float64 {
	return math.Ceil(float64(total) / float64(num))
}
