package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scrappers/internal/domain"
	"time"
)

const (
	perekrestok_payload       = "{\"filter\":{\"category\":%d,\"onlyWithProductReviews\":false, \"features\":%s},\"withBestProductReviews\":false}"
	perekrestok_brand_payload = "{\"filter\":{\"category\":%d},\"page\":%d,\"perPage\":20}"
	perekrestok_api_url       = "https://www.perekrestok.ru/api/customer/1.4.1.0/catalog/product/grouped-feed"
	perekrestok_pork_id       = 139
)

type perekrestokResponse struct {
	Content content `json:"content"`
}

type content struct {
	Items []item `json:"items"`
}

type item struct {
	Products []perekrestokProduct `json:"products"`
}

type perekrestokProduct struct {
	Id    int64               `json:"id"`
	Price perekrestokPriceTag `json:"priceTag"`
	Image perekrestokImage    `json:"image"`
	Title string              `json:"title"`
}

type perekrestokImage struct {
	Url string `json:"cropUrlTemplate"`
}

type perekrestokPriceTag struct {
	Price      int64  `json:"price"`
	GrossPrice *int64 `json:"grossPrice"`
}

type PerekrestokRequester struct {
	metricsService domain.IDomainMetricsService
	envService     domain.IEnvService
}

func NewPerekrestokRequester(metricsService domain.IDomainMetricsService, envService domain.IEnvService) *PerekrestokRequester {
	return &PerekrestokRequester{metricsService: metricsService, envService: envService}
}

func (requester *PerekrestokRequester) Process(id int, features []domain.Feature, delay int) ([]domain.ProductInfo, error) {
	envSnapshot, err := requester.envService.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get environments for perekrestok requester: %s", err)
	}
	env := envSnapshot.GetPerekrestokEnv()

	featuresBytes, err := json.Marshal(features)
	if err != nil {
		return nil, fmt.Errorf("can't marshal features: %s", err)
	}
	featuresString := string(featuresBytes)

	resp, err := doPerekrestokRequest(perekrestok_api_url, fmt.Sprintf(perekrestok_payload, id, featuresString), env)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %+v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %+v", err)
	}

	requester.metricsService.Inc("perekrestok", resp.StatusCode)

	log.Printf("PerekrestokRequest: id=%d, feature=%s; status=%d", id, featuresString, resp.StatusCode)

	var response perekrestokResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshal perekrestok response: %+v", err)
	}

	result := []domain.ProductInfo{}

	for _, item := range response.Content.Items {
		result = append(result, makePerekrestokProductInfos(item.Products)...)
	}

	time.Sleep(time.Duration(delay) * time.Second)

	return result, nil
}

func (requester *PerekrestokRequester) Ping() error {
	envSnapshot, err := requester.envService.Get()
	if err != nil {
		return fmt.Errorf("can't get environments: %s", err)
	}
	env := envSnapshot.GetPerekrestokEnv()

	resp, err := doPerekrestokRequest(perekrestok_api_url, fmt.Sprintf(perekrestok_payload, perekrestok_pork_id, "[]"), env)
	if err != nil {
		return fmt.Errorf("ping perekrestok return err: %s", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("ping perekrestok return not 200 code: %s", resp.Status)
	}

	return nil
}

func doPerekrestokRequest(url, payload string, env *domain.PerekrestokEnv) (*http.Response, error) {
	requestPayload := []byte(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %+v", err)
	}
	req.Header.Add("auth", env.GetAuth())
	req.Header.Add("cookie", env.GetCookie())
	req.Header.Add("content-type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")

	return processRequest(req)
}

func makePerekrestokProductInfos(products []perekrestokProduct) []domain.ProductInfo {
	result := []domain.ProductInfo{}

	for _, product := range products {
		info := domain.ProductInfo{Id: product.Id, Title: product.Title, ShopName: "Перекрёсток", PictureUrl: fmt.Sprintf(product.Image.Url, "400x400")}
		if product.Price.GrossPrice == nil {
			info.PriceDiscount, info.PriceRegular = product.Price.Price, product.Price.Price
		} else {
			info.PriceDiscount, info.PriceRegular = product.Price.Price, *product.Price.GrossPrice
		}
		result = append(result, info)
	}

	return result
}
