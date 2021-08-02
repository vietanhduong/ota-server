package redis

import (
	"encoding/json"
	"time"
)

func (c *Client) StoreWithTTL(key string, value interface{}, ttl time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(key, p, ttl).Err()
}

func (c *Client) Store(key string, value interface{}) error {
	return c.StoreWithTTL(key, value, 0)
}

func (c *Client) GetValue(key string, dest interface{}) error {
	p, err := c.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(p), dest)
}
