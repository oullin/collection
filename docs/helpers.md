# Helpers

The root `collection` package includes a small set of general-purpose helper functions. These are standalone utilities that complement the collection types.

## Value

Returns the given value unchanged. Useful as a no-op callback or for consistent value resolution patterns.

```go
func Value[T any](value T) T
```

```go
v := collection.Value(42) // 42
```

## ValueFunc

Calls the given callback and returns its result. Useful for deferred evaluation.

```go
func ValueFunc[T any](callback func() T) T
```

```go
v := collection.ValueFunc(func() string {
    return computeExpensiveDefault()
})
```

**Why:** Allows lazy computation — the callback is only executed when ValueFunc is called, not when it's defined. Useful for default values that are expensive to compute.

## Head

Returns the first element of a slice and `true`, or the zero value and `false` if the slice is empty.

```go
func Head[T any](items []T) (T, bool)
```

```go
first, ok := collection.Head([]string{"alice", "bob", "charlie"})
// first = "alice", ok = true

first, ok = collection.Head([]string{})
// first = "", ok = false
```

**Why:** A safe way to get the first element without risking a panic on an empty slice. The boolean return lets you distinguish between "empty" and "zero value."

## Last

Returns the last element of a slice and `true`, or the zero value and `false` if the slice is empty.

```go
func Last[T any](items []T) (T, bool)
```

```go
last, ok := collection.Last([]int{10, 20, 30})
// last = 30, ok = true
```

**Why:** Same safety benefits as Head, but for the tail of the slice.

## WhenValue

Returns `value` if `condition` is true, otherwise returns the first default or the zero value.

```go
func WhenValue[T any](condition bool, value T, defaults ...T) T
```

```go
label := collection.WhenValue(isAdmin, "Admin Panel", "Dashboard")
// Returns "Admin Panel" if isAdmin is true, "Dashboard" otherwise

cssClass := collection.WhenValue(isActive, "active")
// Returns "active" if isActive, "" (zero value) otherwise
```

**Why:** A concise inline conditional that avoids verbose if/else blocks when selecting between two values. Especially useful when constructing strings, configs, or template data.

## WhenFunc

Like WhenValue but accepts callbacks for deferred evaluation. Calls the matching callback only when needed.

```go
func WhenFunc[T any](condition bool, callback func() T, defaults ...func() T) T
```

```go
result := collection.WhenFunc(useCache,
    func() []User { return cache.GetUsers() },
    func() []User { return db.QueryUsers() },
)
```

**Why:** Ensures only the needed branch is evaluated. If the cache path is chosen, the database is never queried. Essential when the alternative computations have side effects or are expensive.

## Error Types

The package defines two error types for collection operations that can fail:

### ItemNotFoundError

Returned by `FirstOrFail`, `Sole`, and similar methods when no matching item exists.

```go
type ItemNotFoundError struct {
    Message string
}
```

```go
user, err := users.FirstOrFail(func(u User, _ int) bool {
    return u.Email == "missing@example.com"
})
if err != nil {
    // err is *collection.ItemNotFoundError
    var notFound *collection.ItemNotFoundError
    if errors.As(err, &notFound) {
        log.Println("user not found")
    }
}
```

### MultipleItemsFoundError

Returned by `Sole` when more than one item matches the predicate.

```go
type MultipleItemsFoundError struct {
    Count   int
    Message string
}
```

```go
admin, err := users.Sole(func(u User, _ int) bool {
    return u.Role == "admin"
})
if err != nil {
    var multiple *collection.MultipleItemsFoundError
    if errors.As(err, &multiple) {
        log.Printf("expected 1 admin, found %d", multiple.Count)
    }
}
```

**Why:** Typed errors let you use `errors.As` for precise error handling, and the `Count` field on `MultipleItemsFoundError` gives you diagnostic information without parsing error strings.

## Shared Types

### Numeric

A type constraint for numeric types, used by aggregation functions like `Sum`, `Avg`, `Min`, `Max`.

```go
type Numeric interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
        ~float32 | ~float64
}
```

### Pair

A generic key-value pair, used by `Combine`, `NewMapFromPairs`, and `ToPairs`.

```go
type Pair[K any, V any] struct {
    Key   K
    Value V
}
```

```go
pairs := []collection.Pair[string, int]{
    {Key: "alice", Value: 95},
    {Key: "bob", Value: 87},
}
m := collection.NewMapFromPairs(pairs...)
```
