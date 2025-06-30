package domain

type PerekrestokEnv struct {
	auth   string
	cookie string
}

type PerekrestokCategory struct {
	Id        int       `json:"id"`
	Features  []Feature `json:"feature"`
	Blacklist []int     `json:"blacklist"`
}

type Feature struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewPerekrestokEnv(auth, cookie string) *PerekrestokEnv {
	return &PerekrestokEnv{auth: auth, cookie: cookie}
}

func (env *PerekrestokEnv) GetAuth() string {
	return env.auth
}

func (env *PerekrestokEnv) GetCookie() string {
	return env.cookie
}
