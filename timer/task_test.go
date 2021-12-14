package timer

import (
	"context"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	//d, _ := time.Parse(time.RFC3339, "2021-12-08T06:59:59+08:00")
	d := time.Now()
	getCensus(context.TODO(), d)
}
