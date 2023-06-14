package kustomize

import (
	"fmt"

	"github.com/bigkevmcd/peanut/pkg/kustomize/parser"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/kyaml/filesys"
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
	resMap, err := parser.ParseTreeToResMap(path, filesys.MakeFsOnDisk())
	if err != nil {
		return nil, err
	}
	m := []*unstructured.Unstructured{}
	for _, v := range resMap.AllIds() {
		r, err := resMap.GetById(v)
		if err != nil {
			return nil, fmt.Errorf("getting resource %s: %w", v, err)
		}
		data, err := r.Map()
		if err != nil {
			return nil, fmt.Errorf("getting data: %w", err)
		}
		u := convert(data)
		m = append(m, u)
	}
	return m, nil
}

// convert converts a Kustomize resource into a generic Unstructured resource
// which the gitops engine Sync needs.
func convert(data map[string]interface{}) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: data,
	}
}
