package metrics

import "github.com/argoproj/gitops-engine/pkg/sync/common"

var _ Interface = (*MockMetrics)(nil)

// MockMetrics is a type that provides a simple counter for metrics for test
// purposes.
type MockMetrics struct {
	Synced       int64
	SyncFailed   int64
	Pruned       int64
	PruneSkipped int64
}

// NewMock creates and returns a MockMetrics.
func NewMock() *MockMetrics {
	return &MockMetrics{}
}

// CountFailedAPICall records failed outgoing API calls to upstream services.
func (p *MockMetrics) Record(r []common.ResourceSyncResult) {
	for _, r := range r {
		switch r.Status {
		case common.ResultCodeSynced:
			p.Synced++
		case common.ResultCodeSyncFailed:
			p.SyncFailed++
		case common.ResultCodePruned:
			p.Pruned++
		case common.ResultCodePruneSkipped:
			p.PruneSkipped++
		}
	}
}
