package rio

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/richelieu42/go-scales/src/assertKit"
	"github.com/richelieu42/go-scales/src/idKit"
	"net/http"
	"sync"
	"time"
)

// NewGinHandler
/*
@param listener 不能为nil
*/
func NewGinHandler(listener Listener) (gin.HandlerFunc, error) {
	if err := assertKit.NotNil(listener, "listener"); err != nil {
		return nil, err
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
		defer conn.Close()

		c := &Channel{
			id:       idKit.NewULID(),
			conn:     conn,
			lock:     new(sync.Mutex),
			listener: listener,
		}
		/* 监听: WebSocket客户端主动关闭连接 */
		conn.SetCloseHandler(func(code int, text string) error {
			if Remove(c.id) {
				c.listener.OnCloseByFrontend(c, code, text)
			}

			// 默认的close handler
			message := websocket.FormatCloseMessage(code, text)
			_ = conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
			return nil
		})
		Add(c)
		listener.OnHandshake(c)
		/* 接收WebSocket客户端发来的消息 */
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				if Remove(c.id) {
					if ce, ok := err.(*websocket.CloseError); ok {
						listener.OnCloseByFrontend(c, ce.Code, ce.Text)
					} else {
						listener.OnCloseByBackend(c)
					}
				}
				break
			}
			listener.OnMessage(c, messageType, p)
		}
	}, nil
}
