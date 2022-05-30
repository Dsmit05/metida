package lru

import "testing"

//go test -bench=. -benchmem -benchtime=5x
const lenCache = 10000

// ValS is value lru Cache
type ValS []string

// KeyStrict is key struct lru Cache
type KeyStrict struct {
	ProfileID int
	SiteID    int
}

func (o KeyStrict) Equally(k KeyI) bool {
	return o == k
}

func BenchmarkLruCacheAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache := NewLruCache(lenCache, 0, 0)
		for i := 0; i < lenCache; i++ {
			key := KeyStrict{ProfileID: i, SiteID: i}
			val := ValS{"H", "I"}
			b.StartTimer()
			cache.Add(key, val, 0)
			b.StopTimer()
		}
	}
}

func BenchmarkLruCacheGet(b *testing.B) {
	cache := NewLruCache(0, 0, 0)

	for i := 0; i < lenCache; i++ {
		key := KeyStrict{ProfileID: i, SiteID: i}
		val := ValS{"H", "I"}
		cache.Add(key, val, 0)
	}

	for i := 0; i < b.N; i++ {
		for i := 0; i < lenCache; i++ {
			key := KeyStrict{ProfileID: i, SiteID: i}
			b.StartTimer()
			cache.Get(key)
			b.StopTimer()
		}
	}
}

func BenchmarkLruCacheExist(b *testing.B) {
	cache := NewLruCache(0, 0, 0)

	for i := 0; i < lenCache; i++ {
		key := KeyStrict{ProfileID: i, SiteID: i}
		val := ValS{"H", "I"}
		cache.Add(key, val, 0)
	}

	for i := 0; i < b.N; i++ {
		for i := 0; i < lenCache; i++ {
			key := KeyStrict{ProfileID: i + lenCache/2, SiteID: i + lenCache/2}
			b.StartTimer()
			cache.IsExist(key)
			b.StopTimer()
		}
	}
}
