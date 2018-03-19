package gocache

import (
	"errors"
	"testing"
	"time"
)

var (
	name1 = "storage-cache-1"
	name2 = "storage-cache-2"
)

func TestNewCache(t *testing.T) {
	c, err := NewCache(name1, 0, 0)
	if err == ErrorAlreadyExist {
		c, err = GetCache(name1)
	}
	if err != nil {
		t.Error(err)
	}

	// Check double cache
	_, err = NewCache(name1, DefaultGarbageInterval, DefaultExpiration)
	if err != ErrorAlreadyExist {
		if err == nil {
			err = errors.New("new cache create double cache with name " + name1)
		}
		t.Error(err)
	}

	ccache, ok := c.(*cache)
	if !ok {
		t.Errorf("cache wrong type %T", c)
	}

	if ccache.expiration != DefaultExpiration {
		t.Errorf("cache expiration %d != default %d", ccache.expiration, DefaultExpiration)
	}

	if ccache.garbageInterval != DefaultGarbageInterval {
		t.Errorf("cache garbage interval %d != default %d", ccache.garbageInterval, DefaultGarbageInterval)
	}

	gInterval := time.Millisecond * 100
	exp := time.Second
	c2, err := NewCache(name2, gInterval, exp)
	if err != nil {
		t.Error(err)
	}

	ccache2 := c2.(*cache)
	if ccache2.expiration != exp {
		t.Errorf("cache expiration %d != %d", ccache2.expiration, exp)
	}

	if ccache2.garbageInterval != gInterval {
		t.Errorf("cache garbage interval %d != %d", ccache2.garbageInterval, gInterval)
	}
}

func TestCacheList(t *testing.T) {
	// Check and create cache for multi and alone test
	c1, _ := GetCache(name1)
	if c1 == nil {
		NewCache(name1, 0, 0)
	}

	// Check and create cache for multi and alone test
	c2, _ := GetCache(name2)
	if c2 == nil {
		NewCache(name2, 0, 0)
	}

	list := CacheList()
	var fname1, fname2 bool
	for _, name := range list {
		if name == name1 {
			fname1 = true
		}
		if name == name2 {
			fname2 = true
		}
	}

	if !fname1 {
		t.Error(name1 + " not found")
	}
	if !fname2 {
		t.Error(name1 + " not found")
	}
}

func TestDeleteCache(t *testing.T) {
	c, _ := GetCache(name1)
	if c == nil {
		c, _ = NewCache(name1, 0, 0)
	}

	c.Add("test", 10, 0)

	DeleteCache(name1)

	if c, _ := GetCache(name1); c != nil {
		t.Errorf("cache %s not delete", name1)
	}
}
