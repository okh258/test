package routers

import (
	"github.com/astaxie/beego"
	"test/cmd/controllers"
	"test/cmd/ws"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &ws.WebSocketHandle{}, "get:Join")
}
