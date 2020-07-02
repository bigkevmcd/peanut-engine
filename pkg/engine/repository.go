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

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/bigkevmcd/peanut-engine/pkg/parser"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
	parser     parser.ManifestParser
}

// NewRepository creates and returns a new PeanutRepository.
func NewRepository(cfg GitConfig, p parser.ManifestParser) *PeanutRepository {
	return &PeanutRepository{
		config:     cfg,
		remoteName: defaultRemoteName,
		parser:     p,
	}
}

// TODO: should Open and Clone be package level functions instead of methods?

// Clone clones the configured repository to the provided path.
func (p *PeanutRepository) Clone(repoPath string) error {
	p.repoPath = repoPath

	opts := &git.CloneOptions{
		Auth:          p.config.BasicAuth(),
		RemoteName:    p.remoteName,
		URL:           p.config.RepoURL,
		ReferenceName: plumbing.NewBranchReferenceName(p.config.Branch),
	}
	clone, err := git.PlainClone(p.repoPath, false, opts)
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

// Sync does a Fetch and Pull, and returns the HeadHash.
func (p *PeanutRepository) Sync() (plumbing.Hash, error) {
	err := p.repo.Fetch(&git.FetchOptions{
		Auth:       p.config.BasicAuth(),
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

// ParseManifests parses this repository's path, and returns the parsed
// resources.
// TODO: should this take a path? Is there
func (p *PeanutRepository) ParseManifests() ([]*unstructured.Unstructured, error) {
	res, err := p.parser.Parse(filepath.Join(p.repoPath, p.config.Path))
	if err != nil {
		return nil, err
	}
	for _, v := range res {
		annotations := v.GetAnnotations()
		if annotations == nil {
			annotations = map[string]string{}
		}
		gcm, err := p.GCMark(kube.GetResourceKey(v))
		if err != nil {
			return nil, err
		}
		annotations[annotationGCMark] = gcm
		v.SetAnnotations(annotations)
	}
	return res, nil
}

// IsManaged is used by the cache to determine whether or not a resource is
// a managed resource.
// TODO: is this appropriate for the Repository?
func (p *PeanutRepository) IsManaged(r *cache.Resource) bool {
	gcm, err := p.GCMark(r.ResourceKey())
	if err != nil {
		panic(err)
	}
	return r.Info.(*resourceInfo).gcMark == gcm
}

// GCMark calculates a signature for the resource from the repo URL and path
// along with the GVK.
func (p *PeanutRepository) GCMark(key kube.ResourceKey) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(fmt.Sprintf("%s/%s", p.config.RepoURL, p.config.Path)))
	if err != nil {
		return "", err
	}
	_, err = h.Write([]byte(strings.Join([]string{key.Group, key.Kind, key.Name}, "/")))
	if err != nil {
		return "", err
	}
	return "sha256." + base64.RawURLEncoding.EncodeToString(h.Sum(nil)), nil
}

func upToDate(err error) bool {
	return err == git.NoErrAlreadyUpToDate
}
