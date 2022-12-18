package rio

func BindGroup(c *Channel, group string) {
	managerLock.Lock()
	defer managerLock.Unlock()

	c.group = group
}

func BindUser(c *Channel, user string) {
	managerLock.Lock()
	defer managerLock.Unlock()

	c.user = user
}

func BindBsId(c *Channel, bsId string) {
	managerLock.Lock()
	defer managerLock.Unlock()

	c.bsId = bsId
}
