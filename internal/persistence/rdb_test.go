package persistence

import (
	"os"
	"testing"

	"github.com/Sagor0078/redis-clone/internal/cache"
)

func TestSaveAndLoad(t *testing.T) {
	_ = os.Remove(rdbFile)

	cache.Set("foo", "bar")
	cache.Set("hello", "world")

	Save()

	cache.Delete("foo")
	cache.Delete("hello")

	if _, ok := cache.Get("foo"); ok {
		t.Fatal("Expected foo to be deleted from cache")
	}

	Load()

	if val, ok := cache.Get("foo"); !ok || val != "bar" {
		t.Errorf("Expected foo=bar after load, got %v", val)
	}
	if val, ok := cache.Get("hello"); !ok || val != "world" {
		t.Errorf("Expected hello=world after load, got %v", val)
	}

	_ = os.Remove(rdbFile)
}
