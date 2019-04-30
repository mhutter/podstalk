package main

import (
	"log"

	"github.com/mhutter/podstalk/services"
)

func main() {
	// Configure logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Read configurations
	addr := getAddr()
	namespace := getNamespace()
	kubeconfig := getKubeconfig()

	// Create services
	w := services.NewWatcher(kubeconfig, namespace)
	r := services.NewRegistry()
	s := services.NewServer(addr, r)

	// Start services
	w.Start()
	updates := r.Start(w.Events)

	// Post updates to log
	go func() {
		for e := range updates {
			log.Printf("%-8s - %s", e.Type, e.Pod.Name)
		}
	}()

	s.Start()
}
