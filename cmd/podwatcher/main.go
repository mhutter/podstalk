package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mhutter/podstalk/kube"
	"github.com/mhutter/podstalk/podwatcher"
)

func main() {
	// Configure logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Define & parse flags
	var kubeconfig, namespace string
	kube.CommonFlags(&kubeconfig)
	namespaceFlag(&namespace)
	flag.Parse()

	// Create clientset
	clientset, err := kube.GetClientset(kubeconfig)
	if err != nil {
		log.Fatalln("ERROR configuring Kubernetes client:", err)
	}

	// Create & start watcher
	watcher := podwatcher.New(clientset, namespace)
	watcher.List()

	log.Printf("We have %d pods!", len(watcher.PodNames))
	for _, n := range watcher.PodNames {
		log.Println(" -> ", n)
	}

	log.Println("Goodbye")
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
		flag.StringVar(namespace, "namespace", "",
			"namespace to use")
	}
}
