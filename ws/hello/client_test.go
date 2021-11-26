package hello

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"git.devops.com/wsim/hflib/jwt"
	gwebsocket "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func TestUat(t *testing.T) {
	var wsUrl = "wss://tm-uat.puttinggreenz.com/admin-api-server/ws"
	c1, err := newChatClient(wsUrl,1)
	assert.NoError(t, err)
	assert.NotNil(t, c1)
}

func New(url string, uid int64, header http.Header) (*ChatClient, error) {
	conn, r, err := gwebsocket.DefaultDialer.Dial(url, header)
	if err != nil {
		body, _ := ioutil.ReadAll(r.Body)
		return nil, fmt.Errorf("Fail to connect %w status: %v %v body: %s. Request Header: %v ", err, r.Status, r.Header, body, header)
	}
	c := &ChatClient{uid: uid, conn: conn}
	go c.loopMessage()
	return c, nil
}

func newChatClient(wsUrl string, userId int64) (*ChatClient, error) {
	osType := "golang testing"
	deviceToken := fmt.Sprintf("test_device_%d", userId)
	accessToken, _ := jwt.MakeWaToken(userId, osType, deviceToken)
	header := http.Header{}
	header.Set("x-wa-request-id", fmt.Sprintf("%d", time.Now().UnixNano()))
	header.Set("x-os-type", osType)
	header.Set("x-device-token", deviceToken)
	header.Set("x-access-token", accessToken)
	return New(wsUrl, userId, header)
}

func (c *ChatClient) loopMessage() {
	messageType, buf, err := c.conn.ReadMessage()
	for err == nil {
		fmt.Println("Got", messageType, string(buf))
		messageType, buf, err = c.conn.ReadMessage()
	}
}

type ChatClient struct {
	uid  int64
	conn *gwebsocket.Conn
}

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
