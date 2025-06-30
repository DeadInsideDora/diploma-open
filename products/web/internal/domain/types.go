package domain

type MatchData struct {
	Title    string        `json:"title"`
	Category string        `json:"category"`
	Image    *string       `json:"image"`
	Data     []MasterData  `json:"master_data"`
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

type Category struct {
	Type    string   `json:"type"`
	Filters []Filter `json:"filters"`
}

type Filter struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
