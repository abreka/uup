package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// It's best to leak as little from golang's errors as possible.
// Thus, these errors are sent to the user for bad inputs.
var (
	ErrEmptyInput    = errors.New("user input empty")
	ErrMissingProto  = errors.New("missing proto in uri")
	ErrMissingHost   = errors.New("missing host in uri")
	ErrMissingPort   = errors.New("missing port in uri")
	ErrInvalidFormat = errors.New("user input invalid")
)

type Req struct {
	Host  string
	Proto int
}

const (
	ProtoUDP = iota
	ProtoTCP
	ProtoHTTP
	ProtoHTTPS
	ProtoQuic
)

var ProtoMap = map[string]int{
	"UDP":   ProtoUDP,
	"TCP":   ProtoTCP,
	"HTTP":  ProtoHTTP,
	"HTTPS": ProtoHTTPS,
	"QUIC":  ProtoQuic,
}

func ParseReq(userInput string) (*Req, error) {
	userInput = strings.TrimSpace(userInput)
	if userInput == "" {
		return nil, ErrEmptyInput
	}

	uri, err := url.Parse(userInput)
	if err != nil {
		if strings.Contains(err.Error(), "missing protocol scheme") {
			return nil, ErrMissingProto
		} else {
			return nil, ErrInvalidFormat
		}
	}

	protoStr := strings.ToUpper(uri.Scheme)

	proto, ok := ProtoMap[protoStr]
	if !ok {
		return nil, fmt.Errorf("protocol %s is not supported", protoStr)
	}

	host := strings.ToUpper(uri.Host)
	if host == "" {
		return nil, ErrMissingHost
	}

	if uri.Port() == "" {
		return nil, ErrMissingPort
	}

	return &Req{
		Host:  host,
		Proto: proto,
	}, nil
}

type Resp struct {
	Success bool `json:"success"`
}
