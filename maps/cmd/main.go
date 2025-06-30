package main

import (
	"fmt"
	"log"
	"maps_service/internal/api"
	"net/http"
	"os"
)

type Settings struct {
	ApiKey2Gis string
	Port       string
}

func GetSettings() (*Settings, error) {
	settings := Settings{}

	if apiKey := os.Getenv("API_KEY_2GIS"); apiKey != "" {
		settings.ApiKey2Gis = apiKey
	} else {
		return nil, fmt.Errorf("can't get API_KEY_2GIS env")
	}

	if port := os.Getenv("PORT"); port != "" {
		settings.Port = port
	} else {
		return nil, fmt.Errorf("can't get PORT env")
	}

	return &settings, nil
}

func main() {
	settings, err := GetSettings()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /available-shops", api.CreateAvailableShopsHandler(settings.ApiKey2Gis))
	mux.HandleFunc("POST /shops", api.CreateShopsHandler(settings.ApiKey2Gis))
	mux.HandleFunc("POST /optimal-routes", api.CreateOptimalRoutesHandler(settings.ApiKey2Gis))
	mux.HandleFunc("POST /distance", api.CreateDistanceHandler(settings.ApiKey2Gis))

	cors := api.CorsMiddleware(mux)

	log.Printf("Server starting on :%s\n", settings.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", settings.Port), cors); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
