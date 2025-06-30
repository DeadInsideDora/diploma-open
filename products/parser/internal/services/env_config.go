package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"scrappers/internal/domain"
)

type configData struct {
	LentaDeviceId     string `json:"lentaDeviceId"`
	LentaSessionToken string `json:"lentaSessionToken"`
	PerekrestokAuth   string `json:"perekrestokAuth"`
	PerekrestokCookie string `json:"perekrestokCookie"`
	MagnitCookie      string `json:"magnitCookie"`
}

type EnvConfigService struct {
	path string
}

func NewEnvConfigService(path string) *EnvConfigService {
	return &EnvConfigService{path: path}
}

func (service *EnvConfigService) Get() (*domain.Environments, error) {
	jsonFile, err := os.Open(service.path)
	if err != nil {
		return nil, fmt.Errorf("can't load env config file: %s", err)
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var config configData
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("can't unmarshal config: %s", err)
	}

	return domain.NewEnvironments(domain.NewLentaEnv(config.LentaDeviceId, config.LentaSessionToken), domain.NewPerekrestokEnv(config.PerekrestokAuth, config.PerekrestokCookie), domain.NewMagnitEnv(config.MagnitCookie)), nil
}
