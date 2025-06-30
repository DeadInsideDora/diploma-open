package services

import (
	"encoding/json"
	"fmt"
	"log"
	"optimizer/internal/domain"
	"slices"
)

type OptimizerService struct {
	mapsService     domain.IMapsService
	productsService domain.IProductsService
}

type addressInfo struct {
	Products          []productInfo `json:"products"`
	Duration          int           `json:"duration"`
	TotalPrice        int64         `json:"totalPrice"`
	PricesFilled      int           `json:"pricedFilled"`
	Id                string        `json:"id"`
	StoreName         string        `json:"store"`
	OriginalStoreName string        `json:"origstore"`
	StorePoint        domain.Point  `json:"point"`
	CachedPrice       *int64        `json:"cached"`
	Visited           map[string]struct{}
	Previous          *addressInfo `json:"prev"`
}

type productInfo struct {
	Price     *int64  `json:"price"`
	StoreName *string `json:"store"`
}

type extendedPlace struct {
	ShopInfo domain.Place
	ShopName string
}

type optimizedStores struct {
	Products []productInfo
	Stores   []finalStoreInfo
}

type finalStoreInfo struct {
	StorePoint        domain.Point
	StoreName         string
	OriginalStoreName string
}

func NewOptimizerService(mapsService domain.IMapsService, productsService domain.IProductsService) *OptimizerService {
	return &OptimizerService{mapsService: mapsService, productsService: productsService}
}

func (service *OptimizerService) Get(products []domain.InputProductInfo, discounts []string, userPoint domain.Point, radius int64, exchange int64) (*domain.OptimizerResult, error) {
	discountsMap := createDiscountsMap(discounts)
	amounts := getAmounts(products)

	matchData, err := collectProducts(products, service.productsService)
	if err != nil {
		return nil, err
	}

	{
		bytes, _ := json.Marshal(matchData)
		log.Println(string(bytes))
	}

	shopInfos, err := service.mapsService.GetNearShops(userPoint, radius)
	if err != nil {
		return nil, err
	}

	{
		bytes, _ := json.Marshal(shopInfos)
		log.Println(string(bytes))
	}

	durMatrix, err := createDurMatrix(service.mapsService, shopInfos, userPoint)
	if err != nil {
		return nil, err
	}
	idToShop := createIdToShop(shopInfos)

	states := fordBellman(durMatrix, idToShop, matchData, discountsMap, amounts, exchange)

	best := getBestState(states, durMatrix, exchange)

	if best == nil {
		return nil, fmt.Errorf("can't get optimal route")
	}

	stores := optimizeStores(best)
	mtr, err := service.getTSP(stores, userPoint)
	if err != nil {
		return nil, err
	}

	return collectResult(products, stores, mtr, exchange)
}

func (service *OptimizerService) getTSP(stores optimizedStores, userPoint domain.Point) (*domain.MinTimeRoute, error) {
	points := []domain.Point{userPoint}
	for _, store := range stores.Stores {
		points = append(points, store.StorePoint)
	}

	mtr, err := service.mapsService.GetTSP(points, 0)
	if err != nil {
		return nil, fmt.Errorf("can't get tsp between points: %+v", err)
	}
	if len(mtr.Points) != len(stores.Stores)+2 {
		return nil, fmt.Errorf("unexpected len of tsp points: actual=%d, expected=%d", len(mtr.Points), len(stores.Stores)+2)
	}

	return mtr, nil
}

func getBestState(states map[string]*addressInfo, durMatrix map[string]map[string]int, exchange int64) *addressInfo {
	var info *addressInfo = nil
	var bestStateByFilled *addressInfo = nil
	for id, value := range states {
		if id == USER_POINT_ID {
			continue
		}
		if bestStateByFilled == nil || value.PricesFilled > bestStateByFilled.PricesFilled {
			bestStateByFilled = value
		}
		if value.PricesFilled == len(value.Products) {
			val, ok := durMatrix[id][USER_POINT_ID]
			if !ok {
				log.Printf("warning: no route between shop with id=%s and user point", value.Id)
				continue
			}
			value.Duration += val
			value.CachedPrice = calculatePriceWithExchange(value, exchange)
			if value.CachedPrice != nil && (info == nil || info.CachedPrice == nil || info.CachedPrice != nil && (*value.CachedPrice < *info.CachedPrice || *value.CachedPrice == *info.CachedPrice && value.TotalPrice < info.TotalPrice)) {
				info = value
			}
		}
	}

	if bestStateByFilled == nil {
		log.Println("NO BEST STATE BY PRICES FILLED")
	} else {
		bytes, _ := json.Marshal(*bestStateByFilled)
		log.Printf("BEST STATE BY PRICES FILLED: %s", string(bytes))
	}

	return info
}

func collectPoints(shopInfo domain.ShopInfo) []domain.Point {
	result := []domain.Point{}
	for _, place := range shopInfo.Info {
		result = append(result, place.Point)
	}
	return result
}

func fordBellman(durMatrix map[string]map[string]int, idToShop map[string]extendedPlace, matchData []domain.MatchData, discountsMap map[string]struct{}, amounts []int64, exchange int64) map[string]*addressInfo {
	states := make(map[string]*addressInfo)
	states[USER_POINT_ID] = nil

	for {
		updated := false
		for from, value := range durMatrix {
			cur, ok := states[from]
			if !ok {
				continue
			}

			for to, dur := range value {
				if to == USER_POINT_ID {
					continue
				}
				log.Printf("Possible edge %s -> %s (%d sec)", from, to, dur)
				val, ok := idToShop[to]
				if !ok {
					log.Printf("Warning: no such shop with id=%s", to)
					continue
				}
				if cur != nil {
					if _, ok := cur.Visited[val.ShopName]; ok {

						continue
					}
				}

				addressInfo := buildNewAddressInfo(cur, dur, matchData, discountsMap, val.ShopInfo, val.ShopName, amounts, exchange)
				old, ok := states[to]

				if !ok {
					log.Printf("Updated state for %s (1st statement)", to)
					states[to] = &addressInfo
					updated = true
				} else {
					if addressInfo.CachedPrice == nil && old.CachedPrice == nil && (addressInfo.PricesFilled > old.PricesFilled || addressInfo.PricesFilled == old.PricesFilled && addressInfo.Duration < old.Duration) {
						log.Printf("Updated state for %s (2nd statement)", to)
						states[to] = &addressInfo
						updated = true
					}
					if addressInfo.CachedPrice != nil && old.CachedPrice == nil {
						log.Printf("Updated state for %s (3rd statement)", to)
						states[to] = &addressInfo
						updated = true
					}
					if old.CachedPrice != nil && addressInfo.CachedPrice != nil {
						if *old.CachedPrice > *addressInfo.CachedPrice {
							log.Printf("Updated state for %s (4th statement)", to)
							states[to] = &addressInfo
							updated = true
						}
						if *old.CachedPrice == *addressInfo.CachedPrice && (exchange <= 2000 && addressInfo.TotalPrice < old.TotalPrice || exchange > 2000 && addressInfo.Duration < old.Duration) {
							log.Printf("Updated state for %s (5th statement)", to)
							states[to] = &addressInfo
							updated = true
						}
					}
				}
			}
		}

		if !updated {
			break
		}
	}

	return states
}

func buildNewAddressInfo(info *addressInfo, dur int, matchData []domain.MatchData, discounts map[string]struct{}, shop domain.Place, shopName string, amounts []int64, exchange int64) addressInfo {
	var products []productInfo
	var duration int
	if info == nil {
		duration = 0
		products = make([]productInfo, len(matchData))
	} else {
		duration = info.Duration
		products = info.Products
	}
	result := addressInfo{Duration: duration + dur, Id: shop.Id, StorePoint: shop.Point, OriginalStoreName: shop.Name, StoreName: shopName, Visited: make(map[string]struct{}), Previous: info, PricesFilled: 0, Products: []productInfo{}, TotalPrice: 0}
	for i, product := range matchData {
		productInfo := products[i]
		for _, price := range product.Prices {
			if price.ShopName != shopName {
				continue
			}
			newInfo := getProductInfo(price, shopName, discounts)
			if newInfo.Price != nil && (productInfo.Price == nil || *productInfo.Price > *newInfo.Price) {
				productInfo = newInfo
			}
		}
		if productInfo.Price != nil {
			result.TotalPrice += *productInfo.Price * amounts[i]
			result.PricesFilled += 1
		}
		result.Products = append(result.Products, productInfo)
	}
	result.CachedPrice = calculatePriceWithExchange(&result, exchange)
	if info != nil {
		for key := range info.Visited {
			result.Visited[key] = struct{}{}
		}
	}
	result.Visited[shopName] = struct{}{}
	return result
}

func calculatePriceWithExchange(info *addressInfo, exchange int64) *int64 {
	if info.PricesFilled != len(info.Products) {
		return nil
	}
	summary := info.TotalPrice + int64((float64(info.Duration)/6.)*float64(exchange))
	return &summary
}

func optimizeStores(info *addressInfo) optimizedStores {
	finalVisitedStores := make(map[string]struct{})
	for _, product := range info.Products {
		finalVisitedStores[*product.StoreName] = struct{}{}
	}

	result := optimizedStores{Products: info.Products}
	cur := info

	for cur != nil {
		_, ok := finalVisitedStores[cur.StoreName]
		if ok {
			result.Stores = append(result.Stores, finalStoreInfo{StorePoint: cur.StorePoint, StoreName: cur.StoreName, OriginalStoreName: cur.OriginalStoreName})
		}
		cur = cur.Previous
	}
	slices.Reverse(result.Stores)

	return result
}

func collectResult(products []domain.InputProductInfo, stores optimizedStores, mtr *domain.MinTimeRoute, exchange int64) (*domain.OptimizerResult, error) {
	result := domain.OptimizerResult{}
	points := mtr.Points[1 : len(mtr.Points)-1]
	if len(points) != len(stores.Stores) {
		return nil, fmt.Errorf("unexpected points len")
	}

	for _, storeId := range points {
		store := stores.Stores[storeId-1]
		info := domain.StoreInfo{Products: []domain.OutputProductInfo{}, Store: store.OriginalStoreName, StorePoint: store.StorePoint, Price: 0}
		for i, data := range stores.Products {
			if data.StoreName != nil && *data.StoreName == store.StoreName {
				price := *data.Price * products[i].Amount
				info.Price += price
				result.TotalPrice += price
				info.Products = append(info.Products, domain.OutputProductInfo{Info: createProductInfoWithAmount(&products[i]), Price: price})
			}
		}
		result.Stores = append(result.Stores, info)
	}

	result.Cost = result.TotalPrice + int64((float64(mtr.Duration)/6.)*float64(exchange))
	return &result, nil
}

func createProductInfoWithAmount(product *domain.InputProductInfo) domain.ProductInfoWithAmount {
	return domain.ProductInfoWithAmount{Type: product.Info.Type, Name: product.Info.Name, Url: product.Info.Url, Weighed: product.Info.Weighed, Amount: product.Amount}
}
