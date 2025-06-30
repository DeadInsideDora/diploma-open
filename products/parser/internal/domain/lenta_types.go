package domain

type LentaEnv struct {
	deviceId     string
	sessionToken string
}

type LentaCategory struct {
	Id              int             `json:"id"`
	Multicheckboxes []Multicheckbox `json:"multicheckbox"`
	Blacklist       []int           `json:"blacklist"`
}

type Multicheckbox struct {
	Key      string `json:"key"`
	ValueIds []int  `json:"valueIds"`
}

func NewLentaEnv(deviceId, sessionToken string) *LentaEnv {
	return &LentaEnv{deviceId: deviceId, sessionToken: sessionToken}
}

func (env *LentaEnv) GetDeviceId() string {
	return env.deviceId
}

func (env *LentaEnv) GetSessionToken() string {
	return env.sessionToken
}
