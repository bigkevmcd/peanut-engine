package engine

import (
	"log"
	"testing"

	// "github.com/bigkevmcd/peanut/pkg/gitfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func TestParseManifest(t *testing.T) {
	gt := clonedTree(t)
	c := GitConfig{RepoURL: "https://github.com/bigkevmcd/peanut-engine.git", Branch: "main", Path: "pkg/engine/testdata"}

	m, err := c.parseManifests(gt)
	assertNoError(t, err)

	log.Printf("KEVIN!!!! %#v\n", m)
	t.Fail()
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
