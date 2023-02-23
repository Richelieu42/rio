package manager

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/sliceKit"
)

func PushToBsId(messageType int, bsid string, data []byte) {
	if len(data) == 0 {
		// 无效的推送数据
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

func PushToGroup(messageType int, group string, data []byte, exceptBsIds []string) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	s := mapKit.Get(groupMap, group)

	for _, channel := range s {
		if sliceKit.Contains(exceptBsIds, channel.GetBsId()) {
			continue
		}
		// TODO: 输出推送失败的error
		_ = channel.PushMessage(messageType, data)
	}
}
