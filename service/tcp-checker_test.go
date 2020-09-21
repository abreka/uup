package service

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestService_CheckTCP(t *testing.T) {
	svc := Service{MaxDialTime: time.Second}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	}))

	req, err := ParseReq(strings.Replace(ts.URL, "http://", "tcp://", 1))
	require.NoError(t, err)
	resp, err := svc.Execute(req)
	require.NoError(t, err)
	require.Equal(t, &Resp{Success: true}, resp)

	ts.Close()
	resp, err = svc.Execute(req)
	require.NoError(t, err)
	require.Equal(t, &Resp{Success: false}, resp)
}

func TestService_CheckHTTP(t *testing.T) {
	svc := Service{
		MaxDialTime: time.Second,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	}))
	defer ts.Close()

	req, err := ParseReq(ts.URL)
	require.NoError(t, err)
	resp, err := svc.Execute(req)
	require.NoError(t, err)
	require.Equal(t, &Resp{Success: true}, resp)
}
