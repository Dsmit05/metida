package set

import "sync"

// CacheSet implements set with mutex.
type CacheSet struct {
	mutex *sync.RWMutex
	set   map[string]struct{}
}

// NewCacheSet return new instance of cacheSet.
func NewCacheSet(sizeSet int) *CacheSet {
	return &CacheSet{
		mutex: &sync.RWMutex{},
		set:   make(map[string]struct{}, sizeSet),
	}
}

// IsExist return true if key in map, if not have key return false and add this key
func (o *CacheSet) IsExist(s string) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	_, ok := o.set[s]
	if !ok {
		o.set[s] = struct{}{}
	}

	return ok
}
