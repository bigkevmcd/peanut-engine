package engine

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/resource"

	"github.com/bigkevmcd/peanut/pkg/kustomize/parser"
)

var defaultRefSpecs = []config.RefSpec{
	config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
}

const (
	defaultRemoteName = "origin"
)

// PeanutRepository wraps git.Repository with some high-level functionality.
type PeanutRepository struct {
	config     GitConfig
	repo       *git.Repository
	remoteName string
	repoPath   string
}

// NewRepository creates and returns a new PeanutRepository.
func NewRepository(cfg GitConfig) *PeanutRepository {
	return &PeanutRepository{
		config:     cfg,
		remoteName: defaultRemoteName,
	}
}

// Clone clones the configured repository to the provided path.
func (p *PeanutRepository) Clone(repoPath string) error {
	p.repoPath = repoPath
	clone, err := git.PlainClone(p.repoPath, false, &git.CloneOptions{
		RemoteName:    p.remoteName,
		URL:           p.config.RepoURL,
		ReferenceName: plumbing.NewBranchReferenceName(p.config.Branch),
	})
	if err != nil {
		return fmt.Errorf("failed to clone %s to %s: %w", p.config.RepoURL, p.repoPath, err)
	}
	p.repo = clone
	return nil
}

// Open assumes that the provided path contains a valid Git clone with the
// correct branch.
func (p *PeanutRepository) Open(openPath string) error {
	p.repoPath = openPath
	repo, err := git.PlainOpen(p.repoPath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", p.repoPath, err)
	}
	p.repo = repo
	return nil
}

// HeadHash returns the hash of the head commit of the repository.
func (p *PeanutRepository) HeadHash() (plumbing.Hash, error) {
	ref, err := p.repo.Head()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to get the Head for the Repository: %w", err)
	}
	return ref.Hash(), nil
}

// TreeForHash returns a git tree that can be used to access files.
func (p *PeanutRepository) TreeForHash(h plumbing.Hash) (*object.Tree, error) {
	commit, err := p.repo.CommitObject(h)
	if err != nil {
		return nil, fmt.Errorf("failed to get the CommitObject for %s: %w", h, err)
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get a Tree for the commit %s: %w", h, err)
	}
	return tree, nil
}

// Sync does a Fetch and Pull, and returns the HeadHash.
func (p *PeanutRepository) Sync() (plumbing.Hash, error) {
	err := p.repo.Fetch(&git.FetchOptions{
		RemoteName: p.remoteName,
		RefSpecs:   defaultRefSpecs,
	})
	if err != nil {
		if !upToDate(err) {
			return plumbing.ZeroHash, fmt.Errorf("failed to fetch from the Repository: %w", err)
		}
		return plumbing.ZeroHash, err
	}
	wtree, err := p.repo.Worktree()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to get a Worktree from the Repository: %w", err)
	}
	err = wtree.Pull(&git.PullOptions{
		RemoteName:    p.remoteName,
		ReferenceName: plumbing.NewBranchReferenceName(p.config.Branch),
	})

	if err != nil {
		if !upToDate(err) {
			return plumbing.ZeroHash, fmt.Errorf("failed to pull the Worktree: %w", err)
		}
		return plumbing.ZeroHash, err
	}
	return p.HeadHash()
}

// ParseManifests parses this repository's path, and returns the kustomized
// resources.
func (p *PeanutRepository) ParseManifests(h plumbing.Hash) ([]*unstructured.Unstructured, error) {
	res, err := parser.ParseTreeToResMap(filepath.Join(p.repoPath, p.config.Path), fs.MakeRealFS())
	if err != nil {
		return nil, err
	}
	m := []*unstructured.Unstructured{}
	for _, v := range res {
		annotations := v.GetAnnotations()
		if annotations == nil {
			annotations = map[string]string{}
		}
		u := convert(v)
		annotations[annotationGCMark] = p.GCMark(kube.GetResourceKey(u))
		v.SetAnnotations(annotations)
		m = append(m, u)
	}
	return m, nil
}

func (p *PeanutRepository) IsManaged(r *cache.Resource) bool {
	return r.Info.(*resourceInfo).gcMark == p.GCMark(r.ResourceKey())
}

func (p *PeanutRepository) GCMark(key kube.ResourceKey) string {
	h := sha256.New()
	_, _ = h.Write([]byte(fmt.Sprintf("%s/%s", p.config.RepoURL, p.config.Path)))
	_, _ = h.Write([]byte(strings.Join([]string{key.Group, key.Kind, key.Name}, "/")))
	return "sha256." + base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func convert(r *resource.Resource) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: r.Map(),
	}
}