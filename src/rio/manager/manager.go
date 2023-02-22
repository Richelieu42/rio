package manager

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/rio/src/rio/bean"
	"sync"
)

var (
	// rwLock 读写锁
	rwLock = new(sync.RWMutex)

	// allMap key: id属性
	allMap = make(map[string]*bean.Channel)
	// groupMap key: group属性
	groupMap = make(map[string]*bean.Channel)
	// userMap key: user属性
	userMap = make(map[string]*bean.Channel)
	// bdIdMap key: bsId属性
	bdIdMap = make(map[string]*bean.Channel)
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
func RemoveChannel(channel *bean.Channel) bool {
	rwLock.Lock()
	defer rwLock.Unlock()

	_, ok := mapKit.Remove(allMap, channel.GetId())
	return ok
}
