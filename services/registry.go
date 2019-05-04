package services

import (
	"log"
	"reflect"

	"github.com/mhutter/podstalk"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

// Registry stores infos about pods
type Registry struct {
	Pods    map[types.UID]*podstalk.PodStatus
	Updates chan *podstalk.Event
	Debug   bool
}

// NewRegistry returns a new registry
func NewRegistry() Registry {
	return Registry{
		Pods:    make(map[types.UID]*podstalk.PodStatus),
		Updates: make(chan *podstalk.Event, 10),
		Debug:   false,
	}
}

// Start the registry service. It will handle Events received via "events",
// update its internal state and emmit the Event via "updates" if something
// actually changed.
func (r Registry) Start(events <-chan *podstalk.Event) {
	go r.listen(events)
	log.Println("Registry ready")
}

func (r Registry) listen(events <-chan *podstalk.Event) {
	for e := range events {
		switch e.Type {
		case watch.Added:
			r.Pods[e.Pod.UID] = e.Pod
			r.publish(e)

		case watch.Modified:
			// Only update and publish events if anything has actually changed.
			// Since the controller will publish ALL updates, we might get
			// updates on pods that have only updates on fields not represented
			// in PodStatus
			if !reflect.DeepEqual(r.Pods[e.Pod.UID], e.Pod) {
				r.Pods[e.Pod.UID] = e.Pod
				r.publish(e)
			}

		case watch.Deleted:
			delete(r.Pods, e.Pod.UID)
			r.publish(e)
		}
	}
	log.Println("Registry stopped")
}

func (r Registry) publish(e *podstalk.Event) {
	r.Updates <- e
	if r.Debug {
		log.Printf("%-8s - %s", e.Type, e.Pod.Name)
	}
}

// ListPods returns all pods in the registry as a slice
func (r Registry) ListPods() []podstalk.PodStatus {
	list := make([]podstalk.PodStatus, 0, len(r.Pods))

	for _, p := range r.Pods {
		pod := *p
		list = append(list, pod)
	}

	return list
}
