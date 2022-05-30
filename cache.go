package cache

import (
	"sync"
	"time"
)

type Data struct {
	value    string
	deadline *time.Time
}

type Cache struct {
	mut  sync.Mutex
	data map[string]Data
}

func NewCache() Cache {
	return Cache{
		data: map[string]Data{},
	}
}

func (cach *Cache) Get(key string) (string, bool) {
	cach.mut.Lock()
	defer cach.mut.Unlock()

	if _, ok := cach.data[key]; !ok {
		return "", false
	}

	if cach.data[key].deadline != nil && cach.data[key].deadline.Before(time.Now()) {
		return "", false
	}

	return cach.data[key].value, true
}

func (cach *Cache) Put(key, value string) {
	cach.mut.Lock()
	defer cach.mut.Unlock()

	cach.data[key] = Data{value, nil}
}

func (cach *Cache) Keys() []string {
	cach.mut.Lock()
	defer cach.mut.Unlock()

	var keys []string
	now := time.Now()

	for k, v := range cach.data {
		if v.deadline != nil && v.deadline.Before(now) {
			continue
		}

		keys = append(keys, k)
	}

	return keys
}

func (cach *Cache) PutTill(key, value string, deadline time.Time) {
	cach.mut.Lock()
	defer cach.mut.Unlock()

	cach.data[key] = Data{value, &deadline}
}
