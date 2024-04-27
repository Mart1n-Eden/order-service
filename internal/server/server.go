package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(h http.Handler, addr string) *Server {

	s := &Server{
		httpServer: &http.Server{
			Addr:         ":" + addr,
			Handler:      h,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	return s
}

func (s *Server) Run() {
	s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
}
