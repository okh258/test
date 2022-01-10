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

func TestStringLen(t *testing.T) {
	s := "随手一打，就是标准的十五个字。"
	t.Logf("s len: %v, 占用(%v)个字", len(s), Len(s))
	s = "随手一打,就是标准的十五个字."
	t.Logf("s len: %v, 占用(%v)个字", len(s), Len(s))
	s = "随手一打,就是标准的15个字."
	t.Logf("s len: %v, 占用(%v)个字", len(s), Len(s))
	s = "w shou y da 15. w shou y da 15"
	t.Logf("s len: %v, 占用(%v)个字", len(s), Len(s))
	s = "随手一打, 15 ge zi "
	t.Logf("s len: %v, 占用(%v)个字", len(s), Len(s))
}

func Len(s string) (num float64) {
	for range s {
		num++
	}
	return
}
