# arr -- Generic Slice Utilities

```
import "github.com/gocanto/collection/arr"
```

The `arr` package provides standalone, generic helper functions for working with Go slices. Every function operates on plain slices, returns new slices (the original is never mutated), and can be called without constructing a collection object.

Use `arr` when you need a single operation on a slice and do not need the fluent chaining that `collection.Collection` provides.

---

## Table of Contents

- [Accessible](#accessible)
- [IsList](#islist)
- [First](#first)
- [Last](#last)
- [Take](#take)
- [Only](#only)
- [Except](#except)
- [Flatten](#flatten)
- [Collapse](#collapse)
- [Wrap](#wrap)
- [WrapSlice](#wrapslice)
- [Prepend](#prepend)
- [Push](#push)
- [Shuffle](#shuffle)
- [Random](#random)
- [Sort](#sort)
- [SortDesc](#sortdesc)
- [SortRecursive](#sortrecursive)
- [Where](#where)
- [WhereNotNull](#wherenotnull)
- [Reject](#reject)
- [Partition](#partition)
- [Every](#every)
- [Some](#some)
- [Exists](#exists)
- [Has](#has)
- [HasAny](#hasany)
- [Join](#join)
- [CrossJoin](#crossjoin)
- [Divide](#divide)
- [Map](#map)
- [MapWithKeys](#mapwithkeys)
- [MapSpread](#mapspread)
- [KeyBy](#keyby)
- [Pluck](#pluck)

---

## Accessible

```go
func Accessible(value any) bool
```

Reports whether the given value is non-nil. This is useful as a guard before performing operations on interface values that may be nil.

```go
var s []int
arr.Accessible(s)   // false -- nil slice
arr.Accessible(42)  // true

m := map[string]int{}
arr.Accessible(m)   // true
```

**Why it is useful:** Provides a clean, intention-revealing nil check when working with `any`-typed values received from configuration maps, JSON decoding, or other untyped sources.

---

## IsList

```go
func IsList[T any](items []T) bool
```

Reports whether the given slice is a sequential list. In Go every slice is inherently a sequential, integer-indexed list, so this always returns `true`. It exists for API compatibility with the PHP Laravel collection counterpart.

```go
arr.IsList([]string{"a", "b", "c"}) // true
arr.IsList([]int{})                 // true
```

**Why it is useful:** Allows code ported from Laravel's PHP collections to compile and behave correctly without modification.

---

## First

```go
func First[T any](items []T, callbacks ...func(T, int) bool) (T, bool)
```

Returns the first element of the slice, or the first element that matches the provided callback predicate. The second return value reports whether a match was found.

```go
// Without a callback -- returns the first element.
val, ok := arr.First([]int{10, 20, 30})
// val == 10, ok == true

// With a callback -- returns the first even number.
val, ok = arr.First([]int{1, 3, 4, 6}, func(v int, _ int) bool {
    return v%2 == 0
})
// val == 4, ok == true

// Empty slice.
val, ok = arr.First([]int{})
// val == 0, ok == false
```

**Why it is useful:** Replaces the common `if len(s) > 0 { s[0] }` pattern with a single call that also supports predicate-based searching.

---

## Last

```go
func Last[T any](items []T, callbacks ...func(T, int) bool) (T, bool)
```

Returns the last element of the slice, or the last element that matches the provided callback predicate. The second return value reports whether a match was found.

```go
val, ok := arr.Last([]string{"a", "b", "c"})
// val == "c", ok == true

val, ok = arr.Last([]int{1, 2, 3, 4, 5}, func(v int, _ int) bool {
    return v < 4
})
// val == 3, ok == true
```

**Why it is useful:** Safely retrieves the last element (or last matching element) without manual length checks or reverse iteration.

---

## Take

```go
func Take[T any](items []T, limit int) []T
```

Returns a new slice containing up to `limit` elements. A positive limit takes from the front; a negative limit takes from the end.

```go
arr.Take([]int{1, 2, 3, 4, 5}, 3)   // [1, 2, 3]
arr.Take([]int{1, 2, 3, 4, 5}, -2)  // [4, 5]
arr.Take([]int{1, 2}, 10)           // [1, 2]
```

**Why it is useful:** Provides a single function for both "take first N" and "take last N" operations, handling edge cases automatically.

---

## Only

```go
func Only[T any](items []T, indices []int) []T
```

Returns a new slice containing only the elements at the given indices. Out-of-range indices are silently skipped.

```go
arr.Only([]string{"a", "b", "c", "d"}, []int{0, 2})
// ["a", "c"]

arr.Only([]string{"a", "b"}, []int{0, 5})
// ["a"]  -- index 5 is out of range, skipped
```

**Why it is useful:** Selects specific positions from a slice without manual bounds checking.

---

## Except

```go
func Except[T any](items []T, indices []int) []T
```

Returns a new slice containing all elements except those at the given indices.

```go
arr.Except([]string{"a", "b", "c", "d"}, []int{1, 3})
// ["a", "c"]
```

**Why it is useful:** The inverse of `Only` -- removes specific positions from a slice cleanly.

---

## Flatten

```go
func Flatten[T any](items [][]T) []T
```

Flattens a slice of slices into a single, flat slice.

```go
nested := [][]int{{1, 2}, {3, 4}, {5}}
arr.Flatten(nested)
// [1, 2, 3, 4, 5]
```

**Why it is useful:** Eliminates nested loops when merging grouped results into a single list.

---

## Collapse

```go
func Collapse[T any](items [][]T) []T
```

Merges a slice of slices into a single slice. This is an alias for `Flatten`.

```go
arr.Collapse([][]string{{"a", "b"}, {"c"}})
// ["a", "b", "c"]
```

**Why it is useful:** Provides an alternative name for `Flatten` that may read more naturally when the intent is to "collapse" grouped data.

---

## Wrap

```go
func Wrap[T any](value T) []T
```

Wraps a single value in a one-element slice.

```go
arr.Wrap(42)       // [42]
arr.Wrap("hello")  // ["hello"]
```

**Why it is useful:** Normalizes a single value into a slice so downstream code can always work with `[]T`.

---

## WrapSlice

```go
func WrapSlice[T any](value []T) []T
```

Returns the given slice unchanged. Use this when the value is already a slice and you want to avoid double-wrapping.

```go
s := []int{1, 2, 3}
arr.WrapSlice(s)  // [1, 2, 3] -- same slice returned
```

**Why it is useful:** Pairs with `Wrap` to provide a uniform wrapping API that works whether the input is a single value or already a slice.

---

## Prepend

```go
func Prepend[T any](items []T, value T) []T
```

Inserts a value at the beginning of the slice and returns the new slice.

```go
arr.Prepend([]int{2, 3, 4}, 1)
// [1, 2, 3, 4]
```

**Why it is useful:** Go's built-in `append` adds to the end; `Prepend` provides the inverse operation in a single, readable call.

---

## Push

```go
func Push[T any](items []T, values ...T) []T
```

Appends one or more values to the end of the slice and returns the new slice.

```go
arr.Push([]int{1, 2}, 3, 4, 5)
// [1, 2, 3, 4, 5]
```

**Why it is useful:** A thin wrapper around `append` that reads more naturally in collection-style code.

---

## Shuffle

```go
func Shuffle[T any](items []T) []T
```

Returns a new slice with the elements in random order. The original slice is not modified.

```go
arr.Shuffle([]int{1, 2, 3, 4, 5})
// e.g. [3, 1, 5, 2, 4] -- random each time
```

**Why it is useful:** Provides an immutable shuffle without needing to set up a `rand.Shuffle` call manually.

---

## Random

```go
func Random[T any](items []T, counts ...int) []T
```

Returns a new slice containing `count` randomly selected elements. If `count` is omitted it defaults to 1. If `count` exceeds the slice length, all elements are returned in random order.

```go
arr.Random([]string{"a", "b", "c", "d"}, 2)
// e.g. ["c", "a"]

arr.Random([]int{10, 20, 30})
// e.g. [20]  -- defaults to 1
```

**Why it is useful:** Quickly samples elements from a slice without writing shuffle-and-truncate boilerplate.

---

## Sort

```go
func Sort[T any](items []T, less func(a, b T) bool) []T
```

Returns a new slice sorted using the provided comparison function. The sort is stable (equal elements retain their original order).

```go
arr.Sort([]int{3, 1, 4, 1, 5}, func(a, b int) bool {
    return a < b
})
// [1, 1, 3, 4, 5]

type User struct {
    Name string
    Age  int
}
users := []User{{Name: "Bob", Age: 30}, {Name: "Alice", Age: 25}}
arr.Sort(users, func(a, b User) bool {
    return a.Age < b.Age
})
// [{Alice 25}, {Bob 30}]
```

**Why it is useful:** Returns a new sorted slice instead of sorting in place, which keeps the original data intact.

---

## SortDesc

```go
func SortDesc[T any](items []T, less func(a, b T) bool) []T
```

Returns a new slice sorted in descending order. The `less` function defines the ascending order; `SortDesc` reverses it internally.

```go
arr.SortDesc([]int{3, 1, 4, 1, 5}, func(a, b int) bool {
    return a < b
})
// [5, 4, 3, 1, 1]
```

**Why it is useful:** Saves you from writing a reversed comparator when you need descending order.

---

## SortRecursive

```go
func SortRecursive[T any](items []T, less func(a, b T) bool) []T
```

Sorts a slice using the provided comparison function. For flat slices (which is the only kind in Go generics), this behaves identically to `Sort`.

```go
arr.SortRecursive([]int{5, 2, 8, 1}, func(a, b int) bool {
    return a < b
})
// [1, 2, 5, 8]
```

**Why it is useful:** Exists for API parity with Laravel collections. Use `Sort` directly in new Go code.

---

## Where

```go
func Where[T any](items []T, callback func(T, int) bool) []T
```

Returns a new slice containing only the elements for which the callback returns `true`. The callback receives each element and its index.

```go
evens := arr.Where([]int{1, 2, 3, 4, 5, 6}, func(v int, _ int) bool {
    return v%2 == 0
})
// [2, 4, 6]
```

**Why it is useful:** A standard filter operation with access to both value and index -- the foundation for most data-selection patterns.

---

## WhereNotNull

```go
func WhereNotNull[T comparable](items []T) []T
```

Returns a new slice with all zero-value elements removed. The type must satisfy `comparable`.

```go
arr.WhereNotNull([]string{"a", "", "b", "", "c"})
// ["a", "b", "c"]

arr.WhereNotNull([]int{0, 1, 0, 2, 3})
// [1, 2, 3]
```

**Why it is useful:** Quickly strips empty strings, zero integers, nil pointers, or any other zero-value sentinel from a typed slice.

---

## Reject

```go
func Reject[T any](items []T, callback func(T, int) bool) []T
```

Returns a new slice containing only the elements for which the callback returns `false`. It is the inverse of `Where`.

```go
arr.Reject([]int{1, 2, 3, 4, 5}, func(v int, _ int) bool {
    return v > 3
})
// [1, 2, 3]
```

**Why it is useful:** Makes intent clearer when you want to express "remove items matching X" rather than "keep items not matching X."

---

## Partition

```go
func Partition[T any](items []T, callback func(T, int) bool) ([]T, []T)
```

Splits the slice into two: the first contains elements where the callback returns `true`, the second contains the rest.

```go
pass, fail := arr.Partition([]int{1, 2, 3, 4, 5, 6}, func(v int, _ int) bool {
    return v%2 == 0
})
// pass == [2, 4, 6]
// fail == [1, 3, 5]
```

**Why it is useful:** Performs a single pass to separate elements into two groups, avoiding two separate filter calls.

---

## Every

```go
func Every[T any](items []T, callback func(T, int) bool) bool
```

Reports whether every element in the slice satisfies the callback. Returns `true` for an empty slice.

```go
arr.Every([]int{2, 4, 6}, func(v int, _ int) bool {
    return v%2 == 0
})
// true

arr.Every([]int{2, 3, 6}, func(v int, _ int) bool {
    return v%2 == 0
})
// false
```

**Why it is useful:** Validates that all elements meet a condition without writing a manual loop with an early break.

---

## Some

```go
func Some[T any](items []T, callback func(T, int) bool) bool
```

Reports whether at least one element in the slice satisfies the callback. Returns `false` for an empty slice.

```go
arr.Some([]int{1, 3, 5, 8}, func(v int, _ int) bool {
    return v%2 == 0
})
// true

arr.Some([]int{1, 3, 5}, func(v int, _ int) bool {
    return v%2 == 0
})
// false
```

**Why it is useful:** Provides a short-circuiting existence check -- iteration stops as soon as a match is found.

---

## Exists

```go
func Exists[T any](items []T, index int) bool
```

Reports whether the given index is valid (in bounds) for the slice.

```go
arr.Exists([]string{"a", "b", "c"}, 2)   // true
arr.Exists([]string{"a", "b", "c"}, 5)   // false
arr.Exists([]string{"a", "b", "c"}, -1)  // false
```

**Why it is useful:** A readable bounds check that replaces `index >= 0 && index < len(s)`.

---

## Has

```go
func Has[T any](items []T, indices ...int) bool
```

Reports whether **all** of the given indices are valid for the slice.

```go
arr.Has([]int{10, 20, 30}, 0, 1, 2)  // true
arr.Has([]int{10, 20, 30}, 0, 5)     // false -- 5 is out of range
```

**Why it is useful:** Validates multiple indices in one call before accessing them.

---

## HasAny

```go
func HasAny[T any](items []T, indices ...int) bool
```

Reports whether **at least one** of the given indices is valid for the slice.

```go
arr.HasAny([]int{10, 20, 30}, 5, 1)    // true -- index 1 is valid
arr.HasAny([]int{10, 20, 30}, 5, 10)   // false -- neither is valid
```

**Why it is useful:** A looser existence check than `Has`, useful when any one valid index is sufficient.

---

## Join

```go
func Join(items []string, glue string, finalGlues ...string) string
```

Concatenates string slice elements with a glue string. An optional final glue is placed between the last two elements.

```go
arr.Join([]string{"a", "b", "c"}, ", ")
// "a, b, c"

arr.Join([]string{"a", "b", "c"}, ", ", " and ")
// "a, b and c"

arr.Join([]string{"Go"}, ", ", " and ")
// "Go"
```

**Why it is useful:** Produces human-readable lists like "Alice, Bob and Charlie" without manual string building.

---

## CrossJoin

```go
func CrossJoin[T any](lists ...[]T) [][]T
```

Returns the Cartesian product of the given slices. Each result element is a slice combining one element from each input.

```go
arr.CrossJoin([]int{1, 2}, []int{10, 20})
// [[1, 10], [1, 20], [2, 10], [2, 20]]

arr.CrossJoin([]string{"S", "M"}, []string{"Red", "Blue"})
// [["S", "Red"], ["S", "Blue"], ["M", "Red"], ["M", "Blue"]]
```

**Why it is useful:** Generates all combinations across multiple dimensions (sizes, colors, options) without nested loops.

---

## Divide

```go
func Divide[T any](items []T) ([]int, []T)
```

Splits a slice into two: a slice of indices and a slice of values.

```go
indices, values := arr.Divide([]string{"a", "b", "c"})
// indices == [0, 1, 2]
// values  == ["a", "b", "c"]
```

**Why it is useful:** Separates positional information from data, useful when you need to work with indices and values independently.

---

## Map

```go
func Map[T any, R any](items []T, callback func(T, int) R) []R
```

Applies the callback to each element and returns a new slice of the transformed results. The callback receives each element and its index.

```go
arr.Map([]int{1, 2, 3}, func(v int, _ int) int {
    return v * 2
})
// [2, 4, 6]

arr.Map([]string{"hello", "world"}, func(s string, _ int) int {
    return len(s)
})
// [5, 5]
```

**Why it is useful:** The fundamental transform operation -- converts a slice of one type into a slice of another without manual loop construction.

---

## MapWithKeys

```go
func MapWithKeys[T any, K comparable, V any](items []T, callback func(T) (K, V)) map[K]V
```

Applies the callback to each element, which produces a key-value pair, and collects the results into a map.

```go
type User struct {
    ID   int
    Name string
}

users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
result := arr.MapWithKeys(users, func(u User) (int, string) {
    return u.ID, u.Name
})
// map[1:"Alice" 2:"Bob"]
```

**Why it is useful:** Converts a slice into a lookup map in a single expression, common when indexing records by a unique field.

---

## MapSpread

```go
func MapSpread[T any, R any](items []T, callback func(T, int) R) []R
```

Applies the callback to each element and returns a new slice. In Go this is identical to `Map`.

```go
arr.MapSpread([]int{1, 2, 3}, func(v int, idx int) string {
    return fmt.Sprintf("%d:%d", idx, v)
})
// ["0:1", "1:2", "2:3"]
```

**Why it is useful:** Exists for API parity with Laravel collections. Use `Map` in new Go code.

---

## KeyBy

```go
func KeyBy[T any, K comparable](items []T, keyFunc func(T) K) map[K]T
```

Indexes the slice elements by the key returned from `keyFunc`, producing a map. If duplicate keys exist, the last element wins.

```go
type Product struct {
    SKU  string
    Name string
}

products := []Product{
    {SKU: "A1", Name: "Widget"},
    {SKU: "B2", Name: "Gadget"},
}
bysku := arr.KeyBy(products, func(p Product) string {
    return p.SKU
})
// map["A1":{SKU:"A1", Name:"Widget"} "B2":{SKU:"B2", Name:"Gadget"}]
```

**Why it is useful:** Creates a lookup table from a slice in one call, avoiding manual map construction loops.

---

## Pluck

```go
func Pluck[T any, V any](items []T, valueFunc func(T) V) []V
```

Extracts a single field or derived value from each element and returns the collected values as a slice.

```go
type Employee struct {
    Name   string
    Salary float64
}

employees := []Employee{
    {Name: "Alice", Salary: 90000},
    {Name: "Bob", Salary: 75000},
}

names := arr.Pluck(employees, func(e Employee) string {
    return e.Name
})
// ["Alice", "Bob"]

salaries := arr.Pluck(employees, func(e Employee) float64 {
    return e.Salary
})
// [90000, 75000]
```

**Why it is useful:** Extracts a column of data from a slice of structs, mirroring SQL's `SELECT column` pattern.
