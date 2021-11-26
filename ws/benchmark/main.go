package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"test/model"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	userNum       int           // 用户数
	loginInterval time.Duration // 用户登录时间间隔
	msgInterval   time.Duration // 同一个用户发送消息间隔
)

func init() {
	flag.IntVar(&userNum, "u", 10, "登录用户数")
	flag.DurationVar(&loginInterval, "l", 100*time.Millisecond, "用户陆续登录时间间隔")
	flag.DurationVar(&msgInterval, "m", 10*time.Second, "用户发送消息时间间隔")
}

func main() {
	flag.Parse()

	for i := 0; i < userNum; i++ {
		go UserConnect("user" + strconv.Itoa(i))
		time.Sleep(loginInterval)
	}
	log.Println("user connect ok...")

	select {}
}

func UserConnect(uid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://127.0.0.1:2022/ws?uid=%v", uid), nil)
	if err != nil {
		log.Printf("Dial failed, uid: %v, err:%v\b", uid, err)
		os.Exit(1)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "内部错误！")

	//go sendMessage(conn, uid)

	ctx = context.Background()

	for {
		var message *model.Message
		err = wsjson.Read(ctx, conn, &message)
		if err != nil {
			log.Println("receive msg error:", err)
			return
		}

		//if message.ClientSendTime.IsZero() {
		//	continue
		//}
		if d := time.Now().Sub(message.ClientSendTime); d > 1*time.Second {
			fmt.Printf("uid(%v)接收到服务端响应(%d)：%#v\n", uid, d.Milliseconds(), message)
		}
	}

	conn.Close(websocket.StatusNormalClosure, "")
}

func sendMessage(conn *websocket.Conn, nickname string) {
	ctx := context.Background()
	i := 1
	for {
		msg := map[string]string{
			"content":   "来自" + nickname + "的消息:" + strconv.Itoa(i),
			"send_time": strconv.FormatInt(time.Now().UnixNano(), 10),
		}
		err := wsjson.Write(ctx, conn, msg)
		if err != nil {
			log.Println("send msg error:", err, "nickname:", nickname, "no:", i)
		}
		i++

		time.Sleep(msgInterval)
	}
}
