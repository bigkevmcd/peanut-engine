package engine

import (
	"github.com/bigkevmcd/peanut/pkg/kustomize/parser"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/pkg/fs"
)

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
