package base

import (
	"strconv"
	"testing"
)

func TestStringParseFloat(t *testing.T) {

	totalAmount, err := strconv.ParseFloat("0.01", 64)
	if err != nil {
		t.Fatalf("金额转换失败: %v", err.Error())
		return
	}
	t.Logf("amout: %v", totalAmount)
}
