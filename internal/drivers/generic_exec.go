package drivers

import (
	"context"
	"os"
	"path/filepath"

	snapshotv1alpha1 "github.com/SMLoureiro/vector-snapshot-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type genericExecDriver struct {
	pol snapshotv1alpha1.SnapshotPolicy
	k8s client.Client
}

type podTarget struct{ pod corev1.Pod }
func (p podTarget) Name() string { return p.pod.Namespace + "-" + p.pod.Name }

func NewGenericExecDriver(pol snapshotv1alpha1.SnapshotPolicy, k8s client.Client) EngineDriver {
	return &genericExecDriver{pol: pol, k8s: k8s}
}

func (d *genericExecDriver) DiscoverTargets(ctx context.Context, sel *metav1.LabelSelector) ([]Target, error) {
	var pods corev1.PodList
	opts := []client.ListOption{}
	if sel != nil {
		ls, _ := metav1.LabelSelectorAsSelector(sel)
		opts = append(opts, client.MatchingLabelsSelector{Selector: ls})
	}
	if err := d.k8s.List(ctx, &pods, opts...); err != nil { return nil, err }
	res := make([]Target, 0, len(pods.Items))
	for _, p := range pods.Items {
		res = append(res, podTarget{pod: p})
	}
	return res, nil
}

func (d *genericExecDriver) SnapshotTarget(ctx context.Context, t Target, up Uploader) (string, error) {
	// Stub: create a dummy artifact file to upload.
	tmp := filepath.Join(os.TempDir(), "generic-snapshot-"+t.Name()+".tgz")
	if err := os.WriteFile(tmp, []byte("generic snapshot placeholder"), 0o644); err != nil {
		return "", err
	}
	return up.Upload(ctx, tmp)
}
