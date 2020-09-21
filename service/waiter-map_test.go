package service

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestWaiterMap(t *testing.T) {
	wm := NewWaiterMap()

	addr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 80,
	}

	// It doesn't yet exist.
	w := wm.Lookup(addr)
	assert.Nil(t, w)

	// It does now.
	w1 := wm.getElseCreate(addr.String())
	assert.NotNil(t, w1)
	w = wm.Lookup(addr)
	assert.Equal(t, w, w1)
	assert.Equal(t, 1, w.listeners)

	// It gets the existing entry
	w2 := wm.getElseCreate(addr.String())
	assert.NotNil(t, w2)
	w = wm.Lookup(addr)
	assert.Equal(t, w, w2)
	assert.Equal(t, w1, w2)
	assert.Equal(t, 2, w.listeners)

	// Gets existing entry; but doesn't dispose of it because
	// there are two listeners.
	var called bool
	wm.WithWaiterResp(addr, time.Millisecond, func(w3 *Waiter) {
		assert.Equal(t, w2, w3)
		assert.Equal(t, 3, w3.listeners)
		called = true
	})
	assert.True(t, called)
	assert.NotNil(t, wm.Lookup(addr))
	assert.Equal(t, 2, w.listeners)

	// It automatically dispose otherwise.
	addr2 := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 2),
		Port: 80,
	}
	called = false
	wm.WithWaiterResp(addr2, time.Millisecond, func(w4 *Waiter) {
		assert.Equal(t, 1, w4.listeners)
		called = true
	})
	assert.True(t, called)
	assert.Nil(t, wm.Lookup(addr2))
}


