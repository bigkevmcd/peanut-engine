package engine

import (
	"testing"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestParseManifest(t *testing.T) {
	m := parseManifests(t)

	if l := len(m); l != 3 {
		t.Fatalf("got %d resources, wanted 3", len(m))
	}
}

func TestParseManifestParsesResource(t *testing.T) {
	m := parseManifests(t)

	d := findByKind(m, "Deployment")
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

func TestParseManifestAddsAnnotation(t *testing.T) {
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut-engine.git", Branch: "main", Path: "pkg/engine/testdata"}
	r := testRepository(t, c)
	h, err := r.HeadHash()
	assertNoError(t, err)
	m, err := r.ParseManifests(h)
	assertNoError(t, err)

	d := m[0]
	want := map[string]string{
		annotationGCMark: c.getGCMark(kube.GetResourceKey(d)),
	}

	if diff := cmp.Diff(want, d.GetAnnotations()); diff != "" {
		t.Fatalf("parsed manifest:\n%s", diff)
	}
}

func testRepository(t *testing.T, c GitConfig) *PeanutRepository {
	t.Helper()
	r := NewRepository(c)
	err := r.Open("../..")
	assertNoError(t, err)
	return r
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
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

func parseManifests(t *testing.T) []*unstructured.Unstructured {
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut-engine.git", Branch: "main", Path: "pkg/engine/testdata"}
	r := testRepository(t, c)
	h, err := r.HeadHash()
	assertNoError(t, err)
	m, err := r.ParseManifests(h)
	assertNoError(t, err)
	return m
}
