package concurrentcache

import (
	"sync"
	"time"
)

type Locker interface {
	RLock()
	TryRLock() bool
	RUnlock()
	Lock()
	TryLock() bool
	Unlock()
	RLocker() sync.Locker
}

type ConcurrentCache[KeyType comparable, ValueType any] interface {
	Close()
	AccessRead(callback func(cache map[KeyType]ValueType))
	AccessWrite(callback func(cache map[KeyType]ValueType))
}

type concurrentCache[KeyType comparable, ValueType any] struct {
	sync.RWMutex
	stopChan chan struct{}
	wg       sync.WaitGroup
	cache    map[KeyType]ValueType
	update   func(locker Locker, cache map[KeyType]ValueType)
}

func NewConcurrentCache[KeyType comparable, ValueType any](updateInterval time.Duration, update func(locker Locker, cache map[KeyType]ValueType)) ConcurrentCache[KeyType, ValueType] {
	c := &concurrentCache[KeyType, ValueType]{
		RWMutex:  sync.RWMutex{},
		stopChan: make(chan struct{}),
		wg:       sync.WaitGroup{},
		cache:    make(map[KeyType]ValueType),
		update:   update,
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(updateInterval)
		for {
			select {
			case <-ticker.C:
				c.update(c, c.cache)
			case <-c.stopChan:
				return
			}
		}
	}()

	return c
}

func (c *concurrentCache[KeyType, ValueType]) Close() {
	close(c.stopChan)
	c.wg.Wait()
}

func (c *concurrentCache[KeyType, ValueType]) AccessRead(callback func(cache map[KeyType]ValueType)) {
	c.RLock()
	defer c.RUnlock()
	callback(c.cache)
}

func (c *concurrentCache[KeyType, ValueType]) AccessWrite(callback func(cache map[KeyType]ValueType)) {
	c.Lock()
	defer c.Unlock()
	callback(c.cache)
}
