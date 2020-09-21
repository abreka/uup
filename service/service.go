package service

import (
	"context"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"time"
)

const DefaultMaxDialTime = time.Second * 10

type Service struct {
	tcpDisabled           bool
	r                     *mux.Router
	httpServer            *http.Server
	quickChecker          *QuicChecker
	udpChecker            *UDPChecker
	homeTmpl              *template.Template
	MaxDialTime           time.Duration
	activeProtocols       []string
	inspectXForrwardedFor bool
}

type ServiceOpt func(*Service) error

func New(opts ...ServiceOpt) (*Service, error) {
	svc := &Service{
		MaxDialTime: DefaultMaxDialTime,
	}
	err := svc.LoadHomePage()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err := opt(svc)
		if err != nil {
			return nil, err
		}
	}

	if svc.r == nil {
		svc.r = mux.NewRouter()
	}

	svc.activeProtocols, err = svc.computeActiveProtocols()
	if err != nil {
		return nil, err
	}

	svc.buildRoutes()

	return svc, nil
}

func (svc *Service) ListenAndServe(ctx context.Context) error {
	if svc.udpChecker != nil {
		go svc.udpChecker.StartHandler()
	}

	return nil
}

func (svc *Service) Cleanup() {
	if svc.udpChecker != nil {
		svc.udpChecker.Stop()
	}

	if svc.quickChecker != nil {
		svc.quickChecker.conn.Close()
	}
}

func (svc *Service) Routes() *mux.Router {
	return svc.r
}
