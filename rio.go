package rio

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/richelieu42/go-scales/src/core/errorKit"
	"github.com/richelieu42/rio/src/manager"
	"net/http"
	"time"
)

type (
	Listener interface {
		// OnIllegalRequest 判断请求非法时（不是WebSocket请求）
		OnIllegalRequest(ctx *gin.Context)

		// OnFailureToUpgrade 将连接升级为WebSocket协议失败时
		OnFailureToUpgrade(ctx *gin.Context, err error)

		OnHandshake(c *manager.Channel)

		// OnMessage 收到 WebSocket客户端 发来的消息
		OnMessage(c *manager.Channel, messageType int, data []byte)

		// OnCloseByFrontEnd WebSocket连接断开（因为前端）时
		OnCloseByFrontEnd(c *manager.Channel)
	}
)

// NewGinHandler
/*
@param listener 可以为nil，但不推荐这么干
*/
func NewGinHandler(listener Listener) (gin.HandlerFunc, error) {
	if listener == nil {
		return nil, errorKit.Simple("listener == nil")
	}

	var upgrader = websocket.Upgrader{
		HandshakeTimeout: time.Second * 6,
		CheckOrigin: func(r *http.Request) bool {
			// 允许跨域
			return true
		},
	}

	return func(ctx *gin.Context) {
		// 先判断是不是websocket请求
		if !websocket.IsWebSocketUpgrade(ctx.Request) {
			listener.OnIllegalRequest(ctx)
			return
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Writer.Header())
		if err != nil {
			// 升级为WebSocket协议失败
			listener.OnFailureToUpgrade(ctx, err)
			return
		}
		// ！！！：下面一行代码至关重要，否则会导致WebSocket连接关不掉
		defer conn.Close()

		c := manager.WrapToChannel(conn)
		manager.Add(c)
		c.ReceiveMessages()
	}, nil
}
