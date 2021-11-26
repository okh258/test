package ws

import (
	"context"
	"git.devops.com/wsim/hflib/logs"
	"time"
)

var (
	ticker     *time.Ticker
	htDuration int64 //秒
)

func init() {
	startCheckHt()
}

func startCheckHt() {
	ticker = time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				ReconnectingWebSocket()
			}
		}
	}()
}

func stopCheckHt() {
	if ticker != nil {
		ticker.Stop()
	}
}

// ReconnectingWebSocket 发送心跳
func ReconnectingWebSocket() {
	cxt := context.TODO()
	// 给所有在线用户发送心跳
	for key := range Broadcaster.users {
		err := Broadcaster.users[key].conn.Ping(cxt)
		if err != nil {
			logs.Infof(cxt, "ping failed, err: %v", err)
		}
	}
}
