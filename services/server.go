package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mhutter/podstalk"
)

// Server handles HTTP and Websocket requests
type Server struct {
	http.Server

	Registry Registry
	Debug    bool

	pendingRequests chan struct{}
	upgrader        websocket.Upgrader
	clients         clients
}

type clients map[*websocket.Conn]*sync.Mutex

// NewServer returns a configured server
func NewServer(addr string, reg Registry, basePath string) *Server {
	r := mux.NewRouter()

	s := &Server{
		Registry: reg,
		Debug:    false,
		Server: http.Server{
			Addr:           addr,
			Handler:        r,
			ReadTimeout:    8 * time.Second,
			WriteTimeout:   8 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// We don't really care about the origin, so simply accept
			// the connection
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(clients),
	}

	// API Routes
	api := r.PathPrefix(basePath + "/api").Subrouter()
	api.HandleFunc("/pods", s.ListPods)
	api.HandleFunc("/ws", s.HandleSocket)

	// Serve static files
	fs := http.StripPrefix(basePath, http.FileServer(http.Dir("public")))
	r.PathPrefix(basePath + "/").Handler(fs).Methods("GET")

	return s
}

// Start starts up the HTTP server
func (s *Server) Start() {
	go s.broadcast()
	go s.start()
}

func (s *Server) start() {
	s.pendingRequests = make(chan struct{})
	defer close(s.pendingRequests)

	log.Printf("\x1b[32mListening on %s\x1b[0m", s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("Error from Server:", err)
	}
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() {
	if err := s.Shutdown(context.Background()); err != nil {
		log.Println("Error during HTTP server shutdown:", err)
	}

	// Wait for all requests to finish
	<-s.pendingRequests
	log.Println("Server stopped")
}

// ListPods GET /api/pods - returns an array of all pods
func (s *Server) ListPods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(s.Registry.ListPods())
}

// HandleSocket GET /api/ws - upgrades connection to a WebSocket
func (s *Server) HandleSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	if s.Debug {
		log.Println("New client:", r.RemoteAddr)
	}
	defer conn.Close()

	// Send an ADDED event for all pods currently in the registry
	for _, pod := range s.Registry.ListPods() {
		e := podstalk.Event{
			Type: "ADDED",
			Pod:  &pod,
		}
		conn.WriteJSON(e)
	}
	s.clients[conn] = &sync.Mutex{}

	// Listen on the connection until it is closed
	for {
		if _, _, err := conn.NextReader(); err != nil {
			if err2, ok := err.(*websocket.CloseError); ok {
				if err2.Code == websocket.CloseGoingAway ||
					err2.Code == websocket.CloseNormalClosure {
					// normal close codes, don't care
					break
				}
			}

			log.Println("Client error:", err)
		}
	}

	delete(s.clients, conn)

	if s.Debug {
		log.Println("Client disconnected:", r.RemoteAddr)
	}
}

func (s *Server) broadcast() {
	for e := range s.Registry.Updates {
		for c, l := range s.clients {
			go publish(e, c, l)
		}
	}
}

func publish(e *podstalk.Event, c *websocket.Conn, lock *sync.Mutex) {
	lock.Lock()
	defer lock.Unlock()
	if err := c.WriteJSON(e); err != nil {
		log.Println("Error sending event:", err)
	}
}
