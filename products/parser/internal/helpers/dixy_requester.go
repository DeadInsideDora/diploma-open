package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"scrappers/internal/domain"
	"strconv"
	"strings"
	"time"
)

const (
	dixy_api_url = "https://dixy.ru/ajax/listing-json.php?block=product-list&sid=%d&perPage=30&page=%d"
	dixy_pork_id = 49
)

type dixyResponse struct {
	PagenData dixyPagenData `json:"pagenData"`
	Cards     []dixyCard    `json:"cards"`
}

type dixyCard struct {
	Id            int64   `json:"id,string"`
	Title         string  `json:"title"`
	Url           string  `json:"src"`
	PriceRegular  float64 `json:"oldPriceSimple"`
	PriceDiscount float64 `json:"priceSimple"`
}

func (d *dixyCard) UnmarshalJSON(data []byte) error {
	type Alias dixyCard
	aux := &struct {
		PriceRegular  string `json:"oldPriceSimple"`
		PriceDiscount string `json:"priceSimple"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	regular := strings.ReplaceAll(aux.PriceRegular, " ", "")
	discount := strings.ReplaceAll(aux.PriceDiscount, " ", "")

	var err error
	d.PriceRegular, err = strconv.ParseFloat(regular, 64)
	if err != nil {
		return err
	}

	d.PriceDiscount, err = strconv.ParseFloat(discount, 64)
	if err != nil {
		return err
	}

	return nil
}

type dixyPagenData struct {
	Count int `json:"element_count,string"`
}

type DixyRequester struct {
	metricsService domain.IDomainMetricsService
}

func NewDixyRequester(metricsService domain.IDomainMetricsService) *DixyRequester {
	return &DixyRequester{metricsService: metricsService}
}

func (requester *DixyRequester) Process(id int, filterData []domain.FilterData, delay int) ([]domain.ProductInfo, error) {
	filterDataBytes, err := json.Marshal(filterData)
	if err != nil {
		return nil, fmt.Errorf("can't get json from filterData: %+v", err)
	}
	filterDataString := string(filterDataBytes)

	result := []domain.ProductInfo{}

	total := 0

	for page := 1; len(result) != total || len(result) == 0; page += 1 {
		resp, err := doDixyRequest(id, page, filterDataString)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %+v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %+v", err)
		}

		requester.metricsService.Inc("dixy", resp.StatusCode)

		log.Printf("DixyRequest: id=%d, page=%d, filterData=%s (scrapped=%d, all=%d); status=%d", id, page, filterDataString, len(result), total, resp.StatusCode)

		var response []dixyResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("error unmarshal dixy response body: %+v", err)
		}
		if length := len(response); length != 1 {
			return nil, fmt.Errorf("error get unexpected dixy response size: %d", length)
		}

		result = append(result, makeDixyProductInfos(response[0].Cards)...)
		total = int(response[0].PagenData.Count)

		time.Sleep(time.Duration(delay) * time.Second)
	}

	return result, nil
}

func (requester *DixyRequester) Ping() error {
	resp, err := doDixyRequest(dixy_pork_id, 1, "null")
	if err != nil {
		return fmt.Errorf("ping dixy return err: %s", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("ping dixy return not 200 code: %s", resp.Status)
	}

	return nil
}

func doDixyRequest(id int, page int, filterData string) (*http.Response, error) {
	url := fmt.Sprintf("https://dixy.ru/ajax/listing-json.php?block=product-list&sid=%d&perPage=30&page=%d", id, page)

	var req *http.Request
	var err error

	if filterData != "null" {
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("filterData", filterData)
		err := writer.Close()
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest("POST", url, payload)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Cookie", "DIXY2_SM_FOREVER_BASKET=33820439; DIXY2_SM_d_guid=0%3Asu0h9c%3A; _ym_uid=1743463777915119398; _ym_d=1743463777; _ga=GA1.1.1508804860.1743463777; _userGUID=0:m8xpbggg:Nbzz0wJ3FCLz2m5fDj7enawgMry5TMgN; _userGUID=0:m8xpbggg:Nbzz0wJ3FCLz2m5fDj7enawgMry5TMgN; BX_USER_ID=1ae6a2240e5a5dd3b61870824d75abfb; PHPSESSID=Vb67VItDqkEjLUTaI9cgVYm3VMC0jnyT; tmr_lvid=d694480d5879f9aefcef7903dad39bee; tmr_lvidTS=1748457941720; hide-cookie=true; qrator_msid2=v2.0.1749424742.656.bcf229e4nTQdWhae|enLztBkJgqD9gJfy|+HDPg+e5T57oWv32gRAd7B0Uf1NWDik+cIKYNt/tGudA7m8ZDwqbxVOauG5AgSZFQR9aGtZ/VC17Qb8GEMLC4sz2FOfxT7PPCUgHvoVhaFo=-/YMxFc7z5MyUBHlQWQyXnrAqFwk=; _ym_isad=2; _ym_visorc=w; digi_uc=|v:174784:2000219613:2000437012:2000219633:2000203655:DI00080510:T000024619!174785:2000294696!174846:2000563331:2000005386!174933:2000208710|c:174386:2000563791:2000592212!174867:2000002844!174932:2000548652!174942:2000210260; dSesn=6d419f97-d289-c167-34c0-0c3bbb8aebbb; _dvs=0:mboabzev:5fTwSu8dEud6ZAdbm5QpKTToKzTqm8fw; domain_sid=bw2Vjv28LjcfY9ompH7kX%3A1749424760598; tmr_detect=0%7C1749424761938; _ga_J3JT2KMN08=GS2.1.s1749424759$o16$g0$t1749424777$j42$l0$h0")

	return processRequest(req)
}

func makeDixyProductInfos(cards []dixyCard) []domain.ProductInfo {
	result := []domain.ProductInfo{}

	for _, card := range cards {
		info := domain.ProductInfo{Id: card.Id, Title: card.Title, ShopName: "Дикси", PictureUrl: fmt.Sprintf("https://dixy.ru%s", card.Url)}
		discount := int64(card.PriceDiscount * 100)
		regular := int64(card.PriceRegular * 100)

		if regular == 0 {
			info.PriceDiscount, info.PriceRegular = discount, discount
		} else {
			info.PriceDiscount, info.PriceRegular = discount, regular
		}

		result = append(result, info)
	}

	return result
}
