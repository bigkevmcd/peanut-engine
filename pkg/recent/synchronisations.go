package recent

import (
	"time"

	"container/ring"

	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/go-git/go-git/v5/plumbing"
)

// NewRecentSynchronisations creates and returns a ring buffer of
// synchronisations.
func NewRecentSynchronisations(r *ring.Ring) *RecentSynchronisations {
	return &RecentSynchronisations{recent: r}
}

// RecentSynchronisations represents a ring buffer of recent sync states.
type RecentSynchronisations struct {
	recent *ring.Ring
}

// Add records the details of a synchronisation in the ring.
func (r *RecentSynchronisations) Add(start, end time.Time, sha plumbing.Hash, syncErr error, results []common.ResourceSyncResult) {
	r.recent.Value = Synchronisation{Start: start, End: end, SHA: sha.String(), Error: syncErr, Results: results}
	r.recent = r.recent.Next()
}

// Latest returns the last recorded synchronisation.
func (r *RecentSynchronisations) Latest() Synchronisation {
	return r.recent.Prev().Value.(Synchronisation)
}

// Synchronisation represents a sync run from the gitops engine.
type Synchronisation struct {
	Start   time.Time                   `json:"startTime"`
	End     time.Time                   `json:"endTime"`
	SHA     string                      `json:"sha"`
	Error   error                       `json:"err"`
	Results []common.ResourceSyncResult `json:"results"`
}
