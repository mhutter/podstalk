package services

import (
	"log"
	"net/http"
	"time"
)

// Server handles HTTP and Websocket requests
type Server struct {
	http.Server

	Registry Registry
}

// NewServer returns a configured server
func NewServer(addr string, r Registry) *Server {
	return &Server{
		Registry: r,
		Server: http.Server{
			Addr:           addr,
			Handler:        nil,
			ReadTimeout:    8 * time.Second,
			WriteTimeout:   8 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Start starts up the HTTP server
func (s *Server) Start() {
	log.Printf("\x1b[32mListening on %s\x1b[0m", s.Addr)
	log.Fatal(s.ListenAndServe())
}
