# Collection -- Fluent Generic Collections for Go

```
go get github.com/gocanto/collection
```

**collection** is a Go port of Laravel's [Illuminate\Support\Collection](https://github.com/laravel/framework/tree/master/src/Illuminate/Collections). It preserves the same rich API surface -- `filter`, `map`, `reduce`, `flatMap`, `chunk`, `partition`, and many more -- while redesigning everything idiomatically for Go with generics (1.25+), `iter.Seq` for lazy evaluation, and full type safety.

---

## Why This Library Exists

Go's built-in slice and map operations are intentionally minimal. Common tasks like filtering, mapping, partitioning, and reducing require writing boilerplate loops every time. This library provides:

- **Readable data pipelines** -- chain `.Filter().Map().Take()` instead of nesting loops.
- **Type-safe generics** -- no `interface{}` casting; the compiler catches type errors at build time.
- **Immutable returns** -- methods return new collections, leaving the original untouched.
- **Lazy evaluation** -- process large or infinite datasets without allocating the entire result up front.
- **Standalone utilities** -- use individual functions from `arr` or `kv` without constructing a collection object when a single operation is all you need.

---

## Package Structure

```
github.com/gocanto/collection
    |
    |-- collection/              Core collection package
    |       Collection[T]        Fluent wrapper around []T
    |       MapCollection[K,V]   Ordered map with fluent API
    |       LazyCollection[T]    Lazy sequences backed by iter.Seq[T]
    |       Pair[K,V]            Key-value pair type
    |       Numeric              Type constraint for numeric types
    |
    |-- arr/                     Standalone generic slice utilities
    |       Flatten, Sort, Where, Map, Pluck, Partition, ...
    |
    |-- kv/                      Map utilities with dot-notation support
            Get, Set, Has, Dot, Undot, Only, Query, ...
```

### Root package (`collection`)

The root package provides three core collection types and a set of helper functions.

### `arr` sub-package

Package-level generic functions that operate on plain `[]T` slices. Use `arr` when you need a quick one-off operation and do not need method chaining.

See the full API reference: [arr.md](arr.md)

### `kv` sub-package

Package-level functions for `map[string]any` (including dot-notation traversal of nested maps) and generic map utilities.

See the full API reference: [kv.md](kv.md)

---

## Quick Start

### Installation

```bash
go get github.com/gocanto/collection
```

Requires **Go 1.25** or later.

### Basic Collection Usage

```go
package main

import (
    "fmt"
    "github.com/gocanto/collection/collection"
)

func main() {
    // Create a collection from a slice.
    numbers := collection.Collect([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

    // Filter to even numbers, take the first 3.
    result := numbers.
        Filter(func(n int, _ int) bool { return n%2 == 0 }).
        Take(3)

    fmt.Println(result.All()) // [2, 4, 6]
}
```

### Mapping to a Different Type

Because Go generics do not allow methods to introduce new type parameters, type-changing operations like `Map`, `Reduce`, and `FlatMap` are top-level functions:

```go
names := collection.Collect([]string{"alice", "bob", "charlie"})

lengths := collection.Map(names, func(name string, _ int) int {
    return len(name)
})

fmt.Println(lengths.All()) // [5, 3, 7]
```

### MapCollection

```go
m := collection.NewMap(map[string]int{
    "apples":  5,
    "bananas": 3,
    "oranges": 8,
})

m.Has("apples")          // true
val, ok := m.Get("bananas") // 3, true

// Filter entries.
plenty := m.Filter(func(count int, _ string) bool {
    return count > 4
})
fmt.Println(plenty.All()) // map[apples:5 oranges:8]
```

### LazyCollection

Lazy collections defer computation until results are consumed. They are backed by `iter.Seq[T]` and integrate with Go's `range`-over-function syntax:

```go
// Generate a potentially large range lazily.
lazy := collection.LazyRange(1, 1_000_000).
    Filter(func(n int, _ int) bool { return n%2 == 0 }).
    Take(5)

fmt.Println(lazy.All()) // [2, 4, 6, 8, 10]
// Only 10 integers were evaluated, not 1 million.
```

Use `range` with the iterator directly:

```go
for val := range lazy.Iter() {
    fmt.Println(val)
}
```

### Standalone Slice Utilities (arr)

```go
import "github.com/gocanto/collection/arr"

nums := []int{5, 3, 8, 1, 9, 2}

sorted := arr.Sort(nums, func(a, b int) bool { return a < b })
fmt.Println(sorted) // [1, 2, 3, 5, 8, 9]

first, ok := arr.First(nums, func(v int, _ int) bool { return v > 4 })
fmt.Println(first, ok) // 5 true
```

### Standalone Map Utilities (kv)

```go
import "github.com/gocanto/collection/kv"

config := map[string]any{
    "database": map[string]any{
        "host": "localhost",
        "port": 5432,
    },
}

host := kv.Get(config, "database.host")
fmt.Println(host) // "localhost"

kv.Set(config, "database.name", "mydb")
```

---

## Key Design Decisions

### Go Generics

Every collection type is parameterized with type variables (`Collection[T]`, `MapCollection[K, V]`, `LazyCollection[T]`). This means:

- No type assertions at call sites.
- The compiler verifies type compatibility.
- Operations like `Map` can transform `[]T` into `[]R` where `R` is a completely different type.

Because Go does not allow methods to introduce new type parameters, type-changing operations (`Map`, `FlatMap`, `Reduce`, `MapValues`, `Pipe`, `PipeInto`, `LazyMap`, etc.) are implemented as package-level functions rather than methods.

### `iter.Seq` Integration

`LazyCollection` is built on Go's `iter.Seq[T]` iterator protocol. This means:

- Lazy collections compose naturally with any code that accepts `iter.Seq`.
- You can use `range` directly over `lc.Iter()`.
- Constructors like `NewLazy` accept any `iter.Seq[T]`, so you can wrap database cursors, file readers, or channel consumers.

`Collection` and `MapCollection` also expose `.Iter()` methods that return `iter.Seq[T]` and `iter.Seq2[K, V]` respectively, enabling interop with the broader Go iterator ecosystem.

### Immutable Returns

Methods that transform data (`.Filter()`, `.Take()`, `.Reject()`, `.Flatten()`, etc.) return **new** collections. The original collection is never modified by these operations. Mutating methods like `.Push()`, `.Put()`, `.Pull()`, `.Pop()`, and `.Shift()` are clearly named and documented as in-place operations.

This makes it safe to share collections across goroutines (for read-only access) and to build pipelines without worrying about side effects.

---

## When to Use Each Collection Type

| Type | Use when... |
|------|-------------|
| **`Collection[T]`** | You have a slice of data and want to filter, sort, map, chunk, or aggregate it with fluent method chaining. This is the workhorse for most day-to-day slice operations. |
| **`MapCollection[K, V]`** | You are working with key-value data and need ordered iteration, key-based lookup, filtering, merging, or set operations (diff, intersect, union). Maintains insertion order. |
| **`LazyCollection[T]`** | You are working with large datasets, streams, or infinite sequences where you do not want to materialize the entire result. Operations are deferred until you call `.All()`, `.Eager()`, `.Count()`, or iterate with `range`. |
| **`arr.FuncName`** | You need a single slice operation (one filter, one sort, one map) and do not need chaining. Avoids allocating a collection object. |
| **`kv.FuncName`** | You need to work with `map[string]any` (especially nested maps from JSON/config) using dot-notation paths, or you need a quick map filter/transform. |

### Conversion Between Types

```go
// Slice -> Collection
c := collection.Collect([]int{1, 2, 3})

// Collection -> LazyCollection (not built-in; wrap the slice)
lc := collection.LazyFrom(c.All())

// LazyCollection -> Collection
c2 := lc.Eager()   // or lc.Collect()

// Collection -> MapCollection (via arr.KeyBy or MapWithKeys)
users := collection.Collect([]User{ ... })
indexed := collection.NewMap(arr.KeyBy(users.All(), func(u User) int { return u.ID }))
```

---

## Further Reading

- **[arr.md](arr.md)** -- Full API reference for the `arr` slice utilities package.
- **[kv.md](kv.md)** -- Full API reference for the `kv` map utilities package.
- **[helpers.md](helpers.md)** -- Documentation for root-package helper functions.
