package domain

type IDatabaseService interface {
	GetUserById(int64) (*User, error)
	GetUserByLogin(string) (*User, error)
	CreateUser(string, string, string) error
	UpdateCards(int, []string) error
	UpdateMapInfo(int, MapInfo) error
	UpdateExchange(int, int) error
	Close()
}
