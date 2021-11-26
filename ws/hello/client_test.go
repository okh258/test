package hello

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func TestWsClient_msg(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	//uid := "110"
	uid := "500944337634304"
	conn, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://tm-uat.puttinggreenz.com/admin-api-server/ws?uid=%v", uid), nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close(websocket.StatusInternalError, "内部错误！")

	var message map[string]interface{}
	err = wsjson.Read(ctx, conn, &message)
	if err != nil {
		log.Println("receive msg error:", err)
		return
	}
	fmt.Printf("uid(%v)接收到服务端响应(%s)：%#v\n", uid, "d.Milliseconds()", message)

	conn.Close(websocket.StatusNormalClosure, "")
}

func TestWsClient_string(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	//uid := "110"
	uid := "500944337634304"
	c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://127.0.0.1:2022/ws?uid=%v", uid), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "内部错误！")

	err = wsjson.Write(ctx, c, "Hello WebSocket Server")
	if err != nil {
		panic(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("接收到服务端响应：%v\n", v)

	c.Close(websocket.StatusNormalClosure, "")
}
