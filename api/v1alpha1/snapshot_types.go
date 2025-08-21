package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type SnapshotSpec struct {
    PolicyRef string   `json:"policyRef"`
    Reason    string   `json:"reason,omitempty"`
    Targets   []string `json:"targets,omitempty"`
}

type SnapshotStatus struct {
    Phase   string       `json:"phase,omitempty"` // Pending, Running, Succeeded, Failed
    URIs    []string     `json:"uris,omitempty"`
    Started *metav1.Time `json:"started,omitempty"`
    Ended   *metav1.Time `json:"ended,omitempty"`
    Message string       `json:"message,omitempty"`
    Completed int32      `json:"completed,omitempty"`
    Total     int32      `json:"total,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Snapshot struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   SnapshotSpec   `json:"spec,omitempty"`
    Status SnapshotStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type SnapshotList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Snapshot `json:"items"`
}

func init() {
    SchemeBuilder.Register(&Snapshot{}, &SnapshotList{})
}
