# Vector Snapshot Operator (Scaffold)

A multi-engine snapshot operator scaffold for Kubernetes. Engines supported (as stubs): **Qdrant**, **GenericExec**.
Storage uploader is a local placeholder (copies to `/tmp/vector-snapshots`). Swap with S3/GCS/Azure.

## Quick start (dev)
```bash
go mod tidy
make run
```
This starts the controller against your current kubeconfig. You'll need to apply the CRDs (generated via controller-gen in a full setup) or adapt types during development.

## Sample CRs
See `config/samples/` for a `SnapshotPolicy`, `SnapshotStorage`, and a manual `Snapshot`.

## Next steps
- Implement real Qdrant/Weaviate/Milvus/OpenSearch/Vespa drivers.
- Add S3/GCS/Azure uploaders.
- Add RBAC/CRDs via `controller-gen` or kubebuilder.
- Package a Helm chart.
