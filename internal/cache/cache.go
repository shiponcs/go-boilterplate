package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
)

type store struct {
	client *redis.Client
	log    *zap.Logger
	cfg    *config.Config
}

func NewStore(client *redis.Client, log *zap.Logger, cfg *config.Config) *store {
	return &store{client: client, log: log, cfg: cfg}
}

func (s *store) keyBuilder(placeholder string, values ...any) string {
	return fmt.Sprintf(prefix+placeholder, values...)
}

func (s *store) GetWidget(id uint) (*entity.Widget, error) {
	key := s.keyBuilder(keyWidget, id)
	data, err := s.client.Get(s.client.Context(), key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var w entity.Widget
	if err := json.Unmarshal([]byte(data), &w); err != nil {
		s.log.Error("error unmarshalling widget from cache", zap.Error(err))
		return nil, err
	}
	return &w, nil
}

func (s *store) SetWidget(id uint, widget *entity.Widget, ttl time.Duration) error {
	key := s.keyBuilder(keyWidget, id)
	bytes, err := json.Marshal(widget)
	if err != nil {
		s.log.Error("error marshalling widget for cache", zap.Error(err))
		return err
	}
	if err := s.client.Set(s.client.Context(), key, bytes, ttl).Err(); err != nil {
		s.log.Warn("error setting widget in cache", zap.Error(err))
		return err
	}
	return nil
}

func (s *store) ForgetWidget(id uint) error {
	key := s.keyBuilder(keyWidget, id)
	return s.client.Del(s.client.Context(), key).Err()
}
