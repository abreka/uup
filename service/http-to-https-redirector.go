package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

// Starts a HTTP to HTTPS redirector service on bindAddr.
//
// - The server times out reads and writes longer than timeout.
// - The targetAddr should refer to the https address.
// - Cancelling the context shuts down the redirector.
func StartHTTPToHTTPSRedirector(ctx context.Context, httpBindAddr, targetAddr string, timeout time.Duration) error {
	svcName := zap.String("svc", "http-to-https")

	_, port, err := net.SplitHostPort(targetAddr)
	if err != nil {
		logger.Info(err.Error(), zap.String("binding", httpBindAddr), svcName)
		return err
	}

	srv := &http.Server{
		Addr:         httpBindAddr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Handler:      http.HandlerFunc(MakeHTTPSRedirectHandler(port)),
	}

	logger.Info("starting listener", zap.String("binding", httpBindAddr), svcName)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Error(fmt.Sprintf("ListenAndServe exited: %v", err), svcName)
		}
	}()

	go func() {
		<-ctx.Done()
		logger.Info("shutting down", zap.String("binding", httpBindAddr), svcName)

		// TODO: Read code to see how it uses the context argument.
		err := srv.Shutdown(context.Background())
		logger.Error(fmt.Sprintf("shutdown: %v", err), svcName)
	}()

	return nil
}

func MakeHTTPSRedirectHandler(dstPort string) func(w http.ResponseWriter, r *http.Request) {
	urlRewriter := stdHTTPSURLRewriter

	if dstPort != "443" {
		urlRewriter = makeRobustHTTPSURLRewriter(":" + dstPort)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, urlRewriter(r), http.StatusTemporaryRedirect)
	}
}

func makeRobustHTTPSURLRewriter(destPortWithColonPrefix string) func(r *http.Request) string {
	return func(r *http.Request) string {
		return appendOptionalQuery("https://" + extractHost(r) + destPortWithColonPrefix + r.URL.Path, r)
	}
}

func stdHTTPSURLRewriter(r *http.Request) string {
	return appendOptionalQuery("https://"+extractHost(r)+r.URL.Path, r)
}

func appendOptionalQuery(target string, r *http.Request) string {
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	return target
}

func extractHost(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		return r.Host
	}
	return host
}
