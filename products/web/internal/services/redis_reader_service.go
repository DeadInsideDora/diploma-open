package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"web/internal/domain"

	"github.com/redis/go-redis/v9"
)

type redisReaderService struct {
	client *redis.Client
}

func NewRedisReaderService(addr string) *redisReaderService {
	options := redis.Options{
		Addr: addr,
		DB:   0,
	}

	return &redisReaderService{client: redis.NewClient(&options)}
}

func (service *redisReaderService) ReadByCategory(category string) ([]domain.MatchData, error) {
	ctx := context.Background()
	if err := service.client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("can't connected to redis: %s", err)
	}

	result := []domain.MatchData{}
	var cursor uint64
	for {
		keys, cur, err := service.client.Scan(ctx, cursor, fmt.Sprintf("%s:*", category), 20).Result()
		if err != nil {
			return nil, fmt.Errorf("redis scan return an error: %s", err)
		}
		cursor = cur

		if len(keys) > 0 {
			vals, err := service.client.MGet(ctx, keys...).Result()
			if err != nil {
				return nil, fmt.Errorf("can't get values from redis: %s", err)
			}

			for i, raw := range vals {
				if raw == nil {
					continue
				}

				str, ok := raw.(string)
				if !ok {
					continue
				}

				var obj domain.MatchData
				if err := json.Unmarshal([]byte(str), &obj); err != nil {
					log.Printf("can't unmarshal for key %q: %s\n", keys[i], err)
					continue
				}

				result = append(result, obj)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return result, nil
}

func (service *redisReaderService) ReadByNames(category string, names []string) ([]domain.MatchData, error) {
	ctx := context.Background()
	if err := service.client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("can't connected to redis: %s", err)
	}

	keys := getSearchKeys(category, names)
	vals, err := service.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("can't get data from redis: %s", err)
	}

	result := []domain.MatchData{}
	for i, raw := range vals {
		if raw == nil {
			continue
		}

		str, ok := raw.(string)
		if !ok {
			continue
		}

		var obj domain.MatchData
		if err := json.Unmarshal([]byte(str), &obj); err != nil {
			log.Printf("can't unmarshal for key %q: %s\n", keys[i], err)
			continue
		}

		result = append(result, obj)
	}

	return result, nil
}

func (service *redisReaderService) Close() error {
	return service.client.Close()
}

func getSearchKeys(category string, names []string) []string {
	result := []string{}
	for _, name := range names {
		result = append(result, fmt.Sprintf("%s:%s", category, name))
	}
	return result
}
