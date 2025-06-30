package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"optimizer/internal/domain"
)

type MapsService struct {
	url *url.URL
}

type mapsPayload struct {
	Point  domain.Point `json:"point"`
	Radius int64        `json:"radius"`
}

type nearShopsResponse struct {
	Shops []domain.ShopInfo `json:"shops"`
}

type routesRequest struct {
	From []domain.Point `json:"from"`
	To   []domain.Point `json:"to"`
	Type string         `json:"type"`
}

type routesResponse struct {
	Info []domain.RoutesInfo `json:"info"`
}

type tspPayload struct {
	Points     []domain.Point `json:"points"`
	StartPoint int            `json:"startPoint"`
	ByDistance bool           `json:"byDistance"`
	Algorithm  string         `json:"algorithm"`
}

type tspResponse struct {
	Routes []domain.MinTimeRoute `json:"routes"`
}

func NewMapsService(host string) *MapsService {
	u, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	return &MapsService{url: u}
}

func (service *MapsService) GetNearShops(point domain.Point, radius int64) ([]domain.ShopInfo, error) {
	log.Printf("MapsService.GetNearShops(%+v, %d)", point, radius)

	payload := mapsPayload{Point: point, Radius: radius}
	requestBody, _ := json.Marshal(payload)

	resp, err := http.Post(service.url.JoinPath("shops").String(), "application/json", bytes.NewBuffer(requestBody))
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

	var result nearShopsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %+v", err)
	}

	return result.Shops, nil
}

func (service *MapsService) GetRoutesBetweenAddresses(source, targets []domain.Point, transport string) ([]domain.RoutesInfo, error) {
	log.Printf("MapsService.GetRoutesBetweenAddresses(%+v, %+v, %s)", source, targets, transport)

	result := make([]domain.RoutesInfo, len(source))

	for i := 0; i < len(source); i += 10 {
		for j := 0; j < len(targets); j += 10 {
			sourceSlice := source[i:min(len(source), i+10)]
			targetsSlice := targets[j:min(len(targets), j+10)]

			payload := routesRequest{From: sourceSlice, To: targetsSlice, Type: transport}
			requestBody, _ := json.Marshal(payload)

			resp, err := http.Post(service.url.JoinPath("distance").String(), "application/json", bytes.NewBuffer(requestBody))
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

			var response routesResponse
			if err := json.Unmarshal(body, &response); err != nil {
				return nil, fmt.Errorf("error decoding response: %+v", err)
			}

			for _, routeInfo := range response.Info {
				result[i+routeInfo.From].From = i + routeInfo.From
				for _, route := range routeInfo.Routes {
					result[i+routeInfo.From].Routes = append(result[i+routeInfo.From].Routes, domain.Route{To: j + route.To, Time: route.Time})
				}
			}
		}
	}

	return result, nil
}

func (service *MapsService) GetTSP(points []domain.Point, startPoint int) (*domain.MinTimeRoute, error) {
	log.Printf("MapsService.GetTSP(%+v, %d)", points, startPoint)

	payload := tspPayload{Points: points, StartPoint: startPoint, ByDistance: false, Algorithm: "dp"}
	requestBody, _ := json.Marshal(payload)

	resp, err := http.Post(service.url.JoinPath("optimal-routes").String(), "application/json", bytes.NewBuffer(requestBody))
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

	var response tspResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %+v", err)
	}

	for _, mtr := range response.Routes {
		if mtr.Transport == "walking" {
			return &mtr, nil
		}
	}
	return nil, fmt.Errorf("no min time route for walking type")
}
