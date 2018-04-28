package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	hostname string
	ips      []string
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	responseData := map[string]interface{}{
		"hostname": hostname,
		"ips":      ips,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(responseData)
}

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	interfaces, err := net.Interfaces()
	// handle err
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() != nil {
				ips = append(ips, ip.String())
			}
		}
	}

	for _, element := range os.Environ() {
		fmt.Println(element)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        http.HandlerFunc(infoHandler),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	for _, ip := range ips {
		log.Printf("Listening on http://%s:%s/\n", ip, port)
	}
	log.Fatal(s.ListenAndServe())
}
