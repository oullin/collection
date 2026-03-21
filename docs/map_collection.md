# MapCollection[K, V]

`MapCollection` is an ordered map that preserves insertion order. Unlike Go's built-in `map`, iterating over a `MapCollection` always visits entries in the order they were added. It provides a fluent, chainable API for querying, filtering, transforming, and serializing key-value data.

```go
import "github.com/gocanto/collection"
```

---

## Table of Contents

- [Constructors](#constructors)
- [Query](#query)
- [Retrieval](#retrieval)
- [Search](#search)
- [Mutation](#mutation)
- [Filtering](#filtering)
- [Transformation](#transformation)
- [Sorting](#sorting)
- [Iteration](#iteration)
- [String](#string)
- [Conditional](#conditional)
- [Partitioning](#partitioning)
- [Serialization](#serialization)

---

## Constructors

### NewMap

```go
func NewMap[K comparable, V any](items map[K]V) *MapCollection[K, V]
```

Creates a `MapCollection` from an existing Go map. If `nil` is passed, an empty collection is created.

```go
m := collection.NewMap(map[string]int{
    "apples":  5,
    "oranges": 3,
})
fmt.Println(m.Count()) // 2
```

**When to use:** You already have a `map` and want the ordered, fluent API that `MapCollection` provides.

---

### NewMapFromPairs

```go
func NewMapFromPairs[K comparable, V any](pairs ...Pair[K, V]) *MapCollection[K, V]
```

Creates a `MapCollection` from explicit key-value pairs. Insertion order is determined by the order of the arguments. If a key appears more than once, the last value wins but the key retains its first position.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, int]{Key: "a", Value: 1},
    collection.Pair[string, int]{Key: "b", Value: 2},
    collection.Pair[string, int]{Key: "c", Value: 3},
)
fmt.Println(m.Keys().All()) // [a b c]
```

**When to use:** You need guaranteed insertion order from the start, or you are constructing a map programmatically from pairs.

---

## Query

### All

```go
func (m *MapCollection[K, V]) All() map[K]V
```

Returns a shallow copy of the underlying map. Changes to the returned map do not affect the collection.

```go
m := collection.NewMap(map[string]int{"x": 10, "y": 20})
raw := m.All() // map[string]int{"x": 10, "y": 20}
```

**When to use:** You need to pass the data to code that expects a plain `map`.

---

### Keys

```go
func (m *MapCollection[K, V]) Keys() *Collection[K]
```

Returns all keys as a `Collection`, preserving insertion order.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, int]{Key: "first", Value: 1},
    collection.Pair[string, int]{Key: "second", Value: 2},
)
keys := m.Keys().All() // ["first", "second"]
```

**When to use:** You need to iterate over or manipulate just the keys -- for example, checking for duplicates or sorting them independently.

---

### Values

```go
func (m *MapCollection[K, V]) Values() *Collection[V]
```

Returns all values as a `Collection`, preserving insertion order.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, int]{Key: "a", Value: 100},
    collection.Pair[string, int]{Key: "b", Value: 200},
)
vals := m.Values().All() // [100, 200]
```

**When to use:** You want to aggregate, filter, or transform only the values without caring about keys.

---

### Count

```go
func (m *MapCollection[K, V]) Count() int
```

Returns the number of items in the collection.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3})
fmt.Println(m.Count()) // 3
```

---

### IsEmpty

```go
func (m *MapCollection[K, V]) IsEmpty() bool
```

Reports whether the collection contains no items.

```go
m := collection.NewMap(map[string]int{})
fmt.Println(m.IsEmpty()) // true
```

---

### IsNotEmpty

```go
func (m *MapCollection[K, V]) IsNotEmpty() bool
```

Reports whether the collection contains at least one item.

```go
m := collection.NewMap(map[string]int{"a": 1})
fmt.Println(m.IsNotEmpty()) // true
```

---

### Has

```go
func (m *MapCollection[K, V]) Has(key K) bool
```

Reports whether the given key exists in the collection.

```go
m := collection.NewMap(map[string]int{"name": 42})
fmt.Println(m.Has("name"))  // true
fmt.Println(m.Has("email")) // false
```

**When to use:** Quick key-existence check before performing an operation.

---

### HasAny

```go
func (m *MapCollection[K, V]) HasAny(keys ...K) bool
```

Reports whether any of the given keys exist in the collection.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2})
fmt.Println(m.HasAny("b", "z")) // true
fmt.Println(m.HasAny("x", "z")) // false
```

**When to use:** Validating that at least one of several expected keys is present, such as checking for alternative field names in a config.

---

### ContainsOneItem

```go
func (m *MapCollection[K, V]) ContainsOneItem() bool
```

Reports whether the collection contains exactly one item.

```go
m := collection.NewMap(map[string]int{"only": 1})
fmt.Println(m.ContainsOneItem()) // true
```

---

### ContainsManyItems

```go
func (m *MapCollection[K, V]) ContainsManyItems() bool
```

Reports whether the collection contains more than one item.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2})
fmt.Println(m.ContainsManyItems()) // true
```

---

## Retrieval

### Get

```go
func (m *MapCollection[K, V]) Get(key K, defaults ...V) (V, bool)
```

Returns the value for the given key. The second return value indicates whether the key was found. An optional default value is returned (with `false`) when the key is missing.

```go
m := collection.NewMap(map[string]int{"timeout": 30})

val, ok := m.Get("timeout")
// val=30, ok=true

val, ok = m.Get("retries", 3)
// val=3, ok=false (key missing, default returned)
```

**When to use:** Safe value retrieval with an optional fallback, similar to Python's `dict.get(key, default)`.

---

### GetOrPut

```go
func (m *MapCollection[K, V]) GetOrPut(key K, value V) V
```

Returns the value for the given key if it exists. Otherwise stores and returns the provided default value.

```go
m := collection.NewMap(map[string]int{"hits": 10})

val := m.GetOrPut("hits", 0)    // 10 (existing value)
val = m.GetOrPut("misses", 0)   // 0  (inserted and returned)
fmt.Println(m.Has("misses"))    // true
```

**When to use:** Implementing a "get or initialize" pattern -- for example, building a counter map where missing keys default to zero.

---

### First

```go
func (m *MapCollection[K, V]) First(predicates ...func(V, K) bool) (V, bool)
```

Returns the first value in insertion order. If a predicate is provided, returns the first value satisfying it. The second return value indicates whether a match was found.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, int]{Key: "a", Value: 10},
    collection.Pair[string, int]{Key: "b", Value: 20},
    collection.Pair[string, int]{Key: "c", Value: 30},
)

val, _ := m.First() // 10

val, _ = m.First(func(v int, k string) bool {
    return v > 15
})
// val=20
```

---

### Last

```go
func (m *MapCollection[K, V]) Last(predicates ...func(V, K) bool) (V, bool)
```

Returns the last value in insertion order, or the last value satisfying the optional predicate.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, int]{Key: "a", Value: 10},
    collection.Pair[string, int]{Key: "b", Value: 20},
    collection.Pair[string, int]{Key: "c", Value: 30},
)

val, _ := m.Last() // 30

val, _ = m.Last(func(v int, k string) bool {
    return v < 25
})
// val=20
```

---

## Search

### Contains

```go
func (m *MapCollection[K, V]) Contains(predicate func(V, K) bool) bool
```

Reports whether any item in the collection satisfies the predicate.

```go
m := collection.NewMap(map[string]int{"age": 25, "score": 90})

hasHighScore := m.Contains(func(v int, k string) bool {
    return k == "score" && v >= 80
})
fmt.Println(hasHighScore) // true
```

**When to use:** Checking for the existence of a value that meets a complex condition involving both the key and the value.

---

### Some

```go
func (m *MapCollection[K, V]) Some(predicate func(V, K) bool) bool
```

Alias for `Contains`. Use whichever reads more naturally in your code.

---

### DoesntContain

```go
func (m *MapCollection[K, V]) DoesntContain(predicate func(V, K) bool) bool
```

Reports whether no item in the collection satisfies the predicate. The logical inverse of `Contains`.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2})

noNegatives := m.DoesntContain(func(v int, _ string) bool {
    return v < 0
})
fmt.Println(noNegatives) // true
```

---

### HasSole

```go
func (m *MapCollection[K, V]) HasSole(predicate func(V, K) bool) bool
```

Reports whether exactly one item in the collection satisfies the predicate.

```go
m := collection.NewMap(map[string]string{
    "role": "admin",
    "name": "Alice",
    "dept": "Engineering",
})

onlyOneAdmin := m.HasSole(func(v string, k string) bool {
    return v == "admin"
})
fmt.Println(onlyOneAdmin) // true
```

**When to use:** Asserting uniqueness constraints -- for example, ensuring exactly one entry matches a given criterion.

---

### Search

```go
func (m *MapCollection[K, V]) Search(predicate func(V, K) bool) (K, bool)
```

Returns the key of the first item that satisfies the predicate. The second return value indicates whether a match was found.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, int]{Key: "x", Value: 100},
    collection.Pair[string, int]{Key: "y", Value: 200},
    collection.Pair[string, int]{Key: "z", Value: 300},
)

key, found := m.Search(func(v int, _ string) bool {
    return v == 200
})
fmt.Println(key, found) // y true
```

---

## Mutation

### Put

```go
func (m *MapCollection[K, V]) Put(key K, value V) *MapCollection[K, V]
```

Sets a key-value pair. If the key already exists, the value is updated in place. If the key is new, it is appended to the end of the insertion-order list. Returns the collection for chaining.

```go
m := collection.NewMap(map[string]int{})
m.Put("x", 1).Put("y", 2).Put("z", 3)
fmt.Println(m.Count()) // 3
```

---

### Pull

```go
func (m *MapCollection[K, V]) Pull(key K) (V, bool)
```

Removes an item by key and returns its value. The second return value indicates whether the key was found.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2})

val, ok := m.Pull("a") // val=1, ok=true
fmt.Println(m.Count()) // 1
```

**When to use:** Extracting a value while simultaneously removing it, like popping from a dictionary.

---

### Forget

```go
func (m *MapCollection[K, V]) Forget(keys ...K) *MapCollection[K, V]
```

Removes one or more items by key. Returns the collection for chaining.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3})
m.Forget("a", "c")
fmt.Println(m.Count()) // 1
```

**When to use:** Bulk removal of keys you no longer need, such as stripping internal fields before returning data to a caller.

---

## Filtering

### Only

```go
func (m *MapCollection[K, V]) Only(keys ...K) *MapCollection[K, V]
```

Returns a new collection containing only the specified keys.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4})
subset := m.Only("a", "c")
fmt.Println(subset.All()) // map[a:1 c:3]
```

**When to use:** Picking a subset of fields -- for example, selecting only the columns you need from a record.

---

### Except

```go
func (m *MapCollection[K, V]) Except(keys ...K) *MapCollection[K, V]
```

Returns a new collection containing all items except those with the specified keys.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3})
rest := m.Except("b")
fmt.Println(rest.All()) // map[a:1 c:3]
```

**When to use:** The inverse of `Only` -- hiding sensitive or irrelevant fields.

---

### Filter

```go
func (m *MapCollection[K, V]) Filter(callback func(V, K) bool) *MapCollection[K, V]
```

Returns a new collection containing only the items for which the callback returns `true`.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 5, "c": 10, "d": 15})

big := m.Filter(func(v int, _ string) bool {
    return v >= 10
})
fmt.Println(big.All()) // map[c:10 d:15]
```

**When to use:** General-purpose value/key filtering with a custom predicate.

---

### Reject

```go
func (m *MapCollection[K, V]) Reject(callback func(V, K) bool) *MapCollection[K, V]
```

Returns a new collection excluding items for which the callback returns `true`. The logical inverse of `Filter`.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 5, "c": 10})

small := m.Reject(func(v int, _ string) bool {
    return v >= 10
})
fmt.Println(small.All()) // map[a:1 b:5]
```

---

## Transformation

### MapValues (top-level function)

```go
func MapValues[K comparable, V any, R any](
    m *MapCollection[K, V],
    callback func(V, K) R,
) *MapCollection[K, R]
```

Transforms each value using the callback, returning a new `MapCollection` with the same keys and transformed values. Because Go does not allow methods to introduce new type parameters, this is a package-level function.

```go
prices := collection.NewMap(map[string]float64{
    "apple":  1.20,
    "banana": 0.50,
})

doubled := collection.MapValues(prices, func(v float64, _ string) float64 {
    return v * 2
})
fmt.Println(doubled.All()) // map[apple:2.4 banana:1]
```

**When to use:** Applying a uniform transformation to all values -- currency conversion, unit changes, formatting.

---

### MapFlip (top-level function)

```go
func MapFlip[K comparable, V comparable](
    m *MapCollection[K, V],
) *MapCollection[V, K]
```

Swaps keys and values. The original values become keys and original keys become values. Both `K` and `V` must be `comparable`. If duplicate values exist, the last key wins.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3})
flipped := collection.MapFlip(m)
// flipped: map[1:"a", 2:"b", 3:"c"]
```

**When to use:** Building a reverse lookup table.

---

### Merge

```go
func (m *MapCollection[K, V]) Merge(items map[K]V) *MapCollection[K, V]
```

Returns a new collection with the given map merged in. Existing keys are overwritten by the incoming values. New keys are appended.

```go
base := collection.NewMap(map[string]int{"a": 1, "b": 2})
merged := base.Merge(map[string]int{"b": 20, "c": 30})
fmt.Println(merged.All()) // map[a:1 b:20 c:30]
```

---

### Union

```go
func (m *MapCollection[K, V]) Union(items map[K]V) *MapCollection[K, V]
```

Returns a new collection that is the union of the collection and the given map. Keys already present in the original collection are **not** overwritten.

```go
base := collection.NewMap(map[string]int{"a": 1, "b": 2})
united := base.Union(map[string]int{"b": 20, "c": 30})
fmt.Println(united.All()) // map[a:1 b:2 c:30]  ("b" kept original value)
```

**When to use:** Providing defaults -- the incoming map supplies fallback values for missing keys without overriding existing ones.

---

### Replace

```go
func (m *MapCollection[K, V]) Replace(items map[K]V) *MapCollection[K, V]
```

Alias for `Merge`. Existing keys are overwritten by the incoming values.

---

### DiffKeys

```go
func (m *MapCollection[K, V]) DiffKeys(items map[K]V) *MapCollection[K, V]
```

Returns a new collection containing items whose keys are **not** present in the given map.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3})
diff := m.DiffKeys(map[string]int{"a": 0, "c": 0})
fmt.Println(diff.All()) // map[b:2]
```

**When to use:** Finding entries that are unique to one dataset compared to another.

---

### IntersectByKeys

```go
func (m *MapCollection[K, V]) IntersectByKeys(items map[K]V) *MapCollection[K, V]
```

Returns a new collection containing items whose keys **are** present in the given map.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3})
inter := m.IntersectByKeys(map[string]int{"a": 0, "c": 0})
fmt.Println(inter.All()) // map[a:1 c:3]
```

**When to use:** Restricting a collection to only the keys that exist in another dataset.

---

## Sorting

### MapSortKeys (top-level function)

```go
func MapSortKeys[V any](m *MapCollection[string, V]) *MapCollection[string, V]
```

Returns a new collection with string keys sorted in ascending order.

```go
m := collection.NewMap(map[string]int{"banana": 2, "apple": 1, "cherry": 3})
sorted := collection.MapSortKeys(m)
fmt.Println(sorted.Keys().All()) // [apple banana cherry]
```

> **Note:** This function works only with `string` keys. For other key types, use `SortKeysUsing`.

---

### MapSortKeysDesc (top-level function)

```go
func MapSortKeysDesc[V any](m *MapCollection[string, V]) *MapCollection[string, V]
```

Returns a new collection with string keys sorted in descending order.

```go
m := collection.NewMap(map[string]int{"banana": 2, "apple": 1, "cherry": 3})
sorted := collection.MapSortKeysDesc(m)
fmt.Println(sorted.Keys().All()) // [cherry banana apple]
```

---

### SortKeysUsing

```go
func (m *MapCollection[K, V]) SortKeysUsing(less func(K, K) bool) *MapCollection[K, V]
```

Returns a new collection with keys sorted using the provided comparison function. Works with any comparable key type.

```go
m := collection.NewMap(map[int]string{3: "c", 1: "a", 2: "b"})
sorted := m.SortKeysUsing(func(a, b int) bool {
    return a < b
})
fmt.Println(sorted.Keys().All()) // [1 2 3]
```

**When to use:** Sorting by non-string keys or custom orderings (e.g., case-insensitive sort).

---

## Iteration

### Each

```go
func (m *MapCollection[K, V]) Each(callback func(V, K) bool) *MapCollection[K, V]
```

Iterates over items in insertion order, calling the callback for each key-value pair. Return `false` from the callback to stop iteration early. Returns the collection for chaining.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, int]{Key: "a", Value: 1},
    collection.Pair[string, int]{Key: "b", Value: 2},
    collection.Pair[string, int]{Key: "c", Value: 3},
)

m.Each(func(v int, k string) bool {
    fmt.Printf("%s=%d\n", k, v)
    return true // continue
})
// Output:
// a=1
// b=2
// c=3
```

---

### Tap

```go
func (m *MapCollection[K, V]) Tap(callback func(*MapCollection[K, V])) *MapCollection[K, V]
```

Passes the collection to the callback for side effects (logging, debugging) and returns the collection unchanged, preserving a method chain.

```go
result := collection.NewMap(map[string]int{"a": 1, "b": 2}).
    Tap(func(m *collection.MapCollection[string, int]) {
        fmt.Println("count:", m.Count())
    }).
    Filter(func(v int, _ string) bool { return v > 1 })
```

---

### Iter

```go
func (m *MapCollection[K, V]) Iter() iter.Seq2[K, V]
```

Returns an `iter.Seq2[K, V]` iterator for use with Go 1.23+ range-over-func loops. Entries are yielded in insertion order.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, int]{Key: "x", Value: 10},
    collection.Pair[string, int]{Key: "y", Value: 20},
)

for key, value := range m.Iter() {
    fmt.Printf("%s: %d\n", key, value)
}
// Output:
// x: 10
// y: 20
```

**When to use:** Integrating with standard Go `for range` syntax while preserving insertion order. This is the idiomatic way to iterate a `MapCollection` in Go 1.23+.

---

## String

### Implode

```go
func (m *MapCollection[K, V]) Implode(glue string) string
```

Concatenates all values (in insertion order) into a single string, separated by the given glue.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, string]{Key: "first", Value: "Hello"},
    collection.Pair[string, string]{Key: "second", Value: "World"},
)
fmt.Println(m.Implode(" ")) // "Hello World"
```

---

### Join

```go
func (m *MapCollection[K, V]) Join(glue string, finalGlues ...string) string
```

Like `Implode`, but with an optional final glue placed between the last two items. This is useful for producing human-readable lists.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, string]{Key: "a", Value: "Go"},
    collection.Pair[string, string]{Key: "b", Value: "Rust"},
    collection.Pair[string, string]{Key: "c", Value: "Zig"},
)
fmt.Println(m.Join(", ", ", and ")) // "Go, Rust, and Zig"
```

---

## Conditional

### When

```go
func (m *MapCollection[K, V]) When(
    condition bool,
    callback func(*MapCollection[K, V]) *MapCollection[K, V],
    defaults ...func(*MapCollection[K, V]) *MapCollection[K, V],
) *MapCollection[K, V]
```

Applies the callback if the condition is `true`. An optional default callback is applied when the condition is `false`. Returns the (potentially new) collection.

```go
isAdmin := true

m := collection.NewMap(map[string]int{"base": 100}).
    When(isAdmin, func(m *collection.MapCollection[string, int]) *collection.MapCollection[string, int] {
        m.Put("bonus", 50)
        return m
    })
fmt.Println(m.All()) // map[base:100 bonus:50]
```

**When to use:** Conditionally modifying a collection in a method chain without breaking the flow with `if/else` blocks.

---

### Unless

```go
func (m *MapCollection[K, V]) Unless(
    condition bool,
    callback func(*MapCollection[K, V]) *MapCollection[K, V],
    defaults ...func(*MapCollection[K, V]) *MapCollection[K, V],
) *MapCollection[K, V]
```

The inverse of `When`: applies the callback unless the condition is `true`.

```go
isGuest := false

m := collection.NewMap(map[string]string{"page": "dashboard"}).
    Unless(isGuest, func(m *collection.MapCollection[string, string]) *collection.MapCollection[string, string] {
        m.Put("menu", "full")
        return m
    })
fmt.Println(m.Has("menu")) // true
```

---

## Partitioning

### Every

```go
func (m *MapCollection[K, V]) Every(callback func(V, K) bool) bool
```

Reports whether all items in the collection satisfy the callback.

```go
m := collection.NewMap(map[string]int{"a": 2, "b": 4, "c": 6})

allEven := m.Every(func(v int, _ string) bool {
    return v%2 == 0
})
fmt.Println(allEven) // true
```

**When to use:** Validating that every entry meets a constraint before proceeding.

---

### Partition

```go
func (m *MapCollection[K, V]) Partition(
    callback func(V, K) bool,
) (*MapCollection[K, V], *MapCollection[K, V])
```

Splits the collection into two: one where the callback returns `true`, and one where it returns `false`.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4})

pass, fail := m.Partition(func(v int, _ string) bool {
    return v%2 == 0
})
fmt.Println(pass.All()) // map[b:2 d:4]
fmt.Println(fail.All()) // map[a:1 c:3]
```

**When to use:** Splitting data into two groups based on a condition -- for example, separating valid from invalid entries.

---

## Serialization

### ToJSON

```go
func (m *MapCollection[K, V]) ToJSON() ([]byte, error)
```

Serializes the collection to compact JSON.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2})
data, _ := m.ToJSON()
fmt.Println(string(data)) // {"a":1,"b":2}
```

---

### ToPrettyJSON

```go
func (m *MapCollection[K, V]) ToPrettyJSON() ([]byte, error)
```

Serializes the collection to indented (pretty-printed) JSON with 4-space indentation.

```go
m := collection.NewMap(map[string]int{"a": 1, "b": 2})
data, _ := m.ToPrettyJSON()
fmt.Println(string(data))
// {
//     "a": 1,
//     "b": 2
// }
```

---

### String

```go
func (m *MapCollection[K, V]) String() string
```

Returns the JSON representation of the collection as a string. If serialization fails, returns `"{}"`.

```go
m := collection.NewMap(map[string]int{"x": 42})
fmt.Println(m.String()) // {"x":42}
```

---

### MarshalJSON

```go
func (m *MapCollection[K, V]) MarshalJSON() ([]byte, error)
```

Implements the `json.Marshaler` interface, allowing a `MapCollection` to be used anywhere Go's `encoding/json` expects a marshaler.

```go
m := collection.NewMap(map[string]int{"a": 1})
data, _ := json.Marshal(m) // works because MarshalJSON is implemented
```

---

### UnmarshalJSON

```go
func (m *MapCollection[K, V]) UnmarshalJSON(data []byte) error
```

Implements the `json.Unmarshaler` interface.

```go
m := collection.NewMap(map[string]int{})
_ = json.Unmarshal([]byte(`{"x":10,"y":20}`), m)
fmt.Println(m.All()) // map[x:10 y:20]
```

---

### Copy

```go
func (m *MapCollection[K, V]) Copy() *MapCollection[K, V]
```

Creates a shallow copy of the collection. Mutations to the copy do not affect the original.

```go
original := collection.NewMap(map[string]int{"a": 1})
clone := original.Copy()
clone.Put("b", 2)

fmt.Println(original.Count()) // 1
fmt.Println(clone.Count())    // 2
```

---

### Dump

```go
func (m *MapCollection[K, V]) Dump() *MapCollection[K, V]
```

Prints the underlying map to stdout for debugging purposes. Returns the collection for chaining.

```go
collection.NewMap(map[string]int{"debug": 42}).Dump()
// Output: map[debug:42]
```

---

### ToPairs

```go
func (m *MapCollection[K, V]) ToPairs() *Collection[Pair[K, V]]
```

Converts the collection to a `Collection` of `Pair[K, V]` values, preserving insertion order. This is useful for round-tripping between `MapCollection` and pair-based representations.

```go
m := collection.NewMapFromPairs(
    collection.Pair[string, int]{Key: "a", Value: 1},
    collection.Pair[string, int]{Key: "b", Value: 2},
)

pairs := m.ToPairs()
fmt.Println(pairs.Count()) // 2
// Each pair has .Key and .Value fields
```

**When to use:** Converting to a slice-based representation for operations that require indexing, or serializing ordered key-value data.
