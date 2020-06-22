package metrics

import (
	"strings"
	"testing"

	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var _ Interface = (*PrometheusMetrics)(nil)

func TestRecordWithSynced(t *testing.T) {
	result := []common.ResourceSyncResult{
		{Status: common.ResultCodeSynced},
		// {Status: common.ResultCodeSyncFailed},
		// {Status: common.ResultCodePruned},
		// {Status: common.ResultCodePruneSkipped},
	}
	m := New("testing", prometheus.NewRegistry())
	m.Record(result)

	assertMetricGauged(t, m, result, m.synced, `
# HELP testing_synced Number of resources synced
# TYPE testing_synced gauge
testing_synced 1
`)
}

func TestRecordWithSyncFailed(t *testing.T) {
	m := New("testing", prometheus.NewRegistry())
	result := []common.ResourceSyncResult{
		{Status: common.ResultCodeSyncFailed},
	}

	assertMetricGauged(t, m, result, m.syncFailed, `
# HELP testing_sync_failed Number of resources that failed to sync
# TYPE testing_sync_failed gauge
testing_sync_failed 1
`)
}

func TestRecordWithPruned(t *testing.T) {
	m := New("testing", prometheus.NewRegistry())
	result := []common.ResourceSyncResult{
		{Status: common.ResultCodePruned},
	}

	assertMetricGauged(t, m, result, m.pruned, `
# HELP testing_pruned Number of resources pruned
# TYPE testing_pruned gauge
testing_pruned 1
`)
}

func TestRecordWithPruneSkipped(t *testing.T) {
	m := New("testing", prometheus.NewRegistry())
	result := []common.ResourceSyncResult{
		{Status: common.ResultCodePruneSkipped},
	}

	assertMetricGauged(t, m, result, m.pruneSkipped, `
# HELP testing_prune_skipped Number of resources that the pruning skipped
# TYPE testing_prune_skipped gauge
testing_prune_skipped 1
`)
}

func TestCountError(t *testing.T) {
	m := New("testing", prometheus.NewRegistry())

	m.CountError()

	err := testutil.CollectAndCompare(m.errors, strings.NewReader(`
# HELP testing_errors Count of errors during synchronisation
# TYPE testing_errors counter
testing_errors 1
`))
	if err != nil {
		t.Fatal(err)
	}
}

func assertMetricGauged(t *testing.T, m *PrometheusMetrics, r []common.ResourceSyncResult, g prometheus.Gauge, output string) {
	m.Record(r)
	err := testutil.CollectAndCompare(g, strings.NewReader(output))
	if err != nil {
		t.Fatal(err)
	}
}
