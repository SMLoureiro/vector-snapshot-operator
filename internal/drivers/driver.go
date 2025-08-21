package drivers

import (
	"context"
	"errors"

	snapshotv1alpha1 "github.com/SMLoureiro/vector-snapshot-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Target is an abstract shard/pod/endpoint to snapshot.
type Target interface{ Name() string }

type Uploader interface {
	// Upload a local file path and return a URI (s3://..., gs://..., file://...)
	Upload(ctx context.Context, localPath string) (string, error)
}

type EngineDriver interface {
	DiscoverTargets(ctx context.Context, selector *metav1.LabelSelector) ([]Target, error)
	SnapshotTarget(ctx context.Context, t Target, up Uploader) (string, error)
}

func NewDriver(_ context.Context, pol snapshotv1alpha1.SnapshotPolicy, k8s client.Client) (EngineDriver, error) {
	switch pol.Spec.Engine {
	case "Qdrant":
		return NewQdrantDriver(pol, k8s), nil
	case "GenericExec":
		return NewGenericExecDriver(pol, k8s), nil
	default:
		return nil, errors.New("engine not implemented: " + string(pol.Spec.Engine))
	}
}
