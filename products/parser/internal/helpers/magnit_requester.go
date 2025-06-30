package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scrappers/internal/domain"
	"strconv"
	"time"
)

const (
	magnit_payload  = "{\"sort\":{\"order\":\"desc\",\"type\":\"popularity\"},\"pagination\":{\"limit\":30,\"offset\":%d},\"filters\":%s,\"categories\":%s,\"includeAdultGoods\":true,\"storeCode\":\"784198\",\"storeType\":\"1\",\"catalogType\":\"1\"}"
	magnit_api_url  = "https://magnit.ru/webgate/v2/goods/search"
	magnit_pork_ids = "[4868]"
)

type magnitResponse struct {
	Items      []magnitItem     `json:"items"`
	Pagination magnitPagination `json:"pagination"`
}

type magnitItem struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	Gallery   []magnitGallery `json:"gallery"`
	Price     int64           `json:"price"`
	Promotion magnitPromotion `json:"prices"`
	Weighted  magnitWeighted  `json:"weighted"`
}

type magnitPromotion struct {
	OldPrice *int64 `json:"oldPrice"`
}

type magnitGallery struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type magnitPagination struct {
	More bool `json:"hasMore"`
}

type magnitWeighted struct {
	Weight *int64 `json:"shelfWeight"`
}

type MagnitRequester struct {
	metricsService domain.IDomainMetricsService
	envService     domain.IEnvService
}

func NewMagnitRequester(metricsService domain.IDomainMetricsService, envService domain.IEnvService) *MagnitRequester {
	return &MagnitRequester{metricsService: metricsService, envService: envService}
}

func (requester *MagnitRequester) Process(ids int, filters []domain.MagnitFilter, delay int) ([]domain.ProductInfo, error) {
	envSnapshot, err := requester.envService.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get environments: %s", err)
	}
	env := envSnapshot.GetMagnitEnv()

	idsBytes, err := json.Marshal([]int{ids})
	if err != nil {
		return nil, fmt.Errorf("can't marshal id: %s", err)
	}
	idsString := string(idsBytes)

	filtersBytes, err := json.Marshal(filters)
	if err != nil {
		return nil, fmt.Errorf("can't marshal filters: %s", err)
	}
	filtersString := string(filtersBytes)

	result := []domain.ProductInfo{}

	for {
		resp, err := doMagnitRequest(env, idsString, filtersString, len(result))
		if err != nil {
			return nil, fmt.Errorf("error sending request: %+v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %+v", err)
		}

		requester.metricsService.Inc("magnit", resp.StatusCode)

		log.Printf("MagnitRequest: ids=%s; status=%d", idsString, resp.StatusCode)

		var response magnitResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("error unmarshal magnit response: %+v", err)
		}

		result = append(result, makeMagnitProductInfos(response.Items)...)

		time.Sleep(time.Duration(delay) * time.Second)

		if !response.Pagination.More {
			break
		}
	}

	return result, nil
}

func (requester *MagnitRequester) Ping() error {
	envSnapshot, err := requester.envService.Get()
	if err != nil {
		return fmt.Errorf("can't get environments: %s", err)
	}
	env := envSnapshot.GetMagnitEnv()

	resp, err := doMagnitRequest(env, magnit_pork_ids, "null", 0)
	if err != nil {
		return fmt.Errorf("ping magnit return err: %s", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("ping magnit return not 200 code: %s", resp.Status)
	}

	return nil
}

func doMagnitRequest(env *domain.MagnitEnv, ids string, filters string, offset int) (*http.Response, error) {
	requestPayload := []byte(fmt.Sprintf(magnit_payload, offset, filters, ids))

	req, err := http.NewRequest("POST", magnit_api_url, bytes.NewBuffer(requestPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %+v", err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 YaBrowser/25.4.0.0 Safari/537.36")
	req.Header.Add("Cookie", env.GetCookie())
	req.Header.Add("content-type", "application/json")

	return processRequest(req)
}

func makeMagnitProductInfos(items []magnitItem) []domain.ProductInfo {
	result := []domain.ProductInfo{}

	for _, item := range items {
		url := ""
		for _, gallery := range item.Gallery {
			if gallery.Type == "IMAGE" {
				url = gallery.Url
				break
			}
		}
		i64, err := strconv.ParseInt(item.Id, 10, 64)
		if err != nil {
			log.Printf("MagnitRequester: can't parse id to int64: %s", item.Id)
			continue
		}
		info := domain.ProductInfo{Id: i64, Title: item.Name, ShopName: "Магнит", PictureUrl: url}
		if item.Promotion.OldPrice != nil {
			info.PriceDiscount, info.PriceRegular = item.Price, *item.Promotion.OldPrice
		} else {
			info.PriceDiscount, info.PriceRegular = item.Price, item.Price
		}
	
		if item.Weighted.Weight != nil {
			info.PriceDiscount = info.PriceDiscount * 1000 / *item.Weighted.Weight
			info.PriceRegular = info.PriceRegular * 1000 / *item.Weighted.Weight
		}

		result = append(result, info)
	}

	return result
}
