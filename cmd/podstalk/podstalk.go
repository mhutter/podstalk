package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/mhutter/podstalk"
	"github.com/mhutter/podstalk/middleware"
)

var (
	ips []string
)

func init() {
	collectIPs()
}

func collectIPs() {
	interfaces, err := net.Interfaces()
	podstalk.Check(err)

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		podstalk.Check(err)
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
}

func main() {
	port := podstalk.GetEnvOr("PORT", "8080")
	mux := http.NewServeMux()

	mux.Handle("/", middleware.Chain(
		middleware.NewAccessLogger(),
		podstalk.NewInfoHandler(),
	))

	mux.HandleFunc("/appuioli.png", podstalk.LogoHandler)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	for _, ip := range ips {
		log.Printf("Listening on http://%s:%s/\n", ip, port)
	}
	log.Fatal(s.ListenAndServe())
}
