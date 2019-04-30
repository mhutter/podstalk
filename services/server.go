package services

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server handles HTTP and Websocket requests
type Server struct {
	http.Server

	Registry Registry
}

// NewServer returns a configured server
func NewServer(addr string, reg Registry) *Server {
	r := mux.NewRouter()

	s := &Server{
		Registry: reg,
		Server: http.Server{
			Addr:           addr,
			Handler:        r,
			ReadTimeout:    8 * time.Second,
			WriteTimeout:   8 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	r.HandleFunc("/api/pods", s.ListPods)

	return s
}

// Start starts up the HTTP server
func (s *Server) Start() {
	log.Printf("\x1b[32mListening on %s\x1b[0m", s.Addr)
	log.Fatal(s.ListenAndServe())
}

// ListPods GET /api/pods - returns an array of all pods
func (s *Server) ListPods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(s.Registry.ListPods())
}
