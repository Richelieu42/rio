package manager

import (
	"github.com/richelieu42/go-scales/src/core/mapKit"
	"github.com/richelieu42/go-scales/src/core/strKit"
	"github.com/sirupsen/logrus"
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
	logrus.Infof("Channel(id: %s) is added.", c.id)
}

// RemoveByFrontEnd 前端主动断开连接（关闭浏览器tab、干掉浏览器进程...），会调用此方法.
func RemoveByFrontEnd(id string, code int, text string) {
	reason := strKit.Format("CloseHandler with code(%d) and text(%s)", code, text)
	Remove(id, reason)
}

func Remove(id, reason string) {
	if strKit.IsEmpty(id) {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	if _, ok := mapKit.Remove(all, id); ok {
		logrus.Infof("Channel(id: %s) is removed.", id)
	}
}
