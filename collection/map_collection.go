package collection

import (
	"encoding/json"
	"fmt"
	"iter"
	"os"
	"sort"
	"strings"
)

// MapCollection wraps a map and provides a fluent API for working with
// key-value pairs while maintaining insertion order.
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

// All returns a shallow copy of the underlying map.
func (m *MapCollection[K, V]) All() map[K]V {
	result := make(map[K]V, len(m.items))
	for k, v := range m.items {
		result[k] = v
	}
	return result
}

// Keys returns all keys as a Collection, preserving insertion order.
func (m *MapCollection[K, V]) Keys() *Collection[K] {
	result := make([]K, len(m.keys))
	copy(result, m.keys)
	return Collect(result)
}

// Values returns all values as a Collection, preserving insertion order.
func (m *MapCollection[K, V]) Values() *Collection[V] {
	result := make([]V, 0, len(m.keys))
	for _, k := range m.keys {
		result = append(result, m.items[k])
	}
	return Collect(result)
}

// Count returns the number of items in the collection.
func (m *MapCollection[K, V]) Count() int {
	return len(m.items)
}

// IsEmpty reports whether the collection contains no items.
func (m *MapCollection[K, V]) IsEmpty() bool {
	return len(m.items) == 0
}

// IsNotEmpty reports whether the collection contains at least one item.
func (m *MapCollection[K, V]) IsNotEmpty() bool {
	return len(m.items) > 0
}

// Get returns the value for the given key. The second return value indicates
// whether the key was found. An optional default may be provided.
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

// GetOrPut returns the value for the given key if it exists. Otherwise it
// stores and returns the provided default value.
func (m *MapCollection[K, V]) GetOrPut(key K, value V) V {
	if v, ok := m.items[key]; ok {
		return v
	}
	m.Put(key, value)
	return value
}

// Has reports whether the given key exists in the collection.
func (m *MapCollection[K, V]) Has(key K) bool {
	_, ok := m.items[key]
	return ok
}

// HasAny reports whether any of the given keys exist in the collection.
func (m *MapCollection[K, V]) HasAny(keys ...K) bool {
	for _, k := range keys {
		if m.Has(k) {
			return true
		}
	}
	return false
}

// Put sets the given key-value pair in the collection.
func (m *MapCollection[K, V]) Put(key K, value V) *MapCollection[K, V] {
	if _, exists := m.items[key]; !exists {
		m.keys = append(m.keys, key)
	}
	m.items[key] = value
	return m
}

// Pull removes an item by key and returns its value. The second return value
// indicates whether the key was found.
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
func (m *MapCollection[K, V]) Forget(keys ...K) *MapCollection[K, V] {
	for _, key := range keys {
		m.Pull(key)
	}
	return m
}

// Only returns a new collection containing only the specified keys.
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

// Except returns a new collection containing all items except those with the
// specified keys.
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

// Contains reports whether any item in the collection satisfies the predicate.
func (m *MapCollection[K, V]) Contains(predicate func(V, K) bool) bool {
	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			return true
		}
	}
	return false
}

// Some is an alias for Contains.
func (m *MapCollection[K, V]) Some(predicate func(V, K) bool) bool {
	return m.Contains(predicate)
}

// DoesntContain reports whether no item in the collection satisfies the predicate.
func (m *MapCollection[K, V]) DoesntContain(predicate func(V, K) bool) bool {
	return !m.Contains(predicate)
}

// HasSole reports whether exactly one item in the collection satisfies the predicate.
func (m *MapCollection[K, V]) HasSole(predicate func(V, K) bool) bool {
	count := 0
	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return count == 1
}

// Search returns the key of the first item that satisfies the predicate.
// The second return value indicates whether a match was found.
func (m *MapCollection[K, V]) Search(predicate func(V, K) bool) (K, bool) {
	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			return k, true
		}
	}
	var zero K
	return zero, false
}

// First returns the first value that satisfies the optional predicate. If no
// predicate is given, the first item in insertion order is returned.
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

// Last returns the last value that satisfies the optional predicate. If no
// predicate is given, the last item in insertion order is returned.
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

// Each iterates over items in insertion order, calling the callback for each
// key-value pair. Return false from the callback to stop iteration.
func (m *MapCollection[K, V]) Each(callback func(V, K) bool) *MapCollection[K, V] {
	for _, k := range m.keys {
		if !callback(m.items[k], k) {
			break
		}
	}
	return m
}

// Filter returns a new collection containing only the items for which the
// callback returns true.
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

// Reject returns a new collection containing only the items for which the
// callback returns false.
func (m *MapCollection[K, V]) Reject(callback func(V, K) bool) *MapCollection[K, V] {
	return m.Filter(func(v V, k K) bool {
		return !callback(v, k)
	})
}

// MapValues transforms each value using the callback, returning a new
// MapCollection with the same keys and transformed values.
func MapValues[K comparable, V any, R any](m *MapCollection[K, V], callback func(V, K) R) *MapCollection[K, R] {
	result := make(map[K]R, len(m.items))
	newKeys := make([]K, len(m.keys))
	copy(newKeys, m.keys)
	for _, k := range m.keys {
		result[k] = callback(m.items[k], k)
	}
	return &MapCollection[K, R]{items: result, keys: newKeys}
}

// Every reports whether all items in the collection satisfy the callback.
func (m *MapCollection[K, V]) Every(callback func(V, K) bool) bool {
	for _, k := range m.keys {
		if !callback(m.items[k], k) {
			return false
		}
	}
	return true
}

// Partition splits the collection into two: one where the callback returns
// true, and one where it returns false.
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

// Merge returns a new collection with the given map merged in. Existing keys
// are overwritten by the incoming values.
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

// Union returns a new collection that is the union of the collection and the
// given map. Keys already present in the collection are not overwritten.
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

// Replace returns a new collection with the given map merged in. This is an
// alias for Merge.
func (m *MapCollection[K, V]) Replace(items map[K]V) *MapCollection[K, V] {
	return m.Merge(items)
}

// MapMergeRecursive recursively merges the given map into the collection.
// When both sides have a map value for the same key, they are merged recursively.
func MapMergeRecursive[V any](m *MapCollection[string, V], items map[string]V) *MapCollection[string, V] {
	result := make(map[string]V, len(m.items)+len(items))
	newKeys := make([]string, len(m.keys))
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
	return &MapCollection[string, V]{items: result, keys: newKeys}
}

// DiffKeys returns a new collection containing items whose keys are not
// present in the given map.
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

// DiffKeysUsing returns items whose keys are not considered equal to any key in the given map,
// using the provided comparison function.
func (m *MapCollection[K, V]) DiffKeysUsing(items map[K]V, equals func(K, K) bool) *MapCollection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)
	for _, k := range m.keys {
		found := false
		for otherK := range items {
			if equals(k, otherK) {
				found = true
				break
			}
		}
		if !found {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// IntersectByKeys returns a new collection containing items whose keys are
// present in the given map.
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

// DiffAssoc returns items whose key is not present in the given map or whose value differs.
func MapDiffAssoc[K comparable, V comparable](m *MapCollection[K, V], items map[K]V) *MapCollection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)
	for _, k := range m.keys {
		otherVal, exists := items[k]
		if !exists || m.items[k] != otherVal {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// IntersectAssoc returns items whose key exists in the given map and whose value matches.
func MapIntersectAssoc[K comparable, V comparable](m *MapCollection[K, V], items map[K]V) *MapCollection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)
	for _, k := range m.keys {
		if otherVal, exists := items[k]; exists && m.items[k] == otherVal {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}
	return &MapCollection[K, V]{items: result, keys: newKeys}
}

// MapFlip swaps keys and values, returning a new MapCollection where the
// original values become keys and original keys become values.
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

// MapSortKeys returns a new collection with string keys sorted in ascending
// order.
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

// MapSortKeysDesc returns a new collection with string keys sorted in
// descending order.
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

// SortKeysUsing returns a new collection with keys sorted using the provided
// comparison function.
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

// Implode concatenates all values into a single string, separated by the
// given glue string.
func (m *MapCollection[K, V]) Implode(glue string) string {
	parts := make([]string, 0, len(m.keys))
	for _, k := range m.keys {
		parts = append(parts, fmt.Sprint(m.items[k]))
	}
	return strings.Join(parts, glue)
}

// Join concatenates all values into a string separated by glue. An optional
// final glue is placed between the last two items.
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

// Tap passes the collection to the callback and returns the collection,
// allowing side effects without breaking a method chain.
func (m *MapCollection[K, V]) Tap(callback func(*MapCollection[K, V])) *MapCollection[K, V] {
	callback(m)
	return m
}

// When applies the callback if the condition is true. An optional default
// callback is applied when the condition is false.
func (m *MapCollection[K, V]) When(condition bool, callback func(*MapCollection[K, V]) *MapCollection[K, V], defaults ...func(*MapCollection[K, V]) *MapCollection[K, V]) *MapCollection[K, V] {
	if condition {
		return callback(m)
	}
	if len(defaults) > 0 {
		return defaults[0](m)
	}
	return m
}

// Unless applies the callback unless the condition is true. An optional
// default callback is applied when the condition is true.
func (m *MapCollection[K, V]) Unless(condition bool, callback func(*MapCollection[K, V]) *MapCollection[K, V], defaults ...func(*MapCollection[K, V]) *MapCollection[K, V]) *MapCollection[K, V] {
	return m.When(!condition, callback, defaults...)
}

// ToJSON serializes the collection to JSON.
func (m *MapCollection[K, V]) ToJSON() ([]byte, error) {
	return json.Marshal(m.items)
}

// ToPrettyJSON serializes the collection to indented JSON.
func (m *MapCollection[K, V]) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(m.items, "", "    ")
}

// String returns the JSON representation of the collection.
func (m *MapCollection[K, V]) String() string {
	b, err := m.ToJSON()
	if err != nil {
		return "{}"
	}
	return string(b)
}

// MarshalJSON implements the json.Marshaler interface.
func (m *MapCollection[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.items)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *MapCollection[K, V]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.items)
}

// Copy creates a shallow copy of the collection.
func (m *MapCollection[K, V]) Copy() *MapCollection[K, V] {
	return NewMap(m.All())
}

// Dump prints the underlying map to stdout for debugging purposes.
func (m *MapCollection[K, V]) Dump() *MapCollection[K, V] {
	fmt.Printf("%v\n", m.items)
	return m
}

// DD prints the map collection for debugging and terminates the program.
func (m *MapCollection[K, V]) DD() {
	m.Dump()
	os.Exit(1)
}

// ToPairs converts the collection to a Collection of Pair values, preserving
// insertion order.
func (m *MapCollection[K, V]) ToPairs() *Collection[Pair[K, V]] {
	result := make([]Pair[K, V], 0, len(m.keys))
	for _, k := range m.keys {
		result = append(result, Pair[K, V]{Key: k, Value: m.items[k]})
	}
	return Collect(result)
}

// ContainsOneItem reports whether the collection contains exactly one item.
func (m *MapCollection[K, V]) ContainsOneItem() bool {
	return len(m.items) == 1
}

// ContainsManyItems reports whether the collection contains more than one item.
func (m *MapCollection[K, V]) ContainsManyItems() bool {
	return len(m.items) > 1
}

// Iter returns an iterator over the collection's key-value pairs in insertion
// order.
func (m *MapCollection[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range m.keys {
			if !yield(k, m.items[k]) {
				return
			}
		}
	}
}
