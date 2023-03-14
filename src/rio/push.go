package rio

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/sliceKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
)

func PushToBsid(messageType int, data []byte, bsid string) {
	if strKit.IsEmpty(bsid) || len(data) == 0 {
		return
	}

	rwLock.RLock()
	defer rwLock.RUnlock()

	channel := mapKit.Get(bsIdMap, bsid)
	if channel == nil {
		return
	}
	// TODO: 输出推送失败的error
	_ = channel.PushMessage(messageType, data)
}

func PushToUser(messageType int, data []byte, user string, exceptBsids ...string) {
	if strKit.IsEmpty(user) || len(data) == 0 {
		return
	}

	rwLock.RLock()
	defer rwLock.RUnlock()

	set := mapKit.Get(userMap, user)
	if set == nil {
		return
	}
	set.Each(func(channel *Channel) bool {
		if sliceKit.Contains(exceptBsids, channel.GetBsId()) {
			return false
		}
		// TODO: 输出推送失败的error
		_ = channel.PushMessage(messageType, data)
		return false
	})
}

func PushToGroup(messageType int, data []byte, group string, exceptBsids ...string) {
	if strKit.IsEmpty(group) || len(data) == 0 {
		return
	}

	rwLock.RLock()
	defer rwLock.RUnlock()

	set := mapKit.Get(groupMap, group)
	if set == nil {
		return
	}
	set.Each(func(channel *Channel) bool {
		if sliceKit.Contains(exceptBsids, channel.GetBsId()) {
			return false
		}
		// TODO: 输出推送失败的error
		_ = channel.PushMessage(messageType, data)
		return false
	})
}

func PushToAll(messageType int, data []byte, exceptBsids ...string) {
	if len(data) == 0 {
		return
	}

	rwLock.RLock()
	defer rwLock.RUnlock()

	for _, channel := range allMap {
		if sliceKit.Contains(exceptBsids, channel.GetBsId()) {
			continue
		}

		// TODO: 输出推送失败的error
		_ = channel.PushMessage(messageType, data)
	}
}
