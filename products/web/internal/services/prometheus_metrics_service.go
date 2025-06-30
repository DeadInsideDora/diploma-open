package services

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetricsService struct {
	counter *prometheus.CounterVec
}

func NewPrometheusMetricsService() *PrometheusMetricsService {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "products_web_requests_total",
			Help: "Total number of HTTP requests, labeled by path and status code.",
		},
		[]string{"path", "status"},
	)
	prometheus.MustRegister(counter)

	return &PrometheusMetricsService{counter: counter}
}

func (service *PrometheusMetricsService) Inc(path string, statusCode int) {
	service.counter.WithLabelValues(path, strconv.Itoa(statusCode)).Inc()
}
