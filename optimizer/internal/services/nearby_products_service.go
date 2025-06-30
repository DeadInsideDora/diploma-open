package services

import (
	"fmt"
	"log"

	"optimizer/internal/domain"

	"github.com/oleiade/lane/v2"
)

type NearbyProductsService struct {
	mapsService     domain.IMapsService
	productsService domain.IProductsService
}

type storeInfo struct {
	Prices []*int64
	Id     string
	Dur    int64
}

func NewNearbyProductsService(mapsService domain.IMapsService, productsService domain.IProductsService) *NearbyProductsService {
	return &NearbyProductsService{mapsService: mapsService, productsService: productsService}
}

func (service *NearbyProductsService) Get(products []domain.InputProductInfo, discounts []string, userPoint domain.Point, radius int64, exchange int64) (*domain.OptimizerResult, error) {
	discountsMap := createDiscountsMap(discounts)

	matchData, err := collectProducts(products, service.productsService)
	if err != nil {
		return nil, err
	}
	shopInfos, err := service.mapsService.GetNearShops(userPoint, radius)
	if err != nil {
		return nil, err
	}
	durMatrix, err := createDurMatrix(service.mapsService, shopInfos, userPoint)
	if err != nil {
		return nil, err
	}
	idToShop := createIdToShop(shopInfos)

	var result *storeInfo = nil

	pq := lane.NewMinPriorityQueue[storeInfo, string]()
	pq.Push(storeInfo{Prices: make([]*int64, len(matchData)), Id: USER_POINT_ID, Dur: 0}, "0:"+USER_POINT_ID)

	for !pq.Empty() {
		top, _, ok := pq.Pop()

		log.Printf("State: %s %d %+v", top.Id, top.Dur, top.Prices)
		if !ok {
			continue
		}
		if _, ok := durMatrix[top.Id]; !ok {
			log.Printf("no routes from %s", top.Id)
			continue
		}
		for to, d := range durMatrix[top.Id] {
			info := storeInfo{Prices: top.Prices, Id: to}
			var shopName string
			if val, ok := idToShop[to]; ok {
				shopName = val.ShopName
			} else {
				continue
			}
			log.Printf("Possible edge from %s to %s (shopName=%s)", top.Id, to, shopName)

			for i, data := range matchData {
				for _, price := range data.Prices {
					if price.ShopName == shopName {
						productInfo := getProductInfo(price, shopName, discountsMap)
						if productInfo.Price != nil && (info.Prices[i] == nil || *info.Prices[i] > *productInfo.Price) {
							info.Prices[i] = productInfo.Price
						}
						continue
					}
				}
			}

			count := 0
			for _, price := range info.Prices {
				if price != nil {
					count += 1
				}
			}

			if count == len(matchData) {
				result = &info
				break
			}
			dur := top.Dur + int64(d)
			info.Dur = dur
			pq.Push(info, fmt.Sprintf("%d:%s", dur, info.Id))
		}
		if result != nil {
			break
		}
	}

	if result == nil {
		return nil, fmt.Errorf("can't collect all products in nearby shops")
	}

	if _, ok := durMatrix[result.Id]; !ok {
		return nil, fmt.Errorf("no route from %s", result.Id)
	}
	val := durMatrix[result.Id]
	if _, ok := val[USER_POINT_ID]; !ok {
		return nil, fmt.Errorf("no route from %s to userPoint", result.Id)
	}
	result.Dur += int64(val[USER_POINT_ID])

	return getNearbyProductsResult(products, result, exchange), nil
}

func getNearbyProductsResult(products []domain.InputProductInfo, store *storeInfo, exchange int64) *domain.OptimizerResult {
	var totalPrice int64 = 0
	var cost int64 = 0

	for i, price := range store.Prices {
		totalPrice += *price * products[i].Amount
	}
	cost = totalPrice + int64((float64(store.Dur)/6.)*float64(exchange))
	return &domain.OptimizerResult{TotalPrice: totalPrice, Cost: cost}
}
