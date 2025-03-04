package server

import (
	"github.com/rs/cors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Option func(s *Server)

func WithAddress(address string) Option {
	return func(s *Server) {
		s.Address = address
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) {
		if logger != nil {
			s.Logger = logger
		}
	}
}

func WithMiddleware(middleware func(http.Handler) http.Handler) Option {
	return func(s *Server) {
		if middleware != nil {
			s.Handler = middleware(s.Handler)
		}
	}
}

func WithPathPrefix(prefix string) Option {
	return func(s *Server) {
		s.PathPrefix = prefix
	}
}

func WithCORS(enabled bool) Option {
	return func(s *Server) {
		if enabled {
			s.Handler = cors.Default().Handler(s.Handler)
		}
	}
}

func WithDebug(debug bool) Option {
	return func(s *Server) {
		s.Debug = debug
	}
}

func WithStartupTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		if timeout > 0 {
			s.StartupTimeout = timeout
		}
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		if timeout > 0 {
			s.ShutdownTimeout = timeout
		}
	}
}

func WithStripQueryString(enabled bool) Option {
	return func(s *Server) {
		if enabled {
			s.Handler = stripQueryStringHandler(s.Handler)
		}
	}
}
