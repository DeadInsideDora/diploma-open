package services

import (
	"fmt"
	"slices"
)

func nextPermutation(nums []int) bool {
	n := len(nums) - 1
	if n < 1 {
		return false
	}
	i := n - 1
	for ; nums[i] >= nums[i+1]; i-- {
		if i == 0 {
			return false
		}
	}
	j := n
	for nums[i] >= nums[j] {
		j--
	}
	nums[i], nums[j] = nums[j], nums[i]
	for k, j := i+1, n; k < j; {
		nums[k], nums[j] = nums[j], nums[k]
		k++
		j--
	}
	return true
}

func genSequenceExcludeValue(count, exclude int) []int {
	var permutation []int
	for i := 0; i < count; i++ {
		if i == exclude {
			continue
		}
		permutation = append(permutation, i)
	}

	return permutation
}

func updatePathCost(matrix [][]int, previousCost int, previousValid bool, from, to int) (int, bool) {
	if matrix[from][to] == -1 {
		return -1, false
	} else {
		return previousCost + matrix[from][to], previousValid
	}
}

type TSPBruteforce struct{}

type TSPDynProgramming struct{}

func NewTSPBruteforce() *TSPBruteforce {
	return &TSPBruteforce{}
}

func (tsp *TSPBruteforce) Get(matrix [][]int, startPoint int) (int, []int, error) {
	if len(matrix) <= 1 {
		return -1, nil, fmt.Errorf("incorrect points count")
	}

	if startPoint < 0 || len(matrix) <= startPoint {
		return -1, nil, fmt.Errorf("incorrect start point")
	}

	permutation := genSequenceExcludeValue(len(matrix), startPoint)

	result := -1
	var path []int = nil

	for {
		cost := 0
		valid := true

		firstPoint, lastPoint := permutation[0], permutation[len(permutation)-1]

		cost, valid = updatePathCost(matrix, cost, valid, startPoint, firstPoint)
		cost, valid = updatePathCost(matrix, cost, valid, lastPoint, startPoint)

		for i := 0; i+1 < len(permutation); i++ {
			cost, valid = updatePathCost(matrix, cost, valid, permutation[i], permutation[i+1])
		}

		if valid && (result == -1 || result > cost) {
			result = cost
			path = append([]int(nil), permutation...)
		}

		if !nextPermutation(permutation) {
			break
		}
	}

	if path == nil {
		return -1, nil, fmt.Errorf("no path")
	}

	return result, append(append([]int{startPoint}, path...), startPoint), nil
}

func NewTSPDynProgramming() *TSPDynProgramming {
	return &TSPDynProgramming{}
}

func (tsp *TSPDynProgramming) Get(matrix [][]int, startPoint int) (int, []int, error) {
	if len(matrix) <= 1 {
		return -1, nil, fmt.Errorf("incorrect points conunt")
	}

	if startPoint < 0 || len(matrix) <= startPoint {
		return -1, nil, fmt.Errorf("incorrect start point")
	}

	n := len(matrix)
	dp := createMatrix(n, 1<<n)
	path := createMatrix(n, 1<<n)

	dp[startPoint][1<<startPoint] = 0

	for mask := 0; mask < 1<<n; mask++ {
		for last := 0; last < n; last++ {
			if mask&(1<<last) == 0 {
				continue
			}
			if dp[last][mask] == -1 {
				continue
			}

			for next := 0; next < n; next++ {
				if mask&(1<<next) != 0 {
					continue
				}
				if matrix[last][next] == -1 {
					continue
				}

				cost := dp[next][mask|(1<<next)]
				updatedCost := dp[last][mask] + matrix[last][next]

				if cost == -1 || updatedCost < cost {
					dp[next][mask|(1<<next)] = updatedCost
					path[next][mask|(1<<next)] = last
				}
			}
		}
	}

	result := -1

	for i := 0; i < n; i++ {
		if dp[i][(1<<n)-1] == -1 {
			continue
		}
		if matrix[i][startPoint] == -1 {
			continue
		}
		totalCost := dp[i][(1<<n)-1] + matrix[i][startPoint]
		if result == -1 || totalCost < dp[result][(1<<n)-1]+matrix[result][startPoint] {
			result = i
		}
	}

	resultCost := dp[result][(1<<n)-1] + matrix[result][startPoint]
	resultPath := []int{startPoint}

	mask := (1 << n) - 1

	for result != -1 {
		resultPath = append(resultPath, result)
		oldResult := result
		result = path[result][mask]
		mask ^= 1 << oldResult
	}

	slices.Reverse(resultPath)

	return resultCost, resultPath, nil
}
