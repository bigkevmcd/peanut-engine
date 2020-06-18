package engine

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/bigkevmcd/peanut/pkg/gitfs"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/pkg/resource"

	"github.com/bigkevmcd/peanut/pkg/kustomize/parser"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitConfig is the configuration for the repo to extract resources.
type GitConfig struct {
	RepoURL string
	Branch  string
	Path    string
}

// PeanutConfig configures the engine synchronisation.
type PeanutConfig struct {
	Git        GitConfig
	Prune      bool
	Namespace  string
	Namespaced bool
	Resync     time.Duration
	ClonePath  string
}

func (g *GitConfig) getGCMark(key kube.ResourceKey) string {
	h := sha256.New()
	_, _ = h.Write([]byte(fmt.Sprintf("%s/%s", g.RepoURL, g.Path)))
	_, _ = h.Write([]byte(strings.Join([]string{key.Group, key.Kind, key.Name}, "/")))
	return "sha256." + base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func (s *GitConfig) parseManifests(tree *object.Tree) ([]*unstructured.Unstructured, error) {
	res, err := parser.ParseTreeToResMap(s.Path, gitfs.New(tree))
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
		annotations[annotationGCMark] = s.getGCMark(kube.GetResourceKey(u))
		v.SetAnnotations(annotations)
		m = append(m, u)
	}
	return m, nil
}

func convert(r *resource.Resource) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: r.Map(),
	}
}
