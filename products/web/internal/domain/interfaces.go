package domain

type IReaderService interface {
	ReadByCategory(string) ([]MatchData, error)
	ReadByNames(string, []string) ([]MatchData, error)
	Close() error
}

type ICategoriesService interface {
	Get(bool) ([]Category, error)
}

type IMetricsService interface {
	Inc(string, int)
}
