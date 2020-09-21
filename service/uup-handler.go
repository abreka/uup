package service

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

var InternalErrorJson = []byte(`{"err": "internal error"}`)
var ErrProtocolNotHandled = errors.New("protocol not handled")

func (svc *Service) UUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	req, err := ParseReq(vars["proto"] + "://" + GetRemoteIP(r) + ":" + vars["port"])
	if DidError(w, http.StatusBadRequest, err) {
		return
	}

	resp, err := svc.Execute(req)
	if DidError(w, http.StatusBadRequest, err) {
		return
	}

	b, err := json.Marshal(resp)
	w.Write(b)
}

func (svc *Service) Execute(req *Req) (*Resp, error) {
	switch req.Proto {
	case ProtoTCP, ProtoHTTP, ProtoHTTPS:
		if !svc.tcpDisabled {
			return svc.CheckTCP(req), nil
		}
	case ProtoUDP:
		if svc.udpChecker != nil {
			return svc.CheckUDP(req), nil
		}
	case ProtoQuic:
		if svc.quickChecker != nil {
			return svc.CheckQuic(req), nil
		}
	}
	return nil, ErrProtocolNotHandled
}

func DidError(w http.ResponseWriter, errCode int, err error) bool {
	if err == nil {
		return false
	}

	type ErrMsg struct {
		Error string `json:"err"`
	}

	b, err := json.Marshal(&ErrMsg{Error: err.Error()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(InternalErrorJson)
		return true
	}

	w.WriteHeader(errCode)
	w.Write(b)
	return true
}

