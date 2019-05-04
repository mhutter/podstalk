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
	r.Start(w.Events)
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

	<-done
	log.Println("Goodbye :(")
}
