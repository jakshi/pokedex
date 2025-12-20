package pokecache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	defer cache.Stop()

	if cache == nil {
		t.Fatal("expected cache to not be nil")
	}
	if cache.cache == nil {
		t.Fatal("expected cache.cache map to be initialized")
	}
}

func TestAddGet(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	defer cache.Stop()

	key := "https://pokeapi.co/api/v2/location"
	data := []byte(`{"results": []}`)

	cache.Add(key, data)

	got, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected to find key")
	}
	if string(got) != string(data) {
		t.Fatalf("expected %s, got %s", data, got)
	}
}

func TestGetMissing(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	defer cache.Stop()

	_, ok := cache.Get("nonexistent")
	if ok {
		t.Fatal("expected ok to be false for missing key")
	}
}

func TestDelete(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	defer cache.Stop()

	key := "mykey"
	cache.Add(key, []byte("value"))
	cache.Delete(key)

	_, ok := cache.Get(key)
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestReap(t *testing.T) {
	interval := 50 * time.Millisecond
	cache := NewCache(interval)
	defer cache.Stop()

	cache.Add("key", []byte("value"))

	// Wait for reap to run (at least one interval + buffer)
	time.Sleep(interval * 3)

	_, ok := cache.Get("key")
	if ok {
		t.Fatal("expected entry to be reaped after interval")
	}
}