package engine

import (
	"testing"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestParseManifest(t *testing.T) {
	gt := clonedTree(t)
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut-engine.git", Branch: "main", Path: "pkg/engine/testdata"}

	m, err := c.parseManifests(gt)
	assertNoError(t, err)

	if l := len(m); l != 2 {
		t.Fatalf("got %d resources, wanted 2", len(m))
	}
}

func TestParseManifestParsesResource(t *testing.T) {
	gt := clonedTree(t)
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut-engine.git", Branch: "main", Path: "pkg/engine/testdata"}

	m, err := c.parseManifests(gt)
	assertNoError(t, err)

	r := findByKind(m, "Deployment")
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	}
	if diff := cmp.Diff(gvk, r.GroupVersionKind()); diff != "" {
		t.Errorf("parsed manifest:\n%s", diff)
	}
	if n := r.GetName(); n != "taxi" {
		t.Errorf("GetName() got %s, want %s", n, "taxi")
	}
	if n := r.GetNamespace(); n != "taxi-dev" {
		t.Errorf("GetNamespace() got %s, want %s", n, "taxi-dev")
	}
}

func TestParseManifestAddsAnnotation(t *testing.T) {
	gt := clonedTree(t)
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut-engine.git", Branch: "main", Path: "pkg/engine/testdata"}

	m, err := c.parseManifests(gt)
	assertNoError(t, err)

	r := m[0]
	want := map[string]string{
		annotationGCMark: c.getGCMark(kube.GetResourceKey(r)),
	}

	if diff := cmp.Diff(want, r.GetAnnotations()); diff != "" {
		t.Fatalf("parsed manifest:\n%s", diff)
	}
}

func clonedTree(t *testing.T) *object.Tree {
	t.Helper()
	clone, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	ref, err := clone.Head()
	if err != nil {
		t.Fatal(err)
	}
	commit, err := clone.CommitObject(ref.Hash())
	if err != nil {
		t.Fatal(err)
	}
	tree, err := commit.Tree()
	if err != nil {
		t.Fatal(err)
	}
	assertNoError(t, err)
	return tree
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
