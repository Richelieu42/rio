package manager

import (
	"github.com/gorilla/websocket"
	"github.com/richelieu42/go-scales/src/core/errorKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
	"github.com/richelieu42/go-scales/src/idKit"
	"sync"
	"time"
)

type (
	Channel struct {
		id   string
		conn *websocket.Conn
		// true: 当前websocket连接已经被关闭（断开）
		closed bool

		lock *sync.Mutex

		bsId   string
		group  string
		userId string
	}
)

func (c *Channel) GetId() string {
	return c.id
}

// Close TODO: 后端主动关闭连接
func (c *Channel) Close() {
	c.closed = true

	_ = c.conn.Close()
}

// ReceiveMessages 接收 WebSocket客户端 发来的消息（会阻塞直至连接断开）
func (c *Channel) ReceiveMessages(listener Listener) {
	for {
		if c.closed {
			break
		}
		messageType, p, err := c.conn.ReadMessage()
		if err != nil {
			c.closed = true
			break
		}

		if listener != nil {
			listener.OnMessage(c, messageType, p)
		}
	}
}

// PushMessage 推送消息给 WebSocket客户端
/*
@param messageType websocket.TextMessage || websocket.BinaryMessage
*/
func (c *Channel) PushMessage(messageType int, data []byte) error {
	// 防止panic: concurrent write to websocket connection
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return errorKit.Simple("Channel(id: %s) is already closed", c.id)
	}

	err := c.conn.WriteMessage(messageType, data)
	if err != nil {
		c.closed = true
	}
	return err
}

func WrapToChannel(conn *websocket.Conn, listener Listener) *Channel {
	id := idKit.NewULID()
	c := &Channel{
		id:     id,
		conn:   conn,
		closed: false,
		lock:   new(sync.Mutex),
	}

	conn.SetCloseHandler(func(code int, text string) error {
		c.closed = true

		reason := strKit.Format("closure by frontend with code(%d) and text(%s)", code, text)
		Remove(id, reason)

		if listener != nil {
			listener.OnCloseByFrontend(c, code, text)
		}

		// 默认的close handler
		message := websocket.FormatCloseMessage(code, text)
		_ = conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})
	return c
}
