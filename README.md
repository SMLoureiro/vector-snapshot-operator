# Vector DB Snapshot Operator

An opinionated Kubernetes Operator to orchestrate consistent, automated snapshots for multiple vector databases (Qdrant, Weaviate, Milvus, OpenSearch k‑NN/Vespa, and generic engines) and ship them to object storage with retention policies.

## Goals

- Engine‑agnostic CRDs with per‑engine adapters.
- **Cluster‑aware discovery**: label/annotation selectors; shard‑wise coordination.
- **Pluggable backends**: S3, GCS, Azure Blob, NFS/CSI.
- **Integrity**: checksums + manifest; optional quiesce hooks.
- **Observability**: conditions, events, Prometheus metrics.
- **Safety**: idempotent reconcile, finalizers, concurrency limits, per‑engine backoff.
