package main

import (
	"fmt"
	"time"

	"github.com/Dsmit05/metida/pkg/cache/lru"
)

// ValS is value lru Cache
type ValS []string

// KeyStrict is key struct lru Cache
type KeyStrict struct {
	ProfileID int
	SiteID    int
}

func (o KeyStrict) Equally(k lru.KeyI) bool {
	return o == k
}

func main() {
	setKey := make([]KeyStrict, 0)

	cache := lru.NewLruCache(15, time.Second*5, time.Second*6)

	for i := 0; i < 20; i++ {
		setKey = append(setKey, KeyStrict{ProfileID: i, SiteID: i})
		err := cache.Add(setKey[i], ValS{"Val", "NextVal"}, -1)

		if err != nil {
			fmt.Println(err)
		}
	}

	val, bl := cache.Get(setKey[1])
	fmt.Println(val, bl, cache.Len(), setKey[18])
}
