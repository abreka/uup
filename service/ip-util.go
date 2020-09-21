package service

import (
	"net"
	"net/http"
)

// TODO: This s not robust.
//
// Need X-Forwarded-For support, which has security concerns.
func GetRemoteIP(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	addr, _, _ := net.SplitHostPort(remoteAddr)
	return addr
}


