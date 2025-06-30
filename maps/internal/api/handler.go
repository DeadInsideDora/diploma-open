package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps_service/internal/domain"
	"maps_service/internal/services"
	"net/http"
)

type shopsRequest struct {
	Point  domain.Point `json:"point"`
	Radius int64        `json:"radius"`
}

type shopsResponse struct {
	Shops []domain.ShopInfo `json:"shops"`
}

type tspResp struct {
	Routes []domain.MinTimeRoute `json:"routes"`
}

type distanceRequest struct {
	From []domain.Point `json:"from"`
	To   []domain.Point `json:"to"`
	Type string         `json:"type"`
}

type distanceResponse struct {
	Info []pointInfo `json:"info"`
}

type pointInfo struct {
	From   int     `json:"from"`
	Routes []route `json:"routes"`
}

type route struct {
	To   int `json:"to"`
	Time int `json:"time"`
}

func logHttpError(w http.ResponseWriter, message string, code int) {
	log.Println(message)
	http.Error(w, message, code)
}

func CreateShopsHandler(apiKey string) func(w http.ResponseWriter, r *http.Request) {
	shopInfo2Gis := services.NewShopInfo2GisService(apiKey)
	shopsRequester := services.NewShopsRequester(shopInfo2Gis)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST /shops")

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			logHttpError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logHttpError(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Body: %s", string(body))

		var request shopsRequest
		if err := json.Unmarshal(body, &request); err != nil {
			logHttpError(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(shopsResponse{
			Shops: shopsRequester.GetNearbyShops(request.Point, request.Radius),
		})
	}
}

func CreateAvailableShopsHandler(apiKey string) func(w http.ResponseWriter, r *http.Request) {
	shopInfo2Gis := services.NewShopInfo2GisService(apiKey)
	shopsRequester := services.NewShopsRequester(shopInfo2Gis)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /available-shops")

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		json.NewEncoder(w).Encode(shopsRequester.GetAvailableShops())
	}
}

func CreateOptimalRoutesHandler(apiKey string) func(w http.ResponseWriter, r *http.Request) {
	matrixService2Gis := services.NewMatrix2GisService(apiKey)
	tspBruteforce := services.NewTSPBruteforce()
	tspDynProgramming := services.NewTSPDynProgramming()

	routingServiceBruteforce := services.NewRoutingService(matrixService2Gis, tspBruteforce)
	routingServiceDynProgramming := services.NewRoutingService(matrixService2Gis, tspDynProgramming)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			Points     []domain.Point `json:"points"`
			StartPoint int            `json:"startPoint"`
			ByDistance bool           `json:"byDistance"`
			Algorithm  string         `json:"algorithm"`
		}

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

		if err := json.Unmarshal(body, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("POST /optimal-routes: %s", string(body))

		if !(request.Algorithm == "bruteforce" || request.Algorithm == "dp") {
			http.Error(w, "invalid algorithm for tsp", http.StatusBadRequest)
			return
		}

		routes := func() domain.IRoutingService {
			if request.Algorithm == "bruteforce" {
				return routingServiceBruteforce
			}
			return routingServiceDynProgramming
		}().Get(request.Points, request.StartPoint, request.ByDistance)

		json.NewEncoder(w).Encode(tspResp{Routes: routes})
	}
}

func CreateDistanceHandler(apiKey string) func(http.ResponseWriter, *http.Request) {
	matrixService2Gis := services.NewMatrix2GisService(apiKey)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST /distance")
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

		var request distanceRequest

		if err := json.Unmarshal(body, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		points := []domain.Point{}
		sources := []int{}
		targets := []int{}

		for i := 0; i < len(request.From); i += 1 {
			sources = append(sources, i)
		}
		for i := 0; i < len(request.To); i += 1 {
			targets = append(targets, len(request.From)+i)
		}
		points = append(points, request.From...)
		points = append(points, request.To...)

		_, dur, err := matrixService2Gis.Get(points, sources, targets, request.Type)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		info := []pointInfo{}
		for i := 0; i < len(request.From); i += 1 {
			point := pointInfo{From: i, Routes: []route{}}

			for j := 0; j < len(request.To); j += 1 {
				val := dur[i][len(request.From)+j]
				if val != -1 {
					point.Routes = append(point.Routes, route{To: j, Time: val})
				}
			}

			info = append(info, point)
		}

		if err := json.NewEncoder(w).Encode(distanceResponse{Info: info}); err != nil {
			http.Error(w, fmt.Sprintf("encode error: %s", err), http.StatusInternalServerError)
		}
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
