package manager

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/sliceKit"
	"github.com/richelieu42/rio/src/rio/bean"
)

func BindBsId(channel *bean.Channel, bsId string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	channel.SetBsId(bsId)

	mapKit.Set(bsIdMap, bsId, channel)
}

func BindUser(channel *bean.Channel, user string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	channel.SetUser(user)

	s := mapKit.Get(userMap, user)
	s = sliceKit.Append(s, channel)
	mapKit.Set(userMap, user, s)
}

func BindGroup(channel *bean.Channel, group string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	channel.SetGroup(group)

	s := mapKit.Get(groupMap, group)
	s = sliceKit.Append(s, channel)
	mapKit.Set(groupMap, group, s)
}
