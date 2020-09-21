package uup

import (
	"context"
	"fmt"
	"github.com/abreka/uup/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

const Version = "0.0.1"

var (
	disableTCP          bool
	httpAddr            string
	httpsRedirectorAddr string
	dialTimeout         time.Duration
	udpAddr             string
	udpResendAttempts   int
	udpResendInterval   time.Duration
	quicAddr            string
	tlsCertFilePath     string
	tlsKeyFilePath      string
)

// TODO: http redirect
func RunCLI() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.AddCommand(StartCmd)
	rootCmd.PersistentFlags().StringVar(&tlsCertFilePath, "tls-cert-path", "",
		"path to your TLS certificate",
	)
	rootCmd.PersistentFlags().StringVar(&tlsKeyFilePath, "tls-key-path", "",
		"path to your TLS private key",
	)
	rootCmd.PersistentFlags().BoolVar(&disableTCP, "disable-tcp", false,
		"disable TCP checking",
	)
	rootCmd.PersistentFlags().StringVar(&httpAddr, "http-addr", "0.0.0.0:8080",
		"ip:port to bind http service to",
	)
	rootCmd.PersistentFlags().StringVar(&udpAddr, "udp-addr", "",
		"ip:port to bind udp service to",
	)
	rootCmd.PersistentFlags().DurationVar(&udpResendInterval, "udp-resend-interval", time.Second,
		"time to wait before sending the next udp datagram",
	)
	rootCmd.PersistentFlags().IntVar(&udpResendAttempts, "udp-resend-attempts", 10,
		"number of udp datagrams to send before giving up",
	)
	rootCmd.PersistentFlags().StringVar(&quicAddr, "quic-addr", "",
		"ip:port to bind quic service to",
	)
	rootCmd.PersistentFlags().DurationVar(&dialTimeout, "dial-timeout", time.Second,
		"number of udp datagrams to send before giving up",
	)
	rootCmd.PersistentFlags().StringVar(&httpsRedirectorAddr, "redirect-http-to-https-addr", "",
		"launches a process bound to the address that redirects to 443",
	)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Version: Version,
	Use:     "uup",
	Short:   "a service to check whether your port is internet-accessible",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the uup service",
	Run: func(cmd *cobra.Command, args []string) {
		logger, finalizer := service.InitializeLogger()
		defer finalizer()

		opts, useTLS, clientTimeout := buildOptsFromCLIParams()

		s, err := service.New(opts...)
		if err != nil {
			log.Fatalln(err)
		}
		s.ListenAndServe(context.Background())


		// Disable cookies
		srv := &http.Server{
			Handler:      s.Routes(),
			Addr:         httpAddr,
			WriteTimeout: clientTimeout,
			ReadTimeout:  clientTimeout,
		}

		if useTLS {
			if httpsRedirectorAddr != "" {
				// Technically a bug...
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				err := service.StartHTTPToHTTPSRedirector(ctx, httpsRedirectorAddr, httpAddr, time.Second)
				if err != nil {
					os.Exit(1)
				}
			}
			logger.Info("https server listening", zap.String("addr", httpAddr))
			log.Fatal(srv.ListenAndServeTLS(tlsCertFilePath, tlsKeyFilePath))
		} else {
			logger.Info("http server listening", zap.String("addr", httpAddr))
			log.Fatal(srv.ListenAndServe())
		}
	},
}

func buildOptsFromCLIParams() ([]service.ServiceOpt, bool, time.Duration) {
	var opts []service.ServiceOpt
	var clientTimeout time.Duration

	if udpAddr != "" {
		opts = append(opts, service.WithUDPChecker(udpAddr, udpResendAttempts, udpResendInterval))
		// The udp checker waits for one additional resend interval after the final one.
		clientTimeout = udpResendInterval * time.Duration(udpResendAttempts+1)
	}

	if quicAddr != "" {
		opts = append(opts, service.WithQuicService(quicAddr, dialTimeout))
	}

	useTLS := false
	if tlsCertFilePath != "" || tlsKeyFilePath != "" {
		useTLS = true
		if tlsCertFilePath == "" {
			fmt.Fprintln(os.Stderr, "--tls-cert-path not specified")
			os.Exit(1)
		}
		if tlsKeyFilePath == "" {
			fmt.Fprintln(os.Stderr, "--tls-key-path not specified")
			os.Exit(1)
		}
	}

	if dialTimeout > clientTimeout {
		clientTimeout = dialTimeout
	}
	clientTimeout += 2 * time.Second // Give clients a little wiggle room.

	return opts, useTLS, clientTimeout
}

var checkCmd = &cobra.Command{
	Use:     "uup",
	Short:   "call the uup service",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
