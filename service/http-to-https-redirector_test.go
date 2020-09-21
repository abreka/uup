package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStartHTTPToHTTPSRedirector(t *testing.T) {
	l, finalizer := InitializeLogger()
	assert.NotNil(t, l)
	defer finalizer()

	ctx, cancel := context.WithCancel(context.Background())

	err := StartHTTPToHTTPSRedirector(ctx, ":0", "433", time.Second)
	assert.EqualError(t, err, "address 433: missing port in address")

	err = StartHTTPToHTTPSRedirector(ctx, ":0", ":433", time.Second)
	assert.NoError(t, err)
	cancel()
}


func Test_MakeHTTPSRedirectHandler(t *testing.T) {
	stdHandler := MakeHTTPSRedirectHandler("443")
	robustHandler := MakeHTTPSRedirectHandler("9999")

	// -----------------------------------------------------------------------
	r := httptest.NewRequest("GET", "http://localhost/path?q=hello", nil)

	statusCode, redirect := runRedirectHandler(stdHandler, r)
	require.Equal(t, http.StatusTemporaryRedirect, statusCode)
	require.Equal(t, "https://localhost/path?q=hello", redirect)

	statusCode, redirect = runRedirectHandler(robustHandler, r)
	require.Equal(t, http.StatusTemporaryRedirect, statusCode)
	require.Equal(t, "https://localhost:9999/path?q=hello", redirect)

	// -----------------------------------------------------------------------

	r = httptest.NewRequest("GET", "http://localhost:80/path?q=hello", nil)

	statusCode, redirect = runRedirectHandler(stdHandler, r)
	require.Equal(t, http.StatusTemporaryRedirect, statusCode)
	require.Equal(t, "https://localhost/path?q=hello", redirect)

	statusCode, redirect = runRedirectHandler(robustHandler, r)
	require.Equal(t, http.StatusTemporaryRedirect, statusCode)
	require.Equal(t, "https://localhost:9999/path?q=hello", redirect)

	// -----------------------------------------------------------------------

	r = httptest.NewRequest("GET", "http://localhost:5555/path?q=hello", nil)

	statusCode, redirect = runRedirectHandler(stdHandler, r)
	require.Equal(t, http.StatusTemporaryRedirect, statusCode)
	require.Equal(t, "https://localhost/path?q=hello", redirect)

	statusCode, redirect = runRedirectHandler(robustHandler, r)
	require.Equal(t, http.StatusTemporaryRedirect, statusCode)
	require.Equal(t, "https://localhost:9999/path?q=hello", redirect)
}

func runRedirectHandler(handler http.HandlerFunc, r *http.Request) (int, string) {
	w := httptest.NewRecorder()
	handler(w, r)
	resp := w.Result()
	u, _ := resp.Location()
	return resp.StatusCode, u.String()
}

