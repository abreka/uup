package service

import (
	"context"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net"
	"sync"
	"testing"
	"time"
)

type QuicCheckTestSuite struct {
	suite.Suite
	svc *Service

	serverURL string
	serverDidAccept bool
	serverError     error
	serverCtx context.Context
	serverCancel context.CancelFunc

	serverDone sync.WaitGroup
}

func TestQuicCheckTestSuite(t *testing.T) {
	suite.Run(t, new(QuicCheckTestSuite))
}

func (suite *QuicCheckTestSuite) SetupTest() {
	var err error
	suite.svc, err = New(WithQuicService("0.0.0.0:0", time.Millisecond * 500))
	suite.NoError(err)
	suite.NotNil(suite.svc)

	suite.serverURL = ""
	suite.serverDidAccept = false
	suite.serverError = nil
	suite.serverDone = sync.WaitGroup{}
	suite.serverCtx, suite.serverCancel = context.WithCancel(context.Background())
}

func (suite *QuicCheckTestSuite) TearDownTest() {
	if suite.serverURL != "" {
		suite.serverDone.Wait()
	}
	suite.svc.Cleanup()
	suite.svc = nil
}

func (suite *QuicCheckTestSuite) LaunchServer(tlsCfg *tls.Config, quicCfg *quic.Config) {
	var listening sync.WaitGroup
	suite.serverDone.Add(1)
	listening.Add(1)

	go func() {
		defer suite.serverDone.Done()
		conn, err := net.ListenPacket("udp", "0.0.0.0:0")
		suite.NoError(err)
		defer conn.Close()

		_, port, err := net.SplitHostPort(conn.LocalAddr().String())
		suite.serverURL = "quic://0.0.0.0:" + port

		l, err := quic.Listen(conn, tlsCfg, quicCfg)
		suite.NoError(err)

		listening.Done()

		// Accept once
		sess, err := l.Accept(suite.serverCtx)
		if err != nil {
			suite.serverError = err
			return
		}

		suite.serverDidAccept = true
		_, err = sess.AcceptStream(context.Background())
		suite.EqualError(err, "Application error 0xdeadbeaf: U UP? Haz Goodbye!")
	}()

	listening.Wait()
}

func (suite *QuicCheckTestSuite) Test_SameConfiguration() {
	suite.LaunchServer(suite.svc.quickChecker.tlsConfig, suite.svc.quickChecker.cfg)

	req, err := ParseReq(suite.serverURL)
	suite.NoError(err)
	got := suite.svc.CheckQuic(req)
	suite.True(got.Success)

}

func (suite *QuicCheckTestSuite) Test_DifferentNextProto() {
	tlsCfg := generateTLSConfig(1024)
	tlsCfg.NextProtos = []string{"different-proto"}
	suite.LaunchServer(tlsCfg, suite.svc.quickChecker.cfg)

	req, err := ParseReq(suite.serverURL)
	suite.NoError(err)
	got := suite.svc.CheckQuic(req)
	suite.True(got.Success)

	suite.serverCancel() // Listen never returns
}

func (suite *QuicCheckTestSuite) Test_DifferentRSABits() {
	tlsCfg := generateTLSConfig(2048)
	suite.LaunchServer(tlsCfg, suite.svc.quickChecker.cfg)

	req, err := ParseReq(suite.serverURL)
	suite.NoError(err)
	got := suite.svc.CheckQuic(req)
	suite.True(got.Success)

	suite.serverCancel() // Listen never returns
}


func TestService_NoQuicServiceRunning(t *testing.T) {
	svc, err := New(WithQuicService("0.0.0.0:0", time.Millisecond * 500))
	require.NoError(t, err)
	assert.NotNil(t, svc)

	// UDP but discarded
	req, err := ParseReq("quic://0.0.0.0:9")
	require.NoError(t, err)
	got := svc.CheckQuic(req)
	require.False(t, got.Success)
}
