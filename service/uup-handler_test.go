package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestDidError(t *testing.T) {
	assert.False(t, DidError(nil, http.StatusConflict, nil))

	w := httptest.NewRecorder()
	assert.True(t, DidError(w, http.StatusConflict, fmt.Errorf("an error")))
	resp := w.Result()
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	b := make([]byte, 1024)
	n, err := resp.Body.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, `{"err":"an error"}`, string(b[:n]))
}