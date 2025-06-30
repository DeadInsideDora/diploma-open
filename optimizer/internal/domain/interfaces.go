package domain

type IMapsService interface {
	GetNearShops(Point, int64) ([]ShopInfo, error)
	GetRoutesBetweenAddresses([]Point, []Point, string) ([]RoutesInfo, error)
	GetTSP([]Point, int) (*MinTimeRoute, error)
}

type IProductsService interface {
	GetProducts(string, []string) ([]MatchData, error)
}

type IOptimizerService interface {
	Get([]InputProductInfo, []string, Point, int64, int64) (*OptimizerResult, error)
}

type INearbyProductsService interface {
	Get([]InputProductInfo, []string, Point, int64, int64) (*OptimizerResult, error)
}
