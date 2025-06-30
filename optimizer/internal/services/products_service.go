package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"optimizer/internal/domain"
)

type ProductService struct {
	url *url.URL
}

type productPayload struct {
	Type  string   `json:"type"`
	Names []string `json:"names"`
}

func NewProductService(host string) *ProductService {
	u, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	return &ProductService{url: u}
}

func (service *ProductService) GetProducts(category string, names []string) ([]domain.MatchData, error) {
	payload := productPayload{Type: category, Names: names}

	requestBody, _ := json.Marshal(payload)

	resp, err := http.Post(service.url.JoinPath("products-by-names").String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response body: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status %d: %s", resp.StatusCode, string(body))
	}

	var result []domain.MatchData

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %+v", err)
	}

	return result, nil
}
