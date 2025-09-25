package cache

import (
	"sync"
	"time"

	"github.com/Dekamik/farstu/internal/asserts"
)

type Cache[TObject any] interface {
	Get() (*TObject, error)
}

type cacheImpl[TObject any] struct {
	expiresAt time.Time
	mutex     sync.RWMutex
	object    *TObject
	refresh   func() (*TObject, error)
	ttl       int
}

var _ Cache[any] = &cacheImpl[any]{}

func (c *cacheImpl[TObject]) Get() (*TObject, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()

	if c.object == nil || now.After(c.expiresAt) {
		newObject, err := c.refresh()
		if err != nil {
			return nil, err
		}

		c.object = newObject
		c.expiresAt = now.Add(time.Second * time.Duration(c.ttl))
	}

	return c.object, nil
}

func New[TObject any](ttl int, refresh func() (*TObject, error)) Cache[TObject] {
	asserts.Assert(ttl >= 0, "TTL cannot be negative")

	return &cacheImpl[TObject]{
		ttl:     ttl,
		refresh: refresh,
	}
}
