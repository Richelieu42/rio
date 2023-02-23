package rio

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/sliceKit"
)

func BindBsId(channel *Channel, bsId string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	channel.SetBsId(bsId)

	mapKit.Set(bsIdMap, bsId, channel)
}

func BindUser(channel *Channel, user string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	channel.SetUser(user)

	s := mapKit.Get(userMap, user)
	s = sliceKit.Append(s, channel)
	mapKit.Set(userMap, user, s)
}

func BindGroup(channel *Channel, group string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	channel.SetGroup(group)

	s := mapKit.Get(groupMap, group)
	s = sliceKit.Append(s, channel)
	mapKit.Set(groupMap, group, s)
}
