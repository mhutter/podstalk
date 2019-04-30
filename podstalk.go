package podstalk

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// PodStatus contains some info about a pod
type PodStatus struct {
	Name      string      `json:"name"`
	Namespace string      `json:"namespace"`
	Node      string      `json:"node"`
	Phase     v1.PodPhase `json:"phase"`
	UID       types.UID   `json:"uid"`
}
