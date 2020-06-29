package engine

import (
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GitConfig is the configuration for the repo to extract resources.
type GitConfig struct {
	RepoURL   string
	Branch    string
	Path      string
	AuthToken string
}

// PeanutConfig configures the engine synchronisation.
type PeanutConfig struct {
	Prune      bool
	Namespace  string
	Namespaced bool
	Resync     time.Duration
}

func (c *GitConfig) BasicAuth() *http.BasicAuth {
	if c.AuthToken != "" {
		return &http.BasicAuth{
			Username: "peanut",
			Password: c.AuthToken,
		}
	}
	return nil
}
