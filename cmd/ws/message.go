package ws

type Message struct {
	User    *WsUser     `json:"user"`
	Data    interface{} `json:"data"`
	MsgTime int64       `json:"msg_time"`
}
