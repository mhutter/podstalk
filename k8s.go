package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	tokenPath     = "/run/secrets/kubernetes.io/serviceaccount/token"
	namespacePath = "/run/secrets/kubernetes.io/serviceaccount/namespace"
	caPath        = "/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

// PodList holds a list of pods from the k8s api
type PodList struct {
	Items []struct {
		Metadata struct {
			Name   string            `json:"name"`
			Labels map[string]string `json:"labels"`
		} `json:"metadata"`
	} `json:"items"`
}

type K8sClient struct {
	Server    string
	Token     string
	Namespace string
	client    http.Client
}

func NewClient() *K8sClient {
	var buf []byte

	server := fmt.Sprintf(
		"https://%s:%s",
		os.Getenv("KUBERNETES_SERVICE_HOST"),
		os.Getenv("KUBERNETES_SERVICE_PORT"),
	)

	buf, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil
	}
	token := string(buf)

	buf, err = ioutil.ReadFile(namespacePath)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil
	}
	namespace := string(buf)

	caPool := x509.NewCertPool()
	buf, err = ioutil.ReadFile(caPath)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil
	}
	caPool.AppendCertsFromPEM(buf)
	tlsConfig := &tls.Config{RootCAs: caPool}

	return &K8sClient{
		Server:    server,
		Token:     token,
		Namespace: namespace,
		client: http.Client{
			Transport: &http.Transport{TLSClientConfig: tlsConfig},
			Timeout:   5 * time.Second,
		},
	}
}

func (c *K8sClient) ListPods() (list PodList) {
	var err error
	url := fmt.Sprintf("%s/api/v1/namespaces/%s/pods", c.Server, c.Namespace)
	var req *http.Request
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	var res *http.Response
	if res, err = c.client.Do(req); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	if err = json.NewDecoder(res.Body).Decode(&list); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	return
}
