# Generic Sync Map
[![GoDoc](https://pkg.go.dev/badge/github.com/gin-gonic/gin?status.svg)](https://pkg.go.dev/github.com/min0625/gsyncmap?tab=doc)


This Map is based on `generic` and `sync.Map`.

## Installation
```sh
go get github.com/min0625/gsyncmap
```

## Quick start
```go
package main

import "github.com/min0625/gsyncmap"

func main() {
	var m gsyncmap.Map[string, string]

	m.Store("k1", "v1")
	fmt.Println(m.Load("k1")) // v1 true
	fmt.Println(m.Load("k2")) //  false

	m.Delete("k1")
	fmt.Println(m.Load("k1")) //  false
}
```

## Example
See: [./map_example_test.go](./map_example_test.go)
