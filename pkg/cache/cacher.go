package cache

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	ErrCacheMiss    = errors.New("key is missing")
	ErrInvalidKey   = errors.New("key is invalid")
	ErrInvalidValue = errors.New("value type is invalid")
)

type FetchFunc func() (interface{}, error)

type CloseFn func() error

type Config struct {
	Enabled bool          `json:"enabled" yaml:"enabled"`
	Prefix  string        `json:"prefix" yaml:"prefix"`
	Type    string        `json:"type" yaml:"type"`
	TTL     time.Duration `json:"ttl" yaml:"ttl"`
	Redis   RedisConfig   `json:"redis" yaml:"redis"`
}

type RedisConfig struct {
	Cluster      bool          `json:"cluster" yaml:"cluster"`
	Endpoints    []string      `json:"endpoints" yaml:"endpoints"`
	ReadTimeout  time.Duration `json:"read-timeout" yaml:"read-timeout"`
	WriteTimeout time.Duration `json:"write-timeout" yaml:"write-timeout"`
	DialTimeout  time.Duration `json:"dial-timeout" yaml:"dial-timeout"`
	PoolSize     int           `json:"pool-size" yaml:"pool-size"`
	PoolTimeout  time.Duration `json:"pool-timeout" yaml:"pool-timeout"`
	MaxConnAge   time.Duration `json:"max-conn-age" yaml:"max-conn-age"`
	IdleTimeout  time.Duration `json:"idle-timeout" yaml:"idle-timeout"`
}

//go:generate mockery --name Cacher --filename cache_mock.go
type Cacher interface {
	io.Closer
	// Fetch retrieves the item from the cache. If the item does not exist,
	// calls given FetchFunc to create a new item and sets to the cache.
	Fetch(ctx context.Context, key string, value interface{}, fetchFunc FetchFunc) error

	// Get gets an item for the given computeKey.
	Get(ctx context.Context, key string, value interface{}) error

	// Set adds an item to the cache.
	Set(ctx context.Context, key string, value interface{}) error

	// Exists returns a true if the given computeKey is exists, otherwise false.
	Exists(ctx context.Context, key string) (bool, error)

	// Delete removes an item from the cache.
	Delete(ctx context.Context, key string) error
}

func NewCacher(conf *Config) (Cacher, error) {
	if !conf.Enabled {
		return nil, nil
	}
	switch conf.Type {
	case "redis":
		return newRedisCacher(conf)
	default:
		return nil, fmt.Errorf("unknown cache type: %s", conf.Type)
	}
}
