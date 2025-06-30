package services

import (
	"log"
	"maps_service/internal/domain"
)

var TRANSPORT_TYPES = []string{"walking", "taxi", "driving"}

func genSimpleSequence(n int) []int {
	if n < 1 {
		return nil
	}
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = i
	}
	return result
}

type RoutingService struct {
	matrixService domain.IMatrixService
	tspService    domain.ITSPService
}

func NewRoutingService(matrixService domain.IMatrixService, tspService domain.ITSPService) *RoutingService {
	return &RoutingService{matrixService: matrixService, tspService: tspService}
}

func (routing *RoutingService) Get(points []domain.Point, startPoint int, byDistance bool) []domain.MinTimeRoute {
	result := []domain.MinTimeRoute{}

	for _, transport := range TRANSPORT_TYPES {
		sequence := genSimpleSequence(len(points))
		distMatrix, durMatrix, err := routing.matrixService.Get(points, sequence, sequence, transport)

		if err != nil {
			log.Println(err, transport)
			continue
		}

		dur, path, err := routing.tspService.Get(
			func() [][]int {
				if byDistance {
					return distMatrix
				}
				return durMatrix
			}(),
			startPoint,
		)

		if err != nil {
			log.Println(err, transport)
			continue
		}

		result = append(result, domain.MinTimeRoute{Points: path, Duration: dur, Transport: transport})
	}

	return result
}
