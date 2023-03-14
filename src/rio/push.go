package rio

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/sliceKit"
)

func PushToBsId(messageType int, data []byte, bsid string) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	channel := mapKit.Get(bsIdMap, bsid)
	if channel == nil {
		return
	}
	// TODO: 输出推送失败的error
	_ = channel.PushMessage(messageType, data)
}

func PushToUser(messageType int, data []byte, user string, exceptBsIds ...string) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	set := mapKit.Get(userMap, user)
	if set == nil {
		return
	}
	set.Each(func(channel *Channel) bool {
		if sliceKit.Contains(exceptBsIds, channel.GetBsId()) {
			return false
		}
		// TODO: 输出推送失败的error
		_ = channel.PushMessage(messageType, data)
		return false
	})
}

func PushToGroup(messageType int, data []byte, group string, exceptBsIds ...string) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	set := mapKit.Get(groupMap, group)
	if set == nil {
		return
	}
	set.Each(func(channel *Channel) bool {
		if sliceKit.Contains(exceptBsIds, channel.GetBsId()) {
			return false
		}
		// TODO: 输出推送失败的error
		_ = channel.PushMessage(messageType, data)
		return false
	})
}

func PushToAll(messageType int, data []byte) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	for _, channel := range allMap {
		// TODO: 输出推送失败的error
		_ = channel.PushMessage(messageType, data)
	}
}
