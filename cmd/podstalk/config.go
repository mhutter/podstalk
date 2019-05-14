package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultPort = "8080"
	nsPath      = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func getAddr() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":" + defaultPort
}

func getBasePath() string {
	return os.Getenv("BASE_PATH")
}

func getDebug() bool {
	return os.Getenv("DEBUG") != ""
}

func getNamespace() string {
	if ns := os.Getenv("NAMESPACE"); ns != "" {
		return ns
	}

	if ns, err := ioutil.ReadFile(nsPath); err == nil {
		return string(ns)
	}
	return ""
}

func getKubeconfig() string {
	if kc := os.Getenv("KUBECONFIG"); kc != "" {
		return kc
	}

	if home, err := os.UserHomeDir(); err == nil {
		kcPath := filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(kcPath); err != nil {
			log.Fatalln("Could not determine KUBECONFIG")
		}
		return kcPath
	}

	return ""
}
