# kv -- Map Utilities with Dot-Notation Support

```
import "github.com/gocanto/collection/kv"
```

The `kv` package provides standalone utility functions for working with Go maps. Its headline feature is **dot-notation access** for nested `map[string]any` structures (the kind you get from JSON decoding), but it also includes general-purpose filtering, sorting, transformation, and encoding helpers.

Use `kv` when you need quick map operations without constructing a `collection.MapCollection`.

---

## Table of Contents

**Dot-notation helpers (nested `map[string]any`)**
- [Get](#get)
- [Set](#set)
- [Has](#has)
- [Fill](#fill)
- [Add](#add)
- [Forget](#forget)
- [ForgetMany](#forgetmany)
- [Pull](#pull)
- [HasAll](#hasall)
- [HasAny](#hasany)
- [Dot](#dot)
- [Undot](#undot)

**Map filtering and selection**
- [Only](#only)
- [Except](#except)
- [IsAssoc](#isassoc)

**Encoding**
- [Query](#query)
- [ToCssClasses](#tocssclasses)
- [ToCssStyles](#tocssstyles)

**Sorting**
- [Sort](#sort)
- [SortRecursive](#sortrecursive)

**Map transformation**
- [Map](#map)
- [Where](#where)
- [PrependKeysWith](#prependkeyswith)
- [Replace](#replace)

---

## Dot-Notation Helpers

These functions navigate nested `map[string]any` structures using dot-separated key paths such as `"user.address.city"`. They are particularly useful when working with decoded JSON, YAML, or configuration data.

---

### Get

```go
func Get(target map[string]any, key string, defaults ...any) any
```

Retrieves a value from a nested map using a dot-separated key path. If the key is not found, the first default value is returned (or `nil` if no default is provided). Passing an empty key returns the entire map.

```go
data := map[string]any{
    "user": map[string]any{
        "name": "Alice",
        "address": map[string]any{
            "city": "Portland",
        },
    },
}

kv.Get(data, "user.name")                  // "Alice"
kv.Get(data, "user.address.city")           // "Portland"
kv.Get(data, "user.phone")                  // nil
kv.Get(data, "user.phone", "555-0000")      // "555-0000"
kv.Get(data, "")                            // returns the entire map
```

**Why it is useful:** Eliminates chains of type assertions and existence checks when reading deeply nested configuration or JSON data.

---

### Set

```go
func Set(target map[string]any, key string, value any, overwrite ...bool) map[string]any
```

Sets a value in a nested map using a dot-separated key path, creating intermediate maps as needed. By default, existing values are overwritten. Pass `false` as the last argument to preserve existing values.

```go
data := map[string]any{}

kv.Set(data, "app.name", "MyApp")
// data == {"app": {"name": "MyApp"}}

kv.Set(data, "app.version", "1.0")
// data == {"app": {"name": "MyApp", "version": "1.0"}}

// With overwrite disabled:
kv.Set(data, "app.name", "OtherApp", false)
// data["app"]["name"] is still "MyApp"
```

**Why it is useful:** Builds nested map structures programmatically without manually constructing each intermediate map.

---

### Has

```go
func Has(target map[string]any, key string) bool
```

Reports whether the given dot-notated key exists in the nested map. Returns `false` for an empty key.

```go
data := map[string]any{
    "database": map[string]any{
        "host": "localhost",
    },
}

kv.Has(data, "database.host")  // true
kv.Has(data, "database.port")  // false
kv.Has(data, "cache.driver")   // false
```

**Why it is useful:** Checks for the existence of deeply nested keys in a single call.

---

### Fill

```go
func Fill(target map[string]any, key string, value any) map[string]any
```

Sets the value at the given dot-notated key only if that key does not already exist. Existing values are preserved.

```go
data := map[string]any{
    "app": map[string]any{"name": "MyApp"},
}

kv.Fill(data, "app.name", "Other")     // no change -- key exists
kv.Fill(data, "app.debug", true)       // sets app.debug to true
```

**Why it is useful:** Provides safe default-value injection without accidentally overwriting user-provided configuration.

---

### Add

```go
func Add(target map[string]any, key string, value any) map[string]any
```

An alias for `Fill`. Sets a value only if the key does not already exist.

```go
kv.Add(data, "app.env", "production")
```

**Why it is useful:** Offers an alternative name that reads more naturally in certain contexts ("add this key if missing").

---

### Forget

```go
func Forget(target map[string]any, key string) map[string]any
```

Removes the value at the given dot-notated key from the nested map. If the key does not exist, the map is returned unchanged.

```go
data := map[string]any{
    "user": map[string]any{
        "name":  "Alice",
        "email": "alice@example.com",
    },
}

kv.Forget(data, "user.email")
// data == {"user": {"name": "Alice"}}
```

**Why it is useful:** Deletes nested keys without manually traversing the map hierarchy.

---

### ForgetMany

```go
func ForgetMany(items map[string]any, keys ...string) map[string]any
```

Removes one or more dot-notated keys from the map.

```go
data := map[string]any{
    "a": 1, "b": 2, "c": 3, "d": 4,
}

kv.ForgetMany(data, "b", "d")
// data == {"a": 1, "c": 3}
```

**Why it is useful:** Batch-deletes multiple keys in a single call.

---

### Pull

```go
func Pull(items map[string]any, key string, defaults ...any) any
```

Retrieves the value at the given dot-notated key and then removes it from the map. If the key does not exist, the first default value is returned.

```go
data := map[string]any{"token": "abc123", "user": "alice"}

token := kv.Pull(data, "token")
// token == "abc123"
// data  == {"user": "alice"}

missing := kv.Pull(data, "nope", "default_val")
// missing == "default_val"
```

**Why it is useful:** Combines a "get" and a "delete" into one atomic operation, common when consuming one-time values like tokens or nonces.

---

### HasAll

```go
func HasAll(items map[string]any, keys ...string) bool
```

Reports whether **all** of the given dot-notated keys exist in the map.

```go
data := map[string]any{
    "host": "localhost",
    "port": 5432,
}

kv.HasAll(data, "host", "port")          // true
kv.HasAll(data, "host", "port", "user")  // false
```

**Why it is useful:** Validates that all required configuration keys are present before proceeding.

---

### HasAny

```go
func HasAny(items map[string]any, keys ...string) bool
```

Reports whether **at least one** of the given dot-notated keys exists in the map.

```go
kv.HasAny(data, "host", "missing_key")  // true
kv.HasAny(data, "x", "y", "z")          // false
```

**Why it is useful:** Checks whether any of several alternative keys are available.

---

### Dot

```go
func Dot(items map[string]any, prepend ...string) map[string]any
```

Flattens a nested map into a single-level map with dot-notated keys. An optional prefix string is prepended to every key.

```go
nested := map[string]any{
    "user": map[string]any{
        "name": "Alice",
        "address": map[string]any{
            "city":  "Portland",
            "state": "OR",
        },
    },
}

flat := kv.Dot(nested)
// map[
//   "user.name":          "Alice",
//   "user.address.city":  "Portland",
//   "user.address.state": "OR",
// ]

prefixed := kv.Dot(nested, "config")
// map[
//   "config.user.name":          "Alice",
//   "config.user.address.city":  "Portland",
//   "config.user.address.state": "OR",
// ]
```

**Why it is useful:** Converts hierarchical data into a flat key-value representation, useful for logging, diffing, or storing as environment variables.

---

### Undot

```go
func Undot(items map[string]any) map[string]any
```

Expands a flat map with dot-notated keys back into a nested map structure.

```go
flat := map[string]any{
    "database.host": "localhost",
    "database.port": 5432,
    "app.name":      "MyApp",
}

nested := kv.Undot(flat)
// map[
//   "database": map["host": "localhost", "port": 5432],
//   "app":      map["name": "MyApp"],
// ]
```

**Why it is useful:** The inverse of `Dot` -- reconstructs nested structures from flat key-value stores (environment variables, flat config files).

---

## Map Filtering and Selection

---

### Only

```go
func Only(items map[string]any, keys ...string) map[string]any
```

Returns a new map containing only the specified keys. Keys that do not exist in the source are silently skipped.

```go
data := map[string]any{"a": 1, "b": 2, "c": 3, "d": 4}

kv.Only(data, "a", "c")
// map["a": 1, "c": 3]
```

**Why it is useful:** Picks specific keys from a map to create a subset, commonly used to whitelist fields before serialization.

---

### Except

```go
func Except(items map[string]any, keys ...string) map[string]any
```

Returns a new map with the specified keys removed.

```go
data := map[string]any{"a": 1, "b": 2, "c": 3, "d": 4}

kv.Except(data, "b", "d")
// map["a": 1, "c": 3]
```

**Why it is useful:** The inverse of `Only` -- blacklists keys before passing a map to an API or template.

---

### IsAssoc

```go
func IsAssoc(items map[string]any) bool
```

Reports whether the map is non-empty.

```go
kv.IsAssoc(map[string]any{"key": "val"})  // true
kv.IsAssoc(map[string]any{})              // false
```

**Why it is useful:** A quick check to determine whether a map carries any data before processing it.

---

## Encoding

---

### Query

```go
func Query(items map[string]any) string
```

Encodes the map as a URL query string. Values are converted to strings via `fmt.Sprint`.

```go
params := map[string]any{
    "page":     2,
    "per_page": 25,
    "sort":     "name",
}

kv.Query(params)
// "page=2&per_page=25&sort=name"  (keys are sorted alphabetically)
```

**Why it is useful:** Builds URL query strings from arbitrary maps without manually constructing `url.Values`.

---

### ToCssClasses

```go
func ToCssClasses(classes map[string]bool) string
```

Returns a space-separated string of CSS class names whose boolean values are `true`. Classes are sorted alphabetically.

```go
kv.ToCssClasses(map[string]bool{
    "btn":         true,
    "btn-primary": true,
    "disabled":    false,
    "active":      true,
})
// "active btn btn-primary"
```

**Why it is useful:** Conditionally assembles CSS class lists in server-side template rendering.

---

### ToCssStyles

```go
func ToCssStyles(styles map[string]bool) string
```

Returns a space-separated string of CSS style declarations whose boolean values are `true`. A trailing semicolon is appended to each style if not already present. Styles are sorted alphabetically.

```go
kv.ToCssStyles(map[string]bool{
    "color: red":       true,
    "font-weight: bold": true,
    "display: none":     false,
})
// "color: red; font-weight: bold;"
```

**Why it is useful:** Conditionally builds inline CSS style attributes for server-side HTML generation.

---

## Sorting

---

### Sort

```go
func Sort(items map[string]any) map[string]any
```

Returns a new map constructed by iterating over the keys in sorted (ascending) order.

> **Note:** Go maps do not guarantee iteration order. The returned map is sorted at construction time, but subsequent iteration may not preserve the order.

```go
data := map[string]any{"c": 3, "a": 1, "b": 2}
sorted := kv.Sort(data)
// Constructed in order: a=1, b=2, c=3
```

**Why it is useful:** Produces deterministic key ordering for logging, serialization, or display purposes.

---

### SortRecursive

```go
func SortRecursive(items map[string]any) map[string]any
```

Like `Sort`, but also recursively sorts any nested `map[string]any` values.

```go
data := map[string]any{
    "z": map[string]any{"b": 2, "a": 1},
    "a": "first",
}
sorted := kv.SortRecursive(data)
// Outer keys sorted: a, z
// Inner map {"b": 2, "a": 1} also sorted: a, b
```

**Why it is useful:** Produces fully deterministic output from deeply nested maps, ideal for consistent hashing or snapshot testing.

---

## Map Transformation

---

### Map

```go
func Map[V any, R any](items map[string]V, callback func(V, string) R) map[string]R
```

Applies the callback to each value in the map and returns a new map with the transformed values. The callback receives the value and its key.

```go
prices := map[string]float64{
    "apple":  1.50,
    "banana": 0.75,
}

doubled := kv.Map(prices, func(price float64, _ string) float64 {
    return price * 2
})
// map["apple": 3.0, "banana": 1.5]

labels := kv.Map(prices, func(price float64, fruit string) string {
    return fmt.Sprintf("%s: $%.2f", fruit, price)
})
// map["apple": "apple: $1.50", "banana": "banana: $0.75"]
```

**Why it is useful:** Transforms map values while preserving keys, enabling type-changing transformations via generics.

---

### Where

```go
func Where[V any](items map[string]V, callback func(V, string) bool) map[string]V
```

Returns a new map containing only the entries for which the callback returns `true`.

```go
scores := map[string]int{
    "alice": 92,
    "bob":   67,
    "carol": 85,
}

passing := kv.Where(scores, func(score int, _ string) bool {
    return score >= 80
})
// map["alice": 92, "carol": 85]
```

**Why it is useful:** Filters map entries by value (or key) without constructing a new map manually.

---

### PrependKeysWith

```go
func PrependKeysWith[V any](items map[string]V, prefix string) map[string]V
```

Returns a new map with the given prefix prepended to every key.

```go
env := map[string]string{
    "HOST": "localhost",
    "PORT": "8080",
}

kv.PrependKeysWith(env, "APP_")
// map["APP_HOST": "localhost", "APP_PORT": "8080"]
```

**Why it is useful:** Namespaces map keys in bulk, common when merging configuration from multiple sources.

---

### Replace

```go
func Replace(items map[string]any, replacements ...map[string]any) map[string]any
```

Returns a new map starting with all entries from `items`, then overwrites (or adds) entries from each replacement map in order.

```go
defaults := map[string]any{
    "timeout": 30,
    "retries": 3,
    "debug":   false,
}
overrides := map[string]any{
    "timeout": 60,
    "debug":   true,
}

config := kv.Replace(defaults, overrides)
// map["timeout": 60, "retries": 3, "debug": true]
```

**Why it is useful:** Merges configuration layers (defaults, environment, user overrides) with later maps taking precedence.
