package domain

type IShopInfoService interface {
	Get(shop string, point Point, radius int64) (*ShopInfo, error)
}

type IMatrixService interface {
	Get(points []Point, sources, targets []int, transport string) ([][]int, [][]int, error)
}

type ITSPService interface {
	Get(matrix [][]int, startPoint int) (int, []int, error)
}

type IRoutingService interface {
	Get(points []Point, startPoint int, byDistance bool) []MinTimeRoute
}
