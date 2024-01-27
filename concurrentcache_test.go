package concurrentcache

import (
	"bytes"
	"testing"
	"time"
)

func TestConcurrentCache(t *testing.T) {
	c := NewConcurrentCache[string, []uint8](time.Second, func(locker Locker, cache map[string][]uint8) {
		locker.Lock()
		defer locker.Unlock()
		for key, value := range cache {
			if key == "key" {
				if value[0] == 4 {
					delete(cache, key)
					return
				}
				for i := range value {
					value[i] += value[i]
				}
			}
		}
	})
	defer c.Close()

	c.AccessWrite(func(cache map[string][]uint8) {
		cache["key"] = []uint8{1, 2, 3}
	})

	c.AccessRead(func(cache map[string][]uint8) {
		if value := cache["key"]; !bytes.Equal(value, []uint8{1, 2, 3}) {
			t.Error("Expected [1, 2, 3] but got", value)
		}
	})

	time.Sleep(2500 * time.Millisecond)

	c.AccessRead(func(cache map[string][]uint8) {
		if value := cache["key"]; !bytes.Equal(value, []uint8{4, 8, 12}) {
			t.Error("Expected [4, 8, 12] but got", value)
		}
	})

	time.Sleep(time.Second)

	c.AccessRead(func(cache map[string][]uint8) {
		if value := cache["key"]; value != nil {
			t.Error("Expected nil but got", value)
		}
	})
}
