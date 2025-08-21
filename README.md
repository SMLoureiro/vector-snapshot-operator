# Vector DB Snapshot Operator

An opinionated Kubernetes Operator to orchestrate consistent, automated snapshots for multiple vector databases (Qdrant, Weaviate, Milvus, OpenSearch k‑NN/Vespa, and generic engines) and ship them to object storage with retention policies.

## Goals

- Engine‑agnostic CRDs with per‑engine adapters.
- **Cluster‑aware discovery**: label/annotation selectors; shard‑wise coordination.
- **Pluggable backends**: S3, GCS, Azure Blob, NFS/CSI.
- **Integrity**: checksums + manifest; optional quiesce hooks.
- **Observability**: conditions, events, Prometheus metrics.
- **Safety**: idempotent reconcile, finalizers, concurrency limits, per‑engine backoff.

## Local testing

### 0) Ensure code-gen & CRDs exist
make manifests

### 1) Spin up kind & install CRDs
make dev-setup

### 2) In a second terminal: run the controller locally
make dev-run

### 3) Back in the first terminal: apply sample resources
make dev-apply

### 4) Trigger a manual snapshot and watch
make dev-snapshot
make dev-watch
