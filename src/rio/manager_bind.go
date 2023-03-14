package rio

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/setKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
)

func bindData(channel *Channel, bsId, user, group string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	if strKit.IsNotEmpty(bsId) {
		channel.SetBsId(bsId)

		mapKit.Set(bsIdMap, bsId, channel)
	}

	if strKit.IsNotEmpty(user) {
		channel.SetUser(user)

		set := mapKit.Get(userMap, user)
		if set == nil {
			set = setKit.NewSet[*Channel](false)
			mapKit.Set(userMap, user, set)
		}
		set.Add(channel)
	}

	if strKit.IsNotEmpty(group) {
		channel.SetGroup(group)

		set := mapKit.Get(groupMap, group)
		if set == nil {
			set = setKit.NewSet[*Channel](false)
			mapKit.Set(groupMap, group, set)
		}
		set.Add(channel)
	}
}
