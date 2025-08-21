package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	snapshotv1alpha1 "github.com/SMLoureiro/vector-snapshot-operator/api/v1alpha1"
)

type Uploader interface {
    Upload(ctx context.Context, localPath string) (string, error)
}

type localUploader struct {
    base string
}

func (l *localUploader) Upload(ctx context.Context, localPath string) (string, error) {
    // Copy the file into a local folder to simulate upload
    if _, err := os.Stat(localPath); err != nil {
        return "", err
    }
    os.MkdirAll(l.base, 0o755)
    dst := filepath.Join(l.base, filepath.Base(localPath))
    in, err := os.Open(localPath)
    if err != nil { return "", err }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil { return "", err }
    defer out.Close()
    if _, err := io.Copy(out, in); err != nil { return "", err }
    return "file://" + dst, nil
}

// NewUploader returns a local uploader placeholder. Replace with S3/GCS/Azure implementations.
func NewUploader(ctx context.Context, _ *snapshotv1alpha1.SnapshotStorage) (Uploader, error) {
    return &localUploader{base: "/tmp/vector-snapshots"}, nil
}
