package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"scrappers/internal/domain"
)

type JsonMatcherService struct {
	readerFactory domain.IReaderFactory
	matchData     map[string]domain.MatcherCategoryInfo
}

func NewJsonMatcherService(readerFactory domain.IReaderFactory, path string) (*JsonMatcherService, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't load matcher config file: %s", err)
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var config []domain.MatcherCategoryInfo
	if err := json.Unmarshal([]byte(bytes), &config); err != nil {
		return nil, fmt.Errorf("can't unmarshal config: %s", err)
	}

	matchData := make(map[string]domain.MatcherCategoryInfo)
	for _, category := range config {
		_, ok := matchData[category.Type]
		if ok {
			log.Printf("JsonMatcher: duplicate category %s", category.Type)
		} else {
			matchData[category.Type] = category
		}
	}

	return &JsonMatcherService{readerFactory: readerFactory, matchData: matchData}, nil
}

func (service *JsonMatcherService) Match(productInfos []domain.ProductInfo, categoryName string) []domain.MatchData {
	val, ok := service.matchData[categoryName]
	if !ok {
		log.Printf("JsonMatcher: no matching data for category=%s", categoryName)
		return nil
	}
	reader, err := service.readerFactory.Get()
	if err != nil {
		log.Printf("JsonMatcher: can't get reader: %+v", err)
		return nil
	}
	previousMatch, err := service.loadPreviousMatchData(reader, categoryName)
	if err != nil {
		log.Printf("JsonMatcher: can't load previous match: %+v", err)
		return nil
	}

	storeToIdsMap := createStoreToIdsMap(productInfos)

	result := []domain.MatchData{}

	for _, product := range val.Data {
		prices := []domain.MatchPrices{}

		pictures := make(map[string]string)

		for _, store := range product.Stores {
			if _, ok := storeToIdsMap[store.StoreName]; !ok {
				continue
			}
			productInfo, ok := storeToIdsMap[store.StoreName][store.ProductId]
			if !ok {
				log.Printf("JsonMatcher: no matching info for store=%s, productId=%d", store.StoreName, store.ProductId)
				continue
			}
			if productInfo.Title != store.ProductName {
				log.Printf("JsonMatcher: product with id=%d for store=%s didn't match by name=%s (expected=%s)", store.ProductId, store.StoreName, productInfo.Title, store.ProductName)
			}
			pictures[store.StoreName] = productInfo.PictureUrl
			prices = append(prices, domain.MatchPrices{PriceDiscount: productInfo.PriceDiscount, PriceRegular: productInfo.PriceRegular, ShopName: productInfo.ShopName})
		}

		if len(prices) != 0 {
			val, ok := previousMatch[fmt.Sprintf("%s:%s", categoryName, product.Name)]
			image := getBestPic(pictures)
			if image != nil {
				log.Printf("JsonMatcher: name=%s, pic=%s", product.Name, *image)
			} else {
				log.Printf("JsonMatcher: name=%s, pic=nil", product.Name)
			}
			if !ok {
				result = append(result, domain.MatchData{Title: product.Name, Category: categoryName, Data: product.Data, Prices: prices, Image: image})
			} else {
				result = append(result, updateMatchDataPrices(val, prices, image))
			}
		}
	}

	log.Printf("JsonMatcher: category=%s, len=%d", categoryName, len(result))

	return result
}

func (service *JsonMatcherService) loadPreviousMatchData(reader domain.IReaderService, categoryName string) (map[string]domain.MatchData, error) {
	matchData, err := reader.ReadByCategory(categoryName)
	if err != nil {
		return nil, err
	}
	result := make(map[string]domain.MatchData)
	for _, data := range matchData {
		result[fmt.Sprintf("%s:%s", data.Category, data.Title)] = data
	}
	return result, nil
}

func updateMatchDataPrices(data domain.MatchData, prices []domain.MatchPrices, pic *string) domain.MatchData {
	oldPrices := make(map[string]domain.MatchPrices)
	for _, price := range data.Prices {
		oldPrices[price.ShopName] = price
	}
	for _, price := range prices {
		oldPrices[price.ShopName] = price
	}
	newPrices := []domain.MatchPrices{}
	for _, price := range oldPrices {
		newPrices = append(newPrices, price)
	}

	data.Prices = newPrices
	if pic != nil {
		data.Image = pic
	}
	return data
}

func createStoreToIdsMap(productInfos []domain.ProductInfo) map[string]map[int64]*domain.ProductInfo {
	storeToIdsMap := make(map[string]map[int64]*domain.ProductInfo)
	for _, productInfo := range productInfos {
		if _, ok := storeToIdsMap[productInfo.ShopName]; !ok {
			storeToIdsMap[productInfo.ShopName] = make(map[int64]*domain.ProductInfo)
		}
		val, _ := storeToIdsMap[productInfo.ShopName]
		if _, ok := val[productInfo.Id]; ok {
			log.Printf("JsonMatcher: duplicating id=%d for shop=%s", productInfo.Id, productInfo.ShopName)
			continue
		}
		val[productInfo.Id] = &productInfo
	}

	return storeToIdsMap
}

func getBestPic(pictures map[string]string) *string {
	if val, ok := pictures["Перекрёсток"]; ok {
		return &val
	}
	if val, ok := pictures["Лента"]; ok {
		return &val
	}
	if val, ok := pictures["Дикси"]; ok {
		return &val
	}
	return nil
}
