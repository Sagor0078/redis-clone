package cache

import (
	"sync"
	"time"
)

var (
	store       sync.Map
	expirations sync.Map 
)

func Set(key, value string) {
	store.Store(key, value)
	expirations.Delete(key) 
}

func Get(key string) (string, bool) {
	if exp, ok := expirations.Load(key); ok {
		expTime := exp.(time.Time)
		if time.Now().After(expTime) {
			store.Delete(key)
			expirations.Delete(key)
			return "", false
		}
	}
	val, ok := store.Load(key)
	if !ok {
		return "", false
	}
	return val.(string), true
}

func Delete(key string) {
	store.Delete(key)
	expirations.Delete(key)
}

func SetWithExpiration(key, value string, d time.Duration) {
	Set(key, value)
	expTime := time.Now().Add(d)
	expirations.Store(key, expTime)

	go func() {
		time.Sleep(d)
		if exp, ok := expirations.Load(key); ok {
			if exp.(time.Time).Before(time.Now()) {
				Delete(key)
			}
		}
	}()
}

func TTL(key string) time.Duration {
	_, exists := store.Load(key)
	if !exists {
		return -2
	}

	exp, ok := expirations.Load(key)
	if !ok {
		return -1 
	}

	ttl := time.Until(exp.(time.Time))
	if ttl <= 0 {
		Delete(key)
		return -2
	}
	return ttl
}

func FlushAll() {
	store = sync.Map{}
}
