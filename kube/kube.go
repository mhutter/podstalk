// Package kube contains common helper functions around kubernetes, for example
// to load configurations and such
package kube

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CommonFlags(kubeconfig *string) {
	if home := os.Getenv("HOME"); home != "" {
		flag.StringVar(kubeconfig, "kubeconfig",
			filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(kubeconfig, "kubeconfig", "",
			"absolute path to the kubeconfig file")
	}
}

// GetClientset loads a local kubeconf file or determines its in-cluster
// configuration and returns a fully configured clientset.
func GetClientset(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := getConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func getConfig(kubeconfig string) (cfg *rest.Config, err error) {
	if kubeconfig == "" {
		// Use in-cluster config
		cfg, err = rest.InClusterConfig()
	} else {
		// Load config from kube/conf
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return
}
