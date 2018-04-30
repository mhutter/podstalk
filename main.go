package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// PodInfo contains all interesting information about a pod
type PodInfo struct {
	Name           string
	Namespace      string
	IP             string
	ServiceAccount string
	NodeName       string
	NodeIP         string
	Info           map[string]string
	Now            string
	Title          string
	Siblings       []string
}

var (
	podInfo PodInfo
	ips     []string
	client  *K8sClient
	t       *template.Template
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	p := podInfo
	p.Now = time.Now().Format(time.RFC3339Nano)

	if client != nil {
		for _, pod := range client.ListPods().Items {
			p.Siblings = append(p.Siblings, pod.Metadata.Name)
		}
	}

	err := t.Execute(w, p)
	check(err)
}

func init() {
	t = template.Must(template.New("index").Parse(htmlTemplate))
	podInfo = PodInfo{
		Name:           os.Getenv("POD_NAME"),
		Namespace:      os.Getenv("POD_NAMESPACE"),
		IP:             os.Getenv("POD_IP"),
		ServiceAccount: os.Getenv("POD_SA"),
		NodeName:       os.Getenv("NODE_NAME"),
		NodeIP:         os.Getenv("NODE_IP"),
		Title:          getEnvOr("TITLE", "Podstalk"),
	}
	client = NewClient()
	collectIPs()
	collectEnv()
}

func collectEnv() {
	pattern := regexp.MustCompile("^info_([^=]+)=(.+)$")
	info := map[string]string{}

	for _, line := range os.Environ() {
		l := strings.ToLower(line)
		if m := pattern.FindStringSubmatch(l); m != nil {
			info[m[1]] = m[2]
		}
	}

	podInfo.Info = info
}

func collectIPs() {
	interfaces, err := net.Interfaces()
	check(err)

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		check(err)
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

func getEnvOr(name, fallback string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return fallback
}

func main() {
	port := getEnvOr("PORT", "8080")

	handler := chain(
		http.HandlerFunc(jsonAccessLogger),
		http.HandlerFunc(infoHandler),
	)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	for _, ip := range ips {
		log.Printf("Listening on http://%s:%s/\n", ip, port)
	}
	log.Fatal(s.ListenAndServe())
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
