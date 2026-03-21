# Collection[T] API Reference

`Collection[T]` is a generic wrapper around a Go slice that provides a fluent, chainable API for filtering, sorting, transforming, and aggregating data. It is the core type of the `github.com/gocanto/collection/collection` package.

```go
import "github.com/gocanto/collection/collection"
```

> **A note on top-level functions vs methods.** Go generics do not allow methods to introduce new type parameters. Functions like `Map`, `FlatMap`, `Reduce`, `SortBy`, `Unique`, `Duplicates`, `Pluck`, `GroupBy`, `KeyBy`, `CountBy`, `Zip`, `CrossJoin`, `Combine`, `Collapse`, `Diff`, `Intersect`, `Pipe`, `PipeInto`, `Sum`, `SumBy`, `Avg`, `AvgBy`, `Average`, `Min`, `MinBy`, `Max`, `MaxBy`, `Median`, `MedianBy`, `Mode`, `MapToDictionary`, `MapToGroups`, `MapWithKeys`, and `MapInto` are therefore package-level functions, not methods. Call them as `collection.Map(c, fn)`, not `c.Map(fn)`.

---

## Constructors

### New

```go
func New[T any](items ...T) *Collection[T]
```

Creates a new collection from variadic arguments.

```go
nums := collection.New(1, 2, 3)
fmt.Println(nums.Count()) // 3
```

Use `New` when you have a handful of literal values to wrap.

---

### Collect

```go
func Collect[T any](items []T) *Collection[T]
```

Creates a new collection from an existing slice.

```go
names := []string{"Alice", "Bob", "Carol"}
c := collection.Collect(names)
fmt.Println(c.Count()) // 3
```

Use `Collect` when you already have a slice from another function or API call.

---

### Empty

```go
func Empty[T any]() *Collection[T]
```

Creates an empty collection of the given type.

```go
c := collection.Empty[string]()
fmt.Println(c.IsEmpty()) // true
```

Useful as a starting point when you plan to build up a collection with `Push`.

---

### Wrap

```go
func Wrap[T any](value any) *Collection[T]
```

Wraps a value in a collection. If the value is already a `*Collection[T]`, it is returned as-is. If it is a `[]T`, it is wrapped. Otherwise a single-element collection is created.

```go
c := collection.Wrap[int]([]int{10, 20})
fmt.Println(c.Count()) // 2
```

Useful at API boundaries where the input type is not known statically.

---

### Unwrap

```go
func Unwrap[T any](value any) []T
```

Extracts the underlying slice from a collection or returns the value itself if it is already a slice. Returns `nil` if the value is neither.

```go
c := collection.New(1, 2, 3)
items := collection.Unwrap[int](c)
fmt.Println(len(items)) // 3
```

Useful when you need to pass collection data to functions that expect plain slices.

---

### Times

```go
func Times[T any](number int, callback func(int) T) *Collection[T]
```

Creates a collection by invoking the callback `number` times. The callback receives a 1-based index.

```go
labels := collection.Times(3, func(i int) string {
    return fmt.Sprintf("Item #%d", i)
})
// ["Item #1", "Item #2", "Item #3"]
```

Useful for generating seed data, placeholder records, or numbered sequences.

---

### Range

```go
func Range(from, to int) *Collection[int]
```

Creates a collection of consecutive integers from `from` to `to` (inclusive). If `from > to`, the sequence counts downward.

```go
ascending := collection.Range(1, 5)   // [1, 2, 3, 4, 5]
descending := collection.Range(5, 1)  // [5, 4, 3, 2, 1]
```

Useful for generating index ranges, pagination offsets, or numeric test data.

---

## Query

### All

```go
func (c *Collection[T]) All() []T
```

Returns the underlying slice. This does not copy the data.

```go
c := collection.New("a", "b", "c")
items := c.All()
fmt.Println(items) // [a b c]
```

Use `All` when you need direct access to the underlying data for interop with slice-based APIs.

---

### Count

```go
func (c *Collection[T]) Count() int
```

Returns the number of items in the collection.

```go
c := collection.New(10, 20, 30)
fmt.Println(c.Count()) // 3
```

Use `Count` for size checks, loop bounds, or conditional logic based on collection size.

---

### IsEmpty

```go
func (c *Collection[T]) IsEmpty() bool
```

Reports whether the collection contains no items.

```go
c := collection.Empty[int]()
if c.IsEmpty() {
    fmt.Println("nothing here")
}
```

Use `IsEmpty` for guard clauses before processing a collection.

---

### IsNotEmpty

```go
func (c *Collection[T]) IsNotEmpty() bool
```

Reports whether the collection contains at least one item.

```go
results := collection.New("found one")
if results.IsNotEmpty() {
    fmt.Println("we have results")
}
```

Use `IsNotEmpty` as a readable alternative to `c.Count() > 0`.

---

### ContainsOneItem

```go
func (c *Collection[T]) ContainsOneItem() bool
```

Reports whether the collection contains exactly one item.

```go
c := collection.New("only")
fmt.Println(c.ContainsOneItem()) // true
```

Useful for validating that a query returned a singular result.

---

### ContainsManyItems

```go
func (c *Collection[T]) ContainsManyItems() bool
```

Reports whether the collection contains more than one item.

```go
c := collection.New(1, 2, 3)
fmt.Println(c.ContainsManyItems()) // true
```

Useful when you need to distinguish between single-item and multi-item results.

---

### HasMany

```go
func (c *Collection[T]) HasMany() bool
```

Alias for `ContainsManyItems`.

---

### Has

```go
func (c *Collection[T]) Has(index int) bool
```

Reports whether the given index exists. Negative indices count from the end.

```go
c := collection.New("a", "b", "c")
fmt.Println(c.Has(1))   // true
fmt.Println(c.Has(-1))  // true  (last element)
fmt.Println(c.Has(10))  // false
```

Use `Has` to safely check bounds before calling `Get`.

---

### HasAny

```go
func (c *Collection[T]) HasAny(indices ...int) bool
```

Reports whether any of the given indices exist in the collection.

```go
c := collection.New("a", "b")
fmt.Println(c.HasAny(0, 5, 10)) // true (index 0 exists)
```

Useful for checking whether at least one of several expected positions is populated.

---

### Keys

```go
func (c *Collection[T]) Keys() *Collection[int]
```

Returns a collection of indices (0 through n-1).

```go
c := collection.New("a", "b", "c")
keys := c.Keys() // [0, 1, 2]
```

Useful when you need the index set for further computation.

---

## Retrieval

### First

```go
func (c *Collection[T]) First(predicates ...func(T, int) bool) (T, bool)
```

Returns the first element matching the optional predicate. If no predicate is given, returns the first element. The second return value indicates whether a match was found.

```go
type User struct {
    Name string
    Age  int
}

users := collection.New(
    User{"Alice", 30},
    User{"Bob", 25},
    User{"Carol", 35},
)

// First element overall
first, ok := users.First()
fmt.Println(first.Name, ok) // Alice true

// First user over 28
senior, ok := users.First(func(u User, _ int) bool {
    return u.Age > 28
})
fmt.Println(senior.Name, ok) // Alice true
```

Use `First` for quick lookups without scanning the full collection.

---

### FirstOrFail

```go
func (c *Collection[T]) FirstOrFail(predicates ...func(T, int) bool) (T, error)
```

Like `First`, but returns an `ItemNotFoundError` instead of `false` when no match is found.

```go
user, err := users.FirstOrFail(func(u User, _ int) bool {
    return u.Name == "Dave"
})
if err != nil {
    fmt.Println(err) // "item not found"
}
```

Use `FirstOrFail` when a missing item should be treated as an error rather than silently ignored.

---

### Last

```go
func (c *Collection[T]) Last(predicates ...func(T, int) bool) (T, bool)
```

Returns the last element matching the optional predicate. If no predicate is given, returns the last element.

```go
c := collection.New(1, 2, 3, 4, 5)
last, _ := c.Last()
fmt.Println(last) // 5

lastEven, _ := c.Last(func(n int, _ int) bool {
    return n%2 == 0
})
fmt.Println(lastEven) // 4
```

Use `Last` to find the most recent or final matching element.

---

### Sole

```go
func (c *Collection[T]) Sole(predicates ...func(T, int) bool) (T, error)
```

Returns the only element matching the optional predicate. Returns `ItemNotFoundError` if zero items match, or `MultipleItemsFoundError` if more than one matches.

```go
admins := collection.New(
    User{"Alice", 30},
)
admin, err := admins.Sole()
fmt.Println(admin.Name, err) // Alice <nil>

// If multiple items match:
all := collection.New(User{"A", 1}, User{"B", 2})
_, err = all.Sole()
fmt.Println(err) // "multiple items found: 2 items"
```

Use `Sole` when your logic requires exactly one matching record and any other count is an error.

---

### HasSole

```go
func (c *Collection[T]) HasSole(predicates ...func(T, int) bool) bool
```

Reports whether exactly one item matches the optional predicate.

```go
c := collection.New(1, 2, 3)
fmt.Println(c.HasSole(func(n int, _ int) bool { return n > 2 })) // true (only 3)
```

Use `HasSole` as a boolean check before calling `Sole` to avoid error handling.

---

### Get

```go
func (c *Collection[T]) Get(index int, defaults ...T) (T, bool)
```

Returns the item at the given index. Negative indices count from the end. The second return value indicates whether the index was in bounds. An optional default value can be provided.

```go
c := collection.New("a", "b", "c")
val, ok := c.Get(1)
fmt.Println(val, ok) // b true

val, ok = c.Get(-1)
fmt.Println(val, ok) // c true

val, ok = c.Get(99, "fallback")
fmt.Println(val, ok) // fallback false
```

Use `Get` for safe index-based access with optional defaults.

---

### GetOrPut

```go
func (c *Collection[T]) GetOrPut(index int, value T) T
```

Returns the item at the given index if it exists. Otherwise appends the value to the collection and returns it.

```go
c := collection.New("a", "b")
val := c.GetOrPut(0, "default") // "a" (exists)
val = c.GetOrPut(5, "new")      // "new" (appended)
fmt.Println(c.Count())          // 3
```

Useful for "get or initialize" patterns.

---

## Search

### Contains

```go
func (c *Collection[T]) Contains(predicate func(T, int) bool) bool
```

Reports whether any item satisfies the predicate.

```go
c := collection.New(1, 2, 3, 4, 5)
hasEven := c.Contains(func(n int, _ int) bool {
    return n%2 == 0
})
fmt.Println(hasEven) // true
```

Use `Contains` to check for the existence of an item matching a condition.

---

### Some

```go
func (c *Collection[T]) Some(predicate func(T, int) bool) bool
```

Alias for `Contains`.

---

### DoesntContain

```go
func (c *Collection[T]) DoesntContain(predicate func(T, int) bool) bool
```

Reports whether no item satisfies the predicate. The logical inverse of `Contains`.

```go
c := collection.New(1, 3, 5)
noEvens := c.DoesntContain(func(n int, _ int) bool {
    return n%2 == 0
})
fmt.Println(noEvens) // true
```

Use `DoesntContain` for readable "none match" checks.

---

### Search

```go
func (c *Collection[T]) Search(predicate func(T, int) bool) (int, bool)
```

Returns the index of the first item satisfying the predicate. The second return value indicates whether a match was found.

```go
c := collection.New("apple", "banana", "cherry")
idx, found := c.Search(func(s string, _ int) bool {
    return s == "banana"
})
fmt.Println(idx, found) // 1 true
```

Use `Search` when you need the position of a matching item.

---

### Before

```go
func (c *Collection[T]) Before(predicate func(T, int) bool) (T, bool)
```

Returns the item immediately before the first item matching the predicate.

```go
c := collection.New("a", "b", "c", "d")
before, ok := c.Before(func(s string, _ int) bool {
    return s == "c"
})
fmt.Println(before, ok) // b true
```

Useful for getting the predecessor of a known element, such as finding the previous step in a workflow.

---

### After

```go
func (c *Collection[T]) After(predicate func(T, int) bool) (T, bool)
```

Returns the item immediately after the first item matching the predicate.

```go
c := collection.New("a", "b", "c", "d")
after, ok := c.After(func(s string, _ int) bool {
    return s == "b"
})
fmt.Println(after, ok) // c true
```

Useful for getting the successor of a known element, such as finding the next step in a pipeline.

---

## Mutation

### Push

```go
func (c *Collection[T]) Push(values ...T) *Collection[T]
```

Appends one or more items to the end of the collection. Mutates the collection in place.

```go
c := collection.New(1, 2)
c.Push(3, 4)
fmt.Println(c.All()) // [1 2 3 4]
```

Use `Push` to add items to a growing collection.

---

### Add

```go
func (c *Collection[T]) Add(item T) *Collection[T]
```

Appends a single item. Alias for `Push` with one argument.

```go
c := collection.New("x")
c.Add("y")
fmt.Println(c.All()) // [x y]
```

---

### Prepend

```go
func (c *Collection[T]) Prepend(value T) *Collection[T]
```

Adds an item to the beginning of the collection.

```go
c := collection.New(2, 3)
c.Prepend(1)
fmt.Println(c.All()) // [1 2 3]
```

Use `Prepend` when insertion order matters and you need the new item at the front.

---

### Unshift

```go
func (c *Collection[T]) Unshift(value T) *Collection[T]
```

Alias for `Prepend`.

---

### Pop

```go
func (c *Collection[T]) Pop(counts ...int) (T, bool)
```

Removes and returns the last item from the collection. The second return value indicates whether the collection was non-empty.

```go
c := collection.New(1, 2, 3)
last, ok := c.Pop()
fmt.Println(last, ok) // 3 true
fmt.Println(c.All())  // [1 2]
```

Use `Pop` for stack-like (LIFO) behavior.

---

### PopMany

```go
func (c *Collection[T]) PopMany(count int) *Collection[T]
```

Removes and returns the last `count` items from the collection as a new collection.

```go
c := collection.New(1, 2, 3, 4, 5)
popped := c.PopMany(2)
fmt.Println(popped.All()) // [4 5]
fmt.Println(c.All())      // [1 2 3]
```

Use `PopMany` to remove a batch of items from the end.

---

### Shift

```go
func (c *Collection[T]) Shift() (T, bool)
```

Removes and returns the first item from the collection.

```go
c := collection.New("a", "b", "c")
first, ok := c.Shift()
fmt.Println(first, ok) // a true
fmt.Println(c.All())   // [b c]
```

Use `Shift` for queue-like (FIFO) behavior.

---

### ShiftMany

```go
func (c *Collection[T]) ShiftMany(count int) *Collection[T]
```

Removes and returns the first `count` items as a new collection.

```go
c := collection.New(1, 2, 3, 4, 5)
shifted := c.ShiftMany(2)
fmt.Println(shifted.All()) // [1 2]
fmt.Println(c.All())       // [3 4 5]
```

Use `ShiftMany` to dequeue a batch of items.

---

### Put

```go
func (c *Collection[T]) Put(index int, value T) *Collection[T]
```

Sets the item at the given index to the given value. Does nothing if the index is out of bounds.

```go
c := collection.New("a", "b", "c")
c.Put(1, "B")
fmt.Println(c.All()) // [a B c]
```

Use `Put` for in-place updates at a known position.

---

### Pull

```go
func (c *Collection[T]) Pull(index int) (T, bool)
```

Removes and returns the item at the given index. The second return value indicates whether the index was valid.

```go
c := collection.New(10, 20, 30)
val, ok := c.Pull(1)
fmt.Println(val, ok) // 20 true
fmt.Println(c.All()) // [10 30]
```

Use `Pull` to remove a specific element by position.

---

### Forget

```go
func (c *Collection[T]) Forget(index int) *Collection[T]
```

Removes an item by index, mutating the collection. Returns the collection for chaining.

```go
c := collection.New("a", "b", "c")
c.Forget(1)
fmt.Println(c.All()) // [a c]
```

Use `Forget` when you do not need the removed value.

---

### Transform

```go
func (c *Collection[T]) Transform(callback func(T, int) T) *Collection[T]
```

Applies the callback to each item in place, mutating the collection. Unlike `Map`, the return type must be the same as the input type.

```go
prices := collection.New(10.0, 20.0, 30.0)
prices.Transform(func(p float64, _ int) float64 {
    return p * 1.1 // apply 10% markup
})
fmt.Println(prices.All()) // [11 22 33]
```

Use `Transform` for in-place mutations when you do not need a new type.

---

## Iteration

### Each

```go
func (c *Collection[T]) Each(callback func(T, int) bool) *Collection[T]
```

Iterates over items, passing each item and its index to the callback. Return `false` from the callback to stop early.

```go
collection.New("a", "b", "c").Each(func(item string, i int) bool {
    fmt.Printf("%d: %s\n", i, item)
    return true // continue
})
```

Use `Each` for side effects like logging, sending notifications, or accumulating into external state.

---

### EachSpread

```go
func (c *Collection[T]) EachSpread(callback func(T, int) bool) *Collection[T]
```

Operates the same as `Each` in Go. Exists for API parity with the Laravel collection.

---

### Tap

```go
func (c *Collection[T]) Tap(callback func(*Collection[T])) *Collection[T]
```

Passes the collection to the callback for side effects and returns the collection unchanged.

```go
result := collection.New(1, 2, 3).
    Tap(func(c *collection.Collection[int]) {
        fmt.Println("Count:", c.Count())
    }).
    Filter(func(n int, _ int) bool { return n > 1 })
```

Use `Tap` to inspect or log a collection mid-chain without breaking the fluent flow.

---

### TapEach

```go
func (c *Collection[T]) TapEach(callback func(T, int)) *Collection[T]
```

Calls the callback on each item for side effects, returning the original collection unchanged.

```go
collection.New("order-1", "order-2").TapEach(func(id string, _ int) {
    fmt.Println("Processing:", id)
})
```

Use `TapEach` to perform side effects on each item (logging, metrics) without modifying the collection.

---

### Pipe

```go
func Pipe[T any, R any](c *Collection[T], callback func(*Collection[T]) R) R
```

Passes the collection to the callback and returns the callback's result. The result can be any type.

```go
c := collection.New(1, 2, 3, 4, 5)
sum := collection.Pipe(c, func(c *collection.Collection[int]) int {
    total := 0
    c.Each(func(n int, _ int) bool { total += n; return true })
    return total
})
fmt.Println(sum) // 15
```

Use `Pipe` to break out of the fluent chain with a custom computation that returns a different type.

---

### PipeInto

```go
func PipeInto[T any, R any](c *Collection[T], constructor func(*Collection[T]) R) R
```

Passes the collection to a constructor function and returns the result. Semantically identical to `Pipe` but communicates the intent of constructing a new value.

```go
type Report struct{ Total int }

c := collection.New(100, 200, 300)
report := collection.PipeInto(c, func(c *collection.Collection[int]) Report {
    return Report{Total: c.Count()}
})
```

Use `PipeInto` when you want to construct a domain object from a collection.

---

### PipeThrough

```go
func PipeThrough[T any](c *Collection[T], callbacks ...func(*Collection[T]) *Collection[T]) *Collection[T]
```

Passes the collection through a series of callbacks, returning the final result. Each callback receives and returns a `*Collection[T]`.

```go
result := collection.PipeThrough(
    collection.New(1, 2, 3, 4, 5),
    func(c *collection.Collection[int]) *collection.Collection[int] {
        return c.Filter(func(n int, _ int) bool { return n > 2 })
    },
    func(c *collection.Collection[int]) *collection.Collection[int] {
        return c.Reverse()
    },
)
fmt.Println(result.All()) // [5 4 3]
```

Use `PipeThrough` to build reusable processing pipelines from composable stages.

---

### Iter

```go
func (c *Collection[T]) Iter() iter.Seq[T]
```

Returns a Go 1.23+ iterator (`iter.Seq[T]`) that yields each item. Compatible with range-over-func loops.

```go
c := collection.New("x", "y", "z")
for item := range c.Iter() {
    fmt.Println(item)
}
```

Use `Iter` to integrate collections with the standard iterator protocol.

---

### Iter2

```go
func (c *Collection[T]) Iter2() iter.Seq2[int, T]
```

Returns a Go 1.23+ iterator that yields each index-item pair.

```go
c := collection.New("a", "b", "c")
for i, item := range c.Iter2() {
    fmt.Printf("%d: %s\n", i, item)
}
```

Use `Iter2` when you need both index and value in a range loop.

---

## Filtering

### Filter

```go
func (c *Collection[T]) Filter(callback func(T, int) bool) *Collection[T]
```

Returns a new collection containing only items for which the callback returns `true`.

```go
type Order struct {
    ID     int
    Amount float64
}

orders := collection.New(
    Order{1, 50.0},
    Order{2, 150.0},
    Order{3, 25.0},
)
big := orders.Filter(func(o Order, _ int) bool {
    return o.Amount > 100
})
fmt.Println(big.Count()) // 1
```

Use `Filter` for any "keep matching items" operation.

---

### Reject

```go
func (c *Collection[T]) Reject(callback func(T, int) bool) *Collection[T]
```

Returns a new collection excluding items for which the callback returns `true`. The inverse of `Filter`.

```go
nums := collection.New(1, 2, 3, 4, 5)
odds := nums.Reject(func(n int, _ int) bool {
    return n%2 == 0
})
fmt.Println(odds.All()) // [1 3 5]
```

Use `Reject` when it is more natural to describe what you want to remove rather than what you want to keep.

---

### Where

```go
func (c *Collection[T]) Where(predicate func(T) bool) *Collection[T]
```

Filters items using a predicate that receives only the item (no index). A convenience wrapper around `Filter`.

```go
type Product struct {
    Name   string
    InStock bool
}

products := collection.New(
    Product{"Widget", true},
    Product{"Gadget", false},
    Product{"Doohickey", true},
)
available := products.Where(func(p Product) bool {
    return p.InStock
})
fmt.Println(available.Count()) // 2
```

Use `Where` when you do not need the index in your predicate, for cleaner call sites.

---

### WhereNot

```go
func (c *Collection[T]) WhereNot(predicate func(T) bool) *Collection[T]
```

Filters items using a negative predicate. Keeps items for which the predicate returns `false`.

```go
discontinued := products.WhereNot(func(p Product) bool {
    return p.InStock
})
fmt.Println(discontinued.Count()) // 1
```

Use `WhereNot` as the inverse of `Where`.

---

### Unique

```go
func Unique[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T]
```

Returns a new collection containing only items with distinct keys as determined by the key function.

```go
type User struct {
    Email string
    Name  string
}

users := collection.New(
    User{"a@co.com", "Alice"},
    User{"b@co.com", "Bob"},
    User{"a@co.com", "Alice Duplicate"},
)
unique := collection.Unique(users, func(u User) string {
    return u.Email
})
fmt.Println(unique.Count()) // 2
```

Use `Unique` to deduplicate records by a key field.

---

### Duplicates

```go
func Duplicates[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T]
```

Returns a new collection containing items that appear more than once, as determined by the key function.

```go
tags := collection.New("go", "rust", "go", "python", "rust")
dups := collection.Duplicates(tags, func(s string) string { return s })
fmt.Println(dups.All()) // [go rust]
```

Use `Duplicates` to detect and inspect repeated entries in a dataset.

---

## Transformation

### Map

```go
func Map[T any, R any](c *Collection[T], callback func(T, int) R) *Collection[R]
```

Applies the callback to each item and returns a new collection of the results. The output type can differ from the input type.

```go
type User struct {
    Name string
    Age  int
}

users := collection.New(User{"Alice", 30}, User{"Bob", 25})
names := collection.Map(users, func(u User, _ int) string {
    return u.Name
})
fmt.Println(names.All()) // [Alice Bob]
```

Use `Map` whenever you need to transform items into a different type.

---

### FlatMap

```go
func FlatMap[T any, R any](c *Collection[T], callback func(T, int) []R) *Collection[R]
```

Applies the callback to each item (which returns a slice), then flattens the results into a single collection.

```go
sentences := collection.New("hello world", "foo bar")
words := collection.FlatMap(sentences, func(s string, _ int) []string {
    return strings.Split(s, " ")
})
fmt.Println(words.All()) // [hello world foo bar]
```

Use `FlatMap` when each input item maps to multiple output items.

---

### MapInto

```go
func MapInto[T any, R any](c *Collection[T], constructor func(T) R) *Collection[R]
```

Applies the constructor to each item, returning a new collection. Similar to `Map` but the constructor receives only the item, not the index.

```go
type DTO struct{ Value string }

raw := collection.New("a", "b", "c")
dtos := collection.MapInto(raw, func(s string) DTO {
    return DTO{Value: strings.ToUpper(s)}
})
```

Use `MapInto` for wrapping raw values into domain types.

---

### Reduce

```go
func Reduce[T any, R any](c *Collection[T], callback func(R, T, int) R, initial R) R
```

Iterates over the collection and accumulates a single result.

```go
prices := collection.New(19.99, 29.99, 9.99)
total := collection.Reduce(prices, func(sum float64, price float64, _ int) float64 {
    return sum + price
}, 0.0)
fmt.Println(total) // 59.97
```

Use `Reduce` for computing aggregates, building maps, or any fold operation.

---

### Flatten

```go
func (c *Collection[T]) Flatten() *Collection[T]
```

Returns a shallow copy of the collection. For typed Go slices (which are inherently flat), this simply copies the items.

```go
c := collection.New(1, 2, 3)
flat := c.Flatten()
```

---

### Values

```go
func (c *Collection[T]) Values() *Collection[T]
```

Returns a new collection with items re-indexed (a shallow copy).

```go
c := collection.New("a", "b", "c")
v := c.Values()
fmt.Println(v.All()) // [a b c]
```

Use `Values` after operations that may have changed ordering to get a clean, re-indexed copy.

---

### Reverse

```go
func (c *Collection[T]) Reverse() *Collection[T]
```

Returns a new collection with items in reverse order.

```go
c := collection.New(1, 2, 3)
fmt.Println(c.Reverse().All()) // [3 2 1]
```

Use `Reverse` for displaying results in reverse chronological order or inverting a sorted list.

---

### Shuffle

```go
func (c *Collection[T]) Shuffle() *Collection[T]
```

Returns a new collection with items in random order.

```go
c := collection.New(1, 2, 3, 4, 5)
shuffled := c.Shuffle()
fmt.Println(shuffled.All()) // random order
```

Use `Shuffle` for randomizing quiz questions, playlist order, or A/B test assignments.

---

### Random

```go
func (c *Collection[T]) Random(counts ...int) *Collection[T]
```

Returns a new collection with randomly selected items. Defaults to 1 item if no count is given.

```go
c := collection.New("a", "b", "c", "d", "e")
sample := c.Random(2)
fmt.Println(sample.Count()) // 2
```

Use `Random` for sampling a subset of data.

---

### Flip

```go
func (c *Collection[T]) Flip() *Collection[T]
```

Returns a new collection with items in reverse order. For typed Go slices, this behaves identically to `Reverse`.

```go
c := collection.New(1, 2, 3)
fmt.Println(c.Flip().All()) // [3 2 1]
```

---

## Sorting

### Sort

```go
func (c *Collection[T]) Sort(less func(a, b T) bool) *Collection[T]
```

Returns a new collection sorted using the provided comparison function. Uses a stable sort.

```go
c := collection.New(3, 1, 4, 1, 5)
sorted := c.Sort(func(a, b int) bool { return a < b })
fmt.Println(sorted.All()) // [1 1 3 4 5]
```

Use `Sort` when you need full control over the comparison logic.

---

### SortBy

```go
func SortBy[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) *Collection[T]
```

Returns a new collection sorted in ascending order by the key extracted from each item.

```go
type Task struct {
    Name     string
    Priority int
}

tasks := collection.New(
    Task{"Deploy", 3},
    Task{"Test", 1},
    Task{"Build", 2},
)
sorted := collection.SortBy(tasks, func(t Task) int {
    return t.Priority
})
// [{Test 1} {Build 2} {Deploy 3}]
```

Use `SortBy` for clean, declarative ascending sorts by a single field.

---

### SortByDesc

```go
func SortByDesc[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) *Collection[T]
```

Returns a new collection sorted in descending order by the extracted key.

```go
sorted := collection.SortByDesc(tasks, func(t Task) int {
    return t.Priority
})
// [{Deploy 3} {Build 2} {Test 1}]
```

Use `SortByDesc` for "highest first" or "most recent first" ordering.

---

### SortDesc

```go
func (c *Collection[T]) SortDesc(less func(a, b T) bool) *Collection[T]
```

Returns a new collection sorted in descending order using the provided comparison function. Reverses the sense of the `less` function.

```go
c := collection.New(3, 1, 4, 1, 5)
sorted := c.SortDesc(func(a, b int) bool { return a < b })
fmt.Println(sorted.All()) // [5 4 3 1 1]
```

Use `SortDesc` when you need custom descending sort logic.

---

## Partitioning

### Chunk

```go
func (c *Collection[T]) Chunk(size int) [][]T
```

Breaks the collection into multiple slices of the given size.

```go
c := collection.New(1, 2, 3, 4, 5)
chunks := c.Chunk(2)
// [[1 2] [3 4] [5]]
```

Use `Chunk` for batch processing, pagination, or splitting work across goroutines.

---

### ChunkWhile

```go
func (c *Collection[T]) ChunkWhile(callback func(T, int, []T) bool) [][]T
```

Breaks the collection into groups as long as the callback returns `true`. A new group starts each time the callback returns `false`. The callback receives the current item, its index, and the current chunk.

```go
c := collection.New(1, 1, 2, 2, 3, 3)
chunks := c.ChunkWhile(func(val int, _ int, current []int) bool {
    return val == current[len(current)-1]
})
// [[1 1] [2 2] [3 3]]
```

Use `ChunkWhile` to group consecutive runs of related items.

---

### Split

```go
func (c *Collection[T]) Split(numberOfGroups int) [][]T
```

Splits the collection into the given number of groups, distributing items as evenly as possible.

```go
c := collection.New(1, 2, 3, 4, 5)
groups := c.Split(3)
// [[1 2] [3 4] [5]]
```

Use `Split` for distributing work across a fixed number of workers.

---

### SplitIn

```go
func (c *Collection[T]) SplitIn(numberOfGroups int) [][]T
```

Splits the collection into groups, filling non-terminal groups completely before starting the next.

```go
c := collection.New(1, 2, 3, 4, 5)
groups := c.SplitIn(3)
// [[1 2] [3 4] [5]]
```

Use `SplitIn` when you want consistently full groups with only the last group potentially smaller.

---

### Sliding

```go
func (c *Collection[T]) Sliding(size int, steps ...int) [][]T
```

Returns a sliding window view of the collection. The optional step parameter controls how far the window moves between iterations (default 1).

```go
c := collection.New(1, 2, 3, 4, 5)
windows := c.Sliding(3)
// [[1 2 3] [2 3 4] [3 4 5]]

windows = c.Sliding(3, 2)
// [[1 2 3] [3 4 5]]
```

Use `Sliding` for moving average calculations, n-gram extraction, or any sliding window analysis.

---

### Partition

```go
func (c *Collection[T]) Partition(callback func(T, int) bool) (*Collection[T], *Collection[T])
```

Splits the collection into two: items that pass the predicate and items that do not.

```go
nums := collection.New(1, 2, 3, 4, 5, 6)
evens, odds := nums.Partition(func(n int, _ int) bool {
    return n%2 == 0
})
fmt.Println(evens.All()) // [2 4 6]
fmt.Println(odds.All())  // [1 3 5]
```

Use `Partition` to split data into two categories in a single pass.

---

## Slicing

### Slice

```go
func (c *Collection[T]) Slice(offset int, lengths ...int) *Collection[T]
```

Extracts a portion of the collection starting at the given offset. Negative offsets count from the end. An optional length limits how many items are returned.

```go
c := collection.New("a", "b", "c", "d", "e")
fmt.Println(c.Slice(1, 3).All())  // [b c d]
fmt.Println(c.Slice(-2).All())    // [d e]
```

Use `Slice` for extracting subsections from a collection.

---

### Splice

```go
func (c *Collection[T]) Splice(offset int, lengths ...int) *Collection[T]
```

Removes and returns a slice of items starting at the given offset. Mutates the original collection.

```go
c := collection.New(1, 2, 3, 4, 5)
removed := c.Splice(1, 2)
fmt.Println(removed.All()) // [2 3]
fmt.Println(c.All())       // [1 4 5]
```

Use `Splice` to extract and remove a segment from the middle of a collection.

---

### SpliceReplace

```go
func (c *Collection[T]) SpliceReplace(offset, length int, replacement []T) *Collection[T]
```

Removes a portion at the given offset and replaces it with the provided items. Returns the removed items.

```go
c := collection.New("a", "b", "c", "d")
removed := c.SpliceReplace(1, 2, []string{"X", "Y", "Z"})
fmt.Println(removed.All()) // [b c]
fmt.Println(c.All())       // [a X Y Z d]
```

Use `SpliceReplace` for surgical in-place modifications like replacing elements in a pipeline.

---

### Take

```go
func (c *Collection[T]) Take(limit int) *Collection[T]
```

Returns a new collection with the specified number of items from the front. A negative limit takes from the end.

```go
c := collection.New(1, 2, 3, 4, 5)
fmt.Println(c.Take(3).All())  // [1 2 3]
fmt.Println(c.Take(-2).All()) // [4 5]
```

Use `Take` to limit results, such as "top N" queries.

---

### TakeUntil

```go
func (c *Collection[T]) TakeUntil(callback func(T, int) bool) *Collection[T]
```

Returns items from the start until the callback returns `true` (the matching item is excluded).

```go
c := collection.New(1, 2, 3, 4, 5)
result := c.TakeUntil(func(n int, _ int) bool { return n > 3 })
fmt.Println(result.All()) // [1 2 3]
```

Use `TakeUntil` to collect items up to a boundary condition.

---

### TakeWhile

```go
func (c *Collection[T]) TakeWhile(callback func(T, int) bool) *Collection[T]
```

Returns items from the start as long as the callback returns `true`. Stops at the first `false`.

```go
c := collection.New(1, 2, 3, 4, 5)
result := c.TakeWhile(func(n int, _ int) bool { return n < 4 })
fmt.Println(result.All()) // [1 2 3]
```

Use `TakeWhile` to collect a leading run of items matching a condition.

---

### Skip

```go
func (c *Collection[T]) Skip(count int) *Collection[T]
```

Returns a new collection with the first `count` items removed.

```go
c := collection.New(1, 2, 3, 4, 5)
fmt.Println(c.Skip(2).All()) // [3 4 5]
```

Use `Skip` for pagination offsets or discarding headers.

---

### SkipUntil

```go
func (c *Collection[T]) SkipUntil(callback func(T, int) bool) *Collection[T]
```

Skips items until the callback returns `true`, then returns the rest (including the matching item).

```go
c := collection.New(1, 2, 3, 4, 5)
result := c.SkipUntil(func(n int, _ int) bool { return n == 3 })
fmt.Println(result.All()) // [3 4 5]
```

Use `SkipUntil` to start processing from a known marker.

---

### SkipWhile

```go
func (c *Collection[T]) SkipWhile(callback func(T, int) bool) *Collection[T]
```

Skips items while the callback returns `true`, then returns the rest.

```go
c := collection.New(1, 2, 3, 4, 5)
result := c.SkipWhile(func(n int, _ int) bool { return n < 3 })
fmt.Println(result.All()) // [3 4 5]
```

Use `SkipWhile` to drop a leading run of items matching a condition.

---

### Nth

```go
func (c *Collection[T]) Nth(step int, offsets ...int) *Collection[T]
```

Returns a new collection containing every `step`-th element, starting at an optional offset.

```go
c := collection.New("a", "b", "c", "d", "e", "f")
fmt.Println(c.Nth(2).All())    // [a c e]
fmt.Println(c.Nth(2, 1).All()) // [b d f]
```

Use `Nth` for sampling at regular intervals or selecting alternating elements.

---

### ForPage

```go
func (c *Collection[T]) ForPage(page, perPage int) *Collection[T]
```

Returns a subset of items for the given 1-based page number and page size.

```go
c := collection.Range(1, 20)
page2 := c.ForPage(2, 5)
fmt.Println(page2.All()) // [6 7 8 9 10]
```

Use `ForPage` for pagination.

---

## Combining

### Concat

```go
func (c *Collection[T]) Concat(items []T) *Collection[T]
```

Returns a new collection with the given items appended.

```go
a := collection.New(1, 2)
b := a.Concat([]int{3, 4})
fmt.Println(b.All()) // [1 2 3 4]
```

Use `Concat` to join two datasets.

---

### Merge

```go
func (c *Collection[T]) Merge(items []T) *Collection[T]
```

Returns a new collection with the given items appended. Alias for `Concat`.

---

### Pad

```go
func (c *Collection[T]) Pad(size int, value T) *Collection[T]
```

Pads the collection to the specified length with the given value. A positive size pads on the right; a negative size pads on the left.

```go
c := collection.New(1, 2, 3)
fmt.Println(c.Pad(5, 0).All())  // [1 2 3 0 0]
fmt.Println(c.Pad(-5, 0).All()) // [0 0 1 2 3]
```

Use `Pad` to ensure a collection meets a minimum size.

---

### Multiply

```go
func (c *Collection[T]) Multiply(multiplier int) *Collection[T]
```

Returns a new collection with all items repeated the given number of times.

```go
c := collection.New("a", "b")
fmt.Println(c.Multiply(3).All()) // [a b a b a b]
```

Use `Multiply` for repeating patterns or filling templates.

---

### Zip

```go
func Zip[T any](c *Collection[T], others ...[]T) *Collection[[]T]
```

Merges the collection with each of the given slices element-by-element, producing a collection of grouped slices.

```go
names := collection.New("Alice", "Bob", "Carol")
ages := []string{"30", "25", "35"}
zipped := collection.Zip(names, ages)
// [["Alice" "30"] ["Bob" "25"] ["Carol" "35"]]
```

Use `Zip` to combine parallel arrays into rows.

---

### CrossJoin

```go
func CrossJoin[T any](c *Collection[T], others ...[]T) *Collection[[]T]
```

Returns the cross product (Cartesian product) of the collection with the given slices.

```go
sizes := collection.New("S", "M", "L")
colors := []string{"Red", "Blue"}
combos := collection.CrossJoin(sizes, colors)
// [["S" "Red"] ["S" "Blue"] ["M" "Red"] ["M" "Blue"] ["L" "Red"] ["L" "Blue"]]
```

Use `CrossJoin` for generating all possible combinations, such as product variants.

---

### Combine

```go
func Combine[K any, V any](keys *Collection[K], values []V) *Collection[Pair[K, V]]
```

Pairs keys from the collection with values from the given slice, returning a collection of `Pair[K, V]`.

```go
keys := collection.New("name", "age", "city")
vals := []string{"Alice", "30", "NYC"}
pairs := collection.Combine(keys, vals)
// [{name Alice} {age 30} {city NYC}]
```

Use `Combine` to create key-value associations from two parallel lists.

---

### Collapse

```go
func Collapse[T any](c *Collection[[]T]) *Collection[T]
```

Flattens a collection of slices into a single, flat collection.

```go
nested := collection.New([]int{1, 2}, []int{3, 4}, []int{5})
flat := collection.Collapse(nested)
fmt.Println(flat.All()) // [1 2 3 4 5]
```

Use `Collapse` to flatten one level of nesting.

---

### Diff

```go
func Diff[T comparable](c *Collection[T], items []T) *Collection[T]
```

Returns items in the collection that are not present in the given slice.

```go
c := collection.New(1, 2, 3, 4, 5)
diff := collection.Diff(c, []int{2, 4})
fmt.Println(diff.All()) // [1 3 5]
```

Use `Diff` to find what is unique to the collection.

---

### DiffUsing

```go
func (c *Collection[T]) DiffUsing(items []T, equals func(T, T) bool) *Collection[T]
```

Like `Diff`, but uses a custom equality function instead of `==`.

```go
type Item struct{ ID int; Name string }

existing := collection.New(Item{1, "A"}, Item{2, "B"}, Item{3, "C"})
toRemove := []Item{{2, "B"}}
remaining := existing.DiffUsing(toRemove, func(a, b Item) bool {
    return a.ID == b.ID
})
fmt.Println(remaining.Count()) // 2
```

Use `DiffUsing` when items are not comparable with `==` or when you need custom matching.

---

### Intersect

```go
func Intersect[T comparable](c *Collection[T], items []T) *Collection[T]
```

Returns items present in both the collection and the given slice.

```go
c := collection.New(1, 2, 3, 4, 5)
common := collection.Intersect(c, []int{2, 4, 6})
fmt.Println(common.All()) // [2 4]
```

Use `Intersect` to find shared elements between two datasets.

---

### IntersectUsing

```go
func (c *Collection[T]) IntersectUsing(items []T, equals func(T, T) bool) *Collection[T]
```

Like `Intersect`, but uses a custom equality function.

```go
type Tag struct{ Slug string }

mine := collection.New(Tag{"go"}, Tag{"rust"}, Tag{"python"})
theirs := []Tag{{Slug: "go"}, {Slug: "java"}}
shared := mine.IntersectUsing(theirs, func(a, b Tag) bool {
    return a.Slug == b.Slug
})
fmt.Println(shared.Count()) // 1
```

Use `IntersectUsing` for custom struct comparisons.

---

## String

### Implode

```go
func (c *Collection[T]) Implode(glue string) string
```

Joins all items into a single string separated by the given glue. Each item is converted via `fmt.Sprint`.

```go
c := collection.New("Go", "is", "great")
fmt.Println(c.Implode(" ")) // "Go is great"
```

Use `Implode` for simple string concatenation.

---

### Join

```go
func (c *Collection[T]) Join(glue string, finalGlues ...string) string
```

Joins all items into a string separated by `glue`. An optional final glue is placed between the last two items.

```go
c := collection.New("Alice", "Bob", "Carol")
fmt.Println(c.Join(", ", ", and "))
// "Alice, Bob, and Carol"
```

Use `Join` for human-readable lists with a natural language connector.

---

## Conditional

### When

```go
func (c *Collection[T]) When(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T]
```

Applies the callback if the condition is `true`. An optional default callback is applied when the condition is `false`.

```go
isAdmin := true
c := collection.New("read", "write", "delete").
    When(isAdmin, func(c *collection.Collection[string]) *collection.Collection[string] {
        return c.Push("admin")
    })
```

Use `When` for conditional processing in a fluent chain.

---

### WhenEmpty

```go
func (c *Collection[T]) WhenEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T]
```

Applies the callback when the collection is empty.

```go
c := collection.Empty[string]().
    WhenEmpty(func(c *collection.Collection[string]) *collection.Collection[string] {
        return c.Push("default value")
    })
fmt.Println(c.All()) // [default value]
```

Use `WhenEmpty` to provide fallback data.

---

### WhenNotEmpty

```go
func (c *Collection[T]) WhenNotEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T]
```

Applies the callback when the collection is not empty.

```go
c := collection.New(1, 2, 3).
    WhenNotEmpty(func(c *collection.Collection[int]) *collection.Collection[int] {
        return c.Push(4)
    })
fmt.Println(c.Count()) // 4
```

Use `WhenNotEmpty` to perform processing only when there is data.

---

### Unless

```go
func (c *Collection[T]) Unless(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T]
```

Applies the callback unless the condition is `true`. The inverse of `When`.

```go
isReadOnly := false
c := collection.New("data").
    Unless(isReadOnly, func(c *collection.Collection[string]) *collection.Collection[string] {
        return c.Push("appended")
    })
```

Use `Unless` when the inverted condition reads more naturally.

---

### UnlessEmpty

```go
func (c *Collection[T]) UnlessEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T]
```

Applies the callback unless the collection is empty. Equivalent to `WhenNotEmpty`.

---

### UnlessNotEmpty

```go
func (c *Collection[T]) UnlessNotEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T]
```

Applies the callback unless the collection is not empty. Equivalent to `WhenEmpty`.

---

## Aggregation

### Sum

```go
func Sum[T Numeric](c *Collection[T]) T
```

Returns the sum of all items. `T` must satisfy the `Numeric` constraint (any integer or float type).

```go
c := collection.New(10, 20, 30)
fmt.Println(collection.Sum(c)) // 60
```

Use `Sum` for totaling numeric collections.

---

### SumBy

```go
func SumBy[T any, N Numeric](c *Collection[T], valueFunc func(T) N) N
```

Returns the sum of values extracted from each item by the given function.

```go
type LineItem struct {
    Product string
    Price   float64
}

items := collection.New(
    LineItem{"Widget", 9.99},
    LineItem{"Gadget", 24.99},
)
total := collection.SumBy(items, func(li LineItem) float64 {
    return li.Price
})
fmt.Println(total) // 34.98
```

Use `SumBy` to total a specific field from struct collections.

---

### Avg

```go
func Avg[T Numeric](c *Collection[T]) float64
```

Returns the arithmetic mean of all items.

```go
c := collection.New(10.0, 20.0, 30.0)
fmt.Println(collection.Avg(c)) // 20
```

Use `Avg` to compute the average of a numeric collection.

---

### AvgBy

```go
func AvgBy[T any, N Numeric](c *Collection[T], valueFunc func(T) N) float64
```

Returns the arithmetic mean of values extracted by the given function.

```go
type Score struct {
    Student string
    Points  int
}

scores := collection.New(Score{"A", 80}, Score{"B", 90}, Score{"C", 70})
avg := collection.AvgBy(scores, func(s Score) int { return s.Points })
fmt.Println(avg) // 80
```

Use `AvgBy` to compute averages over a specific field.

---

### Average

```go
func Average[T Numeric](c *Collection[T]) float64
```

Alias for `Avg`.

---

### Min

```go
func Min[T cmp.Ordered](c *Collection[T]) (T, bool)
```

Returns the minimum value. `T` must satisfy `cmp.Ordered`. The second return value indicates whether the collection was non-empty.

```go
c := collection.New(5, 3, 8, 1, 4)
min, ok := collection.Min(c)
fmt.Println(min, ok) // 1 true
```

Use `Min` to find the smallest value in a comparable collection.

---

### MinBy

```go
func MinBy[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) (T, bool)
```

Returns the item with the minimum key as determined by the given function.

```go
type Employee struct {
    Name   string
    Salary float64
}

employees := collection.New(
    Employee{"Alice", 75000},
    Employee{"Bob", 55000},
    Employee{"Carol", 90000},
)
lowest, _ := collection.MinBy(employees, func(e Employee) float64 {
    return e.Salary
})
fmt.Println(lowest.Name) // Bob
```

Use `MinBy` to find the item with the lowest value in a specific field.

---

### Max

```go
func Max[T cmp.Ordered](c *Collection[T]) (T, bool)
```

Returns the maximum value. The second return value indicates whether the collection was non-empty.

```go
c := collection.New(5, 3, 8, 1, 4)
max, ok := collection.Max(c)
fmt.Println(max, ok) // 8 true
```

Use `Max` to find the largest value in a comparable collection.

---

### MaxBy

```go
func MaxBy[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) (T, bool)
```

Returns the item with the maximum key as determined by the given function.

```go
highest, _ := collection.MaxBy(employees, func(e Employee) float64 {
    return e.Salary
})
fmt.Println(highest.Name) // Carol
```

Use `MaxBy` to find the item with the highest value in a specific field.

---

### Median

```go
func Median(c *Collection[float64]) float64
```

Returns the median value. Works on `*Collection[float64]`. Returns 0 for an empty collection.

```go
c := collection.New(1.0, 3.0, 2.0, 5.0, 4.0)
fmt.Println(collection.Median(c)) // 3
```

Use `Median` for statistical analysis where the middle value is more representative than the mean.

---

### MedianBy

```go
func MedianBy[T any](c *Collection[T], valueFunc func(T) float64) float64
```

Returns the median of values extracted from each item by the given function.

```go
type Response struct {
    LatencyMs float64
}

responses := collection.New(
    Response{120}, Response{80}, Response{200}, Response{95}, Response{150},
)
median := collection.MedianBy(responses, func(r Response) float64 {
    return r.LatencyMs
})
fmt.Println(median) // 120
```

Use `MedianBy` for computing median values of struct fields.

---

### Mode

```go
func Mode[T comparable](c *Collection[T]) []T
```

Returns the most frequently occurring values. Returns multiple values in case of a tie.

```go
c := collection.New("a", "b", "a", "c", "b", "a")
fmt.Println(collection.Mode(c)) // [a]
```

Use `Mode` to find the most common elements in a dataset.

---

### CountBy

```go
func CountBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]int
```

Counts how many items produce each key from the given function.

```go
type LogEntry struct {
    Level string
    Msg   string
}

logs := collection.New(
    LogEntry{"error", "timeout"},
    LogEntry{"info", "started"},
    LogEntry{"error", "crash"},
    LogEntry{"info", "finished"},
)
counts := collection.CountBy(logs, func(l LogEntry) string {
    return l.Level
})
fmt.Println(counts) // map[error:2 info:2]
```

Use `CountBy` for frequency analysis and histograms.

---

### GroupBy

```go
func GroupBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]*Collection[T]
```

Groups items by a key returned from the given function, producing a map of collections.

```go
type Order struct {
    Customer string
    Amount   float64
}

orders := collection.New(
    Order{"Alice", 100}, Order{"Bob", 50},
    Order{"Alice", 200}, Order{"Bob", 75},
)
grouped := collection.GroupBy(orders, func(o Order) string {
    return o.Customer
})
fmt.Println(grouped["Alice"].Count()) // 2
```

Use `GroupBy` to organize items into categories.

---

### KeyBy

```go
func KeyBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]T
```

Indexes the collection by a key, producing a map. If keys collide, the later value wins.

```go
type User struct {
    ID   int
    Name string
}

users := collection.New(User{1, "Alice"}, User{2, "Bob"})
byID := collection.KeyBy(users, func(u User) int { return u.ID })
fmt.Println(byID[1].Name) // Alice
```

Use `KeyBy` to build lookup maps from collections.

---

### MapToDictionary

```go
func MapToDictionary[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K][]V
```

Maps each item to a key-value pair and groups values by key. Similar to `GroupBy` but extracts a specific value.

```go
type Sale struct {
    Region string
    Amount float64
}

sales := collection.New(
    Sale{"East", 100}, Sale{"West", 200},
    Sale{"East", 150}, Sale{"West", 50},
)
byRegion := collection.MapToDictionary(sales, func(s Sale) (string, float64) {
    return s.Region, s.Amount
})
fmt.Println(byRegion["East"]) // [100 150]
```

Use `MapToDictionary` to group extracted values by a key.

---

### MapToGroups

```go
func MapToGroups[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K][]V
```

Alias for `MapToDictionary`.

---

### MapWithKeys

```go
func MapWithKeys[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K]V
```

Maps each item to a key-value pair, returning a flat map. If keys collide, the later value wins.

```go
type Config struct {
    Key   string
    Value string
}

configs := collection.New(Config{"host", "localhost"}, Config{"port", "8080"})
m := collection.MapWithKeys(configs, func(c Config) (string, string) {
    return c.Key, c.Value
})
fmt.Println(m["host"]) // localhost
```

Use `MapWithKeys` to convert a collection into a map with custom key-value extraction.

---

### Pluck

```go
func Pluck[T any, V any](c *Collection[T], valueFunc func(T) V) *Collection[V]
```

Extracts a value from each item using the given function, returning a new collection.

```go
type Article struct {
    Title  string
    Author string
}

articles := collection.New(
    Article{"Go Generics", "Alice"},
    Article{"Iterators", "Bob"},
)
titles := collection.Pluck(articles, func(a Article) string {
    return a.Title
})
fmt.Println(titles.All()) // [Go Generics Iterators]
```

Use `Pluck` to extract a single field from each item in a collection of structs.

---

## Serialization

### ToSlice

```go
func (c *Collection[T]) ToSlice() []T
```

Returns a copy of the underlying slice.

```go
c := collection.New(1, 2, 3)
s := c.ToSlice()
s[0] = 99
fmt.Println(c.All()) // [1 2 3] (original unchanged)
```

Use `ToSlice` when you need a defensive copy to avoid shared mutation.

---

### ToJSON

```go
func (c *Collection[T]) ToJSON() ([]byte, error)
```

Serializes the collection to JSON bytes.

```go
c := collection.New("a", "b", "c")
data, _ := c.ToJSON()
fmt.Println(string(data)) // ["a","b","c"]
```

Use `ToJSON` for API responses or data serialization.

---

### ToPrettyJSON

```go
func (c *Collection[T]) ToPrettyJSON() ([]byte, error)
```

Serializes the collection to indented JSON bytes (4-space indent).

```go
c := collection.New(1, 2, 3)
data, _ := c.ToPrettyJSON()
fmt.Println(string(data))
// [
//     1,
//     2,
//     3
// ]
```

Use `ToPrettyJSON` for human-readable output, debugging, or config file generation.

---

### String

```go
func (c *Collection[T]) String() string
```

Returns the JSON string representation of the collection. Returns `"[]"` on encoding error.

```go
c := collection.New(1, 2, 3)
fmt.Println(c.String()) // [1,2,3]
```

Implements `fmt.Stringer` for convenient printing.

---

### MarshalJSON

```go
func (c *Collection[T]) MarshalJSON() ([]byte, error)
```

Implements the `json.Marshaler` interface, allowing collections to be directly marshaled.

```go
type Response struct {
    Items *collection.Collection[string] `json:"items"`
}

resp := Response{Items: collection.New("a", "b")}
data, _ := json.Marshal(resp)
fmt.Println(string(data)) // {"items":["a","b"]}
```

---

### UnmarshalJSON

```go
func (c *Collection[T]) UnmarshalJSON(data []byte) error
```

Implements the `json.Unmarshaler` interface, allowing JSON to be decoded directly into a collection.

```go
var c collection.Collection[int]
json.Unmarshal([]byte("[1,2,3]"), &c)
fmt.Println(c.All()) // [1 2 3]
```

---

### Copy

```go
func (c *Collection[T]) Copy() *Collection[T]
```

Creates a shallow copy of the collection.

```go
original := collection.New(1, 2, 3)
clone := original.Copy()
clone.Push(4)
fmt.Println(original.Count()) // 3 (unchanged)
fmt.Println(clone.Count())    // 4
```

Use `Copy` when you need to modify a collection without affecting the original.

---

### Dump

```go
func (c *Collection[T]) Dump() *Collection[T]
```

Prints the collection items to stdout for debugging. Returns the collection for chaining.

```go
collection.New(1, 2, 3).
    Dump().           // prints: [1 2 3]
    Filter(func(n int, _ int) bool { return n > 1 }).
    Dump()            // prints: [2 3]
```

Use `Dump` for quick debugging in a fluent chain.

---

### Len

```go
func (c *Collection[T]) Len() int
```

Returns the number of items. Identical to `Count` but satisfies the `sort.Interface` convention.

---

## Conversion

### Lazy

```go
func (c *Collection[T]) Lazy() *LazyCollection[T]
```

Converts the collection into a `LazyCollection[T]` for deferred evaluation.

```go
c := collection.New(1, 2, 3, 4, 5)
lazy := c.Lazy()
// No work done until you consume:
result := lazy.Filter(func(n int, _ int) bool { return n > 2 }).All()
fmt.Println(result) // [3 4 5]
```

Use `Lazy` when you want to chain multiple operations on a large dataset without creating intermediate slices.

---

### ToBase

```go
func (c *Collection[T]) ToBase() *Collection[T]
```

Returns the collection itself. Provided for API completeness.

---

### Dot

```go
func (c *Collection[T]) Dot() *Collection[T]
```

Returns a shallow copy. For typed Go slices, dot-notation expansion is not applicable.

---

### Undot

```go
func (c *Collection[T]) Undot() *Collection[T]
```

Returns a shallow copy. For typed Go slices, dot-notation expansion is not applicable.

---

### Ensure

```go
func (c *Collection[T]) Ensure(predicate func(T) bool) error
```

Verifies that all items satisfy the predicate. Returns an error if any item fails the check.

```go
ages := collection.New(18, 25, 16, 30)
err := ages.Ensure(func(age int) bool { return age >= 18 })
if err != nil {
    fmt.Println(err) // "collection item failed ensure check"
}
```

Use `Ensure` for validation, such as asserting all items meet a constraint before processing.

---

### Every

```go
func (c *Collection[T]) Every(callback func(T, int) bool) bool
```

Reports whether all items satisfy the predicate.

```go
c := collection.New(2, 4, 6, 8)
allEven := c.Every(func(n int, _ int) bool { return n%2 == 0 })
fmt.Println(allEven) // true
```

Use `Every` for "all match" checks where you need a boolean result.

---

### Only

```go
func (c *Collection[T]) Only(indices ...int) *Collection[T]
```

Returns a new collection containing only items at the given indices.

```go
c := collection.New("a", "b", "c", "d", "e")
fmt.Println(c.Only(0, 2, 4).All()) // [a c e]
```

Use `Only` to cherry-pick specific positions from a collection.

---

### Except

```go
func (c *Collection[T]) Except(indices ...int) *Collection[T]
```

Returns a new collection excluding items at the given indices.

```go
c := collection.New("a", "b", "c", "d", "e")
fmt.Println(c.Except(1, 3).All()) // [a c e]
```

Use `Except` to remove items at known positions without mutating the original.
