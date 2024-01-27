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

type ConcurrentCache[T any] interface {
	Close()
	AccessRead(callback func(cache T))
	AccessWrite(callback func(cache T))
}

type concurrentCache[T any] struct {
	sync.RWMutex
	stopChan chan struct{}
	wg       sync.WaitGroup
	cache    T
	update   func(locker Locker, cache T)
}

func NewConcurrentCache[T any](cache T, updateInterval time.Duration, update func(locker Locker, cache T)) ConcurrentCache[T] {
	c := &concurrentCache[T]{
		RWMutex:  sync.RWMutex{},
		stopChan: make(chan struct{}),
		wg:       sync.WaitGroup{},
		cache:    cache,
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

func (c *concurrentCache[T]) Close() {
	close(c.stopChan)
	c.wg.Wait()
}

func (c *concurrentCache[T]) AccessRead(callback func(cache T)) {
	c.RLock()
	defer c.RUnlock()
	callback(c.cache)
}

func (c *concurrentCache[T]) AccessWrite(callback func(cache T)) {
	c.Lock()
	defer c.Unlock()
	callback(c.cache)
}
