package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/lucas-clemente/quic-go"
	"math/big"
	"net"
	"strings"
	"time"
)

func (svc *Service) CheckQuic(req *Req) *Resp {
	addr, err := net.ResolveUDPAddr("udp", req.Host)
	if err != nil {
		return &Resp{}
	}


	return svc.quickChecker.Check(addr)
}

type QuicChecker struct {
	conn      net.PacketConn
	err       error
	cfg       *quic.Config
	tlsConfig *tls.Config
}

// Create a new QuicChecker service.
//
// This does not do any certificate validation.
func NewQuicChecker(bindAddr string, maxHandhake time.Duration) (*QuicChecker, error) {
	conn, err := net.ListenPacket("udp", bindAddr)
	if err != nil {
		return nil, err
	}

	return &QuicChecker{
		conn:      conn,
		tlsConfig: generateTLSConfig(1024),
		cfg: &quic.Config{
			HandshakeTimeout:   maxHandhake, // TODO
			KeepAlive:          false,
			MaxIncomingStreams: 1,
		},
	}, nil
}

func (q *QuicChecker) Check(addr *net.UDPAddr) *Resp {
	resp := &Resp{}
	sess, err := quic.Dial(q.conn, addr, addr.IP.String(), q.tlsConfig, q.cfg)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "CRYPTO_ERROR: tls: no application protocol") {
			resp.Success = true
		}
	} else {
		resp.Success = true
		sess.CloseWithError(0xdeadbeaf, "U UP? Haz Goodbye!")
	}

	return resp
}

// TODO: Switch to EDSCA
func generateTLSConfig(bits int) *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err) // TODO
	}
	template := x509.Certificate{
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().AddDate(10, 0, 0), // TODO
		SerialNumber: big.NewInt(1),
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"u-up"},
		InsecureSkipVerify: true,
	}
}
