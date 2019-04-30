package podstalk

import (
	"k8s.io/apimachinery/pkg/watch"
)

// Event is sent when a pod is created, modified or deleted
type Event struct {
	Type watch.EventType
	Pod  *PodStatus
}
