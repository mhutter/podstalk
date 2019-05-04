package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/mhutter/podstalk/services"
)

func main() {
	// Configure logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Booting up")

	// Read configurations
	addr := getAddr()
	namespace := getNamespace()
	kubeconfig := getKubeconfig()

	// Create services
	w := services.NewWatcher(kubeconfig, namespace)
	r := services.NewRegistry()
	s := services.NewServer(addr, r)

	// Start services
	updates := r.Start(w.Events)
	w.Start()
	s.Start()

	done := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Stop the watcher and thus the registry
		log.Println("Stopping services...")
		w.Stop()
		s.Stop()

		close(done)
	}()

	// Post updates to log
	go func() {
		for e := range updates {
			log.Printf("%-8s - %s", e.Type, e.Pod.Name)
		}
	}()

	<-done
	log.Println("Goodbye :(")
}
