package ws

import (
	"context"
	"errors"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"git.devops.com/go/golib/util"
	"git.devops.com/wsim/hflib/logs"
)

type WsUser struct {
	UID            int64         `json:"uid"`
	Nickname       string        `json:"nickname"`
	EnterAt        int64         `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn
}

type UserBaseInfo struct {
	Avatar              string `json:"avatar"`
	UserId              int64  `json:"user_id"`
	Nickname            string `json:"nickname"`
	Gender              int32  `json:"gender"`
	Mobile              string `json:"mobile"`
	Usernumber          int64  `json:"usernumber"`
	Username            string `json:"username"`
	Age                 int32  `json:"age"`
	Place               string `json:"place"`
	Birthday            int64  `json:"birthday"`
	VipLevel            int32  `json:"vip_level"`
	VipExpireTime       int64  `json:"vip_expire_time"`
	UserStatus          int32  `json:"user_status"`
	TextSignature       string `json:"text_signature"`
	CertifyStatus       int32  `json:"certify_status"`
	UserIdentity        int32  `json:"user_identity"`
	RegisterTime        int64  `json:"register_time"`
	LastVipBuyTime      int64  `json:"last_vip_buy_time"`
	RegisterIp          string `json:"register_ip"`
	RegisterDevice      string `json:"register_device"`
	LatLng              string `json:"lat_lng"`
	LoginIp             string `json:"login_ip"`
	LoginTime           int64  `json:"login_time"`
	LoginOS             string `json:"login_os"`
	BelongAgentUid      int64  `json:"belong_agent_uid"`
	BelongAgentNickname string `json:"belong_agent_nickname"`
	InviteCount         int64  `json:"invite_count"`
}

// NewUser 创建一个新用户
// queueLen 发送消息缓冲队列长度
func NewUser(conn *websocket.Conn, userInfo *UserBaseInfo, remoteAddr string, queueLen int) *WsUser {
	msgChan := make(chan *Message, queueLen)
	return &WsUser{
		EnterAt:        util.MicroTime(),
		UID:            userInfo.UserId,
		Nickname:       userInfo.Nickname,
		MessageChannel: msgChan,
		Addr:           remoteAddr,
		conn:           conn,
	}
}

func (u *WsUser) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		err := wsjson.Write(ctx, u.conn, msg)
		if err != nil {
			logs.Errorf(ctx, "send msg failed, uid: %v, msg: %+v", u.UID, msg)
		}
	}
}

func (u *WsUser) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			// 判定连接是否关闭了，正常关闭，不认为是错误
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}

			return err
		}
	}
}

// CloseMessageChannel 避免 goroutine 泄露
func (u *WsUser) CloseMessageChannel() {
	close(u.MessageChannel)
}

func (u *WsUser) NewMsg(data interface{}) *Message {
	return &Message{
		User:    u,
		Data:    data,
		MsgTime: util.MicroTime(),
	}
}
