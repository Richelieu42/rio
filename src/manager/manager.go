package manager

import (
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

func Remove(id, reason string) {
	if strKit.IsEmpty(id) {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	delete(all, id)
}
