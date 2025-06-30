package domain

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Login    string   `json:"login"`
	Password string   `json:"-"`
	Cards    []string `json:"cards"`
	Info     *MapInfo `json:"map_info"`
	Exchange int64    `json:"exchange"`
}

type MapInfo struct {
	Point  *Point `json:"point,omitempty"`
	Radius *int64 `json:"radius,omitempty"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
