package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWaiter_WaitTimeout(t *testing.T) {
	w := NewWaiter()

	assert.False(t, w.Wait(time.Millisecond))
	assert.True(t, w.complete)
	assert.False(t, w.success)
}

func TestWaiter_WaitAlreadyComplete(t *testing.T) {
	w := NewWaiter()
	w.complete = true
	w.success = true

	assert.True(t, w.Wait(time.Minute))
	assert.True(t, w.complete)
	assert.True(t, w.success)
}

func TestWaiter_CompleteAlreadyCompleted(t *testing.T) {
	w := NewWaiter()

	w.Complete(true)
	assert.True(t, w.complete)
	assert.True(t, w.success)

	w.Complete(false)
	assert.True(t, w.complete)
	assert.True(t, w.success)
}

func TestWaiter_MultipleWaiters(t *testing.T) {
	w := NewWaiter()
	assert.Equal(t, w.listeners, 1)

	w.IncListeners()
	assert.False(t, w.complete)
	assert.False(t, w.success)
	assert.Equal(t, w.listeners, 2)

	w.DecListeners()
	assert.False(t, w.complete)
	assert.False(t, w.success)
	assert.Equal(t, w.listeners, 1)

	var isDisposable bool
	w.IfDisposable(func() {
		isDisposable = true
	})
	assert.False(t, isDisposable)

	w.DecListeners()
	w.IfDisposable(func() {
		isDisposable = true
	})
	assert.True(t, isDisposable)
}
