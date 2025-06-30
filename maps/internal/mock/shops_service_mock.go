package mock

import (
	"fmt"
	"maps_service/internal/domain"
)

type MockShopInfo struct {
	shopsInfo map[string]domain.ShopInfo
}

func NewMockShopInfo(shopsInfo map[string]domain.ShopInfo) *MockShopInfo {
	return &MockShopInfo{shopsInfo: shopsInfo}
}

func (mock *MockShopInfo) Get(shop string, _ domain.Point, _ int64) (*domain.ShopInfo, error) {
	val, ok := mock.shopsInfo[shop]
	if ok {
		return &val, nil
	}
	return nil, fmt.Errorf("no such shop: %s", shop)
}
