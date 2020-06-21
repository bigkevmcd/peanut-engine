package metrics

import "github.com/argoproj/gitops-engine/pkg/sync/common"

// Interface implementations provide metrics for the system.
type Interface interface {
	// Record records all the metrics summarised from an engine Sync.
	Record([]common.ResourceSyncResult)
}
