package rio

import (
	"net/http"
)

type (
	Listener interface {
		// OnIllegalRequest 判断请求非法时（不是WebSocket请求）
		OnIllegalRequest(w http.ResponseWriter, r *http.Request)

		// OnFailureToUpgrade 将连接升级为WebSocket协议失败时
		OnFailureToUpgrade(w http.ResponseWriter, r *http.Request, err error)

		OnHandshake(c *Channel)

		// OnMessage 收到 WebSocket客户端 发来的消息
		OnMessage(c *Channel, messageType int, data []byte)

		// OnCloseByFrontend WebSocket客户端主动关闭连接
		OnCloseByFrontend(c *Channel, code int, text string)

		// OnCloseByBackend WebSocket服务端主动关闭连接
		OnCloseByBackend(c *Channel)
	}
)
