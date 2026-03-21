# LazyCollection[T]

`LazyCollection` is a lazily-evaluated sequence backed by `iter.Seq[T]`. Items are computed on demand -- nothing runs until you consume the collection. This makes it ideal for large datasets, pipelines with early termination, and infinite or generated sequences.

Because operations like `Filter`, `Take`, and `Skip` return new lazy collections rather than materializing results, you can compose long pipelines that only process the elements actually needed. Memory usage stays proportional to what you consume, not to the total dataset size.

```go
import "github.com/gocanto/collection"
```

---

## Table of Contents

- [Constructors](#constructors)
- [Materialization](#materialization)
- [Query](#query)
- [Retrieval](#retrieval)
- [Search](#search)
- [Iteration](#iteration)
- [Filtering](#filtering)
- [Transformation](#transformation)
- [Partitioning](#partitioning)
- [Slicing](#slicing)
- [Combining](#combining)
- [Grouping](#grouping)
- [Caching](#caching)
- [String](#string)
- [Conditional](#conditional)
- [Other](#other)

---

## Why Lazy?

Consider processing a million-row CSV. With an eager collection you would load all rows into memory, filter them, then transform them. With a lazy collection, each row flows through the pipeline one at a time -- rows that fail the filter never reach the transform step, and if you only need the first 10 matches, the remaining rows are never even read.

```go
// Only processes rows until 10 matches are found.
// Rows that don't match "active" are never transformed.
results := collection.LazyFrom(millionRows).
    Filter(func(r Row, _ int) bool { return r.Status == "active" }).
    Take(10).
    All()
```

---

## Constructors

### NewLazy

```go
func NewLazy[T any](source iter.Seq[T]) *LazyCollection[T]
```

Creates a `LazyCollection` from an `iter.Seq[T]` iterator function. The source receives a `yield` function; call `yield` for each item and stop generating if `yield` returns `false`.

```go
// A lazy sequence of squares
squares := collection.NewLazy(func(yield func(int) bool) {
    for i := 1; i <= 1000; i++ {
        if !yield(i * i) {
            return
        }
    }
})

first3 := squares.Take(3).All() // [1, 4, 9]
```

**When to use:** You have a custom generator function or an existing `iter.Seq[T]` from another library.

---

### LazyFrom

```go
func LazyFrom[T any](items []T) *LazyCollection[T]
```

Creates a `LazyCollection` from a slice. The slice is iterated lazily -- no copy is made up front.

```go
names := collection.LazyFrom([]string{"Alice", "Bob", "Charlie"})
first, _ := names.First() // "Alice"
```

**When to use:** You have data in a slice but want to apply lazy pipeline operations to avoid intermediate allocations.

---

### LazyEmpty

```go
func LazyEmpty[T any]() *LazyCollection[T]
```

Creates an empty `LazyCollection`. Useful as a starting point or default value.

```go
empty := collection.LazyEmpty[int]()
fmt.Println(empty.IsEmpty()) // true
```

---

### LazyRange

```go
func LazyRange(from, to int) *LazyCollection[int]
```

Creates a lazy collection of sequential integers from `from` to `to` (inclusive). If `from > to`, the sequence counts downward.

```go
ascending := collection.LazyRange(1, 5).All()  // [1, 2, 3, 4, 5]
descending := collection.LazyRange(5, 1).All() // [5, 4, 3, 2, 1]
```

**When to use:** Generating index ranges, pagination offsets, or numeric test data without allocating a slice.

---

### LazyTimes

```go
func LazyTimes[T any](number int, callback func(int) T) *LazyCollection[T]
```

Creates a lazy collection by invoking the callback `n` times. The callback receives a **1-based** index.

```go
labels := collection.LazyTimes(3, func(i int) string {
    return fmt.Sprintf("item-%d", i)
})
fmt.Println(labels.All()) // ["item-1", "item-2", "item-3"]
```

**When to use:** Creating test fixtures, placeholder data, or repeated computed values.

---

## Materialization

Materialization methods force evaluation of the lazy pipeline and produce concrete results.

### All

```go
func (lc *LazyCollection[T]) All() []T
```

Eagerly evaluates the lazy collection and returns all items as a slice.

```go
items := collection.LazyRange(1, 5).All() // [1, 2, 3, 4, 5]
```

---

### Eager

```go
func (lc *LazyCollection[T]) Eager() *Collection[T]
```

Eagerly evaluates the lazy collection and returns a `Collection[T]` -- the eager, feature-rich slice wrapper.

```go
c := collection.LazyFrom([]int{3, 1, 2}).Eager()
// c is now a *Collection[int] with sorting, mapping, etc.
```

**When to use:** You have finished building a lazy pipeline and need the full method set of `Collection` for subsequent operations.

---

### Collect

```go
func (lc *LazyCollection[T]) Collect() *Collection[T]
```

Alias for `Eager`. Converts the lazy collection to an eager `Collection`.

---

### Iter

```go
func (lc *LazyCollection[T]) Iter() iter.Seq[T]
```

Returns the underlying `iter.Seq[T]` iterator for use with Go 1.23+ range-over-func loops.

```go
lazy := collection.LazyFrom([]string{"a", "b", "c"})

for item := range lazy.Iter() {
    fmt.Println(item)
}
// Output:
// a
// b
// c
```

**When to use:** Integrating a lazy pipeline with standard Go `for range` syntax. Because `Iter()` returns an `iter.Seq[T]`, you get first-class compatibility with the Go 1.23+ iterator protocol. Consumers can break out of the loop early and the upstream pipeline stops immediately -- no wasted work.

```go
// Processing stops as soon as we find what we need.
for val := range collection.LazyRange(1, 1_000_000).Iter() {
    if val > 10 {
        break // only 11 items were generated
    }
}
```

---

## Query

### Count

```go
func (lc *LazyCollection[T]) Count() int
```

Returns the total number of items by eagerly evaluating the entire collection.

```go
count := collection.LazyRange(1, 100).Count() // 100
```

> **Note:** This consumes the full sequence. For single-pass iterators, subsequent calls re-evaluate from the source.

---

### IsEmpty

```go
func (lc *LazyCollection[T]) IsEmpty() bool
```

Reports whether the lazy collection contains no items. Only evaluates the first element to determine emptiness.

```go
fmt.Println(collection.LazyEmpty[int]().IsEmpty())   // true
fmt.Println(collection.LazyRange(1, 5).IsEmpty())    // false
```

---

### IsNotEmpty

```go
func (lc *LazyCollection[T]) IsNotEmpty() bool
```

Reports whether the lazy collection contains at least one item.

---

### ContainsOneItem

```go
func (lc *LazyCollection[T]) ContainsOneItem() bool
```

Reports whether the lazy collection contains exactly one item.

```go
fmt.Println(collection.LazyFrom([]int{42}).ContainsOneItem()) // true
```

---

### ContainsManyItems

```go
func (lc *LazyCollection[T]) ContainsManyItems() bool
```

Reports whether the lazy collection contains more than one item.

---

### Has

```go
func (lc *LazyCollection[T]) Has(index int) bool
```

Reports whether the given zero-based index exists in the collection. Iterates up to that index.

```go
lazy := collection.LazyFrom([]string{"a", "b", "c"})
fmt.Println(lazy.Has(2))  // true
fmt.Println(lazy.Has(10)) // false
```

---

### HasAny

```go
func (lc *LazyCollection[T]) HasAny(indices ...int) bool
```

Reports whether any of the given indices exist in the collection.

```go
lazy := collection.LazyFrom([]int{10, 20, 30})
fmt.Println(lazy.HasAny(1, 5)) // true (index 1 exists)
```

---

## Retrieval

### First

```go
func (lc *LazyCollection[T]) First(predicates ...func(T, int) bool) (T, bool)
```

Returns the first element matching the optional predicate. If no predicate is given, returns the first element. The second return value indicates whether a match was found.

```go
lazy := collection.LazyFrom([]int{1, 2, 3, 4, 5})

val, _ := lazy.First() // 1

val, _ = lazy.First(func(v int, _ int) bool {
    return v > 3
})
// val=4
```

**When to use:** Finding the first match without materializing the entire sequence.

---

### FirstOrFail

```go
func (lc *LazyCollection[T]) FirstOrFail(predicates ...func(T, int) bool) (T, error)
```

Like `First`, but returns an `ItemNotFoundError` instead of a boolean when no match exists.

```go
lazy := collection.LazyFrom([]int{1, 2, 3})

val, err := lazy.FirstOrFail(func(v int, _ int) bool {
    return v > 100
})
if err != nil {
    fmt.Println(err) // "item not found"
}
```

**When to use:** When a missing element should be treated as an error rather than a normal condition.

---

### Last

```go
func (lc *LazyCollection[T]) Last(predicates ...func(T, int) bool) (T, bool)
```

Returns the last element matching the optional predicate. Requires a full pass through the sequence.

```go
lazy := collection.LazyFrom([]int{1, 2, 3, 4, 5})

val, _ := lazy.Last() // 5

val, _ = lazy.Last(func(v int, _ int) bool {
    return v%2 == 0
})
// val=4
```

---

### Sole

```go
func (lc *LazyCollection[T]) Sole(predicates ...func(T, int) bool) (T, error)
```

Returns the only element matching the optional predicate. Returns an error if zero or more than one element matches (`ItemNotFoundError` or `MultipleItemsFoundError`).

```go
lazy := collection.LazyFrom([]int{1, 2, 3})

val, err := lazy.Sole(func(v int, _ int) bool {
    return v == 2
})
// val=2, err=nil

_, err = lazy.Sole(func(v int, _ int) bool {
    return v > 1
})
// err: "multiple items found: 2 items"
```

**When to use:** Enforcing that exactly one item matches a condition. Useful for lookups that must be unique.

---

### Get

```go
func (lc *LazyCollection[T]) Get(index int) (T, bool)
```

Returns the item at the given zero-based index. Only iterates up to that index.

```go
lazy := collection.LazyRange(10, 20)

val, ok := lazy.Get(3) // val=13, ok=true
val, ok = lazy.Get(99) // ok=false
```

---

## Search

### Contains

```go
func (lc *LazyCollection[T]) Contains(predicate func(T, int) bool) bool
```

Reports whether any item satisfies the predicate. Stops iteration at the first match.

```go
lazy := collection.LazyFrom([]string{"go", "rust", "zig"})

hasGo := lazy.Contains(func(s string, _ int) bool {
    return s == "go"
})
fmt.Println(hasGo) // true
```

---

### Some

```go
func (lc *LazyCollection[T]) Some(predicate func(T, int) bool) bool
```

Alias for `Contains`.

---

### DoesntContain

```go
func (lc *LazyCollection[T]) DoesntContain(predicate func(T, int) bool) bool
```

Reports whether no item satisfies the predicate. The logical inverse of `Contains`.

```go
lazy := collection.LazyFrom([]int{2, 4, 6})

noOdds := lazy.DoesntContain(func(v int, _ int) bool {
    return v%2 != 0
})
fmt.Println(noOdds) // true
```

---

### HasSole

```go
func (lc *LazyCollection[T]) HasSole(predicates ...func(T, int) bool) bool
```

Reports whether exactly one item matches the optional predicate. If no predicate is given, checks whether the collection has exactly one item total.

```go
lazy := collection.LazyFrom([]int{1, 2, 3, 4, 5})

fmt.Println(lazy.HasSole(func(v int, _ int) bool {
    return v == 3
})) // true

fmt.Println(lazy.HasSole(func(v int, _ int) bool {
    return v > 3
})) // false (two items match: 4 and 5)
```

---

### Search

```go
func (lc *LazyCollection[T]) Search(predicate func(T, int) bool) (int, bool)
```

Returns the index of the first item satisfying the predicate. Stops iteration at the match.

```go
lazy := collection.LazyFrom([]string{"a", "b", "c", "d"})

idx, found := lazy.Search(func(s string, _ int) bool {
    return s == "c"
})
fmt.Println(idx, found) // 2 true
```

---

### Before

```go
func (lc *LazyCollection[T]) Before(predicate func(T, int) bool) (T, bool)
```

Returns the item immediately before the first item matching the predicate. Returns `false` if the match is the first item or no match is found.

```go
lazy := collection.LazyFrom([]int{10, 20, 30, 40})

val, ok := lazy.Before(func(v int, _ int) bool {
    return v == 30
})
fmt.Println(val, ok) // 20 true
```

**When to use:** Finding the predecessor of an element in a sequence -- for example, the previous step in a pipeline.

---

### After

```go
func (lc *LazyCollection[T]) After(predicate func(T, int) bool) (T, bool)
```

Returns the item immediately after the first item matching the predicate. Returns `false` if the match is the last item or no match is found.

```go
lazy := collection.LazyFrom([]int{10, 20, 30, 40})

val, ok := lazy.After(func(v int, _ int) bool {
    return v == 20
})
fmt.Println(val, ok) // 30 true
```

---

## Iteration

### Each

```go
func (lc *LazyCollection[T]) Each(callback func(T, int) bool) *LazyCollection[T]
```

Iterates over items, calling the callback for each one with the item and its zero-based index. Return `false` from the callback to stop early.

```go
collection.LazyFrom([]string{"a", "b", "c"}).Each(func(s string, i int) bool {
    fmt.Printf("[%d] %s\n", i, s)
    return true
})
// Output:
// [0] a
// [1] b
// [2] c
```

> **Note:** `Each` evaluates the sequence. It returns the original lazy collection, not a new one.

---

### Tap

```go
func (lc *LazyCollection[T]) Tap(callback func(*LazyCollection[T])) *LazyCollection[T]
```

Passes the lazy collection to the callback for side effects and returns it unchanged.

```go
result := collection.LazyFrom([]int{1, 2, 3}).
    Tap(func(lc *collection.LazyCollection[int]) {
        fmt.Println("before filter:", lc.Count())
    }).
    Filter(func(v int, _ int) bool { return v > 1 })
```

---

### TapEach

```go
func (lc *LazyCollection[T]) TapEach(callback func(T, int)) *LazyCollection[T]
```

Returns a new lazy collection that calls the callback on each item as it passes through the pipeline. Unlike `Each`, `TapEach` is lazy -- the callback only fires when the resulting collection is consumed.

```go
result := collection.LazyFrom([]int{1, 2, 3}).
    TapEach(func(v int, i int) {
        fmt.Printf("processing %d at index %d\n", v, i)
    }).
    Filter(func(v int, _ int) bool { return v > 1 }).
    All()
// "processing" is printed for each item as it flows through
```

**When to use:** Adding logging or metrics to a lazy pipeline without breaking the chain or forcing evaluation.

---

## Filtering

### Filter

```go
func (lc *LazyCollection[T]) Filter(callback func(T, int) bool) *LazyCollection[T]
```

Returns a new lazy collection containing only items for which the callback returns `true`. The filter is applied lazily -- items are tested one at a time as they are consumed.

```go
evens := collection.LazyRange(1, 100).Filter(func(v int, _ int) bool {
    return v%2 == 0
})
first5 := evens.Take(5).All() // [2, 4, 6, 8, 10]
```

---

### Reject

```go
func (lc *LazyCollection[T]) Reject(callback func(T, int) bool) *LazyCollection[T]
```

Returns a new lazy collection excluding items for which the callback returns `true`. The logical inverse of `Filter`.

```go
noNegatives := collection.LazyFrom([]int{-2, -1, 0, 1, 2}).
    Reject(func(v int, _ int) bool { return v < 0 })
fmt.Println(noNegatives.All()) // [0, 1, 2]
```

---

## Transformation

### LazyMap (top-level function)

```go
func LazyMap[T any, R any](
    lc *LazyCollection[T],
    callback func(T, int) R,
) *LazyCollection[R]
```

Transforms each item using the callback, returning a new lazy collection of the transformed type. Because Go methods cannot introduce new type parameters, this is a package-level function.

```go
names := collection.LazyFrom([]string{"alice", "bob"})

lengths := collection.LazyMap(names, func(s string, _ int) int {
    return len(s)
})
fmt.Println(lengths.All()) // [5, 3]
```

**When to use:** Mapping between types in a lazy pipeline -- for example, converting raw records to domain objects, or extracting a single field.

---

### LazyFlatMap (top-level function)

```go
func LazyFlatMap[T any, R any](
    lc *LazyCollection[T],
    callback func(T, int) []R,
) *LazyCollection[R]
```

Transforms each item into a slice and flattens the results into a single lazy sequence.

```go
sentences := collection.LazyFrom([]string{"hello world", "foo bar"})

words := collection.LazyFlatMap(sentences, func(s string, _ int) []string {
    return strings.Split(s, " ")
})
fmt.Println(words.All()) // ["hello", "world", "foo", "bar"]
```

**When to use:** One-to-many transformations where each input produces a variable number of outputs.

---

### LazyReduce (top-level function)

```go
func LazyReduce[T any, R any](
    lc *LazyCollection[T],
    callback func(R, T, int) R,
    initial R,
) R
```

Reduces the lazy collection to a single value by applying the callback to an accumulator and each item. Eagerly evaluates the entire sequence.

```go
sum := collection.LazyReduce(
    collection.LazyRange(1, 100),
    func(acc int, val int, _ int) int { return acc + val },
    0,
)
fmt.Println(sum) // 5050
```

---

### LazyUnique (top-level function)

```go
func LazyUnique[T any, K comparable](
    lc *LazyCollection[T],
    keyFunc func(T) K,
) *LazyCollection[T]
```

Returns a lazy collection containing only unique items as determined by the key function. The first occurrence of each key is kept.

```go
type User struct {
    ID   int
    Name string
}

users := collection.LazyFrom([]User{
    {1, "Alice"}, {2, "Bob"}, {1, "Alice (dup)"},
})

unique := collection.LazyUnique(users, func(u User) int {
    return u.ID
})
fmt.Println(unique.Count()) // 2
```

**When to use:** Deduplicating a stream by a key without sorting.

---

### LazyPluck (top-level function)

```go
func LazyPluck[T any, V any](
    lc *LazyCollection[T],
    valueFunc func(T) V,
) *LazyCollection[V]
```

Extracts a value from each item, returning a new lazy collection of the extracted values. This is a convenience wrapper around `LazyMap`.

```go
type Product struct {
    Name  string
    Price float64
}

products := collection.LazyFrom([]Product{
    {"Widget", 9.99},
    {"Gadget", 24.99},
})

prices := collection.LazyPluck(products, func(p Product) float64 {
    return p.Price
})
fmt.Println(prices.All()) // [9.99, 24.99]
```

---

## Partitioning

### Chunk

```go
func (lc *LazyCollection[T]) Chunk(size int) [][]T
```

Eagerly evaluates the lazy collection and splits items into groups of the given size. The last group may contain fewer items.

```go
chunks := collection.LazyRange(1, 7).Chunk(3)
// [[1, 2, 3], [4, 5, 6], [7]]
```

**When to use:** Batch processing -- for example, sending items to an API in groups of N.

---

### ChunkWhile

```go
func (lc *LazyCollection[T]) ChunkWhile(
    callback func(T, int, []T) bool,
) [][]T
```

Splits the lazy collection into groups where consecutive items satisfy the callback. A new chunk is started when the callback returns `false`. The callback receives the current item, its index, and the current chunk being built.

```go
// Group consecutive even/odd numbers together
chunks := collection.LazyFrom([]int{2, 4, 6, 1, 3, 8, 10}).
    ChunkWhile(func(val int, _ int, current []int) bool {
        return val%2 == current[0]%2
    })
// [[2, 4, 6], [1, 3], [8, 10]]
```

---

### Nth

```go
func (lc *LazyCollection[T]) Nth(step int, offsets ...int) *LazyCollection[T]
```

Returns every `step`-th element, starting from an optional offset (default 0). The result is lazy.

```go
every3rd := collection.LazyRange(1, 12).Nth(3)
fmt.Println(every3rd.All()) // [1, 4, 7, 10]

// With offset
every3rdFrom1 := collection.LazyRange(1, 12).Nth(3, 1)
fmt.Println(every3rdFrom1.All()) // [2, 5, 8, 11]
```

**When to use:** Sampling data at regular intervals or implementing striped processing.

---

## Slicing

### Take

```go
func (lc *LazyCollection[T]) Take(limit int) *LazyCollection[T]
```

Returns a new lazy collection with at most `limit` items from the beginning. A negative limit takes from the end (which requires eager evaluation of the entire sequence).

```go
first3 := collection.LazyRange(1, 1000).Take(3)
fmt.Println(first3.All()) // [1, 2, 3]

last3 := collection.LazyRange(1, 10).Take(-3)
fmt.Println(last3.All()) // [8, 9, 10]
```

**When to use:** Limiting results, implementing pagination, or previewing the first N items of a large dataset.

---

### TakeUntil

```go
func (lc *LazyCollection[T]) TakeUntil(callback func(T, int) bool) *LazyCollection[T]
```

Returns items from the beginning until the callback returns `true`. The item that triggers `true` is **not** included.

```go
items := collection.LazyFrom([]int{1, 2, 3, 4, 5}).
    TakeUntil(func(v int, _ int) bool { return v == 4 })
fmt.Println(items.All()) // [1, 2, 3]
```

---

### TakeWhile

```go
func (lc *LazyCollection[T]) TakeWhile(callback func(T, int) bool) *LazyCollection[T]
```

Returns items from the beginning while the callback returns `true`. Stops at the first `false`.

```go
items := collection.LazyFrom([]int{2, 4, 6, 7, 8}).
    TakeWhile(func(v int, _ int) bool { return v%2 == 0 })
fmt.Println(items.All()) // [2, 4, 6]
```

---

### TakeUntilTimeout

```go
func (lc *LazyCollection[T]) TakeUntilTimeout(timeout time.Duration) *LazyCollection[T]
```

Returns items until the given duration has elapsed. Useful for processing streams with a time budget.

```go
import "time"

// Process items for at most 2 seconds
results := someLazyStream.TakeUntilTimeout(2 * time.Second).All()
```

**When to use:** Rate-limited consumption, polling with a deadline, or any scenario where processing time matters more than item count.

---

### Skip

```go
func (lc *LazyCollection[T]) Skip(count int) *LazyCollection[T]
```

Returns a new lazy collection that skips the first `count` items.

```go
items := collection.LazyRange(1, 10).Skip(5)
fmt.Println(items.All()) // [6, 7, 8, 9, 10]
```

---

### SkipUntil

```go
func (lc *LazyCollection[T]) SkipUntil(callback func(T, int) bool) *LazyCollection[T]
```

Skips items until the callback returns `true`, then yields that item and all subsequent items.

```go
items := collection.LazyFrom([]int{1, 2, 3, 4, 5}).
    SkipUntil(func(v int, _ int) bool { return v >= 3 })
fmt.Println(items.All()) // [3, 4, 5]
```

---

### SkipWhile

```go
func (lc *LazyCollection[T]) SkipWhile(callback func(T, int) bool) *LazyCollection[T]
```

Skips items while the callback returns `true`, then yields the rest.

```go
items := collection.LazyFrom([]int{1, 2, 3, 4, 5}).
    SkipWhile(func(v int, _ int) bool { return v < 3 })
fmt.Println(items.All()) // [3, 4, 5]
```

---

### Slice

```go
func (lc *LazyCollection[T]) Slice(offset int, lengths ...int) *LazyCollection[T]
```

Returns a subset of the lazy collection starting at `offset` with an optional length. Equivalent to chaining `Skip(offset).Take(length)`.

```go
items := collection.LazyRange(1, 20).Slice(5, 3)
fmt.Println(items.All()) // [6, 7, 8]
```

---

## Combining

### Concat

```go
func (lc *LazyCollection[T]) Concat(items []T) *LazyCollection[T]
```

Returns a new lazy collection with the given items appended after the current sequence.

```go
combined := collection.LazyFrom([]int{1, 2}).Concat([]int{3, 4, 5})
fmt.Println(combined.All()) // [1, 2, 3, 4, 5]
```

---

### Pad

```go
func (lc *LazyCollection[T]) Pad(size int, value T) *LazyCollection[T]
```

Returns a new lazy collection padded to the given size with the specified value.
- **Positive size:** pads at the end.
- **Negative size:** pads at the beginning (requires eager evaluation).

```go
padded := collection.LazyFrom([]int{1, 2}).Pad(5, 0)
fmt.Println(padded.All()) // [1, 2, 0, 0, 0]

leftPadded := collection.LazyFrom([]int{1, 2}).Pad(-5, 0)
fmt.Println(leftPadded.All()) // [0, 0, 0, 1, 2]
```

**When to use:** Ensuring a minimum length -- for example, filling a table row to a fixed column count.

---

## Grouping

### LazyGroupBy (top-level function)

```go
func LazyGroupBy[T any, K comparable](
    lc *LazyCollection[T],
    keyFunc func(T) K,
) map[K]*LazyCollection[T]
```

Groups items by the key returned by the given function. This requires eager evaluation to build the groups. Each group is returned as a `LazyCollection`.

```go
type Order struct {
    Status string
    Total  float64
}

orders := collection.LazyFrom([]Order{
    {"pending", 100}, {"shipped", 200}, {"pending", 50},
})

grouped := collection.LazyGroupBy(orders, func(o Order) string {
    return o.Status
})

pending := grouped["pending"].All()
// [{pending 100}, {pending 50}]
```

---

### LazyKeyBy (top-level function)

```go
func LazyKeyBy[T any, K comparable](
    lc *LazyCollection[T],
    keyFunc func(T) K,
) map[K]T
```

Indexes items by the key returned by the given function. If duplicate keys exist, the later value overwrites the earlier one.

```go
type User struct {
    ID   int
    Name string
}

users := collection.LazyFrom([]User{{1, "Alice"}, {2, "Bob"}})

byID := collection.LazyKeyBy(users, func(u User) int {
    return u.ID
})
fmt.Println(byID[1].Name) // "Alice"
```

---

### LazyCountBy (top-level function)

```go
func LazyCountBy[T any, K comparable](
    lc *LazyCollection[T],
    keyFunc func(T) K,
) map[K]int
```

Counts occurrences of each key returned by the given function.

```go
words := collection.LazyFrom([]string{"go", "rust", "go", "zig", "go"})

counts := collection.LazyCountBy(words, func(s string) string {
    return s
})
fmt.Println(counts) // map[go:3 rust:1 zig:1]
```

---

## Caching

### Remember

```go
func (lc *LazyCollection[T]) Remember() *LazyCollection[T]
```

Returns a lazy collection that caches items on first iteration. Subsequent iterations reuse the cached values instead of re-evaluating the source.

```go
expensive := collection.NewLazy(func(yield func(int) bool) {
    for i := 0; i < 3; i++ {
        fmt.Println("computing", i)
        if !yield(i) {
            return
        }
    }
}).Remember()

expensive.All() // prints "computing 0", "computing 1", "computing 2"
expensive.All() // no output -- results are cached
```

**When to use:** When the source is expensive to evaluate (I/O, network calls, heavy computation) and you need to iterate more than once.

---

### Throttle

```go
func (lc *LazyCollection[T]) Throttle(delay time.Duration) *LazyCollection[T]
```

Returns a new lazy collection that inserts a delay between each yielded item. The delay applies lazily as items are consumed.

```go
import "time"

// Rate-limit API calls to 1 per second
results := collection.LazyFrom(apiRequests).
    Throttle(time.Second).
    All()
```

**When to use:** Rate limiting, gentle polling, or preventing thundering-herd issues when processing items that trigger external calls.

---

## String

### Implode

```go
func (lc *LazyCollection[T]) Implode(glue string) string
```

Joins all items into a single string separated by the given glue. Items are converted to strings using `fmt.Sprint`.

```go
result := collection.LazyFrom([]int{1, 2, 3}).Implode(", ")
fmt.Println(result) // "1, 2, 3"
```

---

### Join

```go
func (lc *LazyCollection[T]) Join(glue string, finalGlues ...string) string
```

Like `Implode`, but with an optional final glue placed between the last two items.

```go
result := collection.LazyFrom([]string{"Go", "Rust", "Zig"}).
    Join(", ", ", and ")
fmt.Println(result) // "Go, Rust, and Zig"
```

---

## Conditional

### When

```go
func (lc *LazyCollection[T]) When(
    condition bool,
    callback func(*LazyCollection[T]) *LazyCollection[T],
    defaults ...func(*LazyCollection[T]) *LazyCollection[T],
) *LazyCollection[T]
```

Applies the callback if the condition is `true`. An optional default callback is applied when the condition is `false`.

```go
includeArchived := false

orders := collection.LazyFrom(allOrders).
    When(includeArchived, func(lc *collection.LazyCollection[Order]) *collection.LazyCollection[Order] {
        return lc // include everything
    }, func(lc *collection.LazyCollection[Order]) *collection.LazyCollection[Order] {
        return lc.Filter(func(o Order, _ int) bool {
            return o.Status != "archived"
        })
    })
```

**When to use:** Conditionally modifying a pipeline without breaking the method chain.

---

### WhenEmpty

```go
func (lc *LazyCollection[T]) WhenEmpty(
    callback func(*LazyCollection[T]) *LazyCollection[T],
    defaults ...func(*LazyCollection[T]) *LazyCollection[T],
) *LazyCollection[T]
```

Applies the callback if the lazy collection is empty.

```go
results := fetchResults().
    WhenEmpty(func(lc *collection.LazyCollection[Result]) *collection.LazyCollection[Result] {
        return collection.LazyFrom(defaultResults)
    })
```

---

### WhenNotEmpty

```go
func (lc *LazyCollection[T]) WhenNotEmpty(
    callback func(*LazyCollection[T]) *LazyCollection[T],
    defaults ...func(*LazyCollection[T]) *LazyCollection[T],
) *LazyCollection[T]
```

Applies the callback if the lazy collection is not empty.

---

### Unless

```go
func (lc *LazyCollection[T]) Unless(
    condition bool,
    callback func(*LazyCollection[T]) *LazyCollection[T],
    defaults ...func(*LazyCollection[T]) *LazyCollection[T],
) *LazyCollection[T]
```

The inverse of `When`: applies the callback unless the condition is `true`.

---

## Other

### Every

```go
func (lc *LazyCollection[T]) Every(callback func(T, int) bool) bool
```

Reports whether all items satisfy the callback. Stops iteration at the first failure.

```go
allPositive := collection.LazyFrom([]int{1, 2, 3, 4}).
    Every(func(v int, _ int) bool { return v > 0 })
fmt.Println(allPositive) // true
```

---

### Dump

```go
func (lc *LazyCollection[T]) Dump() *LazyCollection[T]
```

Prints the items to stdout for debugging. Eagerly evaluates the sequence and returns a new lazy collection backed by the evaluated items.

```go
collection.LazyFrom([]int{1, 2, 3}).Dump()
// Output: [1 2 3]
```

> **Note:** Because `Dump` must evaluate the sequence to print it, the returned collection is backed by the materialized slice, not the original source.
