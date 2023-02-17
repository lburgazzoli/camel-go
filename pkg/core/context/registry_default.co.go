package context

import "sync"

func NewDefaultRegistry() Registry {
	r := defaultRegistry{
		store:     make(map[string]interface{}),
		storeLock: &sync.RWMutex{},
	}

	return &r
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
