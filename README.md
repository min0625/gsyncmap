# Generic Sync Map

[![Go Reference](https://pkg.go.dev/badge/github.com/min0625/gsyncmap.svg)](https://pkg.go.dev/github.com/min0625/gsyncmap)

**English** | [繁體中文](./README.zh-TW.md)

A generic, type-safe wrapper around [`sync.Map`](https://pkg.go.dev/sync#Map) for Go.

## Features

- **Type-safe** — keys and values are statically typed; no `any` casts at call sites.
- **Zero-value ready** — the zero value is usable immediately, just like `sync.Map`.
- **Familiar API** — mirrors the standard `sync.Map` method set.
- **No dependencies** — built entirely on the standard library.

## Installation

```sh
go get github.com/min0625/gsyncmap
```

Requires Go 1.24 or later.

## Types

| Type | Value constraint | `CompareAndDelete` / `CompareAndSwap` |
|------|------------------|----------------------------------------|
| [`Map[Key comparable, Value any]`](./map.go) | `any` | available but **deprecated** (may panic at runtime) |
| [`ComparableMap[Key, Value comparable]`](./comparable_map.go) | `comparable` | safe, panic-free |

### `Map[Key comparable, Value any]`

A generic concurrent map that accepts any value type. Suitable for most use cases.

> **Note:** `CompareAndDelete` and `CompareAndSwap` are available on `Map` but
> deprecated — they panic at runtime if `Value` is not comparable (e.g. slice,
> map, func). Use `ComparableMap` when these operations are needed.

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
	fmt.Println(m.Load("k1"))                       // v2 true
}
```

## Documentation

- API reference: [pkg.go.dev/github.com/min0625/gsyncmap](https://pkg.go.dev/github.com/min0625/gsyncmap)
- Runnable examples: [map_example_test.go](./map_example_test.go)

## License

See [LICENSE](./LICENSE).
