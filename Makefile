# --- tools ---
CONTROLLER_GEN ?= $(shell which controller-gen)
KIND ?= kind
KUBECTL ?= kubectl

# --- cluster names/dirs ---
CLUSTER_NAME ?= vso-dev
CRD_DIR := config/crd/bases

.PHONY: tidy build run
tidy:
	@go mod tidy

build:
	@go build ./...

run:
	@go run ./cmd/manager

# --- code gen & CRDs ---
.PHONY: generate manifests install uninstall
generate:
	@[ -n "$(CONTROLLER_GEN)" ] || (echo "controller-gen not found. install: go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.18.0" && exit 1)
	@$(CONTROLLER_GEN) object:headerFile="" paths=./api/...

manifests: generate
	@[ -n "$(CONTROLLER_GEN)" ] || (echo "controller-gen not found. install: go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.18.0" && exit 1)
	@mkdir -p $(CRD_DIR)
	@$(CONTROLLER_GEN) crd:crdVersions=v1 paths=./api/... output:crd:dir=$(CRD_DIR)

install: manifests
	@$(KUBECTL) apply -f $(CRD_DIR)

uninstall:
	-@$(KUBECTL) delete -f $(CRD_DIR) --ignore-not-found

# --- kind cluster lifecycle ---
.PHONY: kind-up kind-down kind-kubeconfig
kind-up:
	@$(KIND) get clusters | grep -qx $(CLUSTER_NAME) || $(KIND) create cluster --name $(CLUSTER_NAME)
	@echo "✔ kind cluster $(CLUSTER_NAME) ready"

kind-down:
	-@$(KIND) delete cluster --name $(CLUSTER_NAME)

kind-kubeconfig:
	@$(KIND) get kubeconfig --name $(CLUSTER_NAME) > .kubeconfig
	@echo "export KUBECONFIG=$$(pwd)/.kubeconfig"

# --- smoke test workflow ---
.PHONY: dev-setup dev-run dev-apply dev-snapshot dev-watch
dev-setup: kind-up install
	@echo "✔ CRDs installed in kind: $(CLUSTER_NAME)"

dev-run:
	@echo "Note: run this in a separate terminal to stream logs."
	@echo "Artifacts will be saved locally to /tmp/vector-snapshots"
	@echo "Starting controller..."
	@$(MAKE) run

dev-apply:
	@$(KUBECTL) apply -f config/samples/storage_s3.yaml
	@# faster schedule during dev:
	@$(KUBECTL) apply -f config/samples/policy_qdrant.yaml
	@$(KUBECTL) get snapshotpolicies -A

dev-snapshot:
	@$(KUBECTL) apply -f config/samples/manual_snapshot.yaml
	@$(KUBECTL) get snapshots -A

dev-watch:
	@$(KUBECTL) get snapshots -A -w


.PHONY: test
test: manifests
	@go test ./controllers -v -count=1
