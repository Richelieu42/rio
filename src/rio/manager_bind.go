package rio

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/sliceKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
)

func BindData(channel *Channel, bsId, user, group string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	if strKit.IsNotEmpty(bsId) {
		channel.SetBsId(bsId)

		mapKit.Set(bsIdMap, bsId, channel)
	}

	if strKit.IsNotEmpty(user) {
		channel.SetUser(user)

		channels := mapKit.Get(userMap, user)
		channels = sliceKit.Append(channels, channel)
		mapKit.Set(userMap, user, channels)
	}

	if strKit.IsNotEmpty(group) {
		channel.SetGroup(group)

		channels := mapKit.Get(groupMap, group)
		channels = sliceKit.Append(channels, channel)
		mapKit.Set(groupMap, group, channels)
	}
}
