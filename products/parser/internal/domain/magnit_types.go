package domain

type MagnitEnv struct {
	cookie string
}

type MagnitCategory struct {
	Id      int            `json:"id"`
	Filters []MagnitFilter `json:"filters"`
}

type MagnitFilter struct {
	Id       string   `json:"id"`
	Elements []string `json:"elements"`
}

func NewMagnitEnv(cookie string) *MagnitEnv {
	return &MagnitEnv{cookie: cookie}
}

func (env *MagnitEnv) GetCookie() string {
	return env.cookie
}
