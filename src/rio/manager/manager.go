package manager

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/sliceKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
	"github.com/richelieu42/rio/src/rio/bean"
	"sync"
)

var (
	// rwLock 读写锁
	rwLock = new(sync.RWMutex)

	// allMap key: id属性（一对一）
	allMap = make(map[string]*bean.Channel)
	// groupMap key: group属性（一对多）
	groupMap = make(map[string][]*bean.Channel)
	// userMap key: user属性（一对多）
	userMap = make(map[string][]*bean.Channel)
	// bsIdMap key: bsId属性（一对一）
	bsIdMap = make(map[string]*bean.Channel)
)

func AddChannel(channel *bean.Channel) {
	rwLock.Lock()
	defer rwLock.Unlock()

	allMap[channel.GetId()] = channel
}

// RemoveChannel
/*
@return
*/
func RemoveChannel(channel *bean.Channel) {
	rwLock.Lock()
	defer rwLock.Unlock()

	id := channel.GetId()
	_, _ = mapKit.Remove(allMap, id)

	bsId := channel.GetBsId()
	if strKit.IsNotEmpty(bsId) {
		_, _ = mapKit.Remove(bsIdMap, bsId)
	}

	user := channel.GetUser()
	if strKit.IsNotEmpty(user) {
		s := mapKit.Get(userMap, user)
		s, _ = sliceKit.Remove(s, channel)
		mapKit.Set(userMap, user, s)
	}

	group := channel.GetGroup()
	if strKit.IsNotEmpty(group) {
		s := mapKit.Get(groupMap, group)
		s, _ = sliceKit.Remove(s, channel)
		mapKit.Set(groupMap, group, s)
	}
}
