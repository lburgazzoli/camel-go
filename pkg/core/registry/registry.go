package registry

import (
	"sync"

	"github.com/lburgazzoli/camel-go/pkg/api"
)

func NewDefaultRegistry() (api.Registry, error) {
	r := defaultRegistry{
		store:     make(map[string]interface{}),
		storeLock: &sync.RWMutex{},
	}

	return &r, nil
}

type defaultRegistry struct {
	store     map[string]interface{}
	storeLock *sync.RWMutex
}

func (r *defaultRegistry) Set(key string, val interface{}) {
	r.storeLock.Lock()
	defer r.storeLock.Unlock()

	r.store[key] = val
}

func (r *defaultRegistry) Get(key string) (interface{}, bool) {
	r.storeLock.RLock()
	defer r.storeLock.RUnlock()

	answer, ok := r.store[key]

	return answer, ok
}

func (r *defaultRegistry) Del(key string) interface{} {
	r.storeLock.RLock()
	defer r.storeLock.RUnlock()

	answer := r.store[key]

	delete(r.store, key)

	return answer
}

func GetAs[T any](r api.Registry, key string) (T, bool) {
	v1, ok := r.Get(key)
	if !ok {
		var result T
		return result, false
	}

	v2, ok := v1.(T)
	if !ok {
		var result T
		return result, false
	}

	return v2, ok
}
