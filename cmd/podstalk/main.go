package main

import (
	"log"
	"sync"

	"github.com/mhutter/podstalk/server"
	"github.com/mhutter/podstalk/watcher"
)

func main() {
	// Configure logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Read configurations
	port := getPort()
	namespace := getNamespace()
	kubeconfig := getKubeconfig()

	stop := make(chan struct{})
	var wg sync.WaitGroup

	go startWatcher(kubeconfig, namespace, stop, &wg)
	startServer(port)
	wg.Wait()
}

func startServer(port string) {
	s := server.New(":" + port)
	s.Start()
}

func startWatcher(kubeconfig, namespace string, stop <-chan struct{}, wg *sync.WaitGroup) {
	// Create clientset
	clientset, err := getClientset(kubeconfig)
	if err != nil {
		log.Fatalln("ERROR configuring Kubernetes client:", err)
	}

	// Create & start watcher
	wg.Add(1)
	w := watcher.New(clientset, namespace)
	w.Watch(stop)
	wg.Done()
}
