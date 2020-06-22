package engine

import (
	"testing"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestParseManifestAddsAnnotation(t *testing.T) {
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut-engine.git", Branch: "main", Path: "pkg/engine/testdata"}
	r := testRepository(t, c)
	m, err := r.ParseManifests()
	assertNoError(t, err)

	d := m[0]
	want := map[string]string{
		annotationGCMark: r.GCMark(kube.GetResourceKey(d)),
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
