package ws

import (
	"context"
	"fmt"
	"git.devops.com/wsim/hflib/logs"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"testing"
)

var (
	addr = ":2022"
)

func TestWsServer(t *testing.T) {
	RegisterHandle()

	t.Fatal(http.ListenAndServe(addr, nil))
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
		logs.Errorf(req.Context(), "websocket accept failed, err: %v", err)
		return
	}

	// 1. uid 是否合法
	uid := req.FormValue("uid")
	if l := len(uid); l < 4 || l > 64 {
		logs.Infof(req.Context(), "uid illegal, uid: %v", uid)
		wsjson.Write(req.Context(), conn, fmt.Sprintf("uid illegal, uid: %v", uid))
		conn.Close(websocket.StatusUnsupportedData, fmt.Sprintf("uid illegal, uid: %v", uid))
		return
	}
	logs.Infof(req.Context(), "uid len: %v", len(uid))
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
	user := NewUser(conn, uid, req.RemoteAddr, 10)

	// 2. 开启给用户发送消息的协程
	go user.SendMessage(context.TODO())

	// 3. 给当前用户发送 echo 信息
	user.MessageChannel <- user.NewMsg("connect ok...")

	// 4. 将该用户加入广播器的用户中
	Broadcaster.UserEntering(user)
	logs.Infof(req.Context(), "user online, uid: %v", uid)
	// 5. 用户后续处理
	go UserHandle(user, conn, uid)
}

func UserHandle(user *User, conn *websocket.Conn, uid string) {
	ctx := context.TODO()
	// 1. 接收用户消息
	err := user.ReceiveMessage(ctx)

	// 2. 用户离开
	Broadcaster.UserLeaving(user)
	logs.Infof(ctx, "user offline, uid: %v", uid)

	if err == nil {
		closeErr := conn.Close(websocket.StatusNormalClosure, "err is empty")
		if closeErr != nil {
			logs.Errorf(ctx, "conn close failed, err: %v", closeErr)
		}
	} else {
		logs.Errorf(ctx, "read error from client, err: %v", err)
		closeErr := conn.Close(websocket.StatusInternalError, fmt.Sprintf("read error from client, err: %v", err))
		if closeErr != nil {
			logs.Errorf(ctx, "conn close failed, err: %v", closeErr)
		}
	}
}
