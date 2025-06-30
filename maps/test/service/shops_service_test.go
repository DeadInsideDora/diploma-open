package service_test

import (
	"maps_service/internal/domain"
	"maps_service/internal/mock"
	"maps_service/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getShopsData() map[string]domain.ShopInfo {
	return map[string]domain.ShopInfo{
		"Лента": {
			Shop: "Лента",
			Info: []domain.Place{
				{
					Name: "Гипер Лента, гипермаркет",
					Point: domain.Point{
						Lon: 82.876501,
						Lat: 55.009161,
					},
					Schedule: domain.Schedule{
						Monday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "00:00",
									To:   "24:00",
								},
							},
						},
						Tuesday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "00:00",
									To:   "24:00",
								},
							},
						},
						Wednesday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "00:00",
									To:   "24:00",
								},
							},
						},
						Thursday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "00:00",
									To:   "24:00",
								},
							},
						},
						Friday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "00:00",
									To:   "24:00",
								},
							},
						},
						Saturday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "00:00",
									To:   "24:00",
								},
							},
						},
						Sunday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "00:00",
									To:   "24:00",
								},
							},
						},
					},
				},
				{
					Name: "Супер Лента, супермаркет",
					Point: domain.Point{
						Lon: 82.839666,
						Lat: 54.969566,
					},
					Schedule: domain.Schedule{
						Monday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "08:00",
									To:   "23:00",
								},
							},
						},
						Tuesday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "08:00",
									To:   "23:00",
								},
							},
						},
						Wednesday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "08:00",
									To:   "23:00",
								},
							},
						},
						Thursday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "08:00",
									To:   "23:00",
								},
							},
						},
						Friday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "08:00",
									To:   "23:00",
								},
							},
						},
						Saturday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "08:00",
									To:   "23:00",
								},
							},
						},
						Sunday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "08:00",
									To:   "23:00",
								},
							},
						},
					},
				},
				{
					Name: "Лента Зоомаркет",
					Point: domain.Point{
						Lon: 82.876501,
						Lat: 555.009161,
					},
					Schedule: domain.Schedule{
						Monday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "11:00",
									To:   "23:00",
								},
							},
						},
						Tuesday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "11:00",
									To:   "23:00",
								},
							},
						},
						Wednesday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "11:00",
									To:   "23:00",
								},
							},
						},
						Thursday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "11:00",
									To:   "23:00",
								},
							},
						},
						Friday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "11:00",
									To:   "23:00",
								},
							},
						},
						Saturday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "11:00",
									To:   "23:00",
								},
							},
						},
						Sunday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "11:00",
									To:   "23:00",
								},
							},
						},
					},
				},
			},
		},
		"Перекрёсток": {
			Shop: "Перекрёсток",
			Info: []domain.Place{
				{
					Name: "Перекресток, магазин-закусочная",
					Point: domain.Point{
						Lon: 82.876501,
						Lat: 555.009161,
					},
					Schedule: domain.Schedule{
						Monday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "09:00",
									To:   "21:00",
								},
							},
						},
						Tuesday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "09:00",
									To:   "21:00",
								},
							},
						},
						Wednesday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "09:00",
									To:   "21:00",
								},
							},
						},
						Thursday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "09:00",
									To:   "21:00",
								},
							},
						},
						Friday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "09:00",
									To:   "21:00",
								},
							},
						},
						Saturday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "09:00",
									To:   "21:00",
								},
							},
						},
						Sunday: domain.DaySchedule{
							WorkingHours: []domain.WorkingHours{
								{
									From: "09:00",
									To:   "21:00",
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestShopsService(t *testing.T) {
	shopsData := getShopsData()
	mock := mock.NewMockShopInfo(shopsData)

	shopsService := services.NewShopsRequester(mock)

	t.Run(
		"AvailableShops",
		func(t *testing.T) {
			assert.Equal(
				t,
				[]string{
					"Лента",
					"Перекрёсток",
				},
				shopsService.GetAvailableShops(),
				"ShopsService return incorrect available shops",
			)
		},
	)

	t.Run(
		"GetNearbyShops",
		func(t *testing.T) {
			assert.Equal(
				t,
				[]domain.ShopInfo{
					shopsData["Лента"],
					shopsData["Перекрёсток"],
				},
				shopsService.GetNearbyShops(domain.Point{}, 0),
				"ShopsService return incorrect nearby shops",
			)
		},
	)
}

func TestShopsService_NoShops(t *testing.T) {
	shopsData := make(map[string]domain.ShopInfo)
	mock := mock.NewMockShopInfo(shopsData)

	shopsService := services.NewShopsRequester(mock)

	t.Run(
		"GetNearbyShops_NoShops",
		func(t *testing.T) {
			assert.Equal(
				t,
				[]domain.ShopInfo(nil),
				shopsService.GetNearbyShops(domain.Point{}, 0),
				"ShopsService return incorrect nearby shops",
			)
		},
	)
}

func TestShopsService_NoПерекрёсток(t *testing.T) {
	shopsData := getShopsData()
	delete(shopsData, "Перекрёсток")
	mock := mock.NewMockShopInfo(shopsData)

	shopsService := services.NewShopsRequester(mock)

	t.Run(
		"GetNearbyShops_NoShops",
		func(t *testing.T) {
			assert.Equal(
				t,
				[]domain.ShopInfo{
					shopsData["Лента"],
				},
				shopsService.GetNearbyShops(domain.Point{}, 0),
				"ShopsService return incorrect nearby shops",
			)
		},
	)
}
