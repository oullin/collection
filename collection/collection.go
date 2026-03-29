// Package collection provides a fluent, generic wrapper for working with slices of data.
package collection

import (
	"cmp"
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"math/rand/v2"
	"os"
	"slices"
	"sort"
	"strings"
)

// Collection wraps a slice and provides a fluent API for working with arrays of data.
type Collection[T any] struct {
	items []T
}

// New creates a new Collection from the given items.
func New[T any](items ...T) *Collection[T] {
	if items == nil {
		items = make([]T, 0)
	}

	return &Collection[T]{items: items}
}

// Collect creates a new Collection from a slice.
func Collect[T any](items []T) *Collection[T] {
	if items == nil {
		items = make([]T, 0)
	}

	return &Collection[T]{items: items}
}

// Empty creates an empty Collection.
func Empty[T any]() *Collection[T] {
	return &Collection[T]{items: make([]T, 0)}
}

// Wrap wraps the given value in a collection if it is not already one.
func Wrap[T any](value any) *Collection[T] {
	switch v := value.(type) {
	case *Collection[T]:
		return v
	case []T:
		return Collect(v)
	default:
		if item, ok := value.(T); ok {
			return New(item)
		}

		return Empty[T]()
	}
}

// Unwrap returns the underlying items from a Collection, or the value itself if it is a slice.
func Unwrap[T any](value any) []T {
	if c, ok := value.(*Collection[T]); ok {
		return c.All()
	}

	if items, ok := value.([]T); ok {
		return items
	}

	return nil
}

// Times creates a new collection by invoking the callback a given number of times.
func Times[T any](number int, callback func(int) T) *Collection[T] {
	if number < 1 {
		return Empty[T]()
	}

	items := make([]T, number)

	for i := 0; i < number; i++ {
		items[i] = callback(i + 1)
	}

	return Collect(items)
}

// Range creates a collection of consecutive integers from start to end (inclusive).
func Range(from, to int) *Collection[int] {
	if from > to {
		items := make([]int, 0, from-to+1)

		for i := from; i >= to; i-- {
			items = append(items, i)
		}

		return Collect(items)
	}

	items := make([]int, 0, to-from+1)

	for i := from; i <= to; i++ {
		items = append(items, i)
	}

	return Collect(items)
}

// All returns all items in the collection as a slice.
func (c *Collection[T]) All() []T {
	return c.items
}

// Count returns the total number of items in the collection.
func (c *Collection[T]) Count() int {
	return len(c.items)
}

// IsEmpty reports whether the collection contains no items.
func (c *Collection[T]) IsEmpty() bool {
	return len(c.items) == 0
}

// IsNotEmpty reports whether the collection contains at least one item.
func (c *Collection[T]) IsNotEmpty() bool {
	return len(c.items) > 0
}

// ContainsOneItem reports whether the collection contains exactly one item.
func (c *Collection[T]) ContainsOneItem() bool {
	return len(c.items) == 1
}

// ContainsManyItems reports whether the collection contains more than one item.
func (c *Collection[T]) ContainsManyItems() bool {
	return len(c.items) > 1
}

// HasMany is an alias for ContainsManyItems.
func (c *Collection[T]) HasMany() bool {
	return c.ContainsManyItems()
}

// First returns the first element matching the optional predicate.
// If no predicate is provided, the first element is returned.
// The second return value indicates whether a match was found.
func (c *Collection[T]) First(predicates ...func(T, int) bool) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	if len(predicates) == 0 || predicates[0] == nil {
		return c.items[0], true
	}

	predicate := predicates[0]

	for i, item := range c.items {
		if predicate(item, i) {
			return item, true
		}
	}

	var zero T

	return zero, false
}

// FirstOrFail returns the first element matching the optional predicate,
// or an ItemNotFoundError if no match is found.
func (c *Collection[T]) FirstOrFail(predicates ...func(T, int) bool) (T, error) {
	item, ok := c.First(predicates...)

	if !ok {
		var zero T

		return zero, &ItemNotFoundError{}
	}

	return item, nil
}

// Last returns the last element matching the optional predicate.
// If no predicate is provided, the last element is returned.
// The second return value indicates whether a match was found.
func (c *Collection[T]) Last(predicates ...func(T, int) bool) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	if len(predicates) == 0 || predicates[0] == nil {
		return c.items[len(c.items)-1], true
	}

	predicate := predicates[0]

	for i := len(c.items) - 1; i >= 0; i-- {
		if predicate(c.items[i], i) {
			return c.items[i], true
		}
	}

	var zero T

	return zero, false
}

// Sole returns the only element matching the optional predicate.
// It returns an ItemNotFoundError if no items match, or a MultipleItemsFoundError
// if more than one item matches.
func (c *Collection[T]) Sole(predicates ...func(T, int) bool) (T, error) {
	var filtered *Collection[T]

	if len(predicates) == 0 || predicates[0] == nil {
		filtered = c
	} else {
		filtered = c.Filter(predicates[0])
	}

	if filtered.Count() == 0 {
		var zero T

		return zero, &ItemNotFoundError{}
	}

	if filtered.Count() > 1 {
		var zero T

		return zero, &MultipleItemsFoundError{Count: filtered.Count()}
	}

	return filtered.items[0], nil
}

// HasSole reports whether exactly one item matches the optional predicate.
func (c *Collection[T]) HasSole(predicates ...func(T, int) bool) bool {
	var filtered *Collection[T]

	if len(predicates) == 0 || predicates[0] == nil {
		filtered = c
	} else {
		filtered = c.Filter(predicates[0])
	}

	return filtered.Count() == 1
}

// Get returns the item at the given index.
// Negative indices count from the end. The second return value indicates
// whether the index was within bounds. An optional default may be provided.
func (c *Collection[T]) Get(index int, defaults ...T) (T, bool) {
	if index < 0 {
		index = len(c.items) + index
	}

	if index >= 0 && index < len(c.items) {
		return c.items[index], true
	}

	if len(defaults) > 0 {
		return defaults[0], false
	}

	var zero T

	return zero, false
}

// GetOrPut returns the item at the given index if it exists.
// Otherwise it appends the value to the collection and returns it.
func (c *Collection[T]) GetOrPut(index int, value T) T {
	if index >= 0 && index < len(c.items) {
		return c.items[index]
	}

	c.items = append(c.items, value)

	return value
}

// Put sets the item at the given index to the given value.
func (c *Collection[T]) Put(index int, value T) *Collection[T] {
	if index >= 0 && index < len(c.items) {
		c.items[index] = value
	}

	return c
}

// Pull removes and returns the item at the given index.
// The second return value indicates whether the index was valid.
func (c *Collection[T]) Pull(index int) (T, bool) {
	if index < 0 || index >= len(c.items) {
		var zero T

		return zero, false
	}

	item := c.items[index]
	c.items = append(c.items[:index], c.items[index+1:]...)

	return item, true
}

// Contains reports whether any item in the collection satisfies the predicate.
func (c *Collection[T]) Contains(predicate func(T, int) bool) bool {
	for i, item := range c.items {
		if predicate(item, i) {
			return true
		}
	}

	return false
}

// Some is an alias for Contains.
func (c *Collection[T]) Some(predicate func(T, int) bool) bool {
	return c.Contains(predicate)
}

// DoesntContain reports whether no item in the collection satisfies the predicate.
func (c *Collection[T]) DoesntContain(predicate func(T, int) bool) bool {
	return !c.Contains(predicate)
}

// Search returns the index of the first item satisfying the predicate.
// The second return value indicates whether a match was found.
func (c *Collection[T]) Search(predicate func(T, int) bool) (int, bool) {
	for i, item := range c.items {
		if predicate(item, i) {
			return i, true
		}
	}

	return -1, false
}

// Before returns the item immediately before the first item matching the predicate.
func (c *Collection[T]) Before(predicate func(T, int) bool) (T, bool) {
	for i, item := range c.items {
		if predicate(item, i) {
			if i == 0 {
				var zero T

				return zero, false
			}

			return c.items[i-1], true
		}
	}

	var zero T

	return zero, false
}

// After returns the item immediately after the first item matching the predicate.
func (c *Collection[T]) After(predicate func(T, int) bool) (T, bool) {
	for i, item := range c.items {
		if predicate(item, i) {
			if i >= len(c.items)-1 {
				var zero T

				return zero, false
			}

			return c.items[i+1], true
		}
	}

	var zero T

	return zero, false
}

// Push appends one or more items to the end of the collection.
func (c *Collection[T]) Push(values ...T) *Collection[T] {
	c.items = append(c.items, values...)

	return c
}

// Add appends a single item to the end of the collection.
func (c *Collection[T]) Add(item T) *Collection[T] {
	return c.Push(item)
}

// Prepend adds an item to the beginning of the collection.
func (c *Collection[T]) Prepend(value T) *Collection[T] {
	c.items = append([]T{value}, c.items...)

	return c
}

// Unshift is an alias for Prepend.
func (c *Collection[T]) Unshift(value T) *Collection[T] {
	return c.Prepend(value)
}

// Pop removes and returns the last item from the collection.
// The second return value indicates whether the collection was non-empty.
func (c *Collection[T]) Pop(counts ...int) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	count := 1

	if len(counts) > 0 {
		count = counts[0]
	}

	_ = count // For single pop
	item := c.items[len(c.items)-1]
	c.items = c.items[:len(c.items)-1]

	return item, true
}

// PopMany removes and returns the last n items from the collection.
func (c *Collection[T]) PopMany(count int) *Collection[T] {
	if count >= len(c.items) {
		popped := Collect(c.items)
		c.items = make([]T, 0)

		return popped
	}

	idx := len(c.items) - count
	popped := Collect(c.items[idx:])
	c.items = c.items[:idx]

	return popped
}

// Shift removes and returns the first item from the collection.
// The second return value indicates whether the collection was non-empty.
func (c *Collection[T]) Shift() (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	item := c.items[0]
	c.items = c.items[1:]

	return item, true
}

// ShiftMany removes and returns the first n items from the collection.
func (c *Collection[T]) ShiftMany(count int) *Collection[T] {
	if count >= len(c.items) {
		shifted := Collect(c.items)
		c.items = make([]T, 0)

		return shifted
	}

	shifted := Collect(c.items[:count])
	c.items = c.items[count:]

	return shifted
}

// Each iterates over the items, passing each item and its index to the callback.
// Return false from the callback to stop iterating.
func (c *Collection[T]) Each(callback func(T, int) bool) *Collection[T] {
	for i, item := range c.items {
		if !callback(item, i) {
			break
		}
	}

	return c
}

// EachSpread iterates over the collection's items, passing each item value
// into the given callback. In Go this operates the same as Each.
func (c *Collection[T]) EachSpread(callback func(T, int) bool) *Collection[T] {
	return c.Each(callback)
}

// Tap passes the collection to the given callback and returns the collection unchanged.
func (c *Collection[T]) Tap(callback func(*Collection[T])) *Collection[T] {
	callback(c)

	return c
}

// Pipe passes the collection to the given callback and returns the callback's result.
func Pipe[T any, R any](c *Collection[T], callback func(*Collection[T]) R) R {
	return callback(c)
}

// PipeInto passes the collection to the given constructor and returns the result.
func PipeInto[T any, R any](c *Collection[T], constructor func(*Collection[T]) R) R {
	return constructor(c)
}

// PipeThrough passes the collection through a series of callbacks, returning the final result.
func PipeThrough[T any](c *Collection[T], callbacks ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	result := c

	for _, cb := range callbacks {
		result = cb(result)
	}

	return result
}

// Filter returns a new collection containing only items for which the callback returns true.
func (c *Collection[T]) Filter(callback func(T, int) bool) *Collection[T] {
	result := make([]T, 0)

	for i, item := range c.items {
		if callback(item, i) {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// Reject returns a new collection containing only items for which the callback returns false.
func (c *Collection[T]) Reject(callback func(T, int) bool) *Collection[T] {
	return c.Filter(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Map applies the callback to each item and returns a new collection of results.
func Map[T any, R any](c *Collection[T], callback func(T, int) R) *Collection[R] {
	result := make([]R, len(c.items))

	for i, item := range c.items {
		result[i] = callback(item, i)
	}

	return Collect(result)
}

// Transform applies the callback to each item in place, mutating the collection.
func (c *Collection[T]) Transform(callback func(T, int) T) *Collection[T] {
	for i, item := range c.items {
		c.items[i] = callback(item, i)
	}

	return c
}

// FlatMap applies the callback to each item, flattening the resulting slices into a single collection.
func FlatMap[T any, R any](c *Collection[T], callback func(T, int) []R) *Collection[R] {
	result := make([]R, 0)

	for i, item := range c.items {
		result = append(result, callback(item, i)...)
	}

	return Collect(result)
}

// MapInto applies the constructor to each item, returning a new collection of the mapped type.
func MapInto[T any, R any](c *Collection[T], constructor func(T) R) *Collection[R] {
	result := make([]R, len(c.items))

	for i, item := range c.items {
		result[i] = constructor(item)
	}

	return Collect(result)
}

// Reduce iterates over the collection and accumulates a single result using the callback.
func Reduce[T any, R any](c *Collection[T], callback func(R, T, int) R, initial R) R {
	result := initial

	for i, item := range c.items {
		result = callback(result, item, i)
	}

	return result
}

// Flatten returns a shallow copy of the collection.
// For non-nested typed slices this returns the items as-is.
func (c *Collection[T]) Flatten() *Collection[T] {
	return Collect(append([]T{}, c.items...))
}

// Chunk breaks the collection into multiple slices of the given size.
func (c *Collection[T]) Chunk(size int) [][]T {
	if size <= 0 {
		return nil
	}

	chunks := make([][]T, 0)

	for i := 0; i < len(c.items); i += size {
		end := i + size

		if end > len(c.items) {
			end = len(c.items)
		}

		chunk := make([]T, end-i)
		copy(chunk, c.items[i:end])
		chunks = append(chunks, chunk)
	}

	return chunks
}

// ChunkWhile breaks the collection into groups as long as the callback returns true.
// A new group is started each time the callback returns false.
func (c *Collection[T]) ChunkWhile(callback func(T, int, []T) bool) [][]T {
	if len(c.items) == 0 {
		return nil
	}

	chunks := make([][]T, 0)
	current := []T{c.items[0]}

	for i := 1; i < len(c.items); i++ {
		if callback(c.items[i], i, current) {
			current = append(current, c.items[i])
		} else {
			chunks = append(chunks, current)
			current = []T{c.items[i]}
		}
	}

	chunks = append(chunks, current)

	return chunks
}

// Split breaks the collection into the given number of groups.
func (c *Collection[T]) Split(numberOfGroups int) [][]T {
	if len(c.items) == 0 || numberOfGroups <= 0 {
		return nil
	}

	groups := make([][]T, 0, numberOfGroups)
	groupSize := float64(len(c.items)) / float64(numberOfGroups)

	for i := 0; i < numberOfGroups; i++ {
		start := int(math.Round(float64(i) * groupSize))
		end := int(math.Round(float64(i+1) * groupSize))

		if start >= len(c.items) {
			break
		}

		if end > len(c.items) {
			end = len(c.items)
		}

		chunk := make([]T, end-start)
		copy(chunk, c.items[start:end])
		groups = append(groups, chunk)
	}

	return groups
}

// SplitIn splits the collection into groups, filling non-terminal groups completely.
func (c *Collection[T]) SplitIn(numberOfGroups int) [][]T {
	size := int(math.Ceil(float64(len(c.items)) / float64(numberOfGroups)))

	return c.Chunk(size)
}

// Sliding returns a sliding window view of the collection with the given window size and step.
func (c *Collection[T]) Sliding(size int, steps ...int) [][]T {
	step := 1

	if len(steps) > 0 {
		step = steps[0]
	}

	if size <= 0 || step <= 0 || len(c.items) == 0 {
		return nil
	}

	chunks := make([][]T, 0)

	for i := 0; i+size <= len(c.items); i += step {
		chunk := make([]T, size)
		copy(chunk, c.items[i:i+size])
		chunks = append(chunks, chunk)
	}

	return chunks
}

// Slice extracts a portion of the collection starting at the given offset.
// An optional length limits how many items are returned.
func (c *Collection[T]) Slice(offset int, lengths ...int) *Collection[T] {
	if offset < 0 {
		offset = len(c.items) + offset

		if offset < 0 {
			offset = 0
		}
	}

	if offset >= len(c.items) {
		return Empty[T]()
	}

	if len(lengths) > 0 {
		end := offset + lengths[0]

		if end > len(c.items) {
			end = len(c.items)
		}

		result := make([]T, end-offset)
		copy(result, c.items[offset:end])

		return Collect(result)
	}

	result := make([]T, len(c.items)-offset)
	copy(result, c.items[offset:])

	return Collect(result)
}

// Splice removes and returns a slice of items starting at the given offset.
// An optional length limits the number of items removed.
func (c *Collection[T]) Splice(offset int, lengths ...int) *Collection[T] {
	length := len(c.items) - offset

	if len(lengths) > 0 {
		length = lengths[0]
	}

	if offset < 0 {
		offset = len(c.items) + offset

		if offset < 0 {
			offset = 0
		}
	}

	if offset >= len(c.items) {
		return Empty[T]()
	}

	end := offset + length

	if end > len(c.items) {
		end = len(c.items)
	}

	removed := make([]T, end-offset)
	copy(removed, c.items[offset:end])
	c.items = append(c.items[:offset], c.items[end:]...)

	return Collect(removed)
}

// SpliceReplace removes a portion at the given offset and replaces it with the provided items.
// It returns the removed items.
func (c *Collection[T]) SpliceReplace(offset, length int, replacement []T) *Collection[T] {
	if offset < 0 {
		offset = len(c.items) + offset

		if offset < 0 {
			offset = 0
		}
	}

	end := offset + length

	if end > len(c.items) {
		end = len(c.items)
	}

	removed := make([]T, end-offset)
	copy(removed, c.items[offset:end])
	newItems := make([]T, 0, len(c.items)-length+len(replacement))
	newItems = append(newItems, c.items[:offset]...)
	newItems = append(newItems, replacement...)
	newItems = append(newItems, c.items[end:]...)
	c.items = newItems

	return Collect(removed)
}

// Take returns a new collection with the specified number of items from the front.
// A negative limit takes from the end.
func (c *Collection[T]) Take(limit int) *Collection[T] {
	if limit < 0 {
		return c.Slice(limit)
	}

	return c.Slice(0, limit)
}

// TakeUntil returns items from the start until the callback returns true.
func (c *Collection[T]) TakeUntil(callback func(T, int) bool) *Collection[T] {
	result := make([]T, 0)

	for i, item := range c.items {
		if callback(item, i) {
			break
		}

		result = append(result, item)
	}

	return Collect(result)
}

// TakeWhile returns items from the start as long as the callback returns true.
func (c *Collection[T]) TakeWhile(callback func(T, int) bool) *Collection[T] {
	return c.TakeUntil(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Skip returns a new collection with the first n items removed.
func (c *Collection[T]) Skip(count int) *Collection[T] {
	return c.Slice(count)
}

// SkipUntil skips items until the callback returns true, then returns the rest.
func (c *Collection[T]) SkipUntil(callback func(T, int) bool) *Collection[T] {
	result := make([]T, 0)
	found := false

	for i, item := range c.items {
		if !found && callback(item, i) {
			found = true
		}

		if found {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// SkipWhile skips items as long as the callback returns true, then returns the rest.
func (c *Collection[T]) SkipWhile(callback func(T, int) bool) *Collection[T] {
	return c.SkipUntil(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Nth returns a new collection containing every n-th element, starting at an optional offset.
func (c *Collection[T]) Nth(step int, offsets ...int) *Collection[T] {
	offset := 0

	if len(offsets) > 0 {
		offset = offsets[0]
	}

	result := make([]T, 0)

	for i := offset; i < len(c.items); i += step {
		result = append(result, c.items[i])
	}

	return Collect(result)
}

// ForPage returns a subset of items for the given page number and page size.
func (c *Collection[T]) ForPage(page, perPage int) *Collection[T] {
	offset := (page - 1) * perPage

	return c.Slice(offset, perPage)
}

// Values returns a new collection with re-indexed items (a shallow copy).
func (c *Collection[T]) Values() *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)

	return Collect(result)
}

// Keys returns a new Collection[int] containing the indices 0 through n-1.
func (c *Collection[T]) Keys() *Collection[int] {
	keys := make([]int, len(c.items))

	for i := range c.items {
		keys[i] = i
	}

	return Collect(keys)
}

// Reverse returns a new collection with items in reverse order.
func (c *Collection[T]) Reverse() *Collection[T] {
	result := make([]T, len(c.items))

	for i, j := 0, len(c.items)-1; j >= 0; i, j = i+1, j-1 {
		result[i] = c.items[j]
	}

	return Collect(result)
}

// Shuffle returns a new collection with items in random order.
func (c *Collection[T]) Shuffle() *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return Collect(result)
}

// Random returns a new collection with the specified number of randomly selected items.
func (c *Collection[T]) Random(counts ...int) *Collection[T] {
	count := 1

	if len(counts) > 0 {
		count = counts[0]
	}

	shuffled := c.Shuffle()

	if count >= len(shuffled.items) {
		return shuffled
	}

	return Collect(shuffled.items[:count])
}

// Sort returns a new collection sorted using the provided comparison function.
func (c *Collection[T]) Sort(less func(a, b T) bool) *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	sort.SliceStable(result, func(i, j int) bool {
		return less(result[i], result[j])
	})

	return Collect(result)
}

// SortBy returns a new collection sorted in ascending order by the given key function.
func SortBy[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	sort.SliceStable(result, func(i, j int) bool {
		return keyFunc(result[i]) < keyFunc(result[j])
	})

	return Collect(result)
}

// SortByDesc returns a new collection sorted in descending order by the given key function.
func SortByDesc[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	sort.SliceStable(result, func(i, j int) bool {
		return keyFunc(result[i]) > keyFunc(result[j])
	})

	return Collect(result)
}

// SortDesc returns a new collection sorted in descending order using the provided comparison function.
func (c *Collection[T]) SortDesc(less func(a, b T) bool) *Collection[T] {
	return c.Sort(func(a, b T) bool {
		return less(b, a)
	})
}

// Unique returns a new collection containing only items with distinct keys
// as determined by the given key function.
func Unique[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	seen := make(map[K]bool)
	result := make([]T, 0)

	for _, item := range c.items {
		key := keyFunc(item)

		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return Collect(result)
}

// Duplicates returns a new collection containing all duplicate items
// as determined by the given key function.
func Duplicates[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	seen := make(map[K]bool)
	result := make([]T, 0)

	for _, item := range c.items {
		key := keyFunc(item)

		if seen[key] {
			result = append(result, item)
		} else {
			seen[key] = true
		}
	}

	return Collect(result)
}

// Every reports whether all items in the collection satisfy the given predicate.
func (c *Collection[T]) Every(callback func(T, int) bool) bool {
	for i, item := range c.items {
		if !callback(item, i) {
			return false
		}
	}

	return true
}

// Partition splits the collection into two: items that pass the predicate and items that do not.
func (c *Collection[T]) Partition(callback func(T, int) bool) (*Collection[T], *Collection[T]) {
	pass := make([]T, 0)
	fail := make([]T, 0)

	for i, item := range c.items {
		if callback(item, i) {
			pass = append(pass, item)
		} else {
			fail = append(fail, item)
		}
	}

	return Collect(pass), Collect(fail)
}

// Concat returns a new collection with the given items appended.
func (c *Collection[T]) Concat(items []T) *Collection[T] {
	result := make([]T, len(c.items)+len(items))
	copy(result, c.items)
	copy(result[len(c.items):], items)

	return Collect(result)
}

// Merge returns a new collection with the given items merged in (appended).
func (c *Collection[T]) Merge(items []T) *Collection[T] {
	return c.Concat(items)
}

// Pad pads the collection to the specified length with the given value.
// A negative size pads on the left; a positive size pads on the right.
func (c *Collection[T]) Pad(size int, value T) *Collection[T] {
	absSize := size

	if absSize < 0 {
		absSize = -absSize
	}

	if len(c.items) >= absSize {
		result := make([]T, len(c.items))
		copy(result, c.items)

		return Collect(result)
	}

	padCount := absSize - len(c.items)
	padding := make([]T, padCount)

	for i := range padding {
		padding[i] = value
	}

	if size < 0 {
		result := make([]T, absSize)
		copy(result, padding)
		copy(result[padCount:], c.items)

		return Collect(result)
	}

	result := make([]T, absSize)
	copy(result, c.items)
	copy(result[len(c.items):], padding)

	return Collect(result)
}

// Multiply returns a new collection with all items repeated the given number of times.
func (c *Collection[T]) Multiply(multiplier int) *Collection[T] {
	if multiplier <= 0 {
		return Empty[T]()
	}

	result := make([]T, 0, len(c.items)*multiplier)

	for i := 0; i < multiplier; i++ {
		result = append(result, c.items...)
	}

	return Collect(result)
}

// Flip returns a new collection with the item order reversed.
// For typed Go slices this reverses the element order.
func (c *Collection[T]) Flip() *Collection[T] {
	result := make([]T, len(c.items))

	for i, j := 0, len(c.items)-1; j >= 0; i, j = i+1, j-1 {
		result[i] = c.items[j]
	}

	return Collect(result)
}

// Forget removes an item from the collection by index, mutating the collection.
func (c *Collection[T]) Forget(index int) *Collection[T] {
	if index < 0 || index >= len(c.items) {
		return c
	}

	c.items = append(c.items[:index], c.items[index+1:]...)

	return c
}

// Implode joins elements into a string using the given glue, converting each item via fmt.Sprint.
func (c *Collection[T]) Implode(glue string) string {
	parts := make([]string, len(c.items))

	for i, item := range c.items {
		parts[i] = fmt.Sprint(item)
	}

	return strings.Join(parts, glue)
}

// Join joins elements into a string using the given glue,
// with an optional final separator between the last two items.
func (c *Collection[T]) Join(glue string, finalGlues ...string) string {
	if len(c.items) == 0 {
		return ""
	}

	if len(c.items) == 1 {
		return fmt.Sprint(c.items[0])
	}

	parts := make([]string, len(c.items))

	for i, item := range c.items {
		parts[i] = fmt.Sprint(item)
	}

	if len(finalGlues) > 0 && finalGlues[0] != "" {
		finalGlue := finalGlues[0]
		last := parts[len(parts)-1]
		parts = parts[:len(parts)-1]

		return strings.Join(parts, glue) + finalGlue + last
	}

	return strings.Join(parts, glue)
}

// When applies the callback if the condition is true; otherwise it applies the optional default callback.
func (c *Collection[T]) When(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	if condition {
		return callback(c)
	}

	if len(defaults) > 0 {
		return defaults[0](c)
	}

	return c
}

// WhenEmpty applies the callback when the collection is empty.
func (c *Collection[T]) WhenEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(c.IsEmpty(), callback, defaults...)
}

// WhenNotEmpty applies the callback when the collection is not empty.
func (c *Collection[T]) WhenNotEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(c.IsNotEmpty(), callback, defaults...)
}

// Unless applies the callback unless the condition is true.
func (c *Collection[T]) Unless(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(!condition, callback, defaults...)
}

// UnlessEmpty applies the callback unless the collection is empty.
func (c *Collection[T]) UnlessEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.WhenNotEmpty(callback, defaults...)
}

// UnlessNotEmpty applies the callback unless the collection is not empty.
func (c *Collection[T]) UnlessNotEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.WhenEmpty(callback, defaults...)
}

// Zip merges the collection with each of the given slices element-by-element.
func Zip[T any](c *Collection[T], others ...[]T) *Collection[[]T] {
	maxLen := len(c.items)

	for _, o := range others {
		if len(o) > maxLen {
			maxLen = len(o)
		}
	}

	result := make([][]T, maxLen)

	for i := 0; i < maxLen; i++ {
		group := make([]T, 0, 1+len(others))

		if i < len(c.items) {
			group = append(group, c.items[i])
		} else {
			var zero T

			group = append(group, zero)
		}

		for _, o := range others {
			if i < len(o) {
				group = append(group, o[i])
			} else {
				var zero T

				group = append(group, zero)
			}
		}

		result[i] = group
	}

	return Collect(result)
}

// CrossJoin returns the cross product of the collection with the given slices.
func CrossJoin[T any](c *Collection[T], others ...[]T) *Collection[[]T] {
	results := [][]T{{}}
	allLists := append([][]T{c.items}, others...)

	for _, list := range allLists {
		var newResults [][]T

		for _, result := range results {
			for _, item := range list {
				newResult := make([]T, len(result)+1)
				copy(newResult, result)
				newResult[len(result)] = item
				newResults = append(newResults, newResult)
			}
		}

		results = newResults
	}

	return Collect(results)
}

// Combine pairs keys from this collection with values from the given slice,
// returning a collection of Pair values.
func Combine[K any, V any](keys *Collection[K], values []V) *Collection[Pair[K, V]] {
	minLen := len(keys.items)

	if len(values) < minLen {
		minLen = len(values)
	}

	result := make([]Pair[K, V], minLen)

	for i := 0; i < minLen; i++ {
		result[i] = Pair[K, V]{Key: keys.items[i], Value: values[i]}
	}

	return Collect(result)
}

// Collapse flattens a collection of slices into a single, flat collection.
func Collapse[T any](c *Collection[[]T]) *Collection[T] {
	result := make([]T, 0)

	for _, items := range c.items {
		result = append(result, items...)
	}

	return Collect(result)
}

// Diff returns the items in the collection that are not present in the given slice.
func Diff[T comparable](c *Collection[T], items []T) *Collection[T] {
	lookup := make(map[T]bool, len(items))

	for _, item := range items {
		lookup[item] = true
	}

	result := make([]T, 0)

	for _, item := range c.items {
		if !lookup[item] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// DiffUsing returns items not present in the given slice, using a custom equality function.
func (c *Collection[T]) DiffUsing(items []T, equals func(T, T) bool) *Collection[T] {
	result := make([]T, 0)

	for _, item := range c.items {
		found := false

		for _, other := range items {
			if equals(item, other) {
				found = true

				break
			}
		}

		if !found {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// Intersect returns the items present in both the collection and the given slice.
func Intersect[T comparable](c *Collection[T], items []T) *Collection[T] {
	lookup := make(map[T]bool, len(items))

	for _, item := range items {
		lookup[item] = true
	}

	result := make([]T, 0)

	for _, item := range c.items {
		if lookup[item] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// IntersectUsing returns items present in both the collection and the given slice,
// using a custom equality function.
func (c *Collection[T]) IntersectUsing(items []T, equals func(T, T) bool) *Collection[T] {
	result := make([]T, 0)

	for _, item := range c.items {
		for _, other := range items {
			if equals(item, other) {
				result = append(result, item)

				break
			}
		}
	}

	return Collect(result)
}

// ToSlice returns a copy of the underlying slice.
func (c *Collection[T]) ToSlice() []T {
	result := make([]T, len(c.items))
	copy(result, c.items)

	return result
}

// ToJSON serializes the collection to JSON bytes.
func (c *Collection[T]) ToJSON() ([]byte, error) {
	return json.Marshal(c.items)
}

// ToPrettyJSON serializes the collection to indented JSON bytes.
func (c *Collection[T]) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(c.items, "", "    ")
}

// String returns the JSON string representation of the collection.
func (c *Collection[T]) String() string {
	b, err := c.ToJSON()

	if err != nil {
		return "[]"
	}

	return string(b)
}

// MarshalJSON implements the json.Marshaler interface.
func (c *Collection[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.items)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *Collection[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.items)
}

// Copy creates a shallow copy of the collection.
func (c *Collection[T]) Copy() *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)

	return Collect(result)
}

// Len returns the number of items, implementing sort.Interface.
func (c *Collection[T]) Len() int {
	return len(c.items)
}

// TapEach calls the given callback on each item for side effects, returning the original collection.
func (c *Collection[T]) TapEach(callback func(T, int)) *Collection[T] {
	for i, item := range c.items {
		callback(item, i)
	}

	return c
}

// Dump prints the collection items to stdout for debugging.
func (c *Collection[T]) Dump() *Collection[T] {
	fmt.Printf("%v\n", c.items)

	return c
}

// DD prints the collection items for debugging and terminates the program.
func (c *Collection[T]) DD() {
	c.Dump()
	os.Exit(1)
}

// Only returns a new collection containing only items at the given indices.
func (c *Collection[T]) Only(indices ...int) *Collection[T] {
	result := make([]T, 0, len(indices))

	for _, idx := range indices {
		if idx >= 0 && idx < len(c.items) {
			result = append(result, c.items[idx])
		}
	}

	return Collect(result)
}

// Except returns a new collection excluding items at the given indices.
func (c *Collection[T]) Except(indices ...int) *Collection[T] {
	excludeSet := make(map[int]bool, len(indices))

	for _, idx := range indices {
		excludeSet[idx] = true
	}

	result := make([]T, 0)

	for i, item := range c.items {
		if !excludeSet[i] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// Has reports whether the given index exists in the collection.
// Negative indices count from the end.
func (c *Collection[T]) Has(index int) bool {
	if index < 0 {
		index = len(c.items) + index
	}

	return index >= 0 && index < len(c.items)
}

// HasAny reports whether any of the given indices exist in the collection.
func (c *Collection[T]) HasAny(indices ...int) bool {
	for _, idx := range indices {
		if c.Has(idx) {
			return true
		}
	}

	return false
}

// Pluck extracts a value from each item using the given function, returning a new collection.
func Pluck[T any, V any](c *Collection[T], valueFunc func(T) V) *Collection[V] {
	result := make([]V, len(c.items))

	for i, item := range c.items {
		result[i] = valueFunc(item)
	}

	return Collect(result)
}

// GroupBy groups the collection's items by a key returned from the given function.
func GroupBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]*Collection[T] {
	groups := make(map[K]*Collection[T])

	for _, item := range c.items {
		key := keyFunc(item)

		if _, ok := groups[key]; !ok {
			groups[key] = Empty[T]()
		}

		groups[key].Push(item)
	}

	return groups
}

// KeyBy indexes the collection by a key returned from the given function.
func KeyBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]T {
	result := make(map[K]T)

	for _, item := range c.items {
		result[keyFunc(item)] = item
	}

	return result
}

// CountBy counts how many items produce each key from the given function.
func CountBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]int {
	result := make(map[K]int)

	for _, item := range c.items {
		result[keyFunc(item)]++
	}

	return result
}

// MapToDictionary maps each item to a key-value pair and groups values by key.
func MapToDictionary[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K][]V {
	result := make(map[K][]V)

	for _, item := range c.items {
		key, value := callback(item)
		result[key] = append(result[key], value)
	}

	return result
}

// MapToGroups is an alias for MapToDictionary.
func MapToGroups[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K][]V {
	return MapToDictionary(c, callback)
}

// MapWithKeys maps each item to a key-value pair, returning a map.
func MapWithKeys[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K]V {
	result := make(map[K]V)

	for _, item := range c.items {
		key, value := callback(item)
		result[key] = value
	}

	return result
}

// Where filters items using a predicate that receives only the item (no index).
func (c *Collection[T]) Where(predicate func(T) bool) *Collection[T] {
	return c.Filter(func(item T, _ int) bool {
		return predicate(item)
	})
}

// WhereNot filters items using a negative predicate that receives only the item (no index).
func (c *Collection[T]) WhereNot(predicate func(T) bool) *Collection[T] {
	return c.Filter(func(item T, _ int) bool {
		return !predicate(item)
	})
}

// WhereNull returns items whose extracted value equals the zero value.
func WhereNull[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	var zero K
	result := make([]T, 0)

	for _, item := range c.items {
		if keyFunc(item) == zero {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereNotNull returns items whose extracted value is not the zero value.
func WhereNotNull[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	var zero K
	result := make([]T, 0)

	for _, item := range c.items {
		if keyFunc(item) != zero {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereIn returns items whose extracted key value is in the given set.
func WhereIn[T any, K comparable](c *Collection[T], keyFunc func(T) K, values []K) *Collection[T] {
	set := make(map[K]bool, len(values))

	for _, v := range values {
		set[v] = true
	}

	result := make([]T, 0)

	for _, item := range c.items {
		if set[keyFunc(item)] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereNotIn returns items whose extracted key value is not in the given set.
func WhereNotIn[T any, K comparable](c *Collection[T], keyFunc func(T) K, values []K) *Collection[T] {
	set := make(map[K]bool, len(values))

	for _, v := range values {
		set[v] = true
	}

	result := make([]T, 0)

	for _, item := range c.items {
		if !set[keyFunc(item)] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereBetween returns items whose extracted key value is between min and max (inclusive).
func WhereBetween[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K, min, max K) *Collection[T] {
	result := make([]T, 0)

	for _, item := range c.items {
		v := keyFunc(item)

		if v >= min && v <= max {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereNotBetween returns items whose extracted key value is outside the range [min, max].
func WhereNotBetween[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K, min, max K) *Collection[T] {
	result := make([]T, 0)

	for _, item := range c.items {
		v := keyFunc(item)

		if v < min || v > max {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// Dot returns a shallow copy of the collection.
// For typed Go slices, no dot-notation expansion is possible.
func (c *Collection[T]) Dot() *Collection[T] {
	return c.Copy()
}

// Undot returns a shallow copy of the collection.
// For typed Go slices, no dot-notation expansion is possible.
func (c *Collection[T]) Undot() *Collection[T] {
	return c.Copy()
}

// Ensure verifies that all items satisfy the given predicate, returning an error if any do not.
func (c *Collection[T]) Ensure(predicate func(T) bool) error {
	for _, item := range c.items {
		if !predicate(item) {
			return fmt.Errorf("collection item failed ensure check")
		}
	}

	return nil
}

// ToBase returns the collection itself.
func (c *Collection[T]) ToBase() *Collection[T] {
	return c
}

// Lazy returns a new LazyCollection backed by the collection's items.
func (c *Collection[T]) Lazy() *LazyCollection[T] {
	items := c.items

	return NewLazy(func(yield func(T) bool) {
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	})
}

// Median returns the median value of a float64 collection.
func Median(c *Collection[float64]) float64 {
	if len(c.items) == 0 {
		return 0
	}

	sorted := make([]float64, len(c.items))
	copy(sorted, c.items)
	slices.Sort(sorted)
	mid := len(sorted) / 2

	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}

	return sorted[mid]
}

// MedianBy returns the median value extracted from each item by the given function.
func MedianBy[T any](c *Collection[T], valueFunc func(T) float64) float64 {
	return Median(Map(c, func(item T, _ int) float64 {
		return valueFunc(item)
	}))
}

// Mode returns the most frequently occurring values in the collection.
func Mode[T comparable](c *Collection[T]) []T {
	if len(c.items) == 0 {
		return nil
	}

	counts := make(map[T]int)
	maxCount := 0

	for _, item := range c.items {
		counts[item]++

		if counts[item] > maxCount {
			maxCount = counts[item]
		}
	}

	result := make([]T, 0)

	for item, count := range counts {
		if count == maxCount {
			result = append(result, item)
		}
	}

	return result
}

// Sum returns the sum of all items in a numeric collection.
func Sum[T Numeric](c *Collection[T]) T {
	var total T

	for _, item := range c.items {
		total += item
	}

	return total
}

// SumBy returns the sum of values extracted from each item by the given function.
func SumBy[T any, N Numeric](c *Collection[T], valueFunc func(T) N) N {
	var total N

	for _, item := range c.items {
		total += valueFunc(item)
	}

	return total
}

// Avg returns the arithmetic mean of all items in a numeric collection.
func Avg[T Numeric](c *Collection[T]) float64 {
	if len(c.items) == 0 {
		return 0
	}

	return float64(Sum(c)) / float64(len(c.items))
}

// AvgBy returns the arithmetic mean of values extracted from each item by the given function.
func AvgBy[T any, N Numeric](c *Collection[T], valueFunc func(T) N) float64 {
	if len(c.items) == 0 {
		return 0
	}

	return float64(SumBy(c, valueFunc)) / float64(len(c.items))
}

// Average is an alias for Avg.
func Average[T Numeric](c *Collection[T]) float64 {
	return Avg(c)
}

// Min returns the minimum value in an ordered collection.
// The second return value indicates whether the collection was non-empty.
func Min[T cmp.Ordered](c *Collection[T]) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	result := c.items[0]

	for _, item := range c.items[1:] {
		if item < result {
			result = item
		}
	}

	return result, true
}

// MinBy returns the item with the minimum key as determined by the given function.
// The second return value indicates whether the collection was non-empty.
func MinBy[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	result := c.items[0]
	minKey := keyFunc(result)

	for _, item := range c.items[1:] {
		k := keyFunc(item)

		if k < minKey {
			minKey = k
			result = item
		}
	}

	return result, true
}

// Max returns the maximum value in an ordered collection.
// The second return value indicates whether the collection was non-empty.
func Max[T cmp.Ordered](c *Collection[T]) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	result := c.items[0]

	for _, item := range c.items[1:] {
		if item > result {
			result = item
		}
	}

	return result, true
}

// MaxBy returns the item with the maximum key as determined by the given function.
// The second return value indicates whether the collection was non-empty.
func MaxBy[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	result := c.items[0]
	maxKey := keyFunc(result)

	for _, item := range c.items[1:] {
		k := keyFunc(item)

		if k > maxKey {
			maxKey = k
			result = item
		}
	}

	return result, true
}

// Iter returns an iter.Seq[T] that yields each item in the collection.
func (c *Collection[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range c.items {
			if !yield(item) {
				return
			}
		}
	}
}

// Iter2 returns an iter.Seq2[int, T] that yields each index-item pair in the collection.
func (c *Collection[T]) Iter2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, item := range c.items {
			if !yield(i, item) {
				return
			}
		}
	}
}
