package engine

import (
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// GitRepository represents the operations that the engine needs to perform.
type GitRepository interface {
	Clone(string) error
	Open(string) error
	HeadHash() (plumbing.Hash, error)
	Sync() (plumbing.Hash, error)
	ParseManifests() ([]*unstructured.Unstructured, error)
	IsManaged(r *cache.Resource) bool
}
