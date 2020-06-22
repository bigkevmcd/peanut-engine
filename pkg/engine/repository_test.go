package engine

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
	gcm, err := r.GCMark(kube.GetResourceKey(d))
	assertNoError(t, err)
	want := map[string]string{
		annotationGCMark: gcm,
	}

	if diff := cmp.Diff(want, d.GetAnnotations()); diff != "" {
		t.Fatalf("parsed manifest:\n%s", diff)
	}
}

func TestOpen(t *testing.T) {
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut-engine.git", Branch: "main", Path: "pkg/engine/testdata"}
	r := NewRepository(c)
	dir, cleanup := mkTempDir(t)
	t.Cleanup(cleanup)
	err := r.Open(dir)

	if !strings.Contains(err.Error(), `repository does not exist`) {
		t.Fatalf("incorrect error: %s", err)
	}
}

func TestClone(t *testing.T) {
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut.git", Branch: "main", Path: "pkg/engine/testdata"}
	dir, cleanup := mkTempDir(t)
	t.Cleanup(cleanup)
	r := NewRepository(c)

	err := r.Clone(dir)
	assertNoError(t, err)

	want := execGitHead(t, dir)
	got, err := r.HeadHash()

	if want != got.String() {
		t.Fatalf("incorrect git SHA from HeadHash, got %#v, want %#v", got.String(), want)
	}
}

func TestCloneWithMissingSource(t *testing.T) {
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/doesnotexist.git", Branch: "main", Path: "pkg/engine/testdata"}
	dir, cleanup := mkTempDir(t)
	t.Cleanup(cleanup)
	r := NewRepository(c)

	err := r.Clone(dir)

	// Unfortunately this is unable to determine the difference between a
	// non-existent Repo and one that requires authentication.
	if !strings.Contains(err.Error(), `authentication required`) {
		t.Fatalf("incorrect error: %s", err)
	}
}

func TestSync(t *testing.T) {
	t.Skip()
}

func TestHeadHash(t *testing.T) {
	want := execGitHead(t, ".")
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut-engine.git", Branch: "main", Path: "pkg/engine/testdata"}
	r := testRepository(t, c)
	got, err := r.HeadHash()
	assertNoError(t, err)

	if want != got.String() {
		t.Fatalf("incorrect git SHA from HeadHash, got %#v, want %#v", got.String(), want)
	}
}

func TestIsManaged(t *testing.T) {
	t.Skip()
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

func mkTempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", "peanut")
	if err != nil {
		t.Fatal(err)
	}
	return dir, func() {
		assertNoError(t, os.RemoveAll(dir))
	}
}

func execGitHead(t *testing.T, dir string) string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	assertNoError(t, err)
	return strings.TrimSpace(string(out))
}
