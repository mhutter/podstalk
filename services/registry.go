package services

import (
	"log"
	"reflect"

	"github.com/mhutter/podstalk"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

// Registry stores infos about pods
type Registry map[types.UID]*podstalk.PodStatus

// NewRegistry returns a new registry
func NewRegistry() Registry {
	return Registry{}
}

// Start the registry service. It will handle Events received via "events",
// update its internal state and emmit the Event via "updates" if something
// actually changed.
func (r Registry) Start(events <-chan *podstalk.Event) <-chan *podstalk.Event {
	updates := make(chan *podstalk.Event)
	go r.listen(events, updates)
	log.Println("Registry ready")
	return updates
}

func (r Registry) listen(events <-chan *podstalk.Event, updates chan<- *podstalk.Event) {
	for e := range events {
		switch e.Type {
		case watch.Added:
			r[e.Pod.UID] = e.Pod
			updates <- e

		case watch.Modified:
			// Only update and publish events if anything has actually changed.
			// Since the controller will publish ALL updates, we might get
			// updates on pods that have only updates on fields not represented
			// in PodStatus
			if !reflect.DeepEqual(r[e.Pod.UID], e.Pod) {
				r[e.Pod.UID] = e.Pod
				updates <- e
			}

		case watch.Deleted:
			delete(r, e.Pod.UID)
			updates <- e
		}
	}
	log.Println("Registry stopped")
}

// ListPods returns all pods in the registry as a slice
func (r Registry) ListPods() []podstalk.PodStatus {
	list := make([]podstalk.PodStatus, 0, len(r))

	for _, p := range r {
		pod := *p
		list = append(list, pod)
	}

	return list
}
