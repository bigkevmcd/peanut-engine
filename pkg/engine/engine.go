package engine

import (
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/utils/errors"
	"github.com/argoproj/gitops-engine/pkg/utils/io"
)

const (
	annotationGCMark = "gitops-agent.argoproj.io/gc-mark"
)

type resourceInfo struct {
	gcMark string
}

func PeanutSync(clientConfig *rest.Config, peanutConfig *PeanutConfig, resync chan bool) {
	var namespaces []string
	if peanutConfig.Namespaced {
		namespaces = []string{peanutConfig.Namespace}
	}
	clusterCache := cache.NewClusterCache(clientConfig,
		cache.SetNamespaces(namespaces),
		cache.SetPopulateResourceInfoHandler(func(un *unstructured.Unstructured, isRoot bool) (info interface{}, cacheManifest bool) {
			// store gc mark of every resource
			gcMark := un.GetAnnotations()[annotationGCMark]
			info = &resourceInfo{gcMark: un.GetAnnotations()[annotationGCMark]}
			// cache resources that has that mark to improve performance
			cacheManifest = gcMark != ""
			return
		}),
	)
	gitOpsEngine := engine.NewEngine(clientConfig, clusterCache)
	closer, err := gitOpsEngine.Run()
	errors.CheckErrorWithCode(err, errors.ErrorCommandSpecific)

	defer io.Close(closer)

	go func() {
		ticker := time.NewTicker(peanutConfig.Resync)
		for {
			<-ticker.C
			log.Infof("Synchronization triggered by timer")
			resync <- true
		}
	}()
	// for ; true; <-resync {
	// 	target, revision, err := s.parseManifests()
	// 	if err != nil {
	// 		log.Warnf("failed to parse target state: %v", err)
	// 		continue
	// 	}

	// 	result, err := gitOpsEngine.Sync(context.Background(), target, func(r *cache.Resource) bool {
	// 		return r.Info.(*resourceInfo).gcMark == s.getGCMark(r.ResourceKey())
	// 	}, revision, namespace, sync.WithPrune(prune))
	// 	if err != nil {
	// 		log.Warnf("failed to synchronize cluster state: %v", err)
	// 		continue
	// 	}
	// 	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// 	_, _ = fmt.Fprintf(w, "RESOURCE\tRESULT\n")
	// 	for _, res := range result {
	// 		_, _ = fmt.Fprintf(w, "%s\t%s\n", res.ResourceKey.String(), res.Message)
	// 	}
	// 	_ = w.Flush()
	// }
}
