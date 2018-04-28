package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	podInfo = map[string]interface{}{}
	ips     []string
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(podInfo)
}

func init() {
	collectHostname()
	collectIPs()
	collectInfo()
}

func collectInfo() {
	pattern := regexp.MustCompile("^info_([^=]+)=(.+)$")
	info := map[string]string{}

	for _, line := range os.Environ() {
		l := strings.ToLower(line)
		if m := pattern.FindStringSubmatch(l); m != nil {
			info[m[1]] = m[2]
		}
	}

	podInfo["info"] = info
}

func collectIPs() {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

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

	podInfo["ips"] = ips
}

func collectHostname() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	podInfo["hostname"] = hostname
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
