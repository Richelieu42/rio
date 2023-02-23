package rio

import (
	"github.com/gorilla/websocket"
	"github.com/richelieu42/go-scales/src/core/errorKit"
	"github.com/richelieu42/go-scales/src/idKit"
	"sync"
)

type (
	Channel struct {
		// id 唯一id
		id string
		// lock 向前端推送消息时会用到
		lock *sync.Mutex

		conn *websocket.Conn

		bsId      string
		group     string
		user      string
		listener  Listener
		extraData map[string]interface{}

		// closed conn是否已经断掉？
		closed bool
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

func (channel *Channel) SetClosed() {
	channel.closed = true
}

// PushMessage 推送 文本消息 给浏览器
/*
@param messageType websocket.TextMessage || websocket.BinaryMessage
*/
func (channel *Channel) PushMessage(messageType int, data []byte) error {
	if len(data) == 0 {
		// 无效的推送
		return nil
	}

	if channel.closed {
		return errorKit.Simple("conn has already been closed")
	}

	// 防止panic: concurrent write to websocket connection
	channel.lock.Lock()
	defer channel.lock.Unlock()

	if channel.closed {
		return errorKit.Simple("conn has already been closed")
	}

	return channel.conn.WriteMessage(messageType, data)
}

// Close 后端主动关闭连接
func (channel *Channel) Close() {
	if channel.closed {
		return
	}

	channel.lock.Lock()
	defer channel.lock.Unlock()

	if channel.closed {
		return
	}
	channel.SetClosed()

	_ = channel.conn.Close()
	if RemoveChannel(channel) {
		channel.listener.OnCloseByBackend(channel)
	}
}

func (channel *Channel) BindBsid(bsid string) {
	BindBsId(channel, bsid)
}

func (channel *Channel) BindUser(user string) {
	BindUser(channel, user)
}

func (channel *Channel) BindGroup(group string) {
	BindGroup(channel, group)
}
