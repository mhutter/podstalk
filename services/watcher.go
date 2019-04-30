package services

import (
	"log"

	"github.com/mhutter/podstalk"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Watcher implements a simple service that watches for pod events (creation,
// modification, deleteion) and publishes them via a channel.
type Watcher struct {
	corev1.PodsGetter

	Events    chan *podstalk.Event
	namespace string
}

// NewWatcher returns a new, configured Watcher
func NewWatcher(kubeconfig, namespace string) *Watcher {
	// Create clientset
	cs, err := getClientset(kubeconfig)
	if err != nil {
		log.Fatalln("ERROR configuring Kubernetes client:", err)
	}

	return &Watcher{
		PodsGetter: cs.CoreV1(),
		Events:     make(chan *podstalk.Event),
		namespace:  namespace,
	}
}

// Start starts the watcher and returns a chan delivering pod change events.
func (w *Watcher) Start() {
	// start watching
	watcher, err := w.Pods(w.namespace).Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("ERROR watching pods:", err)
	}

	// start actually watching for events
	go func() {
		for e := range watcher.ResultChan() {
			w.handleEvent(e)
		}
	}()
}

func (w *Watcher) handleEvent(e watch.Event) {
	pod := e.Object.(*v1.Pod)
	ps := &podstalk.PodStatus{
		Name:      pod.ObjectMeta.Name,
		Namespace: pod.ObjectMeta.Namespace,
		Node:      pod.Spec.NodeName,
		Phase:     pod.Status.Phase,
		UID:       pod.ObjectMeta.UID,
	}

	w.Events <- &podstalk.Event{Type: e.Type, Pod: ps}
}

// getClientConfig returns a config for the k8s rest client. If kubeconfig is
// an empty string, InClusterConfig is used to retrieve the config, otherwise
// the contents of kubeconfig are used.
func getClientConfig(kubeconfig string) (cfg *rest.Config, err error) {
	if kubeconfig == "" {
		// Use in-cluster config
		cfg, err = rest.InClusterConfig()
	} else {
		// Load config from kube/conf
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return
}

// getClientset loads a local kubeconf file or determines its in-cluster
// configuration and returns a fully configured clientset.
func getClientset(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := getClientConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
