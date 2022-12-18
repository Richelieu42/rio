package rio

func BingGroup(c *Channel, group string) {
	managerLock.Lock()
	defer managerLock.Unlock()

	c.group = group
}

func BingUser(c *Channel, user string) {
	managerLock.Lock()
	defer managerLock.Unlock()

	c.user = user
}

func BingBsId(c *Channel, bsId string) {
	managerLock.Lock()
	defer managerLock.Unlock()

	c.bsId = bsId
}
