package collection

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// MapCollection wraps a map and provides a fluent API for working with key-value pairs.
// This mirrors PHP's associative array behavior in Laravel Collections.
// Equivalent to: Collection with string keys in PHP
type MapCollection[K comparable, V any] struct {
	items map[K]V
	keys  []K // maintains insertion order
}

// NewMap creates a new MapCollection from a map.
func NewMap[K comparable, V any](items map[K]V) *MapCollection[K, V] {
	if items == nil {
		items = make(map[K]V)
	}
	keys := make([]K, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	return &MapCollection[K, V]{items: items, keys: keys}
}

// NewMapFromPairs creates a MapCollection from key-value pairs.
func NewMapFromPairs[K comparable, V any](pairs ...Pair[K, V]) *MapCollection[K, V] {
	items := make(map[K]V, len(pairs))
	keys := make([]K, 0, len(pairs))
	for _, p := range pairs {
		if _, exists := items[p.Key]; !exists {
			keys = append(keys, p.Key)
		}
		items[p.Key] = p.Value
	}
	return &MapCollection[K, V]{items: items, keys: keys}
}

// All returns the underlying map.
// Equivalent to: $collection->all()
func (m *MapCollection[K, V]) All() map[K]V {
	result := make(map[K]V, len(m.items))
	for k, v := range m.items {
		result[k] = v
	}
	return result
}

// Keys returns all keys as a Collection.
// Equivalent to: $collection->keys()
func (m *MapCollection[K, V]) Keys() *Collection[K] {
	result := make([]K, len(m.keys))
	copy(result, m.keys)
	return Collect(result)
}

// Values returns all values as a Collection.
// Equivalent to: $collection->values()
func (m *MapCollection[K, V]) Values() *Collection[V] {
	result := make([]V, 0, len(m.keys))
	for _, k := range m.keys {
		result = append(result, m.items[k])
	}
	return Collect(result)
}

// Count returns the number of items.
// Equivalent to: $collection->count()
func (m *MapCollection[K, V]) Count() int {
	return len(m.items)
}

// IsEmpty determines if the map collection is empty.
// Equivalent to: $collection->isEmpty()
func (m *MapCollection[K, V]) IsEmpty() bool {
	return len(m.items) == 0
}

// IsNotEmpty determines if the map collection is not empty.
// Equivalent to: $collection->isNotEmpty()
func (m *MapCollection[K, V]) IsNotEmpty() bool {
	return len(m.items) > 0
}

// Get returns the value for the given key.
// Equivalent to: $collection->get($key, $default)
func (m *MapCollection[K, V]) Get(key K, defaults ...V) (V, bool) {
	if v, ok := m.items[key]; ok {
		return v, true
	}
	if len(defaults) > 0 {
		return defaults[0], false
	}
	var zero V
	return zero, false
}

// GetOrPut returns the value for the given key, or puts a default if not present.
// Equivalent to: $collection->getOrPut($key, $value)
func (m *MapCollection[K, V]) GetOrPut(key K, value V) V {
	if v, ok := m.items[key]; ok {
		return v
	}
	m.Put(key, value)
	return value
}

// Has determines if a key exists.
// Equivalent to: $collection->has($key)
func (m *MapCollection[K, V]) Has(key K) bool {
	_, ok := m.items[key]
	return ok
}

// HasAny determines if any of the given keys exist.
// Equivalent to: $collection->hasAny($keys)
func (m *MapCollection[K, V]) HasAny(keys ...K) bool {
	for _, k := range keys {
		if m.Has(k) {
			return true
		}
	}
	return false
}

// Put sets the given key-value pair.
// Equivalent to: $collection->put($key, $value)
func (m *MapCollection[K, V]) Put(key K, value V) *MapCollection[K, V] {
	if _, exists := m.items[key]; !exists {
		m.keys = append(m.keys, key)
	}
	m.items[key] = value
	return m
}

// Pull removes and returns an item by key.
// Equivalent to: $collection->pull($key)
func (m *MapCollection[K, V]) Pull(key K) (V, bool) {
	v, ok := m.items[key]
	if ok {
		delete(m.items, key)
		for i, k := range m.keys {
			if k == key {
				m.keys = append(m.keys[:i], m.keys[i+1:]...)
				break
			}
		}
	}
	return v, ok
}

// Forget removes one or more items by key.
// Equivalent to: $collection->forget($keys)
func (m *MapCollection[K, V]) Forget(keys ...K) *MapCollection[K, V] {
	for _, key := range keys {
		m.Pull(key)
	}
	return m
}

// Only returns items with only the specified keys.
// Equivalent to: $collection->only($keys)
func (m *MapCollection[K, V]) Only(keys ...K) *MapCollection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0, len(keys))
	for _, k := range keys {
		if v, ok := m.items[k]; ok {
			result[k] = v
			newKeys = append(newKeys, k)
		}
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// Except returns all items except those with the specified keys.
// Equivalent to: $collection->except($keys)
func (m *MapCollection[K, V]) Except(keys ...K) *MapCollection[K, V] {
	excludeSet := make(map[K]bool, len(keys))
	for _, k := range keys {
		excludeSet[k] = true
	}
	result := make(map[K]V)
	newKeys := make([]K, 0)
	for _, k := range m.keys {
		if !excludeSet[k] {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// Contains determines if the map contains a value matching the predicate.
// Equivalent to: $collection->contains($callback)
func (m *MapCollection[K, V]) Contains(predicate func(V, K) bool) bool {
	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			return true
		}
	}
	return false
}

// DoesntContain is the inverse of Contains.
// Equivalent to: $collection->doesntContain($callback)
func (m *MapCollection[K, V]) DoesntContain(predicate func(V, K) bool) bool {
	return !m.Contains(predicate)
}

// Search searches for a value and returns its key.
// Equivalent to: $collection->search($value)
func (m *MapCollection[K, V]) Search(predicate func(V, K) bool) (K, bool) {
	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			return k, true
		}
	}
	var zero K
	return zero, false
}

// First returns the first value matching the predicate.
// Equivalent to: $collection->first($callback)
func (m *MapCollection[K, V]) First(predicates ...func(V, K) bool) (V, bool) {
	if len(m.keys) == 0 {
		var zero V
		return zero, false
	}
	if len(predicates) == 0 || predicates[0] == nil {
		return m.items[m.keys[0]], true
	}
	predicate := predicates[0]
	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			return m.items[k], true
		}
	}
	var zero V
	return zero, false
}

// Last returns the last value matching the predicate.
// Equivalent to: $collection->last($callback)
func (m *MapCollection[K, V]) Last(predicates ...func(V, K) bool) (V, bool) {
	if len(m.keys) == 0 {
		var zero V
		return zero, false
	}
	if len(predicates) == 0 || predicates[0] == nil {
		return m.items[m.keys[len(m.keys)-1]], true
	}
	predicate := predicates[0]
	for i := len(m.keys) - 1; i >= 0; i-- {
		k := m.keys[i]
		if predicate(m.items[k], k) {
			return m.items[k], true
		}
	}
	var zero V
	return zero, false
}

// Each iterates over items. Return false to stop.
// Equivalent to: $collection->each($callback)
func (m *MapCollection[K, V]) Each(callback func(V, K) bool) *MapCollection[K, V] {
	for _, k := range m.keys {
		if !callback(m.items[k], k) {
			break
		}
	}
	return m
}

// Filter returns items that pass the given truth test.
// Equivalent to: $collection->filter($callback)
func (m *MapCollection[K, V]) Filter(callback func(V, K) bool) *MapCollection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)
	for _, k := range m.keys {
		v := m.items[k]
		if callback(v, k) {
			result[k] = v
			newKeys = append(newKeys, k)
		}
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// Reject returns items that don't pass the given truth test.
// Equivalent to: $collection->reject($callback)
func (m *MapCollection[K, V]) Reject(callback func(V, K) bool) *MapCollection[K, V] {
	return m.Filter(func(v V, k K) bool {
		return !callback(v, k)
	})
}

// MapValues transforms values using a callback.
// Equivalent to: $collection->map($callback)
func MapValues[K comparable, V any, R any](m *MapCollection[K, V], callback func(V, K) R) *MapCollection[K, R] {
	result := make(map[K]R, len(m.items))
	newKeys := make([]K, len(m.keys))
	copy(newKeys, m.keys)
	for _, k := range m.keys {
		result[k] = callback(m.items[k], k)
	}
	return &MapCollection[K, R]{items: result, keys: newKeys}
}

// Every determines if all items pass the given test.
// Equivalent to: $collection->every($callback)
func (m *MapCollection[K, V]) Every(callback func(V, K) bool) bool {
	for _, k := range m.keys {
		if !callback(m.items[k], k) {
			return false
		}
	}
	return true
}

// Partition separates items into two map collections.
// Equivalent to: $collection->partition($callback)
func (m *MapCollection[K, V]) Partition(callback func(V, K) bool) (*MapCollection[K, V], *MapCollection[K, V]) {
	pass := make(map[K]V)
	passKeys := make([]K, 0)
	fail := make(map[K]V)
	failKeys := make([]K, 0)
	for _, k := range m.keys {
		v := m.items[k]
		if callback(v, k) {
			pass[k] = v
			passKeys = append(passKeys, k)
		} else {
			fail[k] = v
			failKeys = append(failKeys, k)
		}
	}
	return &MapCollection[K, V]{items: pass, keys: passKeys},
		&MapCollection[K, V]{items: fail, keys: failKeys}
}

// Merge merges the given map into the collection.
// Equivalent to: $collection->merge($items)
func (m *MapCollection[K, V]) Merge(items map[K]V) *MapCollection[K, V] {
	result := make(map[K]V, len(m.items)+len(items))
	newKeys := make([]K, len(m.keys))
	copy(newKeys, m.keys)
	for k, v := range m.items {
		result[k] = v
	}
	for k, v := range items {
		if _, exists := result[k]; !exists {
			newKeys = append(newKeys, k)
		}
		result[k] = v
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// Union creates a union of the map collection with the given map.
// Equivalent to: $collection->union($items)
func (m *MapCollection[K, V]) Union(items map[K]V) *MapCollection[K, V] {
	result := make(map[K]V, len(m.items)+len(items))
	newKeys := make([]K, len(m.keys))
	copy(newKeys, m.keys)
	for k, v := range m.items {
		result[k] = v
	}
	for k, v := range items {
		if _, exists := result[k]; !exists {
			result[k] = v
			newKeys = append(newKeys, k)
		}
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// Replace replaces items in the collection with the given map.
// Equivalent to: $collection->replace($items)
func (m *MapCollection[K, V]) Replace(items map[K]V) *MapCollection[K, V] {
	return m.Merge(items)
}

// DiffKeys returns items whose keys are not present in the given map.
// Equivalent to: $collection->diffKeys($items)
func (m *MapCollection[K, V]) DiffKeys(items map[K]V) *MapCollection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)
	for _, k := range m.keys {
		if _, exists := items[k]; !exists {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// IntersectByKeys returns items whose keys are present in the given map.
// Equivalent to: $collection->intersectByKeys($items)
func (m *MapCollection[K, V]) IntersectByKeys(items map[K]V) *MapCollection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)
	for _, k := range m.keys {
		if _, exists := items[k]; exists {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// Flip swaps keys and values (requires V to be comparable).
// Equivalent to: $collection->flip()
func MapFlip[K comparable, V comparable](m *MapCollection[K, V]) *MapCollection[V, K] {
	result := make(map[V]K, len(m.items))
	newKeys := make([]V, 0, len(m.keys))
	for _, k := range m.keys {
		v := m.items[k]
		if _, exists := result[v]; !exists {
			newKeys = append(newKeys, v)
		}
		result[v] = k
	}
	return &MapCollection[V, K]{items: result, keys: newKeys}
}

// SortKeys sorts the collection by its keys.
// Equivalent to: $collection->sortKeys()
func MapSortKeys[V any](m *MapCollection[string, V]) *MapCollection[string, V] {
	result := make(map[string]V, len(m.items))
	for k, v := range m.items {
		result[k] = v
	}
	newKeys := make([]string, len(m.keys))
	copy(newKeys, m.keys)
	sort.Strings(newKeys)
	return &MapCollection[string, V]{items: result, keys: newKeys}
}

// SortKeysDesc sorts the collection by its keys in descending order.
// Equivalent to: $collection->sortKeysDesc()
func MapSortKeysDesc[V any](m *MapCollection[string, V]) *MapCollection[string, V] {
	result := make(map[string]V, len(m.items))
	for k, v := range m.items {
		result[k] = v
	}
	newKeys := make([]string, len(m.keys))
	copy(newKeys, m.keys)
	sort.Sort(sort.Reverse(sort.StringSlice(newKeys)))
	return &MapCollection[string, V]{items: result, keys: newKeys}
}

// SortKeysUsing sorts the collection by its keys using a custom comparison.
// Equivalent to: $collection->sortKeysUsing($callback)
func (m *MapCollection[K, V]) SortKeysUsing(less func(K, K) bool) *MapCollection[K, V] {
	result := make(map[K]V, len(m.items))
	for k, v := range m.items {
		result[k] = v
	}
	newKeys := make([]K, len(m.keys))
	copy(newKeys, m.keys)
	sort.SliceStable(newKeys, func(i, j int) bool {
		return less(newKeys[i], newKeys[j])
	})
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// Implode joins values into a string with a glue.
// Equivalent to: $collection->implode($glue)
func (m *MapCollection[K, V]) Implode(glue string) string {
	parts := make([]string, 0, len(m.keys))
	for _, k := range m.keys {
		parts = append(parts, fmt.Sprint(m.items[k]))
	}
	return strings.Join(parts, glue)
}

// Join joins values with glue and an optional final glue.
// Equivalent to: $collection->join($glue, $finalGlue)
func (m *MapCollection[K, V]) Join(glue string, finalGlues ...string) string {
	parts := make([]string, 0, len(m.keys))
	for _, k := range m.keys {
		parts = append(parts, fmt.Sprint(m.items[k]))
	}
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	if len(finalGlues) > 0 && finalGlues[0] != "" {
		last := parts[len(parts)-1]
		rest := parts[:len(parts)-1]
		return strings.Join(rest, glue) + finalGlues[0] + last
	}
	return strings.Join(parts, glue)
}

// Tap passes the map collection to the callback.
// Equivalent to: $collection->tap($callback)
func (m *MapCollection[K, V]) Tap(callback func(*MapCollection[K, V])) *MapCollection[K, V] {
	callback(m)
	return m
}

// When applies the callback if condition is true.
// Equivalent to: $collection->when($condition, $callback, $default)
func (m *MapCollection[K, V]) When(condition bool, callback func(*MapCollection[K, V]) *MapCollection[K, V], defaults ...func(*MapCollection[K, V]) *MapCollection[K, V]) *MapCollection[K, V] {
	if condition {
		return callback(m)
	}
	if len(defaults) > 0 {
		return defaults[0](m)
	}
	return m
}

// Unless applies the callback unless the condition is true.
// Equivalent to: $collection->unless($condition, $callback, $default)
func (m *MapCollection[K, V]) Unless(condition bool, callback func(*MapCollection[K, V]) *MapCollection[K, V], defaults ...func(*MapCollection[K, V]) *MapCollection[K, V]) *MapCollection[K, V] {
	return m.When(!condition, callback, defaults...)
}

// ToJSON serializes the map collection to JSON.
// Equivalent to: $collection->toJson()
func (m *MapCollection[K, V]) ToJSON() ([]byte, error) {
	return json.Marshal(m.items)
}

// ToPrettyJSON serializes to indented JSON.
// Equivalent to: $collection->toPrettyJson()
func (m *MapCollection[K, V]) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(m.items, "", "    ")
}

// String returns the JSON representation.
// Equivalent to: $collection->__toString()
func (m *MapCollection[K, V]) String() string {
	b, err := m.ToJSON()
	if err != nil {
		return "{}"
	}
	return string(b)
}

// MarshalJSON implements json.Marshaler.
func (m *MapCollection[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.items)
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *MapCollection[K, V]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.items)
}

// Copy creates a shallow copy.
func (m *MapCollection[K, V]) Copy() *MapCollection[K, V] {
	return NewMap(m.All())
}

// Dump prints the map for debugging.
// Equivalent to: $collection->dump()
func (m *MapCollection[K, V]) Dump() *MapCollection[K, V] {
	fmt.Printf("%v\n", m.items)
	return m
}

// ToPairs converts the map to a Collection of Pairs.
func (m *MapCollection[K, V]) ToPairs() *Collection[Pair[K, V]] {
	result := make([]Pair[K, V], 0, len(m.keys))
	for _, k := range m.keys {
		result = append(result, Pair[K, V]{Key: k, Value: m.items[k]})
	}
	return Collect(result)
}

// ContainsOneItem determines if there is exactly one item.
// Equivalent to: $collection->containsOneItem()
func (m *MapCollection[K, V]) ContainsOneItem() bool {
	return len(m.items) == 1
}

// ContainsManyItems determines if there are many items.
// Equivalent to: $collection->containsManyItems()
func (m *MapCollection[K, V]) ContainsManyItems() bool {
	return len(m.items) > 1
}
