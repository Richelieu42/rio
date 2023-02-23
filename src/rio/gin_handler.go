package rio

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/richelieu42/go-scales/src/core/errorKit"
	"github.com/richelieu42/rio/src/rio/bean"
	"github.com/richelieu42/rio/src/rio/consts"
	"github.com/richelieu42/rio/src/rio/manager"
	"net/http"
	"time"
)

// NewGinHandler
/*
@param listener 不能为nil
*/
func NewGinHandler(listener bean.Listener) (gin.HandlerFunc, error) {
	if listener == nil {
		return nil, errorKit.Simple("param listener mustn't be nil")
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

		// Upgrade（升级为WebSocket协议）
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Writer.Header())
		if err != nil {
			listener.OnFailureToUpgrade(ctx, err)
			return
		}
		defer conn.Close()

		channel := bean.NewChannel(conn, listener)
		/* 监听: WebSocket客户端主动关闭连接 */
		conn.SetCloseHandler(func(code int, text string) error {
			channel.SetClosed()

			if manager.RemoveChannel(channel) {
				channel.GetListener().OnCloseByFrontend(channel, code, text)
			}

			// 默认的close handler
			message := websocket.FormatCloseMessage(code, text)
			_ = conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
			return nil
		})
		manager.AddChannel(channel)
		listener.OnHandshake(channel)

		/* 绑定数据（通过url参数） */
		if bsid, ok := ctx.GetQuery(consts.KeyBsId); ok {
			channel.BindBsid(bsid)
		}
		if user, ok := ctx.GetQuery(consts.KeyUser); ok {
			channel.BindUser(user)
		}
		if group, ok := ctx.GetQuery(consts.KeyGroup); ok {
			channel.BindGroup(group)
		}

		/* 接收WebSocket客户端发来的消息 */
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				channel.SetClosed()

				if manager.RemoveChannel(channel) {
					if closeErr, ok := err.(*websocket.CloseError); ok {
						listener.OnCloseByFrontend(channel, closeErr.Code, closeErr.Text)
					} else {
						listener.OnCloseByBackend(channel)
					}
				}
				break
			}
			listener.OnMessage(channel, messageType, data)
		}
	}, nil
}
