package routes

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/http/websocket"
	websocket2 "github.com/qbhy/goal-piplin/app/websocket"
)

func WebSocket(router contracts.HttpRouter, engine contracts.HttpEngine) {
	engine.Static("/", "/")

	router.Get("/ws-demo", websocket.New(websocket2.DemoController{}))

	router.Get("/ws", websocket.Default(func(frame contracts.WebSocketFrame) {

		fmt.Println("收到消息", frame.RawString(), frame.Connection().Fd())
		_ = frame.Send("来自服务器的回复1")
		_ = frame.SendBytes([]byte("来自服务器的回复2"))

	}))

}
