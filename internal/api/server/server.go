package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/musicmash/auth/internal/log"
)

type Server struct {
	server *http.Server
	Addr   string
}

type Options struct {
	IP           string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func New(router http.Handler, opts *Options) *Server {
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", opts.IP, opts.Port),
		Handler:      router,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		IdleTimeout:  opts.IdleTimeout,
	}
	return &Server{server: &server, Addr: server.Addr}
}

func (s *Server) ListenAndServe() error {
	log.Infof("server is ready to handle requests at: %v", s.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) SetKeepAlivesEnabled(v bool) {
	s.server.SetKeepAlivesEnabled(v)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
