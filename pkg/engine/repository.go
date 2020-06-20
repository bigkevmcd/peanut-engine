package engine

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type gitRepository interface {
	Clone() error
}

type PeanutRepository struct {
	config GitConfig
	repo   *git.Repository
}

func NewRepository(cfg GitConfig) *PeanutRepository {
	return &PeanutRepository{
		config: cfg,
	}
}

func (p *PeanutRepository) Clone(clonePath string) error {
	repo, err := cloneRepository(p.config, clonePath)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	p.repo = repo
	return nil
}

func (p *PeanutRepository) Open(openPath string) error {
	repo, err := git.PlainOpen(openPath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", openPath, err)
	}
	p.repo = repo
	return nil
}

func (p *PeanutRepository) HeadHash() (plumbing.Hash, error) {
	ref, err := p.repo.Head()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to get the Head for the Repository: %w", err)
	}
	return ref.Hash(), nil
}

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

func (p *PeanutRepository) Sync() (plumbing.Hash, error) {
	err := p.repo.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
		},
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
		RemoteName:    remoteName,
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
