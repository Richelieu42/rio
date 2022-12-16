package rio

import (
	"github.com/gin-gonic/gin"
	"github.com/richelieu42/rio/src/manager"
)

type (
	Listener interface {
		// OnIllegalRequest 判断请求非法时（不是WebSocket请求）
		OnIllegalRequest(ctx *gin.Context)

		// OnFailureToUpgrade 将连接升级为WebSocket协议失败时
		OnFailureToUpgrade(ctx *gin.Context)

		OnHandshake(c *manager.Channel)

		// OnMessage 收到 WebSocket客户端 发来的消息
		OnMessage(c *manager.Channel, messageType int, data []byte)

		// OnCloseByFrontEnd WebSocket连接断开（因为前端）时
		OnCloseByFrontEnd(c *manager.Channel)
	}
)

// NewGinHandler
/*
@param



*/
func NewGinHandler(listener Listener) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
