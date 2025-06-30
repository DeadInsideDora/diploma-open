package service_test

import (
	"maps_service/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTSPService_EqualBruteforceAndDynProgramming(t *testing.T) {
	matrix := [][]int{
		{0, 1, 4},
		{8, 0, 3},
		{2, 5, 0},
	}
	tspDynProgramming := services.NewTSPDynProgramming()
	tspBruteforce := services.NewTSPBruteforce()

	t.Run(
		"Equal answers for tsp bruteforce and tsp dyn programming",
		func(t *testing.T) {
			for startPoint := 0; startPoint < 3; startPoint += 1 {
				bruteforceCost, bruteforcePath, bruteforceError := tspBruteforce.Get(matrix, startPoint)
				dynProgrammingCost, dynProgrammingPath, dynProgrammingError := tspDynProgramming.Get(matrix, startPoint)

				assert.NoError(t, bruteforceError)
				assert.NoError(t, dynProgrammingError)

				assert.Equal(t, bruteforceCost, dynProgrammingCost)
				assert.Equal(t, bruteforcePath, dynProgrammingPath)
			}
		},
	)
}
