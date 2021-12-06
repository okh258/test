package base

import (
	"testing"
)

func TestGetRate(t *testing.T) {
	t.Logf("num: %v\n", GetRate(55, 100, true))
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
