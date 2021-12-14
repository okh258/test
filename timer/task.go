package timer

import (
	"fmt"
	"github.com/astaxie/beego/toolbox"
	"time"
)

func StartTimer() {
	// 获取统计数量, 每日 6,12,23 各执行一次
	getCensusTask := toolbox.NewTask("getCensus", "59 59 6,12,15,16,23 * * *", TestTYT())
	toolbox.AddTask("getCensus", getCensusTask)
	toolbox.StartTask()
}

func TestTYT() func() error {
	return func() error {
		fmt.Printf("start task, now: %v", time.Now().Format("2006-01-02 15:04:05"))
		//os.Exit(0)
		return nil
	}
}
