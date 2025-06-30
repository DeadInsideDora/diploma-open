package mock

import (
	"fmt"
	"maps_service/internal/domain"
)

type ReturnData struct {
	Distance [][]int
	Duration [][]int
	E        error
}

type MockMatrix struct {
	data map[string]ReturnData
}

func NewMockMatrix(data map[string]ReturnData) *MockMatrix {
	return &MockMatrix{data: data}
}

func (mock *MockMatrix) Get(_ []domain.Point, _, _ []int, transport string) ([][]int, [][]int, error) {
	val, ok := mock.data[transport]
	if ok {
		return val.Distance, val.Duration, val.E
	}
	return nil, nil, fmt.Errorf("no matrix for transport: %s", transport)
}
