package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"scrappers/internal/domain"
)

type LocalConfigService struct {
	path string
}

func NewLocalConfigService(path string) *LocalConfigService {
	return &LocalConfigService{path: path}
}

func (service *LocalConfigService) Get() (*domain.Config, error) {
	jsonFile, err := os.Open(service.path)
	if err != nil {
		return nil, fmt.Errorf("can't load config file: %s", err)
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var config domain.Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("can't unmarshal config: %s", err)
	}

	return &config, nil
}
