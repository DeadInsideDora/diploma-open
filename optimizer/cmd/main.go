package main

import (
	"fmt"
	"log"
	"net/http"
	"optimizer/internal/handler"
	"optimizer/internal/services"
	"os"
)

type Settings struct {
	Products string
	Maps     string
	Port     string
}

func GetSettings() (*Settings, error) {
	settings := Settings{}

	if products := os.Getenv("PRODUCTS_SERVICE_URL"); products != "" {
		settings.Products = products
	} else {
		return nil, fmt.Errorf("can't get PRODUCTS_SERVICE_URL env")
	}

	if maps := os.Getenv("MAPS_SERVICE_URL"); maps != "" {
		settings.Maps = maps
	} else {
		return nil, fmt.Errorf("can't get MAPS_SERVICE_URL env")
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

	productService := services.NewProductService(settings.Products)
	mapsService := services.NewMapsService(settings.Maps)
	optimizer := services.NewOptimizerService(mapsService, productService)
	nearbyProducts := services.NewNearbyProductsService(mapsService, productService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", handler.CreateProductsHandler(optimizer))
	mux.HandleFunc("POST /nearby-products", handler.CreateNearbyProductsHandler(nearbyProducts))

	handlerWithCORS := handler.CorsMiddleware(mux)

	log.Printf("Server starting on :%s", settings.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", settings.Port), handlerWithCORS); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
