//go:build tools
// +build tools

package v1alpha1

// This file pins controller-gen via go.mod when using `go generate` with -tags tools.
import _ "sigs.k8s.io/controller-tools/cmd/controller-gen"
