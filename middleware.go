package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		buf, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("ERROR: %s\n, err")
			return
		}

		fmt.Println(string(buf))

		if nc != nil {
			err = nc.Publish("visitor", buf)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
			}
		}
	}()
}

func chain(before, after http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		before.ServeHTTP(w, r)
		after.ServeHTTP(w, r)
	})
}
