package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Logger interface {
	Error(...interface{})
	Info(...interface{})
	Debug(...interface{})
}

// HTTP represents the http server, all config concerning http is done in here
type HTTP struct {
	ln     net.Listener
	server *http.Server
	addr   string
	Domain string
	logger Logger
}

func NewHTTP(opts ...Option) HTTP {
	// init defaults
	s := HTTP{
		server: &http.Server{
			Handler:      http.HandlerFunc(helloWorld),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
		addr:   "localhost:80",
		logger: deadendLogger{},
	}

	for _, opt := range opts {
		opt(&s)
	}
	return s
}

// UseTLS returns true if the cert & key file are specified.
func (h *HTTP) UseTLS() bool {
	return h.Domain != ""
}

// Scheme returns the URL scheme for the server.
func (h *HTTP) Scheme() string {
	if h.UseTLS() {
		return "https"
	}
	return "http"
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (h *HTTP) Port() int {
	if h.ln == nil {
		return 0
	}
	return h.ln.Addr().(*net.TCPAddr).Port
}

// URL returns the local base URL of the running server.
func (h *HTTP) URL() string {
	scheme, port := h.Scheme(), h.Port()

	// Use localhost unless a domain is specified.
	domain := "localhost"
	if h.Domain != "" {
		domain = h.Domain
	}

	// Return without port if using standard porth.
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", h.Scheme(), domain)
	}
	return fmt.Sprintf("%s://%s:%d", h.Scheme(), domain, h.Port())
}

// Run runs a http server listening on the port configured.
// this ends when ListenAndServe receives an error, or the context is cancelled
func (h *HTTP) Run(ctx context.Context) (err error) {
	if h.ln, err = net.Listen("tcp", h.addr); err != nil {
		return err
	}
	errCh := make(chan error, 1)
	go func() {
		h.logger.Debug("server starting up...")
		// hang until the server is closed
		err = h.server.Serve(h.ln)
		h.logger.Debug("server shutting down...")
		// we do not want a http.ErrHTTPClosed error to be sent
		if err == http.ErrServerClosed {
			errCh <- nil
		}
		// send any other error, including nil
		errCh <- err
	}()
	// wait until the above go routine closes
	select {
	case <-ctx.Done():
		return h.shutdown(ctx)
	case err = <-errCh:
		return err
	}
}

// Shutdown runs the sver shutdown method which gracefully finished off any current requests,
// and prevents new ones from coming in
func (h *HTTP) shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
