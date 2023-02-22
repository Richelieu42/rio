package bean

import (
	"github.com/gorilla/websocket"
	"github.com/richelieu42/go-scales/src/idKit"
	"github.com/richelieu42/go-scales/src/jsonKit"
	"github.com/richelieu42/rio/src/rio/manager"
	"sync"
)

type (
	Channel struct {
		id   string
		conn *websocket.Conn

		// lock 向前端推送消息时会加解锁
		lock *sync.Mutex

		bsId      string
		group     string
		user      string
		listener  Listener
		extraData map[string]interface{}
	}
)

func NewChannel(conn *websocket.Conn, listener Listener) *Channel {
	return &Channel{
		id:       idKit.NewULID(),
		conn:     conn,
		lock:     new(sync.Mutex),
		listener: listener,
	}
}

func (channel *Channel) GetId() string {
	return channel.id
}

func (channel *Channel) GetGroup() string {
	return channel.group
}

func (channel *Channel) SetGroup(group string) {
	channel.group = group
}

func (channel *Channel) GetUser() string {
	return channel.user
}

func (channel *Channel) SetUser(user string) {
	channel.user = user
}

func (channel *Channel) GetBsId() string {
	return channel.bsId
}

func (channel *Channel) SetBsId(bsId string) {
	channel.bsId = bsId
}

func (channel *Channel) GetExtraData() map[string]interface{} {
	return channel.extraData
}

func (channel *Channel) SetExtraData(extraData map[string]interface{}) {
	channel.extraData = extraData
}

func (channel *Channel) GetListener() Listener {
	return channel.listener
}

// Close 后端主动关闭连接
func (channel *Channel) Close() {
	_ = channel.conn.Close()
	if manager.RemoveChannel(channel) {
		channel.listener.OnCloseByBackend(channel)
	}
}

// PushMessage 推送消息给WebSocket客户端.
/*
@param messageType websocket.TextMessage || websocket.BinaryMessage
*/
func (channel *Channel) PushMessage(messageType int, data []byte) error {
	// 防止panic: concurrent write to websocket connection
	channel.lock.Lock()
	defer channel.lock.Unlock()

	return channel.conn.WriteMessage(messageType, data)
}

// PushJson 先序列化为json字符串，再推送给WebSocket客户端.
/*
@param messageType websocket.TextMessage || websocket.BinaryMessage
*/
func (channel *Channel) PushJson(messageType int, obj interface{}) error {
	data, err := jsonKit.Marshal(obj)
	if err != nil {
		return err
	}

	return channel.PushMessage(messageType, data)
}
