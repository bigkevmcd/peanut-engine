package engine

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
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

func (g *GitConfig) getGCMark(key kube.ResourceKey) string {
	h := sha256.New()
	_, _ = h.Write([]byte(fmt.Sprintf("%s/%s", g.RepoURL, g.Path)))
	_, _ = h.Write([]byte(strings.Join([]string{key.Group, key.Kind, key.Name}, "/")))
	return "sha256." + base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
