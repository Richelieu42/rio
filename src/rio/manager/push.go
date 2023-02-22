package manager

func PushToBsId(bsid, text string) {
	rwLock.RLock()
	defer rwLock.RUnlock()

}

func PushToGroup(group, text string, exceptBsIds []string) {
	rwLock.RLock()
	defer rwLock.RUnlock()

}
