package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:validation:Enum=Qdrant;Weaviate;Milvus;OpenSearch;Vespa;GenericExec
type EngineType string

// Alias Kubernetes' LabelSelector so drivers can import it via our API package.
type LabelSelector = metav1.LabelSelector


type QdrantConfig struct {
    BaseURL string `json:"baseURL"`
    APIKeySecretRef *SecretKeyRef `json:"apiKeySecretRef,omitempty"`
    PerCollection bool `json:"perCollection,omitempty"`
    Collections []string `json:"collections,omitempty"`
}

type GenericExecConfig struct {
    Command []string `json:"command"`
    Paths   []string `json:"paths"`
}

type SnapshotPolicySpec struct {
    Schedule  string         `json:"schedule"`
    Retention int32          `json:"retention"`
    Selector  *metav1.LabelSelector `json:"selector,omitempty"`
    StorageRef string        `json:"storageRef"`
    Engine    EngineType     `json:"engine"`

    Qdrant  *QdrantConfig     `json:"qdrant,omitempty"`
    Generic *GenericExecConfig `json:"generic,omitempty"`

    MaxConcurrent int32 `json:"maxConcurrent,omitempty"`
    ShardTimeoutSeconds int32 `json:"shardTimeoutSeconds,omitempty"`
}

type SnapshotPolicyStatus struct {
    Conditions []metav1.Condition `json:"conditions,omitempty"`
    LastRun    *metav1.Time       `json:"lastRun,omitempty"`
    NextRun    *metav1.Time       `json:"nextRun,omitempty"`
    Successes  int64              `json:"successes,omitempty"`
    Failures   int64              `json:"failures,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=snappol
type SnapshotPolicy struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   SnapshotPolicySpec   `json:"spec,omitempty"`
    Status SnapshotPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type SnapshotPolicyList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []SnapshotPolicy `json:"items"`
}

func init() {
    SchemeBuilder.Register(&SnapshotPolicy{}, &SnapshotPolicyList{})
}
