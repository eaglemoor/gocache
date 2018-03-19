package gocache

import (
	"errors"
	"testing"
	"time"
)

var (
	name1 = "cache-1"
	name2 = "cache-2"
)

func TestNewCache(t *testing.T) {
	c, err := NewCache(name1, 0, 0)
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
