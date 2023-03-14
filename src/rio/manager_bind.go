package rio

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
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

		set := mapKit.Get(userMap, user)
		set.Add(channel)
	}

	if strKit.IsNotEmpty(group) {
		channel.SetGroup(group)

		set := mapKit.Get(groupMap, group)
		set.Add(channel)
	}
}
