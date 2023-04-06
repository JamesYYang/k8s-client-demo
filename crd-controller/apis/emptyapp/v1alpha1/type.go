package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EmptyApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmptyAppSpec   `json:"spec"`
	Status EmptyAppStatus `json:"status"`
}

type EmptyAppSpec struct {
	ImageName string `json:"imageName"`
	Replicas  int32  `json:"replicas"`
	SvcPort   int32  `json:"svcPort"`
}

type EmptyAppStatus struct {
	AvailableReplicas int32  `json:"availableReplicas"`
	ClusterIP         string `json:"clusterIP"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EmptyAppList is a list of EmptyApp resources
type EmptyAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []EmptyApp `json:"items"`
}
