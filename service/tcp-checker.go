package service

import "net"

func (svc *Service) CheckTCP(req *Req) *Resp {
	resp := &Resp{}

	conn, err := net.DialTimeout("tcp", req.Host, svc.MaxDialTime)
	if err != nil {
		return resp
	}

	resp.Success = true
	conn.Close()
	return resp
}