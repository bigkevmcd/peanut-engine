package kustomize

import (
	"github.com/bigkevmcd/peanut/pkg/kustomize/parser"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/resource"
)

// New creates and returns a new KustomizeParser.
func New() *KustomizeParser {
	return &KustomizeParser{}
}

// KustomizeParser is an implementation of the ManifestParser that can parse
// from Kustomize definitions.
type KustomizeParser struct {
}

// Parse is an implementation of ManifestParser.
func (k *KustomizeParser) Parse(path string) ([]*unstructured.Unstructured, error) {
	res, err := parser.ParseTreeToResMap(path, fs.MakeRealFS())
	if err != nil {
		return nil, err
	}
	m := []*unstructured.Unstructured{}
	for _, v := range res {
		u := convert(v)
		m = append(m, u)
	}
	return m, nil
}

// convert converts a Kustomize resource into a generic Unstructured resource
// which the gitops engine Sync needs.
func convert(r *resource.Resource) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: r.Map(),
	}
}
