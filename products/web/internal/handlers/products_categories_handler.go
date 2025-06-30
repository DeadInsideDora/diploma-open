package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"web/internal/domain"
)

type requestProductsCategories struct {
	Type    string          `json:"type"`
	Filters []domain.Filter `json:"filters"`
}

func RegisterProductsByCategoriesHandler(mux *http.ServeMux, path string, metricsService domain.IMetricsService, categoriesService domain.ICategoriesService, readerService domain.IReaderService) {
	categories, err := categoriesService.Get(false)
	if err != nil {
		panic(err)
	}

	mux.HandleFunc(
		fmt.Sprintf("POST %s", path),
		func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			var request requestProductsCategories
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				WriteHttpError(path, metricsService, w, fmt.Sprintf("invalid JSON: %s", err), http.StatusBadRequest)
				return
			}

			if !checkProductType(categories, request.Type) {
				WriteHttpError(path, metricsService, w, "invalid product type", http.StatusBadRequest)
				return
			}

			products, err := readerService.ReadByCategory(request.Type)

			if err != nil {
				WriteHttpError(path, metricsService, w, fmt.Sprintf("can't get products from redis: %s", err), http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(w).Encode(filterProducts(products, request.Filters)); err != nil {
				WriteHttpError(path, metricsService, w, fmt.Sprintf("encode error: %s", err), http.StatusInternalServerError)
			} else {
				metricsService.Inc(path, http.StatusOK)
			}
		},
	)
}

func checkProductType(categories []domain.Category, categoryType string) bool {
	for _, category := range categories {
		if category.Type == categoryType {
			return true
		}
	}

	return false
}

func filterProducts(products []domain.MatchData, filters []domain.Filter) []domain.MatchData {
	result := []domain.MatchData{}

	for _, product := range products {
		passed := true

		for _, filter := range filters {
			for _, data := range product.Data {
				if data.Key == filter.Name {
					passed = passed && slices.Contains(filter.Values, data.Value)
					if !passed {
						break
					}
				}
			}
			if !passed {
				break
			}
		}

		if passed {
			result = append(result, product)
		}
	}

	return result
}
