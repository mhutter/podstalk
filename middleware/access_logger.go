package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mhutter/podstalk"
	nats "github.com/nats-io/go-nats"
)

type AccessLogger struct {
	nc *nats.Conn
}

func NewAccessLogger() *AccessLogger {
	return &AccessLogger{
		nc: connectNats(),
	}
}

func (al *AccessLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := podstalk.Visit{
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

	if al.nc != nil {
		err = al.nc.Publish("visit", buf)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
	}
}

func connectNats() *nats.Conn {
	nurl := podstalk.GetEnvOr("NATS_URL", "nats://localhost:4222")
	nc, err := nats.Connect(nurl)

	if err != nil {
		log.Printf("WARN: Could not connect to NATS: %s\n", err)
		return nil
	}
	return nc
}
