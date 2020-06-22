package metrics

import (
	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetrics is a wrapper around Prometheus metrics for counting
// events in the system.
type PrometheusMetrics struct {
	synced       prometheus.Gauge
	syncFailed   prometheus.Gauge
	pruned       prometheus.Gauge
	pruneSkipped prometheus.Gauge
	errors       prometheus.Counter
}

// New creates and returns a PrometheusMetrics initialised with prometheus
// gauges.
func New(ns string, reg prometheus.Registerer) *PrometheusMetrics {
	pm := &PrometheusMetrics{}
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}

	pm.synced = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: ns,
		Name:      "synced",
		Help:      "Number of resources synced",
	})

	pm.syncFailed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: ns,
		Name:      "sync_failed",
		Help:      "Number of resources that failed to sync",
	})

	pm.pruned = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: ns,
		Name:      "pruned",
		Help:      "Number of resources pruned",
	})

	pm.pruneSkipped = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: ns,
		Name:      "prune_skipped",
		Help:      "Number of resources that the pruning skipped",
	})

	pm.errors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: ns,
		Name:      "errors",
		Help:      "Count of errors during synchronisation",
	})

	reg.MustRegister(pm.synced)
	reg.MustRegister(pm.syncFailed)
	reg.MustRegister(pm.pruned)
	reg.MustRegister(pm.pruneSkipped)
	reg.MustRegister(pm.errors)
	return pm
}

// Record is an implementation of the metrics Interface.
func (p *PrometheusMetrics) Record(r []common.ResourceSyncResult) {
	var synced, syncFailed, pruned, pruneSkipped float64
	for _, r := range r {
		switch r.Status {
		case common.ResultCodeSynced:
			synced++
		case common.ResultCodeSyncFailed:
			syncFailed++
		case common.ResultCodePruned:
			pruned++
		case common.ResultCodePruneSkipped:
			pruneSkipped++
		}
	}
	p.synced.Set(synced)
	p.syncFailed.Set(syncFailed)
	p.pruned.Set(pruned)
	p.pruneSkipped.Set(pruneSkipped)
}

// CountError counts the number of errors during synchronisation.
func (m *PrometheusMetrics) CountError() {
	m.errors.Inc()
}
