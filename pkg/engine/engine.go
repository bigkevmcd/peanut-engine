package engine

import (
	"context"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/argoproj/gitops-engine/pkg/utils/errors"
	"github.com/argoproj/gitops-engine/pkg/utils/io"
)

const (
	annotationGCMark = "gitops-agent.argoproj.io/gc-mark"
)

type resourceInfo struct {
	gcMark string
}

// StartPeanutSync starts watching the configured Git repository, and
// synchronising the resources.
func StartPeanutSync(clientConfig *rest.Config, peanutConfig PeanutConfig, resync chan bool, done <-chan struct{}) error {
	repo, err := cloneRepository(peanutConfig.Git)
	if err != nil {
		return err
	}
	currentSHA, err := headHash(repo)
	if err != nil {
		return err
	}
	log.Printf("Starting synchronisation from commit: %s", currentSHA)

	namespaces := []string{}
	if peanutConfig.Namespaced {
		namespaces = []string{peanutConfig.Namespace}
	}
	clusterCache := createClusterCache(namespaces, clientConfig)
	gitOpsEngine := engine.NewEngine(clientConfig, clusterCache)
	closer, err := gitOpsEngine.Run()
	errors.CheckErrorWithCode(err, errors.ErrorCommandSpecific)
	defer io.Close(closer)

	go func() {
		ticker := time.NewTicker(peanutConfig.Resync)
		for {
			<-ticker.C
			resync <- true
		}
	}()
	isManaged := isManagedChecker(peanutConfig.Git)
	for {
		select {
		case <-resync:
			log.Printf("Starting Synchronisation")
			newSHA, err := fetchRepository(repo)
			if err != nil {
				log.Errorf("failed to fetch updates to the respository: %s", err)
				continue
			}
			if newSHA != currentSHA {
				if newSHA != plumbing.ZeroHash {
					log.Printf("New commit detected: previous SHA %s, new SHA %s\n", currentSHA, newSHA)
					currentSHA = newSHA
				}
			}
			workTree, err := treeForHash(repo, currentSHA)
			if err != nil {
				log.Errorf("failed to calculate the tree for the hash: %s", err)
				continue
			}
			targets, err := peanutConfig.Git.parseManifests(workTree)
			if err != nil {
				log.Errorf("Failed to synchronize cluster state: %s", err)
			}
			result, err := gitOpsEngine.Sync(
				context.Background(), targets, isManaged,
				currentSHA.String(), peanutConfig.Namespace,
				sync.WithPrune(peanutConfig.Prune))
			if err != nil {
				log.Printf("Failed to synchronize cluster state: %v", err)
				continue
			}
			log.Printf("Applied %#v\n", summariseResults(result))
		case <-done:
			log.Println("Terminating synchronisation")
			return nil
		}
	}
}

func isManagedChecker(gc GitConfig) func(r *cache.Resource) bool {
	return func(r *cache.Resource) bool {
		return r.Info.(*resourceInfo).gcMark == gc.getGCMark(r.ResourceKey())
	}
}

func infoHandler(un *unstructured.Unstructured, isRoot bool) (interface{}, bool) {
	// store gc mark of every resource
	gcMark := un.GetAnnotations()[annotationGCMark]
	info := &resourceInfo{gcMark: un.GetAnnotations()[annotationGCMark]}
	// cache resources that has that mark to improve performance
	cacheManifest := gcMark != ""
	return info, cacheManifest
}

func createClusterCache(namespaces []string, clientConfig *rest.Config) cache.ClusterCache {
	return cache.NewClusterCache(
		clientConfig,
		cache.SetNamespaces(namespaces),
		cache.SetPopulateResourceInfoHandler(infoHandler),
	)
}

// TODO: return better errors here.
func cloneRepository(o GitConfig) (*git.Repository, error) {
	clone, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:   o.RepoURL,
		Depth: 1,
	})
	if err != nil {
		return nil, err
	}
	return clone, err
}

func fetchRepository(r *git.Repository) (plumbing.Hash, error) {
	err := r.Fetch(&git.FetchOptions{RemoteName: "origin"})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Println("No changes detected")
			return plumbing.ZeroHash, nil
		}
		return plumbing.ZeroHash, err
	}
	log.Println("Changes detected")
	return headHash(r)
}

func headHash(r *git.Repository) (plumbing.Hash, error) {
	ref, err := r.Head()
	if err != nil {
		return plumbing.ZeroHash, err
	}
	return ref.Hash(), nil
}

func treeForHash(r *git.Repository, h plumbing.Hash) (*object.Tree, error) {
	commit, err := r.CommitObject(h)
	if err != nil {
		return nil, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	return tree, nil
}

type summary struct {
	Synced       int64
	SyncFailed   int64
	Pruned       int64
	PruneSkipped int64
}

func summariseResults(results []common.ResourceSyncResult) summary {
	s := summary{}
	for _, r := range results {
		switch r.Status {
		case common.ResultCodeSynced:
			s.Synced++
		case common.ResultCodeSyncFailed:
			s.SyncFailed++
		case common.ResultCodePruned:
			s.Pruned++
		case common.ResultCodePruneSkipped:
			s.PruneSkipped++
		}
	}
	return s
}
