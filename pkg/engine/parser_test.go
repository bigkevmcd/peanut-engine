package engine

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestKustomizationParse(t *testing.T) {
	k := &KustomizeParser{}

	res, err := k.Parse("testdata")
	if err != nil {
		t.Fatal(err)
	}

	if l := len(res); l != 3 {
		t.Fatalf("got %d, want 3", l)
	}
}

func TestKustomizationParseExtractsResources(t *testing.T) {
	k := &KustomizeParser{}

	res, err := k.Parse("testdata")
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
