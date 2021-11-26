package base

import (
	"fmt"
	"testing"
	"time"
)

func TestForChan(t *testing.T) {
	ch := make(chan int) //创建一个无缓存channel

	//新建一个goroutine
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(200*time.Millisecond)
			ch <- i // 往通道写数据
		}
		// 不需要再写数据时，关闭channel
		close(ch)
		//ch <- 666 // 关闭channel后无法再发送数据, panic

	}() //别忘了()

	for num := range ch {
		fmt.Println("num = ", num)
	}
	var n = <-ch // 关闭channel后再次读取数据, 读取零值
	fmt.Println("num = ", n)
}
