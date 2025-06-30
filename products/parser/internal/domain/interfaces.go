package domain

type IEnvService interface {
	Get() (*Environments, error)
}

type IConfigService interface {
	Get() (*Config, error)
}

type IWriterService interface {
	Write([]MatchData) error
	Close() error
}

type IReaderService interface {
	ReadByCategory(string) ([]MatchData, error)
	Close() error
}

type IMatcherService interface {
	Match([]ProductInfo, string) []MatchData
}

type IDomainMetricsService interface {
	Inc(string, int)
}

type IWriterFactory interface {
	Get() (IWriterService, error)
}

type IReaderFactory interface {
	Get() (IReaderService, error)
}
