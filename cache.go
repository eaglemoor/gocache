package gocache

import (
	"errors"
	"sync"
	"time"
)

const (
	// NoExpiration Cache without life expiration
	NoExpiration time.Duration = -1

	// DefaultExpiration Default Expiration time
	DefaultExpiration time.Duration = time.Minute * 5

	// DefaultGarbageInterval Default Garbage Clear Interval
	DefaultGarbageInterval time.Duration = time.Millisecond * 500
)

var (
	// ErrorAlreadyExist Cache this used key allready exist
	ErrorAlreadyExist = errors.New("item already exist")

	// ErrorNotFound Can't found cache with used key
	ErrorNotFound = errors.New("not found")
)

// ICache interface
type ICache interface {
	Add(key, value interface{}, d time.Duration) error
	Set(key, value interface{}, d time.Duration) error
	Update(key, value interface{}, d time.Duration) error
	Get(key ...interface{}) map[interface{}]interface{}
}

var _ ICache = &cache{}

// Cache storage
type cache struct {
	items sync.Map

	garbageInterval time.Duration
	expiration      time.Duration

	garbageTicker *time.Ticker
	stop          chan bool
}

type cacheItem struct {
	item interface{}
	exp  int64
}

func (ci cacheItem) Expired() bool {
	if ci.exp == 0 {
		return false
	}

	return time.Now().UnixNano() > ci.exp
}

func (c *cache) Add(key interface{}, value interface{}, d time.Duration) error {
	if _, ok := c.items.Load(key); ok == true {
		return ErrorAlreadyExist
	}

	return c.Set(key, value, d)
}

func (c *cache) Set(key interface{}, value interface{}, d time.Duration) error {
	var exp int64
	if d > 0 {
		exp = time.Now().Add(d).UnixNano()
	} else if d == 0 {
		exp = time.Now().Add(c.expiration).UnixNano()
	}

	c.items.Store(key, cacheItem{
		item: value,
		exp:  exp,
	})
	return nil
}

func (c *cache) Update(key interface{}, value interface{}, d time.Duration) error {
	if _, ok := c.items.Load(key); ok == false {
		return ErrorNotFound
	}

	return c.Set(key, value, d)
}

func (c *cache) Get(keys ...interface{}) map[interface{}]interface{} {
	result := make(map[interface{}]interface{}, len(keys))

	for _, key := range keys {
		if item, ok := c.items.Load(key); ok == true && item.(cacheItem).Expired() == false {
			result[key] = item.(cacheItem).item
		}
	}

	return result
}

func (c *cache) runGarbage() {
	c.stopGarbage()
	c.garbageTicker = time.NewTicker(c.garbageInterval)
	c.stop = make(chan bool, 1)

	go func() {
		for {
			select {
			case <-c.garbageTicker.C:
				c.items.Range(func(key, value interface{}) bool {
					if value.(cacheItem).Expired() == true {
						c.items.Delete(key)
					}
					return true
				})
			// Canceller channel for close gorutine
			case <-c.stop:
				return
			}
		}
	}()
}

func (c *cache) stopGarbage() {
	if c.garbageTicker != nil {
		c.garbageTicker.Stop()
		c.stop <- true
	}
}
