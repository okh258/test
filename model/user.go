package model

import (
	"context"
	"errors"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"test/util"
	"time"
)

// NewUser 创建一个新用户
// queueLen 发送消息缓冲队列长度
func NewUser(conn *websocket.Conn, uid string, remoteAddr string, queueLen int) *User {
	msgChan := make(chan *Message, queueLen)
	return &User{
		UID:            uid,
		EnterAt:        util.MicroTime(),
		Addr:           remoteAddr,
		MessageChannel: msgChan,
		conn:           conn,
	}
}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		err := wsjson.Write(ctx, u.conn, msg)
		if err != nil {
			log.Printf("send msg failed, uid: %v, msg: %+v\n", u.UID, msg)
		} else {
			log.Printf("%v send msg: %+v\n", u.UID, msg)
		}
	}
}

func (u *User) ReceiveMessage(ctx context.Context) error {
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
		log.Printf("%v receive msg: %+v\n", u.UID, receiveMsg)
		//Broadcaster.Broadcast(&Message{
		//	Content: receiveMsg["content"],
		//	MsgTime: util.MicroTime(),
		//})
		//for _, user := range Broadcaster.users {
		//	if user.UID == u.UID {
		//		continue
		//	}
		//	user.MessageChannel <- user.NewMsg(receiveMsg["content"])
		//	break
		//}
	}
}

// CloseMessageChannel 避免 goroutine 泄露
func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

func (u *User) NewMsg(msg string) *Message {
	return &Message{
		User:    u,
		Content: msg,
		MsgTime: util.MicroTime(),
	}
}

type User struct {
	UID            string        `json:"uid"`
	EnterAt        int64         `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn
}

type Message struct {
	User    *User  `json:"user"`
	Type    int    `json:"type"`
	Content string `json:"content"`
	MsgTime int64  `json:"msg_time"`

	ClientSendTime time.Time `json:"client_send_time"`

	Users map[string]*User `json:"users"`
}
