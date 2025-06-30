package main

import (
	"fmt"
	"net/http"
	"os"
	"web/internal/handlers"
	"web/internal/services"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Settings struct {
	RedisAddr string
	Port      string
}

func GetSettings() (*Settings, error) {
	settings := Settings{}

	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		settings.RedisAddr = redisAddr
	} else {
		return nil, fmt.Errorf("can't get REDIS_ADDR env")
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

	categoriesService := services.NewCategoriesService(os.Args[1])
	readerService := services.NewRedisReaderService(settings.RedisAddr)
	metricsService := services.NewPrometheusMetricsService()

	mux := http.NewServeMux()

	handlers.RegisterCategoriesHandler(mux, "/categories", metricsService, categoriesService)
	handlers.RegisterProductsByCategoriesHandler(mux, "/products-by-category", metricsService, categoriesService, readerService)
	handlers.RegisterProductsByNamesHandler(mux, "/products-by-names", metricsService, categoriesService, readerService)
	mux.Handle("/metrics", promhttp.Handler())

	handlerWithCORS := handlers.CorsMiddleware(mux)

	fmt.Printf("Server on :%s\n", settings.Port)
	if err := http.ListenAndServe(fmt.	Sprintf(":%s", settings.Port), handlerWithCORS); err != nil {
		readerService.Close()
		panic(err)
	}
}
