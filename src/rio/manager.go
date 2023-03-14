package rio

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
	"sync"
)

var (
	// rwLock （整体的）读写锁
	rwLock = new(sync.RWMutex)

	// allMap key: id属性（一对一）
	allMap = make(map[string]*Channel)
	// groupMap key: group属性（一对多）
	groupMap = make(map[string]mapset.Set[*Channel])
	// userMap key: user属性（一对多）
	userMap = make(map[string]mapset.Set[*Channel])
	// bsIdMap key: bsId属性（一对一）
	bsIdMap = make(map[string]*Channel)
)

func AddChannel(channel *Channel) {
	rwLock.Lock()
	defer rwLock.Unlock()

	allMap[channel.GetId()] = channel
}

// RemoveChannel
/*
@return 是否移除成功？（以免多次移除）
*/
func RemoveChannel(channel *Channel) (flag bool) {
	rwLock.Lock()
	defer rwLock.Unlock()

	id := channel.GetId()
	_, flag = mapKit.Remove(allMap, id)

	bsId := channel.GetBsId()
	if strKit.IsNotEmpty(bsId) {
		_, _ = mapKit.Remove(bsIdMap, bsId)
	}

	user := channel.GetUser()
	if strKit.IsNotEmpty(user) {
		set := mapKit.Get(userMap, user)
		set.Remove(channel)
	}

	group := channel.GetGroup()
	if strKit.IsNotEmpty(group) {
		set := mapKit.Get(groupMap, group)
		set.Remove(channel)
	}

	return flag
}
