package bean

import (
	"github.com/gin-gonic/gin"
)

type (
	Listener interface {
		// OnIllegalRequest 判断请求非法时（不是WebSocket请求）
		OnIllegalRequest(ctx *gin.Context)

		// OnFailureToUpgrade 将连接升级为WebSocket协议失败时
		OnFailureToUpgrade(ctx *gin.Context, err error)

		OnHandshake(c *Channel)

		// OnMessage 收到 WebSocket客户端 发来的消息
		OnMessage(c *Channel, messageType int, data []byte)

		// OnCloseByFrontend WebSocket客户端主动关闭连接
		OnCloseByFrontend(c *Channel, code int, text string)

		// OnCloseByBackend WebSocket服务端主动关闭连接
		OnCloseByBackend(c *Channel)
	}
)
