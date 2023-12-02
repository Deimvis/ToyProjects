package quiet_hn

import (
	"sync"
	"time"
)

const defaultTTL = 60 // seconds

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
}

type valueWithTTL[V any] struct {
	Value          V
	ExpirationTime time.Time
}

type CacheWithTTL[K comparable, V any] struct {
	data    map[K]valueWithTTL[V]
	ttl     int64 // in seconds
	rwmutex *sync.RWMutex
}

func NewCacheWithTTL[K comparable, V any]() CacheWithTTL[K, V] {
	return CacheWithTTL[K, V]{data: make(map[K]valueWithTTL[V]), ttl: defaultTTL, rwmutex: &sync.RWMutex{}}
}

func (c CacheWithTTL[K, V]) Get(key K) (V, bool) {
	c.rwmutex.RLock()
	defer c.rwmutex.RUnlock()
	v, ok := c.data[key]
	if !ok || time.Since(v.ExpirationTime) > 0 {
		var emptyValue V
		return emptyValue, false
	}
	return v.Value, true
}

func (c CacheWithTTL[K, V]) Put(key K, value V) {
	c.rwmutex.Lock()
	defer c.rwmutex.Unlock()
	c.data[key] = valueWithTTL[V]{
		Value:          value,
		ExpirationTime: time.Now().Add(time.Duration(c.ttl) * time.Second),
	}
}
