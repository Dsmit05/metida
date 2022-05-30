package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ServiceMetrics struct {
	httpRequestDurations   *prometheus.HistogramVec
	httpRequestCounters    *prometheus.CounterVec
	dbRequestErrorCounters prometheus.Counter
}

func NewServiceMetrics() *ServiceMetrics {
	httpRequestCounters := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_counters",
		Help: "The total number of request",
	}, []string{"method", "path", "code"})

	httpRequestDurations := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_durations",
		Help: "The total time durations of requests",
	}, []string{"method", "path", "code"})

	dbRequestErrorCounters := promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_error_counters",
		Help: "The total number of errors events",
	})

	return &ServiceMetrics{
		httpRequestDurations,
		httpRequestCounters,
		dbRequestErrorCounters}
}

// IncDbError увеличивает количество ошибок у бд
func (o *ServiceMetrics) IncDbError() {
	o.dbRequestErrorCounters.Inc()
}
