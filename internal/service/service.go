package service

import (
	"encoding/json"
	"order-service/internal/cache"
)

type Service struct {
	cache *cache.Cache
}

func NewService(c *cache.Cache) *Service {
	return &Service{cache: c}
}

func (s *Service) GetOrder(id string) (json.RawMessage, bool) {
	return s.cache.Get(id)
}
