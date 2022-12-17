package rio

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
	"sync"
)

var (
	// 读写锁
	lock = new(sync.RWMutex)

	all = make(map[string]*Channel)
)

func Add(c *Channel) {
	lock.Lock()
	defer lock.Unlock()

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

	lock.Lock()
	defer lock.Unlock()

	c, ok := mapKit.Remove(all, id)
	if ok {
		c.closed = true
	}
	return ok
}
