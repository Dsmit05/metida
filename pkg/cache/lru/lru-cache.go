package lru

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Internal cache errors
var (
	ErrKeyNotExist       = errors.New("key does not exist")
	ErrKeyAlreadyExist   = errors.New("key already exists")
	ErrExpirationInvalid = errors.New("invalid expiration")
)

// Cache is the main cache type.
type Cache struct {
	len               int        // len cache data.
	cap               int        // maximum cache capacity.
	mx                sync.Mutex // mu is the mutex variable to prevent race conditions.
	lst               *list.List // doubly linked list.
	defaultExpiration time.Duration
}

// KeyI is key Cache, your key is to implement the method Equally(KeyI)
// this method should be able to equal your structure
type KeyI interface {
	Equally(KeyI) bool
}

// ValI is value Cache
type ValI interface{}

// unit Internal cache structure
type unit struct {
	Key        KeyI
	Val        ValI
	Expiration int64
}

// NewLruCache create cache with parameters
// args:
// -cap: capacity cache; if cap <=0: cap = ∞
// -defaultExpiration: default lifetime unit; if defaultExpiration <=0 defaultExpiration: ∞
// -cleanupInterval: сache clearing interval; if cleanupInterval <=0, not auto clearing
// return:
// *Cache: Initialized Cache
func NewLruCache(cap int, defaultExpiration, cleanupInterval time.Duration) *Cache {
	lst := list.New()

	if defaultExpiration < 0 {
		defaultExpiration = 0
	}

	lruCache := &Cache{
		cap:               cap,
		mx:                sync.Mutex{},
		lst:               lst,
		defaultExpiration: defaultExpiration,
	}

	go lruCache.clearExpiredDataWithInterval(cleanupInterval)

	return lruCache
}

// Add adding unit in cache
// args:
// -key: type KeyI
// -val: value struct type ValI
// -exp: lifetime unit:
//		0: ∞
//	   -1: use defaultExpiration
// return:
// - error: key creation error
func (c *Cache) Add(key KeyI, val ValI, exp time.Duration) error {
	_, found := c.get(key)
	if found {
		return ErrKeyAlreadyExist
	}

	if exp < -1 {
		return ErrExpirationInvalid
	}

	var expiration int64

	switch exp {
	case 0:
		expiration = 0
	case -1:
		expiration = time.Now().Add(c.defaultExpiration).UnixNano()
	default:
		expiration = time.Now().Add(exp).UnixNano()
	}

	item := unit{
		Key:        key,
		Val:        val,
		Expiration: expiration,
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	if c.Cap() > 0 && c.Len() == c.Cap() {
		lruKey := c.getLRU()
		c.delete(lruKey.Key)
	}

	c.lst.PushFront(item)
	c.len++

	return nil
}

// Get return value with changing order
func (c *Cache) Get(key KeyI) (ValI, bool) {
	if c.Len() == 0 {
		return nil, false
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	val, found := c.get(key)
	if val == nil {
		return nil, found
	}

	e := c.lst.Remove(val)
	c.lst.PushFront(e) // insert the item to the top

	return val.Value.(unit).Val, found
}

// IsExist check element in the cache, without changing order.
func (c *Cache) IsExist(key KeyI) bool {
	if c.Len() == 0 {
		return false
	}

	c.mx.Lock()
	_, found := c.get(key)

	if found {
		c.mx.Unlock()
		return true
	}
	c.mx.Unlock()

	return false
}

// Clear deleting all elements
func (c *Cache) Clear() {
	c.mx.Lock()
	c.clear()
	c.mx.Unlock()
}

// Peek return value of the key without changing the order.
func (c *Cache) Peek(key KeyI) (ValI, bool) {
	if c.Len() == 0 {
		return nil, false
	}

	c.mx.Lock()
	val, found := c.get(key)
	c.mx.Unlock()

	if !found {
		return nil, found
	}

	return val.Value.(unit).Val, found
}

// Len return cache length.
func (c *Cache) Len() int {
	return c.len
}

// Cap return cache capacity.
func (c *Cache) Cap() int {
	return c.cap
}

// Replace changing the key value taking into account the order of elements.
func (c *Cache) Replace(key KeyI, val ValI) error {
	c.mx.Lock()
	e, found := c.get(key)

	if !found {
		c.mx.Unlock()
		return ErrKeyNotExist
	}

	e.Value = unit{
		Key:        key,
		Val:        val,
		Expiration: e.Value.(unit).Expiration,
	}
	c.mx.Unlock()

	return nil
}

// ClearExpiredData deleting elements data with expired lifetime.
func (c *Cache) ClearExpiredData() {
	c.mx.Lock()
	l := c.Len()
	c.mx.Unlock()

	if l == 0 {
		return
	}

	c.mx.Lock()
	now := time.Now().UnixNano()
	c.clearExpiredData(now)
	c.mx.Unlock()
}

// UpdateValue updating the lifetime and/or key value
// -exp: lifetime unit:
//		0: ∞
//	   -1: not update lifetime
// to update only the date, use val = nil
func (c *Cache) UpdateValue(key KeyI, val ValI, exp time.Duration) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	newExpTime := time.Now().Add(exp).Unix()
	_, err := c.update(key, val, newExpTime)

	return err
}

// get iterable the cache, return true if the element is found
func (c *Cache) get(key KeyI) (*list.Element, bool) {
	for e := c.lst.Front(); e != nil; e = e.Next() {
		if e.Value.(unit).Key.Equally(key) {
			return e, true
		}
	}

	return nil, false
}

// clearExpiredDataWithInterval starts clearing the cache
func (c *Cache) clearExpiredDataWithInterval(cleanupInterval time.Duration) {
	if cleanupInterval <= 0 {
		return
	}

	ticker := time.NewTicker(cleanupInterval)

	for {
		<-ticker.C
		c.ClearExpiredData()
	}
}

// delete remove data from the list, and reduce length
func (c *Cache) delete(key KeyI) {
	v, found := c.get(key)
	if !found {
		return
	}

	c.lst.Remove(v)
	c.len--
}

// getLRU get the last element of a doubly linked list, or zero if the list is empty.
func (c *Cache) getLRU() unit {
	return c.lst.Back().Value.(unit)
}

// clear iterate over a doubly linked list and remove elements.
func (c *Cache) clear() {
	var next *list.Element
	for e := c.lst.Front(); e != nil; e = next {
		next = e.Next()
		c.lst.Remove(e)
		c.len--
	}
}

// clearExpiredData clearing expired data.
func (c *Cache) clearExpiredData(now int64) {
	var next *list.Element
	for e := c.lst.Front(); e != nil; e = next {
		next = e.Next()

		if exp := e.Value.(unit).Expiration; exp != 0 && exp < now {
			c.lst.Remove(e)
			c.len--
		}
	}
}

// update changing lifetime and/or key value and move it to the top of the list.
func (c *Cache) update(key KeyI, val ValI, exp int64) (unit, error) {
	var next *list.Element
	for e := c.lst.Front(); e != nil; e = next {
		next = e.Next()

		if k := e.Value.(unit).Key; k == key {
			if val == nil {
				val = e.Value.(unit).Val
			}

			if exp == -1 {
				exp = e.Value.(unit).Expiration
			}

			c.lst.Remove(e)

			newItem := unit{
				Key:        key,
				Val:        val,
				Expiration: exp,
			}
			c.lst.PushFront(newItem)

			return newItem, nil
		}
	}

	return unit{}, ErrKeyNotExist
}
