package cache

import (
	"container/list"
	"sync"
)

type entry struct {
	key   string
	value string
}

type LRUCache struct {
	capacity int
	mu       sync.Mutex
	ll       *list.List
	cache    map[string]*list.Element
}

var lru *LRUCache

func InitLRU(cap int) {
	lru = &LRUCache{
		capacity: cap,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func Set(key, value string) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if el, ok := lru.cache[key]; ok {
		lru.ll.MoveToFront(el)
		el.Value.(*entry).value = value
		return
	}

	if lru.ll.Len() >= lru.capacity {
		back := lru.ll.Back()
		if back != nil {
			evicted := back.Value.(*entry)
			delete(lru.cache, evicted.key)
			lru.ll.Remove(back)
		}
	}

	e := &entry{key, value}
	el := lru.ll.PushFront(e)
	lru.cache[key] = el
}

func Get(key string) (string, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if el, ok := lru.cache[key]; ok {
		lru.ll.MoveToFront(el)
		return el.Value.(*entry).value, true
	}
	return "", false
}

func Delete(key string) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if el, ok := lru.cache[key]; ok {
		delete(lru.cache, key)
		lru.ll.Remove(el)
	}
}

func FlushAll() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.ll.Init()
	lru.cache = make(map[string]*list.Element)
}
