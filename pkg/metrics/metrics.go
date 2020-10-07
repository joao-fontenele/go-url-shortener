package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// These are prometheus metrics
var (
	HTTPRequestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Counter for total requests received",
		},
		[]string{"route", "status"},
	)

	HTTPRequestsDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"route", "status"},
	)

	DAOFindResultCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dao_find_result_total",
			Help: "Total cache hit/miss",
		},
		[]string{"name", "result"},
	)

	DAOOperationsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dao_operations_total",
			Help: "Total DAO operations by result (error/success)",
		},
		[]string{"name", "result"},
	)

	DAOOperationsDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dao_operation_duration_seconds",
			Help:    "Duration of DAO operations in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"name", "operation"},
	)
)

// Init register metrics to prometheus register
func Init() {
	prometheus.MustRegister(
		HTTPRequestsCounter,
		HTTPRequestsDurationHistogram,
		DAOFindResultCounter,
		DAOOperationsCounter,
		DAOOperationsDurationHistogram,
	)
}
