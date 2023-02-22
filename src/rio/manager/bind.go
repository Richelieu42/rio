package manager

import "github.com/richelieu42/rio/src/rio/bean"

func BindGroup(channel *bean.Channel, group string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	channel.SetGroup(group)
}

func BindUser(channel *bean.Channel, user string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	channel.SetUser(user)
}

func BindBsId(channel *bean.Channel, bsId string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	channel.SetBsId(bsId)
}
