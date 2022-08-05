// Code generated by skv2. DO NOT EDIT.

// Definitions for the Kubernetes types
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +genclient:noStatus

// GroupVersionKind for Probe
var ProbeGVK = schema.GroupVersionKind{
	Group:   "probes.bumblebee.io",
	Version: "v1alpha1",
	Kind:    "Probe",
}

// Probe is the Schema for the probe API
type Probe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ProbeSpec `json:"spec,omitempty"`
}

// GVK returns the GroupVersionKind associated with the resource type.
func (Probe) GVK() schema.GroupVersionKind {
	return ProbeGVK
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProbeList contains a list of Probe
type ProbeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Probe `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Probe{}, &ProbeList{})
}
