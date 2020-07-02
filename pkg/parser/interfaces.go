package parser

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ManifestParser parses a path with manifests into resources.
type ManifestParser interface {
	Parse(string) ([]*unstructured.Unstructured, error)
}
