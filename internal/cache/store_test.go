package cache

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	Set("foo", "bar")

	val, ok := Get("foo")
	if !ok {
		t.Errorf("Expected key 'foo' to be found")
	}
	if val != "bar" {
		t.Errorf("Expected value 'bar', got '%s'", val)
	}
}

func TestDelete(t *testing.T) {
	Set("deleteMe", "bye")
	Delete("deleteMe")

	_, ok := Get("deleteMe")
	if ok {
		t.Errorf("Expected key 'deleteMe' to be deleted")
	}
}

func TestSetWithExpiration(t *testing.T) {
	SetWithExpiration("temp", "value", 100*time.Millisecond)

	val, ok := Get("temp")
	if !ok || val != "value" {
		t.Errorf("Expected key 'temp' to be set")
	}

	time.Sleep(150 * time.Millisecond)

	_, ok = Get("temp")
	if ok {
		t.Errorf("Expected key 'temp' to expire and be deleted")
	}
}
