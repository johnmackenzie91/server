package server

import (
	"net/http"
)

type Option func(s *HTTP)

func WithAddress(address string) Option {
	return func(s *HTTP) {
		s.addr = address
	}
}

func WithHandler(handler http.Handler) Option {
	return func(s *HTTP) {
		s.server.Handler = handler
	}
}

func WithLogger(logger Logger) Option {
	return func(s *HTTP) {
		s.logger = logger
	}
}

func WithServer(svr *http.Server) Option {
	return func(s *HTTP) {
		s.server = svr
	}
}
