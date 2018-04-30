package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type access struct {
	Host       string `json:"host"`
	Method     string `json:"method"`
	Proto      string `json:"proto"`
	RemoteAddr string `json:"remote_addr"`
	Path       string `json:"path"`
	UserAgent  string `json:"user_agent"`
}

func jsonAccessLogger(w http.ResponseWriter, r *http.Request) {
	go func() {
		msg := access{
			Host:       r.Host,
			Method:     r.Method,
			Proto:      r.Proto,
			RemoteAddr: r.RemoteAddr,
			Path:       r.URL.Path,
			UserAgent:  r.Header.Get("User-Agent"),
		}
		json.NewEncoder(os.Stdout).Encode(msg)
	}()
}

func chain(before, after http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		before.ServeHTTP(w, r)
		after.ServeHTTP(w, r)
	})
}
