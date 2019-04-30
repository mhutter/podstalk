package watcher

import (
	"log"
	"reflect"

	"github.com/mhutter/podstalk"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// Watcher watches pods and notifies others about new and changed pods
type Watcher struct {
	corev1.PodsGetter

	Namespace string
	Registry  Registry
}

// Registry stores infos about pods
type Registry map[types.UID]*podstalk.PodStatus

// New creates a new Watcher
func New(cs *kubernetes.Clientset, namespace string) *Watcher {
	return &Watcher{
		PodsGetter: cs.CoreV1(),
		Namespace:  namespace,
		Registry:   Registry{},
	}
}

// Watch watch for pod updates in namespace
func (p *Watcher) Watch(done <-chan struct{}) {
	watcher, err := p.Pods(p.Namespace).Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("ERROR watching pods:", err)
	}

	for {
		select {
		case <-done:
			log.Println("Goodbye")
			watcher.Stop()
			return
		case e := <-watcher.ResultChan():
			p.handleEvent(e)
		}
	}
}

func (p *Watcher) handleEvent(e watch.Event) {
	pod := e.Object.(*v1.Pod)
	uid := pod.ObjectMeta.UID
	ps := &podstalk.PodStatus{
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

func (p *Watcher) publishEvent(typ watch.EventType, ps *podstalk.PodStatus) {
	// TODO: Publish somewhere
	log.Printf("%-8s - %#v\n", typ, ps)
}
