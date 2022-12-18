package rio

import (
	"github.com/gorilla/websocket"
	"github.com/richelieu42/go-scales/src/jsonKit"
	"sync"
)

type (
	Channel struct {
		id   string
		conn *websocket.Conn

		lock *sync.Mutex

		bsId      string
		group     string
		userId    string
		listener  Listener
		extraData map[string]interface{}
	}
)

func (c *Channel) GetId() string {
	return c.id
}

// Close 后端主动关闭连接
func (c *Channel) Close() {
	_ = c.conn.Close()
	if Remove(c.id) {
		c.listener.OnCloseByBackend(c)
	}
}

// PushMessage 推送消息给WebSocket客户端.
/*
@param messageType websocket.TextMessage || websocket.BinaryMessage
*/
func (c *Channel) PushMessage(messageType int, data []byte) error {
	// 防止panic: concurrent write to websocket connection
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.conn.WriteMessage(messageType, data)
}

// PushJson 先序列化为json字符串，再推送给WebSocket客户端.
/*
@param messageType websocket.TextMessage || websocket.BinaryMessage
*/
func (c *Channel) PushJson(messageType int, obj interface{}) error {
	data, err := jsonKit.Marshal(obj)
	if err != nil {
		return err
	}

	return c.PushMessage(messageType, data)
}
