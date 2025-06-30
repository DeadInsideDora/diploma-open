package service_test

import (
	"maps_service/internal/domain"
	"maps_service/internal/mock"
	"maps_service/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutingService_WithErrorsMatrix(t *testing.T) {
	mock := mock.NewMockMatrix(map[string]mock.ReturnData{})
	tspBruteforce := services.NewTSPBruteforce()

	service := services.NewRoutingService(mock, tspBruteforce)

	t.Run(
		"Matrix return errors all time",
		func(t *testing.T) {
			assert.Equal(
				t,
				[]domain.MinTimeRoute{},
				service.Get(
					[]domain.Point{
						{
							Lon: 82.876501,
							Lat: 55.009161,
						},
						{
							Lon: 82.839666,
							Lat: 54.969566,
						},
					},
					0,
					false,
				),
				"routing service should return empty array without errors",
			)
		},
	)
}

func TestRoutingService_Correct(t *testing.T) {
	mock := mock.NewMockMatrix(map[string]mock.ReturnData{
		"walking": {
			Distance: [][]int{
				{0, 1, 4},
				{8, 0, 3},
				{2, 5, 0},
			},
			Duration: [][]int{
				{0, 8, 2},
				{9, 0, 1},
				{7, 7, 0},
			},
			E: nil,
		},
	})
	tspBruteforce := services.NewTSPBruteforce()

	service := services.NewRoutingService(mock, tspBruteforce)

	t.Run(
		"Correct Routing by distance",
		func(t *testing.T) {
			assert.Equal(
				t,
				[]domain.MinTimeRoute{
					{
						Points: []domain.Point{
							{
								Lon: 82.876501,
								Lat: 55.009161,
							},
							{
								Lon: 82.839666,
								Lat: 54.969566,
							},
							{
								Lon: 82.876501,
								Lat: 555.009161,
							},
							{
								Lon: 82.876501,
								Lat: 55.009161,
							},
						},
						Duration:  6,
						Transport: "walking",
					},
				},
				service.Get(
					[]domain.Point{
						{
							Lon: 82.876501,
							Lat: 55.009161,
						},
						{
							Lon: 82.839666,
							Lat: 54.969566,
						},
						{
							Lon: 82.876501,
							Lat: 555.009161,
						},
					},
					0,
					true,
				),
				"routing service should return correct data for walking",
			)
		},
	)

	t.Run(
		"Correct Routing by duration",
		func(t *testing.T) {
			assert.Equal(
				t,
				[]domain.MinTimeRoute{
					{
						Points: []domain.Point{
							{
								Lon: 82.876501,
								Lat: 55.009161,
							},
							{
								Lon: 82.839666,
								Lat: 54.969566,
							},
							{
								Lon: 82.876501,
								Lat: 55.009161,
							},
							{
								Lon: 82.876501,
								Lat: 55.009161,
							},
						},
						Duration:  16,
						Transport: "walking",
					},
				},
				service.Get(
					[]domain.Point{
						{
							Lon: 82.876501,
							Lat: 55.009161,
						},
						{
							Lon: 82.839666,
							Lat: 54.969566,
						},
						{
							Lon: 82.876501,
							Lat: 55.009161,
						},
					},
					0,
					false,
				),
				"routing service should return correct data for walking",
			)
		},
	)
}
