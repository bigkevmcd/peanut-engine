package engine

import (
	"time"
)

// GitConfig is the configuration for the repo to extract resources.
type GitConfig struct {
	RepoURL string
	Branch  string
	Path    string
}

// PeanutConfig configures the engine synchronisation.
type PeanutConfig struct {
	Prune      bool
	Namespace  string
	Namespaced bool
	Resync     time.Duration
}
