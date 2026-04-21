# Generic Sync Map
[![Go Reference](https://pkg.go.dev/badge/github.com/min0625/gsyncmap.svg)](https://pkg.go.dev/github.com/min0625/gsyncmap)

A generic, type-safe wrapper around [`sync.Map`](https://pkg.go.dev/sync#Map).

## Installation
```sh
go get github.com/min0625/gsyncmap
```

## Types

### `Map[Key comparable, Value any]`

A generic concurrent map that accepts any value type. Suitable for most use cases.

> **Note:** `CompareAndDelete` and `CompareAndSwap` are available on `Map` but deprecated.
> They may panic at runtime if `Value` is not a comparable type (e.g. slice, map, func).
> Use `ComparableMap` instead when these operations are needed.

### `ComparableMap[Key, Value comparable]`

A drop-in replacement for `Map` that requires the value type to be comparable.
Provides safe `CompareAndDelete` and `CompareAndSwap` without risk of runtime panic.

## Quick start

```go
package main

import (
	"fmt"

	"github.com/min0625/gsyncmap"
)

func main() {
	var m gsyncmap.Map[string, string]

	m.Store("k1", "v1")
	fmt.Println(m.Load("k1")) // v1 true
	fmt.Println(m.Load("k2")) //  false

	m.Delete("k1")
	fmt.Println(m.Load("k1")) //  false
}
```

When `CompareAndDelete` or `CompareAndSwap` is needed, use `ComparableMap`:

```go
package main

import (
	"fmt"

	"github.com/min0625/gsyncmap"
)

func main() {
	var m gsyncmap.ComparableMap[string, string]

	m.Store("k1", "v1")
	fmt.Println(m.CompareAndSwap("k1", "v1", "v2")) // true
	fmt.Println(m.Load("k1"))                        // v2 true
}
```

## Examples
See: [map_example_test.go](./map_example_test.go)

