package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"optimizer/internal/domain"
)

type request struct {
	Products      []domain.InputProductInfo `json:"products"`
	DiscountCards []string                  `json:"discount_cards"`
	UserPoint     domain.Point              `json:"point"`
	Radius        int64                     `json:"radius"`
	Exchange      int64                     `json:"exchange"`
}

func CreateProductsHandler(optimizer domain.IOptimizerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req request
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := optimizer.Get(req.Products, req.DiscountCards, req.UserPoint, req.Radius, req.Exchange)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func CreateNearbyProductsHandler(nearbyProducts domain.INearbyProductsService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req request
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := nearbyProducts.Get(req.Products, req.DiscountCards, req.UserPoint, req.Radius, req.Exchange)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
