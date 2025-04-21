package cache

import (
	"sync"
	"time"
)

var store sync.Map

func Set(key, value string) {
	store.Store(key, value)
}

func Get(key string) (string, bool) {
	val, ok := store.Load(key)
	if !ok {
		return "", false
	}
	return val.(string), true
}

func Delete(key string) {
	store.Delete(key)
}

func SetWithExpiration(key, value string, d time.Duration) {
	Set(key, value)
	go func() {
		time.Sleep(d)
		Delete(key)
	}()
}
