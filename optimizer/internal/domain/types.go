package domain

type Point struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type MatchData struct {
	Title    string        `json:"title"`
	Category string        `json:"category"`
	Prices   []MatchPrices `json:"prices"`
}

type MatchPrices struct {
	PriceDiscount int64  `json:"price_discount"`
	PriceRegular  int64  `json:"price_regular"`
	ShopName      string `json:"shop_name"`
}

type MasterData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ShopInfo struct {
	Info []Place `json:"info"`
	Shop string  `json:"shop"`
}

type Place struct {
	Name  string `json:"name"`
	Point Point  `json:"point"`
	Id    string `json:"id"`
}

type RoutesInfo struct {
	From   int     `json:"from"`
	Routes []Route `json:"routes"`
}

type Route struct {
	To   int `json:"to"`
	Time int `json:"time"`
}

type ProductInfo struct {
	Type 	string 	`json:"type"`
	Name 	string 	`json:"name"`
	Url	 	string 	`json:"url"`
	Weighed bool	`json:"isWeighed"`
}

type InputProductInfo struct {
	Info   ProductInfo `json:"info"`
	Amount int64       `json:"amount"`
}

type ProductInfoWithAmount struct {
	Type 	string 	`json:"type"`
	Name 	string 	`json:"name"`
	Url	 	string 	`json:"url"`
	Weighed bool	`json:"isWeighed"`
	Amount	int64   `json:"amount"`
}

type OutputProductInfo struct {
	Info  ProductInfoWithAmount	`json:"info"`
	Price int64					`json:"price"`
}

type OptimizerResult struct {
	Stores     []StoreInfo `json:"stores"`
	TotalPrice int64       `json:"price"`
	Cost       int64       `json:"cost"`
}

type StoreInfo struct {
	Products   []OutputProductInfo `json:"products"`
	Store      string              `json:"store"`
	StorePoint Point               `json:"point"`
	Price      int64               `json:"total_price"`
}

type MinTimeRoute struct {
	Points    []int  `json:"points"`
	Duration  int    `json:"duration"`
	Transport string `json:"transport"`
}
