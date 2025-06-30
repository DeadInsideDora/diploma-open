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
			Name: "products_parser_requests_total",
			Help: "Total number of HTTP requests, labeled by domain and status code.",
		},
		[]string{"domain", "status"},
	)
	prometheus.MustRegister(counter)

	return &PrometheusMetricsService{counter: counter}
}

func (service *PrometheusMetricsService) Inc(domain string, statusCode int) {
	service.counter.WithLabelValues(domain, strconv.Itoa(statusCode)).Inc()
}
