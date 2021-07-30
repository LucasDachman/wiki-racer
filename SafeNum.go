package main

import "sync"

type SafeNum struct {
	mtx sync.RWMutex
	num int
}

func (sn *SafeNum) Write(val int) {
	sn.mtx.Lock()
	defer sn.mtx.Unlock()
	sn.num = val
}

func (sn *SafeNum) Read() int {
	sn.mtx.RLock()
	defer sn.mtx.RUnlock()
	return sn.num
}
