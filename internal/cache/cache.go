package cache

import "time"

type Cache[TObject any] interface {
	GetOrRefresh() (*TObject, error)
}

type cacheImpl[TObject any] struct {
	expiresAt   time.Time
	object      *TObject
	refresh     func() (*TObject, error)
	ttl         int
}

var _ Cache[any] = cacheImpl[any]{}

func (c cacheImpl[TObject]) GetOrRefresh() (*TObject, error) {
	now := time.Now()

	if c.object == nil || c.expiresAt.Before(now) {
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
	return cacheImpl[TObject]{
		ttl:         ttl,
		refresh:     refresh,
	}
}
