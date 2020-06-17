package engine

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/bigkevmcd/peanut/pkg/gitfs"
	"github.com/bigkevmcd/peanut/pkg/kustomize/parser"
	"github.com/go-git/go-git/v5/plumbing/object"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
	log.Printf("KEVIN!!!!! %#v\n", res)

	// var res []*unstructured.Unstructured
	// for i := range s.paths {
	// 	if err := filepath.Walk(filepath.Join(s.repoPath, s.paths[i]), func(path string, info os.FileInfo, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if info.IsDir() {
	// 			return nil
	// 		}
	// 		if ext := filepath.Ext(info.Name()); ext != ".json" && ext != ".yml" && ext != ".yaml" {
	// 			return nil
	// 		}
	// 		data, err := ioutil.ReadFile(path)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		items, err := kube.SplitYAML(string(data))
	// 		if err != nil {
	// 			return fmt.Errorf("failed to parse %s: %v", path, err)
	// 		}
	// 		res = append(res, items...)
	// 		return nil
	// 	}); err != nil {
	// 		return nil, "", err
	// 	}
	// }
	// for i := range res {
	// 	annotations := res[i].GetAnnotations()
	// 	if annotations == nil {
	// 		annotations = make(map[string]string)
	// 	}
	// 	annotations[annotationGCMark] = s.getGCMark(kube.GetResourceKey(res[i]))
	// 	res[i].SetAnnotations(annotations)
	// }
	// return res, revision, nil
	return nil, nil
}
