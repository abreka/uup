package service

import (
	"sync"
	"time"
)

type Waiter struct {
	listeners int
	cond    *sync.Cond
	success bool
	complete bool
	sync.Mutex
}

// Create a new waiter for UPD checking.
func NewWaiter() *Waiter {
	w := &Waiter{listeners: 1}
	w.cond = sync.NewCond(w)
	return w
}

func (w *Waiter) IncListeners() {
	w.Lock()
	defer w.Unlock()

	w.listeners += 1
}

func (w *Waiter) DecListeners() {
	w.Lock()
	defer w.Unlock()

	w.listeners -= 1
}

func (w *Waiter) IfDisposable(f func()) {
	w.Lock()
	defer w.Unlock()
	if w.listeners == 0 {
		f()
	}
}

// Complete the waiter with the success condition.
//
// Multiple calls have no effect.
func (w *Waiter) Complete(success bool) {
	w.Lock()

	// Only complete the waiter once.
	if !w.complete {
		w.complete = true
		w.success = success
		w.cond.Broadcast()
	}

	w.Unlock()
}

// Wait until another routine completes or the waiter times out.
//
// Return the success condition (false if timeout or the check failed).
func (w *Waiter) Wait(timeout time.Duration) bool {
	go func() {
		time.Sleep(timeout)
		w.Complete(false)
	}()

	w.Lock()
	for !w.complete {
		w.cond.Wait()
	}
	w.Unlock()

	return w.success
}
