# Concurrent Cache - Concurrent-Safe Cache for Go

Concurrent Cache is a Go package providing a concurrent-safe cache with read and write access.

## Features

- Safe concurrent read and write access to a cache.
- Periodic updates to the cache with user-defined logic.

## Installation

```bash
go get -u github.com/NIR3X/concurrentcache
```

## Usage

```go
package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/NIR3X/concurrentcache"
)

func main() {
	// Create a new concurrent cache with a one-second update interval.
	c := concurrentcache.NewConcurrentCache[map[string][]uint8](make(map[string][]uint8), time.Second, func(locker concurrentcache.Locker, cache map[string][]uint8) {
		locker.Lock()
		defer locker.Unlock()
		// Your custom update logic here
	})

	// Close the cache when done to stop the update goroutine.
	defer c.Close()

	// AccessWrite example
	c.AccessWrite(func(cache map[string][]uint8) {
		cache["key"] = []uint8{1, 2, 3}
	})

	// AccessRead example
	c.AccessRead(func(cache map[string][]uint8) {
		if value := cache["key"]; !bytes.Equal(value, []uint8{1, 2, 3}) {
			fmt.Println("Expected [1, 2, 3] but got", value)
		}
	})
}
```

## License

[![GNU AGPLv3 Image](https://www.gnu.org/graphics/agplv3-155x51.png)](https://www.gnu.org/licenses/agpl-3.0.html)

This program is Free Software: You can use, study share and improve it at your
will. Specifically you can redistribute and/or modify it under the terms of the
[GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.html) as
published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
