package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps_service/internal/domain"
	"net/http"
)

const api_2gis_routing = "https://routing.api.2gis.com/get_dist_matrix?key=%s&version=6.0.0"

type matrixReq struct {
	Points    []domain.Point `json:"points"`
	Sources   []int          `json:"sources"`
	Targets   []int          `json:"targets"`
	Transport string         `json:"transport"`
}

type route struct {
	Distance int    `json:"distance"`
	Duration int    `json:"duration"`
	Source   int    `json:"source_id"`
	Target   int    `json:"target_id"`
	Status   string `json:"status"`
}

type matrixResp struct {
	Routes []route `json:"routes"`
}

type Matrix2GisService struct {
	apiKey string
}

func NewMatrix2GisService(apiKey string) *Matrix2GisService {
	return &Matrix2GisService{apiKey: apiKey}
}

func createMatrix(n, m int) [][]int {
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, m)
		for j := range matrix[i] {
			matrix[i][j] = -1
		}
	}
	return matrix
}

func (matrixService *Matrix2GisService) Get(points []domain.Point, sources, targets []int, transport string) (distanceMatrix [][]int, durationMatrix [][]int, returnErr error) {
	distanceMatrix, durationMatrix = nil, nil

	if len(points) < 2 {
		returnErr = fmt.Errorf("invalid points count")
		return
	}

	request := matrixReq{
		Points:    points,
		Sources:   sources,
		Targets:   targets,
		Transport: transport,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		returnErr = fmt.Errorf("error encoding request body: %+v", err)
		return
	}

	resp, err := http.Post(fmt.Sprintf(api_2gis_routing, matrixService.apiKey), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		returnErr = fmt.Errorf("error sending request: %+v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		returnErr = fmt.Errorf("error reading response body: %+v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		returnErr = fmt.Errorf("api returned status %d: %s", resp.StatusCode, string(body))
		return
	}

	var responseMatrix matrixResp

	err = json.Unmarshal(body, &responseMatrix)
	if err != nil {
		returnErr = fmt.Errorf("error decoding response: %+v", err)
		return
	}

	n := len(points)

	distanceMatrix = createMatrix(n, n)
	durationMatrix = createMatrix(n, n)

	for _, route := range responseMatrix.Routes {
		if route.Status == "OK" {
			distanceMatrix[route.Source][route.Target] = route.Distance
			durationMatrix[route.Source][route.Target] = route.Duration
		}
	}

	return
}
