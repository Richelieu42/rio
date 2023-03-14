package rio

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/richelieu42/go-scales/src/core/errorKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
	"github.com/richelieu42/go-scales/src/http/httpKit"
	"github.com/richelieu42/rio/src/rio/consts"
	"net/http"
	"time"
)

var (
	// upgrader 并发安全的
	upgrader = websocket.Upgrader{
		HandshakeTimeout: time.Second * 6,
		CheckOrigin: func(r *http.Request) bool {
			// 允许跨域
			return true
		},
	}
)

// NewHttpHandler
/*
@param listener 不能为nil
*/
func NewHttpHandler(listener Listener) (func(w http.ResponseWriter, r *http.Request), error) {
	if listener == nil {
		return nil, errorKit.Simple("param listener mustn't be nil")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// 先判断是不是websocket请求
		if !websocket.IsWebSocketUpgrade(r) {
			listener.OnIllegalRequest(w, r)
			return
		}

		// Upgrade（升级为WebSocket协议）
		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			listener.OnFailureToUpgrade(w, r, err)
			return
		}
		defer conn.Close()

		channel := NewChannel(conn, listener)
		/* 监听: WebSocket客户端主动关闭连接 */
		conn.SetCloseHandler(func(code int, text string) error {
			channel.SetClosed()

			if RemoveChannel(channel) {
				channel.GetListener().OnCloseByFrontend(channel, code, text)
			}

			// 默认的close handler
			message := websocket.FormatCloseMessage(code, text)
			_ = conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
			return nil
		})
		AddChannel(channel)
		listener.OnHandshake(channel)

		/* 绑定数据（通过url参数，有的话） */
		bsid := httpKit.GetUrlParam(r, consts.KeyBsid)
		user := httpKit.GetUrlParam(r, consts.KeyUser)
		group := httpKit.GetUrlParam(r, consts.KeyGroup)
		if !strKit.IsAllEmpty(bsid, user, group) {
			channel.BindData(bsid, user, group)
		}

		/* 接收WebSocket客户端发来的消息 */
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				channel.SetClosed()

				if RemoveChannel(channel) {
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

// NewGinHandler
/*
@param listener 不能为nil
*/
func NewGinHandler(listener Listener) (gin.HandlerFunc, error) {
	httpHandler, err := NewHttpHandler(listener)
	if err != nil {
		return nil, err
	}

	return func(ctx *gin.Context) {
		httpHandler(ctx.Writer, ctx.Request)
	}, nil
}
