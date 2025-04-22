package persistence

import (
	"encoding/gob"
	"os"
	"time"

	"github.com/Sagor0078/redis-clone/internal/cache"
)

const rdbFile = "dump.rdb"

func SavePeriodically(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			Save()
		}
	}()
}

func Save() {
	m := make(map[string]string)
	cache.Range(func(k, v string) {
		m[k] = v
	})
	file, err := os.Create(rdbFile)
	if err != nil {
		return
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	_ = enc.Encode(m)
}

func Load() {
	file, err := os.Open(rdbFile)
	if err != nil {
		return
	}
	defer file.Close()

	m := make(map[string]string)
	dec := gob.NewDecoder(file)
	err = dec.Decode(&m)
	if err != nil {
		return
	}
	for k, v := range m {
		cache.Set(k, v)
	}
}
