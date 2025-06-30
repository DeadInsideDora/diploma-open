package domain

type Environments struct {
	lenta       *LentaEnv
	perekrestok *PerekrestokEnv
	magnit      *MagnitEnv
}

type Config struct {
	Categories []ParseCategory `json:"categories"`
	DelayInfo  Delay           `json:"delayInfo"`
}

type ParseCategory struct {
	Type        string                `json:"type"`
	Lenta       []LentaCategory       `json:"lenta"`
	Perekrestok []PerekrestokCategory `json:"perekrestok"`
	Dixy        []DixyCategory        `json:"dixy"`
	Magnit      []MagnitCategory      `json:"magnit"`
}

type Delay struct {
	PingDelay       int `json:"pingDelay"`
	NewRequestDelay int `json:"newRequestDelay"`
}

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

type MatcherCategoryInfo struct {
	Type string `json:"type"`
	Data []MatcherProductInfo
}

type MatcherProductInfo struct {
	Name   string             `json:"name"`
	Data   []MasterData       `json:"master_data"`
	Stores []ProductStoreInfo `json:"stores"`
}

type ProductStoreInfo struct {
	StoreName   string `json:"store_name"`
	ProductName string `json:"product_name"`
	ProductId   int64  `json:"product_id"`
}

type MasterData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProductInfo struct {
	Id            int64
	Title         string
	PriceDiscount int64
	PriceRegular  int64
	PictureUrl    string
	ShopName      string
}

func NewEnvironments(lenta *LentaEnv, perekrestok *PerekrestokEnv, magnit *MagnitEnv) *Environments {
	return &Environments{lenta: lenta, perekrestok: perekrestok, magnit: magnit}
}

func (env *Environments) GetLentaEnv() *LentaEnv {
	return env.lenta
}

func (env *Environments) GetPerekrestokEnv() *PerekrestokEnv {
	return env.perekrestok
}

func (env *Environments) GetMagnitEnv() *MagnitEnv {
	return env.magnit
}
