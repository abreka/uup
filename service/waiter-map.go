package service

import (
	"net"
	"sync"
	"time"
)

type WaiterMap struct {
	waiters map[string]*Waiter
	sync.RWMutex
}

func NewWaiterMap() *WaiterMap {
	return &WaiterMap{
		waiters: make(map[string]*Waiter),
	}
}

// TODO: Golang to upper?
func (wm *WaiterMap) Lookup(addr net.Addr) *Waiter{
	addrStr := addr.String()

	wm.RLock()
	defer wm.RUnlock()
	return wm.waiters[addrStr]
}

func (wm *WaiterMap) getElseCreate(addr string) *Waiter {
	wm.Lock()
	defer wm.Unlock()

	w, found := wm.waiters[addr]
	if found {
		w.IncListeners()
		return w
	}

	w = NewWaiter()
	wm.waiters[addr] = w
	return w
}

// TODO: does a goroutine stop running if the request stops?
func (wm *WaiterMap) WithWaiterResp(addr net.Addr, timeout time.Duration, f func (w *Waiter)) bool {
	addrStr := addr.String()
	w := wm.getElseCreate(addrStr)
	go f(w)
	res := w.Wait(timeout)

	w.DecListeners()
	w.IfDisposable(func() {
		// TODO: THIS IS TRICKY
		wm.Lock()
		defer wm.Unlock()
		delete(wm.waiters, addrStr)
	})

	return res
}

