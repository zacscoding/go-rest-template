package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/zacscoding/go-rest-template/internal/config"
)

//go:generate mockery --name Provider --filename provider.go
type Provider interface {
	// RecordApiCount increases count of api request with given code, method, path labels
	RecordApiCount(code int, method, path string)

	// RecordApiLatency observes given elapsed mills with given code, method, path labels
	RecordApiLatency(code int, method, path string, elapsed time.Duration)

	// RecordCache increases count of cache request with given key, hit
	RecordCache(key string, hit bool)
}

type provider struct {
	namespace string
	subsystem string

	apiMetricsProvider   apiMetricsProvider
	cacheMetricsProvider cacheMetricsProvider
}

type apiMetricsProvider struct {
	requestCounter *prometheus.CounterVec
	requestLatency *prometheus.SummaryVec
}

type cacheMetricsProvider struct {
	cacheTotalCounter *prometheus.CounterVec
	cacheHitCounter   *prometheus.CounterVec
}

// NewProvider returns a new Provider with given conf config.Config.
func NewProvider(conf *config.Config) Provider {
	var (
		ns = conf.Metric.Namespace
		ss = conf.Metric.Subsystem
	)

	p := provider{
		namespace: ns,
		subsystem: ss,
		apiMetricsProvider: apiMetricsProvider{
			requestCounter: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: ns,
					Subsystem: ss,
					Name:      "api_request_count",
					Help:      "Total count of request",
				},
				[]string{"code", "method", "path"},
			),
			requestLatency: promauto.NewSummaryVec(
				prometheus.SummaryOpts{
					Namespace: ns,
					Subsystem: ss,
					Name:      "api_request_latency",
					Help:      "Elapsed time of request",
				},
				[]string{"code", "method", "path"},
			),
		},
		cacheMetricsProvider: cacheMetricsProvider{
			cacheTotalCounter: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: ns,
					Subsystem: ss,
					Name:      "cache_total",
					Help:      "Total count of cache requests",
				},
				[]string{"key"},
			),
			cacheHitCounter: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: ns,
					Subsystem: ss,
					Name:      "cache_hit",
					Help:      "Total cache hit count",
				},
				[]string{"key"},
			),
		},
	}
	return &p
}

func (p *provider) RecordApiCount(code int, method, path string) {
	p.apiMetricsProvider.requestCounter.WithLabelValues(strconv.Itoa(code), method, path).Inc()
}

func (p *provider) RecordApiLatency(code int, method, path string, elapsed time.Duration) {
	mills := float64(elapsed.Milliseconds())
	p.apiMetricsProvider.requestLatency.WithLabelValues(strconv.Itoa(code), method, path).Observe(mills)
}

func (p *provider) RecordCache(key string, hit bool) {
	p.cacheMetricsProvider.cacheTotalCounter.WithLabelValues(key).Inc()
	if hit {
		p.cacheMetricsProvider.cacheHitCounter.WithLabelValues(key).Inc()
	}
}
