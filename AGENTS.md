# Agent Guidelines

This repository provides `gsyncmap` — a generic, type-safe wrapper around [`sync.Map`](https://pkg.go.dev/sync#Map) for Go.

## Repository Structure

| File | Description |
|------|-------------|
| `map.go` | `Map[Key, Value]` — generic concurrent map, value type can be `any` |
| `comparable_map.go` | `ComparableMap[Key, Value]` — variant requiring `comparable` value type |
| `map_test.go` | Unit tests |
| `map_example_test.go` | Runnable example tests (shown in pkg.go.dev) |

## Development Commands

```sh
# Run tests (with race detector)
make test
# go test -v -race -failfast ./...

# Run linter
make lint
# golangci-lint run -v ./...

# Run linter with auto-fix
make fix
# golangci-lint run -v --fix ./...

# Run all checks
make check
```

## Code Conventions

- **Go version**: 1.24+, module path `github.com/min0625/gsyncmap`
- All exported types and methods must have doc comments.
- `Map` wraps `sync.Map` via a direct type alias (`type Map[K, V] sync.Map`); avoid adding fields.
- `CompareAndDelete` / `CompareAndSwap` on `Map` are intentionally **deprecated** — do not un-deprecate them. Direct users to `ComparableMap` instead.
- Keep the two types (`Map` and `ComparableMap`) in sync: any new method added to one should be added to the other when applicable.
- New example functions must follow `Example*` naming conventions and include expected output comments so they serve as testable examples.

## Testing Guidelines

- All tests must pass with `-race` flag enabled.
- Cover both the normal (value found) and zero-value (key absent) code paths for every method.
- Do not add external test dependencies; use only the standard library.
