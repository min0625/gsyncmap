# Generic Sync Map

[![Go Reference](https://pkg.go.dev/badge/github.com/min0625/gsyncmap.svg)](https://pkg.go.dev/github.com/min0625/gsyncmap)

[English](./README.md) | **繁體中文**

針對 Go 的 [`sync.Map`](https://pkg.go.dev/sync#Map) 所打造、泛型且型別安全的封裝。

## 特色

- **型別安全** — 鍵與值皆為靜態型別；呼叫端不需要 `any` 型別轉換。
- **零值即可用** — 零值可立即使用，與 `sync.Map` 一致。
- **熟悉的 API** — 對應標準函式庫 `sync.Map` 的方法集。
- **零相依** — 完全建構於標準函式庫之上。

## 安裝

```sh
go get github.com/min0625/gsyncmap
```

需要 Go 1.24 或更新版本。

## 型別

| 型別 | 值的限制 | `CompareAndDelete` / `CompareAndSwap` |
|------|----------|----------------------------------------|
| [`Map[Key comparable, Value any]`](./map.go) | `any` | 可用但**已棄用**（執行期可能 panic） |
| [`ComparableMap[Key, Value comparable]`](./comparable_map.go) | `comparable` | 安全、不會 panic |

### `Map[Key comparable, Value any]`

接受任意值型別的泛型並行 map，適用於大多數情境。

> **注意：** `Map` 上的 `CompareAndDelete` 與 `CompareAndSwap` 雖然可用，但已棄用
> ——當 `Value` 不是 comparable（例如 slice、map、func）時，會在執行期 panic。
> 需要這些操作時請改用 `ComparableMap`。

### `ComparableMap[Key, Value comparable]`

可直接替換 `Map` 的版本，但要求值型別為 comparable。
提供安全的 `CompareAndDelete` 與 `CompareAndSwap`，不會有執行期 panic 的風險。

## 快速開始

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

當需要 `CompareAndDelete` 或 `CompareAndSwap` 時，請改用 `ComparableMap`：

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

## 文件

- API 參考文件：[pkg.go.dev/github.com/min0625/gsyncmap](https://pkg.go.dev/github.com/min0625/gsyncmap)
- 可執行範例：[map_example_test.go](./map_example_test.go)

## 授權

請參閱 [LICENSE](./LICENSE)。
