package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"web/internal/domain"
)

type CategoriesService struct {
	path   string
	cached *cachedCategories
}

type cachedCategories struct {
	Categories []domain.Category
}

func NewCategoriesService(path string) *CategoriesService {
	return &CategoriesService{path: path}
}

func (service *CategoriesService) Get(needToUpdate bool) ([]domain.Category, error) {
	if !needToUpdate && service.cached != nil {
		return service.cached.Categories, nil
	}

	jsonFile, err := os.Open(service.path)
	if err != nil {
		return nil, fmt.Errorf("can't load categories info file: %s", err)
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var categories []domain.Category
	if err := json.Unmarshal(bytes, &categories); err != nil {
		return nil, fmt.Errorf("can't unmarshal categories: %s", err)
	}

	service.cached = &cachedCategories{Categories: categories}
	return categories, nil
}
