package main

import (
	"encoding/json"
	"log"
	"runtime"
	"strings"

	"github.com/mhutter/podstalk"
	nats "github.com/nats-io/go-nats"
)

func handleVisit(m *nats.Msg) {
	v := &podstalk.Visit{}
	json.Unmarshal(m.Data, v)
	log.Printf(
		"New visit from %s using a %s operating system!",
		v.RemoteAddr,
		extractOS(v.UserAgent),
	)
}

func main() {
	nurl := podstalk.GetEnvOr("NATS_URL", nats.DefaultURL)
	nc, err := nats.Connect(nurl)
	podstalk.Check(err)

	nc.Subscribe(podstalk.Topic, handleVisit)

	log.Printf(
		"Connected to NATS at %s and listening to the '%s' topic!",
		nurl,
		podstalk.Topic,
	)

	runtime.Goexit()
}

func extractOS(ua string) string {
	os := "unknown"
	ua = strings.ToLower(ua)

	switch {
	case strings.Contains(ua, "linux"):
		os = "Linux"

	case strings.Contains(ua, "mac os"):
		os = "Mac"

	case strings.Contains(ua, "windows"):
		os = "Windows"
	}

	return os
}
