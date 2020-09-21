package service

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRemoteIP_Bare(t *testing.T) {
	var remoteIP string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIP = GetRemoteIP(r)
	}))

	resp, err := http.Get(ts.URL)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "127.0.0.1", remoteIP)
}

