package services

import (
	"encoding/json"
	"fmt"
	"log"
	"optimizer/internal/domain"
)

const USER_POINT_ID string = "USER_POINT_ID_UNIQUE_DATA_FOR_MAPPING"

func getProductInfo(matchPrice domain.MatchPrices, shopName string, discounts map[string]struct{}) productInfo {
	price := matchPrice.PriceRegular
	if _, ok := discounts[shopName]; ok {
		price = matchPrice.PriceDiscount
	}
	return productInfo{Price: &price, StoreName: &shopName}
}

func groupProductsByCategory(products []domain.InputProductInfo) map[string][]string {
	result := make(map[string][]string)

	for _, product := range products {
		_, ok := result[product.Info.Type]
		if !ok {
			result[product.Info.Type] = []string{}
		}

		result[product.Info.Type] = append(result[product.Info.Type], product.Info.Name)
	}

	return result
}

func createDiscountsMap(discounts []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, storeName := range discounts {
		result[storeName] = struct{}{}
	}

	return result
}

func createIdToShop(shopInfos []domain.ShopInfo) map[string]extendedPlace {
	result := make(map[string]extendedPlace)
	for _, shopInfo := range shopInfos {
		for _, shop := range shopInfo.Info {
			_, ok := result[shop.Id]
			if ok {
				log.Printf("Duplicating shop with id=%s", shop.Id)
				continue
			}
			result[shop.Id] = extendedPlace{ShopInfo: shop, ShopName: shopInfo.Shop}
		}
	}
	return result
}

func createDurMatrix(mapsService domain.IMapsService, shopInfos []domain.ShopInfo, userPoint domain.Point) (map[string]map[string]int, error) {
	result := make(map[string]map[string]int)
	for i := 0; i < len(shopInfos); i += 1 {
		for j := 0; j < len(shopInfos); j += 1 {
			if i == j {
				continue
			}
			source := collectPoints(shopInfos[i])
			targets := collectPoints(shopInfos[j])

			routes, err := mapsService.GetRoutesBetweenAddresses(source, targets, "walking")
			if err != nil {
				return nil, err
			}

			for _, route := range routes {
				if _, ok := result[shopInfos[i].Info[route.From].Id]; !ok {
					result[shopInfos[i].Info[route.From].Id] = make(map[string]int)
				}
				for _, to := range route.Routes {
					result[shopInfos[i].Info[route.From].Id][shopInfos[j].Info[to.To].Id] = to.Time
				}
			}
		}
	}

	getDurMatrixFromUserPoint(result, mapsService, shopInfos, userPoint)

	return result, nil
}

func getDurMatrixFromUserPoint(result map[string]map[string]int, mapsService domain.IMapsService, shopInfos []domain.ShopInfo, userPoint domain.Point) error {
	for _, shopInfo := range shopInfos {
		points := collectPoints(shopInfo)

		routes, err := mapsService.GetRoutesBetweenAddresses([]domain.Point{userPoint}, points, "walking")
		if err != nil {
			return err
		}
		for _, route := range routes {
			if _, ok := result[USER_POINT_ID]; !ok {
				result[USER_POINT_ID] = make(map[string]int)
			}
			for _, to := range route.Routes {
				result[USER_POINT_ID][shopInfo.Info[to.To].Id] = to.Time
			}
		}

		routes, err = mapsService.GetRoutesBetweenAddresses(points, []domain.Point{userPoint}, "walking")
		if err != nil {
			return err
		}
		for _, route := range routes {
			if _, ok := result[shopInfo.Info[route.From].Id]; !ok {
				result[shopInfo.Info[route.From].Id] = make(map[string]int)
			}
			for _, to := range route.Routes {
				result[shopInfo.Info[route.From].Id][USER_POINT_ID] = to.Time
			}
		}
	}
	return nil
}

func getAmounts(products []domain.InputProductInfo) []int64 {
	result := []int64{}
	for _, product := range products {
		result = append(result, product.Amount)
	}
	return result
}

func collectProducts(products []domain.InputProductInfo, productsService domain.IProductsService) ([]domain.MatchData, error) {
	result := make([]domain.MatchData, len(products))
	productsToIndex := make(map[string]int)
	for i, product := range products {
		productsToIndex[product.Info.Type+":"+product.Info.Name] = i
	}

	for key, values := range groupProductsByCategory(products) {
		data, err := productsService.GetProducts(key, values)
		if err != nil {
			return nil, fmt.Errorf("can't collect products: %s", err)
		}
		for _, match := range data {
			result[productsToIndex[match.Category+":"+match.Title]] = match
		}
	}

	return result, nil
}

func logMatchData(matchData []domain.MatchData) {
	bytes, _ := json.Marshal(matchData)
	log.Printf("MatchData: %s", string(bytes))
}
