package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

var logger *zap.Logger

func WithRouter(r *mux.Router) ServiceOpt {
	return func(svc *Service) error {
		svc.r = r
		return nil
	}
}

func WithUDPChecker(bindAddr string, resendAttempts int, sleepBetweenResend time.Duration) ServiceOpt {
	return func(svc *Service) error {
		u, err := NewUDPChecker(bindAddr, resendAttempts, sleepBetweenResend)
		if err != nil {
			return err
		}
		svc.udpChecker = u
		return nil
	}
}

func WithQuicService(bindAddr string, maxDialTime time.Duration) ServiceOpt {
	return func(svc *Service) error {
		quicService, err := NewQuicChecker(bindAddr, maxDialTime)
		if err != nil {
			return err
		}
		svc.quickChecker = quicService
		return nil
	}
}

func DisableTCP() ServiceOpt {
	return func(svc *Service) error {
		svc.tcpDisabled = true
		return nil
	}
}

func (svc *Service) buildRoutes() error {
	svc.r.HandleFunc("/", svc.Home)
	svc.r.HandleFunc("/uup/{proto}/{port:[0-9]+}", svc.UUpHandler).Methods(http.MethodPost)
	return nil
}

func (svc *Service) computeActiveProtocols() ([]string, error) {
	var activeProtocols []string

	if !svc.tcpDisabled {
		activeProtocols = append(activeProtocols, "HTTP")
		activeProtocols = append(activeProtocols, "HTTPS")
	}
	if svc.quickChecker != nil {
		activeProtocols = append(activeProtocols, "QUIC")
	}
	if svc.udpChecker != nil {
		activeProtocols = append(activeProtocols, "UDP")
	}
	if !svc.tcpDisabled {
		activeProtocols = append(activeProtocols, "TCP")
	}

	if len(activeProtocols) == 0 {
		return nil, fmt.Errorf("configuration implies no active protocols")
	}

	return activeProtocols, nil
}

func InitializeLogger() (*zap.Logger, func() error) {
	var err error
	logger, err = zap.NewProduction() // logger is a global variable.

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to initialize logger: %v\n", err)
		os.Exit(1)
	}

	return logger, logger.Sync
}
