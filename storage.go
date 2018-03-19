package gocache

import (
	"sync"
	"time"
)

var (
	cacheList map[string]ICache
	cacheLock sync.RWMutex
)

func init() {
	cacheList = make(map[string]ICache)
}

// NewCache Create new cache instance
func NewCache(name string, garbageInterval, expiration time.Duration) (ICache, error) {
	cexist, err := GetCache(name)
	if err != nil {
		return nil, err
	}
	if cexist != nil {
		return nil, ErrorAlreadyExist
	}

	c := &cache{
		garbageInterval: DefaultGarbageInterval,
		expiration:      DefaultExpiration,
	}

	if garbageInterval > 0 {
		c.garbageInterval = garbageInterval
	}

	if expiration > 0 {
		c.expiration = expiration
	}

	cacheLock.Lock()
	cacheList[name] = c
	cacheLock.Unlock()

	c.runGarbage()

	return c, nil
}

// GetCache Cache instance
func GetCache(name string) (ICache, error) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()

	return cacheList[name], nil
}

// CacheList Cache names
func CacheList() []string {
	keys := make([]string, len(cacheList))
	cacheLock.RLock()
	defer cacheLock.RUnlock()

	for key := range cacheList {
		keys = append(keys, key)
	}

	return keys
}

// DeleteCache Delete cache with all data
func DeleteCache(name string) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if c, ok := cacheList[name]; ok {
		c.(*cache).stopGarbage()
	}
	delete(cacheList, name)
}
