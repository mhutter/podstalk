package podwatcher

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// PodWatcher watches pods and notifies others about new and changed pods
type PodWatcher struct {
	corev1.PodsGetter

	Namespace string
	PodNames  []string
}

// New creates a new PodWatcher
func New(cs *kubernetes.Clientset, namespace string) *PodWatcher {
	return &PodWatcher{
		PodsGetter: cs.CoreV1(),
		Namespace:  namespace,
		PodNames:   []string{},
	}
}

// List all pods in the namespace
func (p *PodWatcher) List() {
	pods, err := p.Pods(p.Namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("ERROR listing pods:", err)
	}

	podNames := make([]string, 0, len(pods.Items))
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.ObjectMeta.Name)
	}
	p.PodNames = podNames
}
