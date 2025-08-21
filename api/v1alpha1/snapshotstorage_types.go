package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:validation:Enum=S3;GCS;AzureBlob;NFS
type BackendType string

type S3Spec struct {
    Bucket string `json:"bucket"`
    Prefix string `json:"prefix,omitempty"`
    Region string `json:"region,omitempty"`
    CredentialsSecretRef *SecretKeyRef `json:"credentialsSecretRef,omitempty"`
}

type GCSSpec struct {
    Bucket string `json:"bucket"`
    Prefix string `json:"prefix,omitempty"`
    ServiceAccountSecretRef *SecretKeyRef `json:"serviceAccountSecretRef,omitempty"`
}

type AzureBlobSpec struct {
    Container string `json:"container"`
    Prefix    string `json:"prefix,omitempty"`
    ConnectionStringSecretRef *SecretKeyRef `json:"connectionStringSecretRef,omitempty"`
}

type NFSSpec struct {
    Server string `json:"server"`
    Path   string `json:"path"`
}

type SecretKeyRef struct {
    Name string `json:"name"`
    Key  string `json:"key"`
}

type SnapshotStorageSpec struct {
    Type  BackendType   `json:"type"`
    S3    *S3Spec       `json:"s3,omitempty"`
    GCS   *GCSSpec      `json:"gcs,omitempty"`
    Azure *AzureBlobSpec `json:"azure,omitempty"`
    NFS   *NFSSpec      `json:"nfs,omitempty"`
}

type SnapshotStorageStatus struct {
    Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=snstore
type SnapshotStorage struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   SnapshotStorageSpec   `json:"spec,omitempty"`
    Status SnapshotStorageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type SnapshotStorageList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []SnapshotStorage `json:"items"`
}

func init() {
    SchemeBuilder.Register(&SnapshotStorage{}, &SnapshotStorageList{})
}
