package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"service/internal/api"
	"service/internal/services"
)

type Settings struct {
	DatabaseUrl string
}

func GetSettings() (*Settings, error) {
	settings := Settings{}

	if databaseUrl := os.Getenv("DATABASE_URL"); databaseUrl != "" {
		settings.DatabaseUrl = databaseUrl
	} else {
		return nil, fmt.Errorf("can't get DATABASE_URL env")
	}

	return &settings, nil
}

func main() {
	settings, err := GetSettings()
	if err != nil {
		panic(err)
	}

	dbService, err := services.NewPgxService(settings.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("POST /register", api.RegisterHandler(dbService))
	mux.Handle("POST /login", api.LoginHandler(dbService))
	mux.Handle("GET /auth/me", api.AuthMeHandler(dbService))
	mux.Handle("GET /data", api.GetDataHandler(dbService))
	mux.Handle("PUT /cards", api.UpdateCardsHandler(dbService))
	mux.Handle("PUT /map-info", api.UpdateMapInfoHandler(dbService))
	mux.Handle("PUT /exchange", api.UpdateExchangeHandler(dbService))

	cors := api.CorsMiddleware(mux)

	log.Fatal(http.ListenAndServe(":8080", cors).Error())
}
