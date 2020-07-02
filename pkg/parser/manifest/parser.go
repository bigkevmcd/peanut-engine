package manifest

import (
	"github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// New creates and returns a new ManifestivalParser.
func New() *ManifestivalParser {
	return &ManifestivalParser{}
}

// ManifestivalParser is an implementation of the ManifestParser that can parse
// from Kustomize definitions.
type ManifestivalParser struct {
}

// Parse is an implementation of ManifestParser.
func (k *ManifestivalParser) Parse(path string) ([]*unstructured.Unstructured, error) {
	m, err := manifestival.ManifestFrom(manifestival.Path(path))
	if err != nil {
		return nil, err
	}
	parsed := m.Resources()
	res := make([]*unstructured.Unstructured, len(parsed))
	for i := range parsed {
		res[i] = &parsed[i]
	}
	return res, nil
}
