package rio

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
	"sync"
)

var (
	// 读写锁
	managerLock = new(sync.RWMutex)

	all = make(map[string]*Channel)
)

func Add(c *Channel) {
	managerLock.Lock()
	defer managerLock.Unlock()

	all[c.id] = c
}

// Remove
/*
@return
*/
func Remove(id string) bool {
	if strKit.IsEmpty(id) {
		return false
	}

	managerLock.Lock()
	defer managerLock.Unlock()

	_, ok := mapKit.Remove(all, id)
	return ok
}
