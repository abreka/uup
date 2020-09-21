package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

// TODO: expand tsetes and implement repeat UDP
func TestService_CheckUDP(t *testing.T) {
	svc, err := New(WithUDPChecker("127.0.0.1:0", 1, time.Second))
	svc.MaxDialTime = time.Millisecond * 100
	require.NoError(t, err)
	require.NotNil(t, svc)

	var handlerStarted, handlerDone sync.WaitGroup
	handlerStarted.Add(1)
	handlerDone.Add(1)
	go func() {
		defer handlerDone.Done()
		handlerStarted.Done()
		svc.udpChecker.StartHandler()
	}()

	handlerStarted.Wait()

	req, err := ParseReq("udp://" + svc.udpChecker.conn.LocalAddr().String())
	require.NoError(t, err)
	resp := svc.CheckUDP(req)
	assert.True(t, resp.Success)

	// Port should fail in general.
	req2, err := ParseReq("udp://localhost:9")  // Discard protocol, if any
	require.NoError(t, err)
	resp = svc.CheckUDP(req2)
	assert.False(t, resp.Success)

	// Can't resolve error
	req3, err := ParseReq("udp://thishostdoesnotexist:9")  // Discard protocol, if any
	require.NoError(t, err)
	resp = svc.CheckUDP(req3)
	assert.False(t, resp.Success)

	// Stops the listener
	err = svc.udpChecker.Stop()
	require.NoError(t, err)

	// Won't send
	resp = svc.CheckUDP(req)
	assert.False(t, resp.Success)
	handlerDone.Wait()
}
