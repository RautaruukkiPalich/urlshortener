package apiserver

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestMetrics = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace: "urlsh",
		Subsystem: "http",
		Name: "request",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{
		"status",
	},
)

func observeRequest(d time.Duration, status int) {
	requestMetrics.WithLabelValues(
		strconv.Itoa(status),
	).Observe(float64(d.Milliseconds()))
}