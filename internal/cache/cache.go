package cache

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"
)

type Cache struct {
	pool *cache.Cache
}

func NewCache() *Cache {
	c := cache.New(cache.NoExpiration, 0)

	return &Cache{pool: c}
}

func (c *Cache) Add(id string, content json.RawMessage) {
	c.pool.SetDefault(id, content)
}

func (c *Cache) Get(id string) (json.RawMessage, bool) {
	content, ok := c.pool.Get(id)
	if !ok {
		return nil, ok
	}

	rawContent, ok := content.(json.RawMessage)
	if !ok {
		return nil, false
	}

	return rawContent, ok
}

func (c *Cache) Fill(ids []string, contents []json.RawMessage) {
	for i := range contents {
		c.pool.SetDefault(ids[i], contents[i])
	}
}
