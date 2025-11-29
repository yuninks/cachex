package cachex_test

import (
	"testing"
	"time"

	"github.com/yuninks/cachex"
)

func TestCacheSetGet(t *testing.T) {
	c := cachex.NewCache()
	defer c.Close()

	c.Set("test", "test", 100*time.Millisecond)

	da, err := c.Get("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if da != "test" {
		t.Errorf("Expected 'test', got %v", da)
	}
}

func TestCacheExpiration(t *testing.T) {
	c := cachex.NewCache()
	defer c.Close()

	c.Set("test", "test", 100*time.Millisecond)
	time.Sleep(150 * time.Millisecond)
	da, err := c.Get("test")
	if err == nil {
		t.Error("Expected error for expired cache")
	}
	if da != nil {
		t.Errorf("Expected nil for expired cache, got %v", da)
	}
}

func TestCacheDelete(t *testing.T) {
	c := cachex.NewCache()
	defer c.Close()

	c.Set("test", "value", time.Hour)
	c.Delete("test")
	_, err := c.Get("test")
	if err == nil {
		t.Error("Expected error after delete")
	}
}

func TestCacheClear(t *testing.T) {
	c := cachex.NewCache()
	defer c.Close()

	c.Set("key1", "value1", time.Hour)
	c.Set("key2", "value2", time.Hour)
	c.Clear()

	_, err1 := c.Get("key1")
	_, err2 := c.Get("key2")
	if err1 == nil || err2 == nil {
		t.Error("Expected errors after clear")
	}
}

func TestCacheDefaultExpiration(t *testing.T) {
	c := cachex.NewCache()
	defer c.Close()

	c.Set("test", "value", 0)
	da, err := c.Get("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if da != "value" {
		t.Errorf("Expected 'value', got %v", da)
	}
}

func TestCacheNonExistentKey(t *testing.T) {
	c := cachex.NewCache()
	defer c.Close()

	_, err := c.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent key")
	}
}
