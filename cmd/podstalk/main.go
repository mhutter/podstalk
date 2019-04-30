package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/mhutter/podstalk/kube"
	"github.com/mhutter/podstalk/watcher"
)

func main() {
	// Configure logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Define & parse flags
	var kubeconfig, namespace string
	kube.CommonFlags(&kubeconfig)
	namespaceFlag(&namespace)
	flag.Parse()

	stop := make(chan struct{})
	var wg sync.WaitGroup

	startWatcher(kubeconfig, namespace, stop, &wg)
	wg.Wait()
}

func namespaceFlag(namespace *string) {
	path := "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	if file, err := os.Open(path); err == nil {
		defer file.Close()
		buf, _ := ioutil.ReadAll(file)
		defaultNamespace := strings.Trim(string(buf), " \n")
		flag.StringVar(namespace, "namespace", defaultNamespace,
			"(optional) namespace to use")
	} else {
		flag.StringVar(namespace, "namespace", os.Getenv("NAMESPACE"),
			"namespace to use")
	}
}

func startWatcher(kubeconfig, namespace string, stop <-chan struct{}, wg *sync.WaitGroup) {
	// Create clientset
	clientset, err := kube.GetClientset(kubeconfig)
	if err != nil {
		log.Fatalln("ERROR configuring Kubernetes client:", err)
	}

	// Create & start watcher
	wg.Add(1)
	w := watcher.New(clientset, namespace)
	w.Watch(stop)
	wg.Done()
}
