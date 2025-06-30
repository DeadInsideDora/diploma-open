package handlers

import (
	"net/http"
	"web/internal/domain"
)

func WriteHttpError(path string, metricsService domain.IMetricsService, w http.ResponseWriter, message string, status int) {
	metricsService.Inc(path, status)
	http.Error(w, message, status)
}
