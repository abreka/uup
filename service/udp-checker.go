package service

import (
	"log"
	"net"
	"time"
)

func (svc *Service) CheckUDP(req *Req) *Resp {
	addr, err := net.ResolveUDPAddr("udp", req.Host)
	if err != nil {
		return &Resp{}
	}

	return svc.udpChecker.Check(addr, svc.MaxDialTime)
}

type UDPChecker struct {
	conn                net.PacketConn
	waiterMap           *WaiterMap
	sendPacket          []byte
	resendAttempts      int
	sleepBetweenResends time.Duration
	err                 error
}

func NewUDPChecker(bindAddr string, resendAttempts int, sleepBetweenResend time.Duration) (*UDPChecker, error) {
	conn, err := net.ListenPacket("udp", bindAddr)
	if err != nil {
		return nil, err
	}

	return &UDPChecker{
		conn:       conn,
		sendPacket: []byte("U UP?"),
		waiterMap:  NewWaiterMap(),
		resendAttempts: resendAttempts,
		sleepBetweenResends: sleepBetweenResend,
	}, nil
}

// TODO: k-attempts
func (u *UDPChecker) Check(addr net.Addr, timeout time.Duration) *Resp {
	resp := &Resp{}

	success := u.waiterMap.WithWaiterResp(addr, timeout, func(w *Waiter) {
		_, err := u.conn.WriteTo(u.sendPacket, addr)
		if err != nil {
			w.Complete(false)
		}
	})

	resp.Success = success
	return resp
}

func (u *UDPChecker) StartHandler() {
	b := make([]byte, 1024)
	for {
		n, addr, err := u.conn.ReadFrom(b)
		if err != nil {
			u.err = err
			return
		}

		log.Println(addr, b[:n])
		w := u.waiterMap.Lookup(addr)
		if w != nil {
			w.Complete(true)
		}
	}
}

func (u *UDPChecker) Stop() error {
	return u.conn.Close()
}

