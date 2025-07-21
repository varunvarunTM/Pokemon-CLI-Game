package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntries map[string]cacheEntry
	mu sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) *Cache {

	var ca Cache
	ca.cacheEntries = make(map[string]cacheEntry)
	ca.interval = interval

	go ca.reapLoop()

	return &ca

}

func (c *Cache) Add( key string , val []byte ) {
	currentTime := time.Now()
	c.mu.Lock()

	entry := cacheEntry{
		createdAt: currentTime,
		val: val,
	}

	c.cacheEntries[key] = entry

	c.mu.Unlock()
}

func (c *Cache) Get( key string ) ( []byte , bool ) {
	c.mu.Lock()
	entry,ok := c.cacheEntries[key]
	if !ok {
		c.mu.Unlock()
		return nil , false
	}
	c.mu.Unlock()
	return entry.val, true
}

func (c *Cache) reapLoop () {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	tickerChan := ticker.C
	for {
		select {
		case <- tickerChan:
			c.mu.Lock()
			keysToDelete := []string{}
			for key,value:= range c.cacheEntries {
				if time.Since(value.createdAt)  > c.interval {
					keysToDelete = append(keysToDelete,key)
				}
			}
			for _,key := range keysToDelete{
				delete(c.cacheEntries,key)
			}
			c.mu.Unlock()
		}
	}
}