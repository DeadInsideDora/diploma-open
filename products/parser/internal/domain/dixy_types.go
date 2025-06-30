package domain

type DixyCategory struct {
	Id      int          `json:"id"`
	Filters []FilterData `json:"filterData"`
}

type FilterData struct {
	FacetId string `json:"facet_id"`
	Value   string `json:"value"`
	Type    string `json:"type"`
	Id      string `json:"id"`
}
