package ws

import (
	"context"

	"git.devops.com/wsim/hflib/logs"
)

// broadcaster 广播
type broadcaster struct {
	// 所有在线用户
	users map[int64]*WsUser

	// 所有 channel 统一管理，可以避免外部乱用
	enteringChannel chan *WsUser
	leavingChannel  chan *WsUser
	messageChannel  chan *Message

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*WsUser
}

var Broadcaster = &broadcaster{
	users: make(map[int64]*WsUser),

	enteringChannel: make(chan *WsUser),
	leavingChannel:  make(chan *WsUser),
	messageChannel:  make(chan *Message, MessageQueueLen),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*WsUser),
}

// Start 启动广播
func (b *broadcaster) Start() {
	logs.Infof(context.Background(), "broadcaster start... ")
	for {
		select {
		case user := <-b.enteringChannel:
			// 新用户进入
			b.users[user.UID] = user
		case user := <-b.leavingChannel:
			// 用户离开
			delete(b.users, user.UID)
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:
			if len(b.users) < 1 {
				logs.Warnf(context.Background(), "online user num: %v", len(b.users))
				continue
			}
			// 给所有在线用户发送消息
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
		case <-b.requestUsersChannel:
			userList := make([]*WsUser, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) UserEntering(u *WsUser) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *WsUser) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	if len(b.messageChannel) >= MessageQueueLen {
		logs.Errorf(context.Background(), "broadcast queue is full...")
	}
	b.messageChannel <- msg
}

func (b *broadcaster) GetUserList() []*WsUser {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
