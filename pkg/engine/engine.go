package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/argoproj/gitops-engine/pkg/utils/errors"
	"github.com/argoproj/gitops-engine/pkg/utils/io"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"

	"github.com/bigkevmcd/peanut-engine/pkg/metrics"
	"github.com/bigkevmcd/peanut-engine/pkg/recent"
)

const (
	annotationGCMark = "gitops-agent.argoproj.io/gc-mark"
)

type resourceInfo struct {
	gcMark string
}

// StartPeanutSync starts watching the configured Git repository, and
// synchronising the resources.
func StartPeanutSync(clientConfig *rest.Config, config PeanutConfig, peanutRepo GitRepository, met metrics.Interface, syncs *recent.RecentSynchronisations, resync chan bool, done <-chan struct{}) error {
	currentSHA, err := peanutRepo.HeadHash()
	if err != nil {
		return fmt.Errorf("failed to get the head hash: %w", err)
	}
	log.Infof("Starting synchronisation from commit: %s", currentSHA)

	namespaces := []string{}
	if config.Namespaced {
		namespaces = []string{config.Namespace}
	}

	clusterCache := createClusterCache(namespaces, clientConfig)
	gitOpsEngine := engine.NewEngine(clientConfig, clusterCache)
	closer, err := gitOpsEngine.Run()
	errors.CheckErrorWithCode(err, errors.ErrorCommandSpecific)
	defer io.Close(closer)

	go func() {
		ticker := time.NewTicker(config.Resync)
		for {
			<-ticker.C
			resync <- true
		}
	}()

	for {
		select {
		case <-resync:
			log.Infof("Starting Synchronisation from %s", currentSHA)
			start := time.Now()
			newSHA, err := peanutRepo.Sync()
			if err != nil && err != git.NoErrAlreadyUpToDate {
				met.CountError()
				log.Errorf("Failed to fetch updates to the repository: %s", err)
				continue
			}
			if newSHA != currentSHA {
				if newSHA != plumbing.ZeroHash {
					log.Infof("New commit detected: previous SHA %s, new SHA %s", currentSHA, newSHA)
					currentSHA = newSHA
				}
			}
			targets, err := peanutRepo.ParseManifests()
			if err != nil {
				met.CountError()
				log.Errorf("Failed to synchronize cluster state: %s", err)
			}

			result, err := gitOpsEngine.Sync(
				context.Background(), targets, peanutRepo.IsManaged,
				currentSHA.String(), config.Namespace,
				sync.WithPrune(config.Prune))

			syncs.Add(start, time.Now(), currentSHA, err, result)

			if err != nil {
				met.CountError()
				log.Infof("Failed to synchronize cluster state: %v", err)
				continue
			}
			met.Record(result)
		case <-done:
			log.Println("Terminating synchronisation")
			return nil
		}
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
