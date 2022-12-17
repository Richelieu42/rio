package rio

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// NewGinHandler
/*
@param listener 可以为nil，但不推荐这么干
*/
func NewGinHandler(listener Listener) (gin.HandlerFunc, error) {
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
			if listener != nil {
				listener.OnIllegalRequest(ctx)
			}
			return
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Writer.Header())
		if err != nil {
			// 升级为WebSocket协议失败
			if listener != nil {
				listener.OnFailureToUpgrade(ctx, err)
			}
			return
		}
		defer conn.Close()

		c := WrapToChannel(conn, listener)
		Add(c)
		if listener != nil {
			listener.OnHandshake(c)
		}

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				if Remove(c.id) {
					if listener != nil {
						listener.OnCloseByBackend(c)
					}
				}
				break
			}

			if listener != nil {
				listener.OnMessage(c, messageType, p)
			}
		}

	}, nil
}
