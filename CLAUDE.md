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
# go test -race -failfast ./...

# Run linter (verifies config, then lints)
make lint
# golangci-lint config verify
# golangci-lint run --new-from-rev=HEAD ./...

# Run linter with auto-fix
make fix
# go mod tidy
# golangci-lint run --new-from-rev=HEAD --fix ./...

# Check that go.mod/go.sum are tidy
make check-tidy
# go mod tidy -diff

# Run all checks (check-tidy + lint + test)
make check
```

> **Note:** `lint`/`fix` use `--new-from-rev=$(NEW_FROM_REV)` (default `HEAD`), so they
> only report issues on changed lines. To lint the whole tree, pass an empty revision,
> e.g. `make lint NEW_FROM_REV=`.

Tool versions are pinned in `mise.toml` (see it for the exact Go and golangci-lint
versions). The linter uses the golangci-lint v2 config schema (`.golangci.yaml`,
`version: "2"`).

## Code Conventions

- **Go version**: see the `go` directive in `go.mod` (toolchain pinned in `mise.toml`); module path `github.com/min0625/gsyncmap`
- All exported types and methods must have doc comments.
- `Map` is a defined type over `sync.Map` (`type Map[Key comparable, Value any] sync.Map`), accessed through the unexported `syncMap()` helper; do not add fields. The unexported `zero[T]()` helper returns the zero `Value` on miss.
- `ComparableMap[Key, Value comparable]` embeds `Map` and only overrides `CompareAndDelete` / `CompareAndSwap` with panic-free implementations; all other methods are inherited.
- `CompareAndDelete` / `CompareAndSwap` on `Map` are intentionally **deprecated** — do not un-deprecate them. Direct users to `ComparableMap` instead.
- Keep the two types in sync: any new method added to `Map` should be considered for `ComparableMap` (and vice versa) when applicable.
- New example functions must follow `Example*` naming conventions and include expected output comments so they serve as testable examples.

## Testing Guidelines

- All tests must pass with `-race` flag enabled.
- Cover both the normal (value found) and zero-value (key absent) code paths for every method.
- Do not add external test dependencies; use only the standard library.
