package main

import (
	"log"
	"test/model"
)

var MessageQueueLen = 1024

// broadcaster 广播器
type broadcaster struct {
	// 所有在线用户
	users map[string]*model.User

	// 所有 channel 统一管理，可以避免外部乱用
	enteringChannel chan *model.User
	leavingChannel  chan *model.User
	messageChannel  chan *model.Message

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*model.User
}

var Broadcaster = &broadcaster{
	users: make(map[string]*model.User),

	enteringChannel: make(chan *model.User),
	leavingChannel:  make(chan *model.User),
	messageChannel:  make(chan *model.Message, MessageQueueLen),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*model.User),
}

// Start 启动广播器
func (b *broadcaster) Start() {
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
			// 给所有在线用户发送消息
			for _, user := range b.users {
				if msg.User != nil && user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
		case <-b.requestUsersChannel:
			userList := make([]*model.User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) UserEntering(u *model.User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *model.User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *model.Message) {
	if len(b.messageChannel) >= MessageQueueLen {
		log.Println("broadcast queue is full...")
	}
	b.messageChannel <- msg
}

func (b *broadcaster) GetUserList() []*model.User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
