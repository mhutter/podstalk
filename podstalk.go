package podstalk

import (
	"log"
	"reflect"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// PodWatcher watches pods and notifies others about new and changed pods
type PodWatcher struct {
	corev1.PodsGetter

	Namespace string
	Registry  Registry
}

// PodStatus contains some info about a pod
type PodStatus struct {
	Name      string      `json:"name"`
	Namespace string      `json:"namespace"`
	Node      string      `json:"node"`
	Phase     v1.PodPhase `json:"phase"`
	UID       types.UID   `json:"uid"`
}

// Registry stores infos about pods
type Registry map[types.UID]*PodStatus

// New creates a new PodWatcher
func New(cs *kubernetes.Clientset, namespace string) *PodWatcher {
	return &PodWatcher{
		PodsGetter: cs.CoreV1(),
		Namespace:  namespace,
		Registry:   Registry{},
	}
}

// Watch watch for pod updates in namespace
func (p *PodWatcher) Watch() {
	watcher, err := p.Pods(p.Namespace).Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("ERROR watching pods:", err)
	}

	for e := range watcher.ResultChan() {
		p.handleEvent(e)
	}
}

func (p *PodWatcher) handleEvent(e watch.Event) {
	pod := e.Object.(*v1.Pod)
	uid := pod.ObjectMeta.UID
	ps := &PodStatus{
		Name:      pod.ObjectMeta.Name,
		Namespace: pod.ObjectMeta.Namespace,
		Node:      pod.Spec.NodeName,
		Phase:     pod.Status.Phase,
		UID:       pod.ObjectMeta.UID,
	}

	switch e.Type {
	case watch.Added:
		p.Registry[uid] = ps
		p.publishEvent(e.Type, ps)

	case watch.Modified:
		if !reflect.DeepEqual(p.Registry[uid], ps) {
			p.Registry[uid] = ps
			p.publishEvent(e.Type, ps)
		}

	case watch.Deleted:
		delete(p.Registry, uid)
		p.publishEvent(e.Type, ps)
	}
}

func (p *PodWatcher) publishEvent(typ watch.EventType, ps *PodStatus) {
	// TODO: Publish somewhere
	log.Printf("%-8s - %#v\n", typ, ps)
}
