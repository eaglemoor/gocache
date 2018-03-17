package gocache

import (
	"fmt"
	"testing"
	"time"
)

func TestCache_AddGetExpired(t *testing.T) {
	cache := NewCache("TestCache_Add", time.Millisecond*500, 0)

	for i := 0; i < 1000; i++ {
		go cache.Add(fmt.Sprint("TestCache_Add_", i), i, time.Second*5)
	}

	key100 := "TestCache_Add_100"
	add100, ok := cache.Get(key100)[key100]
	if !ok {
		t.Errorf("Can't find TestCache_Add_100")
		t.Fail()
	}
	if add100.(int) != 100 {
		t.Errorf("%s %q != %d", key100, add100, 100)
	}

	time.Sleep(time.Second * 2)
	_, ok = cache.Get(key100)[key100]
	if !ok {
		t.Errorf("%s not exist after 2sec", key100)
	}

	time.Sleep(time.Second * 5)

	_, ok = cache.Get(key100)[key100]
	if ok {
		t.Errorf("%s exist after 7sec", key100)
	}
}

func TestCache_Garbage(t *testing.T) {
	//cache500 := NewCache("TestCache_500", time.Millisecond*500, 0)
}
