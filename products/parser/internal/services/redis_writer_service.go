package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"scrappers/internal/domain"

	"github.com/redis/go-redis/v9"
)

type redisWriterService struct {
	client *redis.Client
}

type RedisWriterFactory struct {
	addr string
}

func newRedisWriterService(addr string) *redisWriterService {
	options := redis.Options{
		Addr: addr,
		DB:   0,
	}

	return &redisWriterService{client: redis.NewClient(&options)}
}

func (service *redisWriterService) Write(matchData []domain.MatchData) error {
	log.Printf("RedisWriterService: writing %d", len(matchData))
	ctx := context.Background()
	if err := service.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("can't connected to redis: %s", err)
	}

	for _, match := range matchData {
		data, err := json.Marshal(match)
		if err != nil {
			return fmt.Errorf("can't marshal data to string: %s", err)
		}

		log.Printf("%s:%s : %s", match.Category, match.Title, data)

		if err := service.client.Set(ctx, fmt.Sprintf("%s:%s", match.Category, match.Title), data, 0).Err(); err != nil {
			return fmt.Errorf("can't write data to redis: %s", err)
		}
	}

	return nil
}

func (service *redisWriterService) Close() error {
	return service.client.Close()
}

func NewRedisWriterFactory(addr string) *RedisWriterFactory {
	return &RedisWriterFactory{addr: addr}
}

func (factory *RedisWriterFactory) Get() (domain.IWriterService, error) {
	return newRedisWriterService(factory.addr), nil
}
