package base

import (
	"fmt"
	"test/util"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	a := util.Timestamp()
	d := time.Unix(a/1000, 0)
	t.Logf("d: %v", d.Format("2006-01-02 15:04:05"))
	t.Logf("d: %v", d.Format("2006年01月02日15时04分05秒"))
}

func TestGetZero(t *testing.T) {
	d := time.Now()
	newTime := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	fmt.Printf("date: %v, int: %v\n", newTime.Format("2006-01-02 15:04:05"), newTime.UnixNano()/1000000)
}

func TestMicroTime(t *testing.T) {
	fmt.Printf("is same day: %v\n", util.MicroTime())
}

func TestInSameDay(t *testing.T) {
	d, _ := time.Parse("2006-01-02 15:04:05", "2021-12-03 18:04:05")
	flag := InSameDay(d, time.Now())
	fmt.Printf("is same day: %v\n", flag)
}

func InSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}
