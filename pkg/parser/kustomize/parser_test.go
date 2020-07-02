package kustomize

import (
	"strings"
	"testing"

	"github.com/bigkevmcd/peanut-engine/pkg/parser"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ parser.ManifestParser = (*KustomizeParser)(nil)

func TestKustomizationParse(t *testing.T) {
	k := &KustomizeParser{}

	res, err := k.Parse("../../testdata")
	if err != nil {
		t.Fatal(err)
	}

	if l := len(res); l != 3 {
		t.Fatalf("got %d, want 3", l)
	}
}

func TestKustomizationParseExtractsResources(t *testing.T) {
	k := &KustomizeParser{}

	res, err := k.Parse("../../testdata")
	if err != nil {
		t.Fatal(err)
	}

	d := findByKind(res, "Deployment")
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	}
	if diff := cmp.Diff(gvk, d.GroupVersionKind()); diff != "" {
		t.Errorf("parsed manifest:\n%s", diff)
	}
	if n := d.GetName(); n != "taxi" {
		t.Errorf("GetName() got %s, want %s", n, "taxi")
	}
	if n := d.GetNamespace(); n != "taxi-dev" {
		t.Errorf("GetNamespace() got %s, want %s", n, "taxi-dev")
	}
}

func TestKustomizationParseWithFailure(t *testing.T) {
	k := &KustomizeParser{}

	_, err := k.Parse("../..")
	if !strings.Contains(err.Error(), `unable to find one of 'kustomization.yaml', 'kustomization.yml' or 'Kustomization' in directory`) {
		t.Fatalf("incorrect error: %#v", err)
	}
}

func findByKind(r []*unstructured.Unstructured, k string) *unstructured.Unstructured {
	for _, v := range r {
		if v.GetKind() == k {
			return v
		}
	}
	return nil
}
