package drivers

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    snapshotv1alpha1 "github.com/SMLoureiro/vector-snapshot-operator/api/v1alpha1"
    "sigs.k8s.io/controller-runtime/pkg/client"
)

type qdrantDriver struct {
    pol snapshotv1alpha1.SnapshotPolicy
    k8s client.Client
}

type simpleTarget struct{ name string }
func (s simpleTarget) Name() string { return s.name }

func NewQdrantDriver(pol snapshotv1alpha1.SnapshotPolicy, k8s client.Client) EngineDriver {
    return &qdrantDriver{pol: pol, k8s: k8s}
}

func (d *qdrantDriver) DiscoverTargets(ctx context.Context, _ *snapshotv1alpha1.LabelSelector) ([]Target, error) {
    if d.pol.Spec.Qdrant == nil || d.pol.Spec.Qdrant.BaseURL == "" {
        return nil, fmt.Errorf("qdrant.baseURL not set")
    }
    return []Target{simpleTarget{name: d.pol.Spec.Qdrant.BaseURL}}, nil
}

func (d *qdrantDriver) SnapshotTarget(ctx context.Context, t Target, up Uploader) (string, error) {
    // MVP: create a tiny dummy file to simulate a snapshot artifact.
    tmp := filepath.Join(os.TempDir(), "qdrant-snapshot-"+sanitize(t.Name())+".tgz")
    if err := os.WriteFile(tmp, []byte("qdrant snapshot placeholder"), 0o644); err != nil {
        return "", err
    }
    return up.Upload(ctx, tmp)
}

func sanitize(s string) string {
    b := make([]rune, 0, len(s))
    for _, r := range s {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
            b = append(b, r)
        }
    }
    if len(b) == 0 { return "target" }
    return string(b)
}
