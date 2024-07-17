package metrics

import (
	compbasemetrics "k8s.io/component-base/metrics"
)

var (
	cceNamespace = "cce_apiserver"
	subsystem    = "watch_cacher"

	// CCEWatchCacheReadWait watch latency between apiserver and etcd
	CCEWatchCacheReadWait = compbasemetrics.NewHistogramVec(
		&compbasemetrics.HistogramOpts{
			Namespace:      cceNamespace,
			Subsystem:      subsystem,
			Name:           "watch_wait_seconds",
			Help:           "Histogram of time spent waiting for a watch cache to become fresh.this is cce metrics",
			StabilityLevel: compbasemetrics.ALPHA,
			Buckets:        []float64{0.005, 0.025, 0.05, 0.1, 0.2, 0.4, 0.6, 0.8, 1.0, 1.25, 1.5, 2, 3, 5, 10, 30, 60, 180, 300},
		}, []string{"resource"})

	SourceCacher                 = "cacher"
	SourceStorage                = "storage"
	CCEWatchCacheResourceVersion = compbasemetrics.NewGaugeVec(
		&compbasemetrics.GaugeOpts{
			Namespace:      cceNamespace,
			Subsystem:      subsystem,
			Name:           "resource_version",
			Help:           "Differ from storage and apiserver resource version in watch cacher by resource type.",
			StabilityLevel: compbasemetrics.ALPHA,
		},
		[]string{"resource", "source"},
	)
)
