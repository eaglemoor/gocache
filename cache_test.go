package gocache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCache_AddGetExpired(t *testing.T) {
	c, err := NewCache(t.Name(), time.Millisecond*500, 0)
	if err != nil {
		t.Fatal(err)
	}

	var w sync.WaitGroup

	for i := 0; i < 1000; i++ {
		// ignore error
		w.Add(1)
		go func(i int) {
			c.Add(fmt.Sprintf("TestCache_Add_%d", i), i, time.Second*5)
			w.Done()
		}(i)
	}
	w.Wait()

	key100 := "TestCache_Add_100"
	add100 := c.Get(key100)
	if len(add100) < 1 || add100[0] == nil {
		t.Errorf("Can't find %s", key100)
	}
	if add100[0].(int) != 100 {
		t.Errorf("%s %+v != %d", key100, add100[0], 100)
	}

	time.Sleep(time.Second * 2)
	add100 = c.Get(key100)
	if len(add100) < 1 || add100[0] == nil {
		t.Errorf("%s not exist after 2sec", key100)
	}

	time.Sleep(time.Second * 5)

	add100 = c.Get(key100)
	if len(add100) < 1 && add100[0] != nil {
		t.Errorf("%s exist after 7sec", key100)
	}
}

func TestCache_Add(t *testing.T) {
	c, err := NewCache(t.Name(), time.Millisecond*500, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Cache work
	err = c.Add("test-1", 1, time.Second*2)
	if err != nil {
		t.Error("Add test-1")
	}

	// More data
	for i := 2; i < 100; i++ {
		if err = c.Add(fmt.Sprint("test-", i), i, time.Second*5); err != nil {
			t.Error(err)
		}
	}

	// Check duplicate error
	err = c.Add("test-1", 1, NoExpiration)
	if err == nil {
		t.Error("test-1 not duplicate")
	}

	// Wait garbage cache "test-1"
	time.Sleep(time.Second * 3)

	err = c.Add("test-1", 1, NoExpiration)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 2)
	test1 := c.Get("test-1")
	if len(test1) < 1 {
		t.Error("test-1 not found")
	} else {
		if test1[0].(int) != 1 {
			t.Errorf("%+v != 1", test1[0])
		}
	}
}

func TestCache_Update(t *testing.T) {
	c, err := NewCache(t.Name(), time.Millisecond*500, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	tkey := "test-1"

	// Cache work
	err = c.Add(tkey, 1, time.Second)
	if err != nil {
		t.Error(err)
	}

	// Update
	err = c.Update(tkey, 10, NoExpiration)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 2)

	v := c.Get(tkey)
	if len(v) < 1 {
		t.Error(tkey + " not found")
	} else {
		if v[0].(int) != 10 {
			t.Errorf("%s wrong value %v != %d", tkey, v[0], 10)
		}
	}

	// Fake update
	err = c.Update("test-2", 12, NoExpiration)
	if err == nil {
		t.Error("Update unknown record")
	}
}

func TestCache_Set(t *testing.T) {
	c := cache{
		garbageInterval: time.Millisecond * 500,
		expiration:      DefaultExpiration,
	}
	c.runGarbage()

	tkey := "test-1"

	var err error

	if err = c.Add(tkey, 10, time.Hour); err != nil {
		t.Error(err)
	}

	if err = c.Set(tkey, 15, NoExpiration); err != nil {
		t.Error(err)
	}

	item, exist := c.items.Load(tkey)
	if !exist {
		t.Error(tkey + " not found")
	}

	if citem, ok := item.(cacheItem); !ok {
		t.Errorf("item is not cacheItem %T", item)
	} else {
		if citem.exp != 0 {
			t.Errorf("Wrong expired time after set = %v", citem.exp)
		}
	}
}

func TestCache_StopGarbage(t *testing.T) {
	c := cache{
		garbageInterval: time.Millisecond * 500,
		expiration:      DefaultExpiration,
	}
	c.runGarbage()

	c.Add("test", 10, 0)

	time.Sleep(time.Second * 1)

	if c.garbageTicker == nil {
		t.Fatal("garbage not start")
	}

	// test double run
	c.runGarbage()
	time.Sleep(time.Second * 2)

	// test double stop
	c.stopGarbage()
	time.Sleep(time.Second)
	c.stopGarbage()
}
