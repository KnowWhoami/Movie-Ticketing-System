package cache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	c := NewCache(InMemoryCache)
	if c == nil {
		t.Fatal("NewCache() returned nil")
	}
}

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache(InMemoryCache)

	c.Set("key1", "value1", time.Minute)
	val, ok := c.Get("key1")
	if !ok {
		t.Fatal("Get() returned not found after Set()")
	}
	if val.(string) != "value1" {
		t.Errorf("Get() = %v, want %v", val, "value1")
	}
}

func TestCache_GetMissing(t *testing.T) {
	c := NewCache(InMemoryCache)

	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("Get() returned found for a key that was never set")
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCache(InMemoryCache)

	c.Set("key1", "value1", time.Minute)
	c.Delete("key1")
	_, ok := c.Get("key1")
	if ok {
		t.Error("Get() returned found after Delete()")
	}
}

func TestCache_DeleteNonExistent(t *testing.T) {
	c := NewCache(InMemoryCache)
	// Should not panic
	c.Delete("does_not_exist")
}

func TestCache_Overwrite(t *testing.T) {
	c := NewCache(InMemoryCache)

	c.Set("key1", "first", time.Minute)
	c.Set("key1", "second", time.Minute)
	val, ok := c.Get("key1")
	if !ok {
		t.Fatal("Get() returned not found after overwrite Set()")
	}
	if val.(string) != "second" {
		t.Errorf("Get() = %v, want %v", val, "second")
	}
}

func TestCache_SetStruct(t *testing.T) {
	type payload struct {
		ID   int
		Name string
	}
	c := NewCache(InMemoryCache)

	p := &payload{ID: 1, Name: "test"}
	c.Set("obj", p, time.Minute)
	val, ok := c.Get("obj")
	if !ok {
		t.Fatal("Get() returned not found for struct value")
	}
	got := val.(*payload)
	if got.ID != p.ID || got.Name != p.Name {
		t.Errorf("Get() = %v, want %v", got, p)
	}
}
