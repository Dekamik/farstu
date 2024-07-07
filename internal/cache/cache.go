package cache

import "time"

type Cache[TObject any, TRefreshArgs any] interface {
	GetOrRefresh() (*TObject, error)
}

type cacheImpl[TObject any, TRefreshArgs any] struct {
	expiresAt   time.Time
	object      *TObject
	refresh     func(TRefreshArgs) (*TObject, error)
	refreshArgs TRefreshArgs
	ttl         int
}

var _ Cache[any, any] = cacheImpl[any, any]{}

func (c cacheImpl[TObject, TRefreshArgs]) GetOrRefresh() (*TObject, error) {
	now := time.Now()

	if c.object == nil || c.expiresAt.Before(now) {
		newObject, err := c.refresh(c.refreshArgs)
		if err != nil {
			return nil, err
		}

		c.object = newObject
		c.expiresAt = now.Add(time.Second * time.Duration(c.ttl))
	}

	return c.object, nil
}

func New[TObject any, TRefreshArgs any](ttl int, refresh func(TRefreshArgs) (*TObject, error), args TRefreshArgs) Cache[TObject, TRefreshArgs] {
	return cacheImpl[TObject, TRefreshArgs]{
		ttl:         ttl,
		refresh:     refresh,
		refreshArgs: args,
	}
}
