package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps_service/internal/domain"
	"net/http"
	"strings"
)

const api_2gis_shops string = "https://catalog.api.2gis.com/3.0/items?q=%s&point=%f,%f&radius=%d&fields=items.point,items.org,items.rubrics,items.schedule&key=%s"

var shops = []string{
	"Лента",
	"Перекрёсток",
	"Дикси",
	"Магнит",
}

type result struct {
	Items []placeWithDebug `json:"items"`
}

type placeWithDebug struct {
	Name     string          `json:"name"`
	Point    domain.Point    `json:"point"`
	Id       string          `json:"id"`
	Rubrics  []rubric        `json:"rubrics"`
	Org      org             `json:"org"`
	Schedule domain.Schedule `json:"schedule"`
}

type rubric struct {
	Name string `json:"name"`
}

type org struct {
	Name string `json:"name"`
}

type response2Gis struct {
	Result result `json:"result"`
}

type ShopInfo2GisService struct {
	apiKey string
}

func NewShopInfo2GisService(apiKey string) *ShopInfo2GisService {
	return &ShopInfo2GisService{apiKey: apiKey}
}

func (shopInfo *ShopInfo2GisService) Get(shop string, point domain.Point, radius int64) (*domain.ShopInfo, error) {
	url := fmt.Sprintf(api_2gis_shops, shop, point.Lon, point.Lat, radius, shopInfo.apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %+v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %+v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %+v", err)
	}

	var response response2Gis
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshal 2Gis response")
	}

	return &domain.ShopInfo{
		Info: filterAndTransform(response.Result.Items, shop),
		Shop: shop,
	}, nil
}

type ShopsRequester struct {
	shopInfoService domain.IShopInfoService
}

func NewShopsRequester(shopInfoService domain.IShopInfoService) *ShopsRequester {
	return &ShopsRequester{shopInfoService: shopInfoService}
}

func (sr *ShopsRequester) GetAvailableShops() []string {
	return shops[:]
}

func (sr *ShopsRequester) GetNearbyShops(point domain.Point, radius int64) []domain.ShopInfo {
	var shops []domain.ShopInfo

	for _, shop := range sr.GetAvailableShops() {
		log.Printf("Process shop: %s, for point = (%f, %f), radius = %d", shop, point.Lon, point.Lat, radius)

		shopInfo, err := sr.shopInfoService.Get(shop, point, radius)

		log.Printf("ShopInfo: %+v, err: %+v", shopInfo, err)
		if err == nil {
			shops = append(shops, *shopInfo)
		} else {
			log.Println(err)
		}
	}

	return shops
}

func filterAndTransform(places []placeWithDebug, shop string) []domain.Place {
	result := []domain.Place{}

	for _, place := range places {
		if checkOrg(place.Org, shop) && checkRubrics(place.Rubrics) {
			result = append(result, domain.Place{Name: place.Name, Point: place.Point, Id: place.Id, Schedule: place.Schedule})
		}
	}

	return result
}

func checkRubrics(rubrics []rubric) bool {
	for _, rubric := range rubrics {
		if rubric.Name == "Супермаркеты" || rubric.Name == "Гипермаркеты" {
			return true
		}
	}
	return false
}

func checkOrg(organization org, shop string) bool {
	return strings.Contains(organization.Name, shop)
}
