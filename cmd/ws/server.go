package ws

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"git.devops.com/wsim/hflib/logs"
)

type WebSocketHandle struct {
	beego.Controller
}

var (
	MessageQueueLen     = 1024
	UserMessageQueueLen = 64
)

//// StartWsServer 开启 WebSocket 服务
//func StartWsServer() {
//	RegisterHandle()
//
//	err := http.ListenAndServe(addr, nil)
//	if err != nil {
//		logs.Errorf(context.Background(), "start websocket failed, err: %v", err)
//		os.Exit(1)
//	}
//}

//func RegisterHandle(w http.ResponseWriter, req *http.Request) {
//	// 广播消息处理
//	go Broadcaster.Start()
//
//	// http.HandleFunc("/", homeHandleFunc)
//	//http.HandleFunc("/ws", WebSocketHandleFunc)
//}

func init() {
	go Broadcaster.Start()
}

func (c *WebSocketHandle) Join() {
	req, w := c.Ctx.Request, c.Ctx.ResponseWriter
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		logs.Errorf(req.Context(), "websocket accept failed, err: %v", err)
		return
	}

	// 1. uid 是否合法
	uid := req.FormValue("uid")
	if l := len(uid); l < 15 || l > 64 {
		logs.Infof(req.Context(), "uid len illegal, uid: %v", uid)
		writeErr := wsjson.Write(req.Context(), conn, fmt.Sprintf("uid len illegal, uid: %v", uid))
		if writeErr != nil {
			logs.Errorf(req.Context(), "websocket write json failed, err: %v", writeErr)
		}
		closeErr := conn.Close(websocket.StatusUnsupportedData, fmt.Sprintf("uid len illegal, uid: %v", uid))
		if closeErr != nil {
			logs.Errorf(req.Context(), "conn close failed, err: %v", closeErr)
		}
		return
	}
	userInfo, flag := UserAuth(uid, req)
	// 2. 校验用户权限
	if !flag {
		closeErr := conn.Close(websocket.StatusUnsupportedData, fmt.Sprintf("uid illegal, uid: %v", uid))
		if closeErr != nil {
			logs.Errorf(req.Context(), "conn close failed, err: %v", closeErr)
		}
	}
	// 3. 接收用户, 并创建相关资源
	AcceptUser(conn, req, userInfo)
}

// UserAuth 用户鉴权
func UserAuth(uidStr string, req *http.Request) (*UserBaseInfo, bool) {
	if len(uidStr) < 1 {
		return nil, false
	}
	uid, err := strconv.ParseInt(uidStr, 0, 64)
	if err != nil {
		logs.Errorf(req.Context(), "uid parseInt failed, uid: %v, err: %v", uid, err)
		return nil, false
	}

	//userInfo, err := serviceses.NewUserService().GetUserInfo(req.Context(), uid)
	//if err != nil {
	//	logs.Errorf(req.Context(), "get user info failed, err: %v", err)
	//	return nil, false
	//}
	return &UserBaseInfo{
		UserId:   500944337634304,
		Nickname: "test",
	}, true
}

// AcceptUser 接收用户
func AcceptUser(conn *websocket.Conn, req *http.Request, userInfo *UserBaseInfo) {
	// 1. 创建用户
	user := NewUser(conn, userInfo, req.RemoteAddr, UserMessageQueueLen)

	// 2. 开启给用户发送消息的协程
	go user.SendMessage(req.Context())

	// 3. 给当前用户发送 echo 信息
	user.MessageChannel <- user.NewMsg("connect ok.")

	// 4. 将该用户加入广播器的用户中
	Broadcaster.UserEntering(user)
	logs.Infof(req.Context(), "user online, uid: %v", userInfo.UserId)

	// 5. 接收用户消息
	err := user.ReceiveMessage(req.Context())

	// 6. 用户离开
	Broadcaster.UserLeaving(user)
	logs.Infof(req.Context(), "user offline, uid: %v", userInfo.UserId)

	if err == nil {
		closeErr := conn.Close(websocket.StatusNormalClosure, "err is empty")
		if closeErr != nil {
			logs.Errorf(req.Context(), "conn close failed, err: %v", closeErr)
		}
	} else {
		logs.Errorf(req.Context(), "read error from client, err: %v", err)
		closeErr := conn.Close(websocket.StatusInternalError, fmt.Sprintf("read error from client, err: %v", err))
		if closeErr != nil {
			logs.Errorf(req.Context(), "conn close failed, err: %v", closeErr)
		}
	}
}
