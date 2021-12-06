package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"test/models"
)

var (
	addr = ":2022"
)

func main() {
	RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}

func RegisterHandle() {
	// 广播消息处理
	go Broadcaster.Start()

	// http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Printf("websocket accept failed, err: %v\n", err)
		return
	}

	// 1. uid 是否合法
	uid := req.FormValue("uid")
	if l := len(uid); l < 4 || l > 64 {
		log.Printf("uid illegal, uid: %v\n", uid)
		wsjson.Write(req.Context(), conn, fmt.Sprintf("uid illegal, uid: %v", uid))
		conn.Close(websocket.StatusUnsupportedData, fmt.Sprintf("uid illegal, uid: %v", uid))
		return
	}
	log.Printf("uid len: %v\n", len(uid))
	// 2. 校验用户权限
	if !UserAuth(uid) {
		conn.Close(websocket.StatusUnsupportedData, fmt.Sprintf("uid illegal, uid: %v", uid))
	}
	// 3. 接收用户, 并创建相关资源
	AcceptUser(conn, uid, req)
}

// 用户鉴权
func UserAuth(uid string) bool {
	if len(uid) > 0 {
		return true
	}
	return false
}

// 接收用户
func AcceptUser(conn *websocket.Conn, uid string, req *http.Request) {
	// 1. 创建用户
	user := models.NewUser(conn, uid, req.RemoteAddr, 10)

	// 2. 开启给用户发送消息的协程
	go user.SendMessage(context.TODO())

	// 3. 给当前用户发送 echo 信息
	user.MessageChannel <- user.NewMsg("connect ok...")

	// 4. 将该用户加入广播器的用户中
	Broadcaster.UserEntering(user)
	log.Printf("user online, uid: %v\n", uid)
	// 5. 用户后续处理
	go UserHandle(user, conn, uid)
}

func UserHandle(user *models.User, conn *websocket.Conn, uid string) {
	ctx := context.TODO()
	// 1. 接收用户消息
	err := user.ReceiveMessage(ctx)

	// 2. 用户离开
	Broadcaster.UserLeaving(user)
	log.Printf("user offline, uid: %v\n", uid)

	if err == nil {
		closeErr := conn.Close(websocket.StatusNormalClosure, "err is empty")
		if closeErr != nil {
			log.Printf("conn close failed, err: %v\n", closeErr)
		}
	} else {
		log.Printf("read error from client, err: %v\n", err)
		closeErr := conn.Close(websocket.StatusInternalError, fmt.Sprintf("read error from client, err: %v", err))
		if closeErr != nil {
			log.Printf("conn close failed, err: %v\n", closeErr)
		}
	}
}
