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
	lenta_payload = "{\"filters\":{\"range\":[],\"checkbox\":[],\"multicheckbox\":%s},\"categoryId\":%d,\"sort\":{\"type\":\"popular\",\"order\":\"desc\"},\"limit\":20,\"offset\":%d}"
	lenta_api_url = "https://lenta.com/api-gateway/v1/catalog/items"
	lenta_pork_id = 651
)

type lentaResponse struct {
	Items []lentaItem `json:"items"`
	Total int64       `json:"total"`
}

type lentaItem struct {
	Id     int64             `json:"id"`
	Name   string            `json:"name"`
	Images []lentaImage      `json:"images"`
	Prices lentaProductPrice `json:"prices"`
	Weight lentaWeight       `json:"weight"`
}

type lentaProductPrice struct {
	Cost        int64 `json:"cost"`
	CostRegular int64 `json:"costRegular"`
}

type lentaImage struct {
	Original string `json:"original"`
}

type lentaWeight struct {
	Net int64 `json:"net"`
}

type LentaRequester struct {
	metricsService domain.IDomainMetricsService
	envService     domain.IEnvService
}

func NewLentaRequester(metricsService domain.IDomainMetricsService, envService domain.IEnvService) *LentaRequester {
	return &LentaRequester{metricsService: metricsService, envService: envService}
}

func (requester *LentaRequester) Process(id int, multicheckboxes []domain.Multicheckbox, delay int) ([]domain.ProductInfo, error) {
	envSnapshot, err := requester.envService.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get environments: %s", err)
	}
	env := envSnapshot.GetLentaEnv()

	checkboxesBytes, err := json.Marshal(multicheckboxes)
	if err != nil {
		return nil, fmt.Errorf("can't marshal multicheckboxes: %s", err)
	}
	checkboxes := string(checkboxesBytes)

	result := []domain.ProductInfo{}
	total := 0

	for len(result) != total || len(result) == 0 {
		resp, err := doLentaRequest(env, checkboxes, id, len(result))
		if err != nil {
			return nil, fmt.Errorf("error sending request: %+v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %+v", err)
		}

		requester.metricsService.Inc("lenta", resp.StatusCode)

		log.Printf("LentaRequest: ids=%d, multicheckboxes=%s, id=%d; status=%d", id, checkboxes, id, resp.StatusCode)

		var response lentaResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("error unmarshal lenta response: %+v", err)
		}

		result = append(result, makeLentaProductInfos(response.Items)...)
		total = int(response.Total)

		time.Sleep(time.Duration(delay) * time.Second)
	}

	return result, nil
}

func (requester *LentaRequester) Ping() error {
	envSnapshot, err := requester.envService.Get()
	if err != nil {
		return fmt.Errorf("can't get environments: %s", err)
	}
	env := envSnapshot.GetLentaEnv()

	resp, err := doLentaRequest(env, "[]", lenta_pork_id, 0)
	if err != nil {
		return fmt.Errorf("ping lenta return err: %s", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("ping lenta return not 200 code: %s", resp.Status)
	}

	return nil
}

func doLentaRequest(env *domain.LentaEnv, checkboxes string, id int, offset int) (*http.Response, error) {
	requestPayload := []byte(fmt.Sprintf(lenta_payload, checkboxes, id, offset))

	req, err := http.NewRequest("POST", lenta_api_url, bytes.NewBuffer(requestPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %+v", err)
	}
	req.Header.Add("deviceid", env.GetDeviceId())
	req.Header.Add("sessiontoken", env.GetSessionToken())
	req.Header.Add("x-platform", "omniweb")
	req.Header.Add("x-retail-brand", "lo")
	req.Header.Add("content-type", "application/json")

	return processRequest(req)
}

func makeLentaProductInfos(items []lentaItem) []domain.ProductInfo {
	result := []domain.ProductInfo{}

	for _, item := range items {
		url := ""
		if len(item.Images) != 0 {
			url = item.Images[0].Original
		}
		info := domain.ProductInfo{Id: item.Id, Title: item.Name, ShopName: "Лента", PictureUrl: url, PriceDiscount: item.Prices.Cost, PriceRegular: item.Prices.CostRegular}
		if item.Weight.Net != 0 {
			info.PriceDiscount = info.PriceDiscount * 1000 / item.Weight.Net
			info.PriceRegular = info.PriceRegular * 1000 / item.Weight.Net
		}

		result = append(result, info)
	}

	return result
}
