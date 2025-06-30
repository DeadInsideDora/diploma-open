package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"scrappers/internal/helpers"
	"scrappers/internal/services"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Settings struct {
	Port      string
	RedisAddr string
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

	configService := services.NewLocalConfigService(os.Args[1])
	envService := services.NewEnvConfigService(os.Args[2])

	writerFactory := services.NewRedisWriterFactory(settings.RedisAddr)
	readerFactory := services.NewRedisReaderFactory(settings.RedisAddr)

	matcher, err := services.NewJsonMatcherService(readerFactory, os.Args[3])
	if err != nil {
		log.Fatalf("can't get matcher: %s", err)
	}

	metricsService := services.NewPrometheusMetricsService()

	lenta := helpers.NewLentaRequester(metricsService, envService)
	perekrestok := helpers.NewPerekrestokRequester(metricsService, envService)
	dixy := helpers.NewDixyRequester(metricsService)
	magnit := helpers.NewMagnitRequester(metricsService, envService)

	go services.NewExecutorService(configService, matcher, writerFactory, lenta, perekrestok, dixy, magnit).Start()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	log.Printf("Listening on :%s\n", settings.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", settings.Port), mux); err != nil {
		panic(err)
	}
}
