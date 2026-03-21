# Collection

A Go port of Laravel's [Illuminate Collections](https://laravel.com/docs/collections) -- fluent, type-safe generic collections powered by Go generics.

```bash
go get github.com/gocanto/collection
```

Requires **Go 1.25** or later.

## Overview

This library is a port of Laravel's [Illuminate\Support\Collection](https://github.com/laravel/framework/tree/master/src/Illuminate/Collections) to Go. It preserves the same rich API surface -- `filter`, `map`, `reduce`, `flatMap`, `chunk`, `partition`, and many more -- while redesigning everything idiomatically for Go with generics, `iter.Seq` for lazy evaluation, and full type safety.

- **Readable data pipelines** -- chain `.Filter().Map().Take()` instead of nesting loops.
- **Type-safe generics** -- no `interface{}` casting; the compiler catches type errors at build time.
- **Immutable returns** -- methods return new collections, leaving the original untouched.
- **Lazy evaluation** -- process large or infinite datasets without allocating the entire result up front.
- **Standalone utilities** -- use individual functions from `arr` or `kv` without constructing a collection object.

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/gocanto/collection/collection"
)

func main() {
    numbers := collection.Collect([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

    result := numbers.
        Filter(func(n int, _ int) bool { return n%2 == 0 }).
        Take(3)

    fmt.Println(result.All()) // [2, 4, 6]
}
```

## Packages

| Package | Description |
|---------|-------------|
| `collection` | Core collection types: `Collection[T]`, `MapCollection[K,V]`, `LazyCollection[T]`, and helper functions. |
| `collection/arr` | Standalone generic slice utilities: `Flatten`, `Sort`, `Where`, `Map`, `Pluck`, `Partition`, etc. |
| `collection/kv` | Map utilities with dot-notation support: `Get`, `Set`, `Has`, `Dot`, `Undot`, `Only`, `Query`, etc. |

## Documentation

Full API documentation is available in the [docs](docs/) directory:

- [Overview](docs/overview.md) -- Architecture, design decisions, and usage guide.
- [Collection](docs/collection.md) -- `Collection[T]` API reference.
- [MapCollection](docs/map_collection.md) -- `MapCollection[K,V]` API reference.
- [LazyCollection](docs/lazy_collection.md) -- `LazyCollection[T]` API reference.
- [arr](docs/arr.md) -- Slice utilities API reference.
- [kv](docs/kv.md) -- Map utilities API reference.
- [Helpers](docs/helpers.md) -- Root-package helper functions.

## License

MIT -- see [LICENSE](LICENSE) for details.
