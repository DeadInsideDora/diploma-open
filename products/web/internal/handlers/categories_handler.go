package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"web/internal/domain"
)

func RegisterCategoriesHandler(mux *http.ServeMux, path string, metricsService domain.IMetricsService, categories domain.ICategoriesService) {
	mux.HandleFunc(
		fmt.Sprintf("GET %s", path),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			categories, err := categories.Get(false)
			if err != nil {
				WriteHttpError(path, metricsService, w, fmt.Sprintf("can't get categories info: %v", err), http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(w).Encode(categories); err != nil {
				WriteHttpError(path, metricsService, w, fmt.Sprintf("encode error: %s", err), http.StatusInternalServerError)
			} else {
				metricsService.Inc(path, http.StatusOK)
			}
		},
	)
}
