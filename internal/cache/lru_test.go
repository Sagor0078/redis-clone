package cache

import "testing"

func TestLRUEviction(t *testing.T) {
	InitLRU(2)
	Set("a", "1")
	Set("b", "2")
	Get("a")      // "a" becomes MRU
	Set("c", "3") // should evict "b"

	if _, ok := Get("b"); ok {
		t.Errorf("Expected b to be evicted")
	}
	if _, ok := Get("a"); !ok {
		t.Errorf("Expected a to be present")
	}
}
