package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"web/internal/domain"
)

type requestProductsNames struct {
	Type  string   `json:"type"`
	Names []string `json:"names"`
}

func RegisterProductsByNamesHandler(mux *http.ServeMux, path string, metricsService domain.IMetricsService, categoriesService domain.ICategoriesService, readerService domain.IReaderService) {
	categories, err := categoriesService.Get(false)
	if err != nil {
		panic(err)
	}

	mux.HandleFunc(
		fmt.Sprintf("POST %s", path),
		func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			var request requestProductsNames
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				WriteHttpError(path, metricsService, w, fmt.Sprintf("invalid JSON: %s", err), http.StatusBadRequest)
				return
			}

			if !checkProductType(categories, request.Type) {
				WriteHttpError(path, metricsService, w, "invalid product type", http.StatusBadRequest)
				return
			}

			products, err := readerService.ReadByNames(request.Type, request.Names)

			if err != nil {
				WriteHttpError(path, metricsService, w, fmt.Sprintf("can't get products from redis: %s", err), http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(w).Encode(products); err != nil {
				WriteHttpError(path, metricsService, w, fmt.Sprintf("encode error: %s", err), http.StatusInternalServerError)
			} else {
				metricsService.Inc(path, http.StatusOK)
			}
		},
	)
}
