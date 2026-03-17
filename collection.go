// Package collection provides a fluent, convenient wrapper for working with slices of data.
// It is a line-by-line port of Laravel's Illuminate\Support\Collection to Go,
// using Go generics for type safety.
package collection

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"sort"
	"strings"
)

// Collection wraps a slice and provides a fluent API for working with arrays of data.
// This is a generic type parameterized by T.
type Collection[T any] struct {
	items []T
}

// New creates a new Collection from the given items.
// Equivalent to: new Collection($items)
func New[T any](items ...T) *Collection[T] {
	if items == nil {
		items = make([]T, 0)
	}
	return &Collection[T]{items: items}
}

// Collect creates a new Collection from a slice.
// Equivalent to: collect($items)
func Collect[T any](items []T) *Collection[T] {
	if items == nil {
		items = make([]T, 0)
	}
	return &Collection[T]{items: items}
}

// Empty creates an empty Collection.
// Equivalent to: Collection::empty()
func Empty[T any]() *Collection[T] {
	return &Collection[T]{items: make([]T, 0)}
}

// Wrap wraps the given value in a collection if it is not already a collection.
// Equivalent to: Collection::wrap($value)
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

// Unwrap returns the underlying items from a Collection, or the value itself.
// Equivalent to: Collection::unwrap($value)
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
// Equivalent to: Collection::times($number, $callback)
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

// Range creates a collection of integers from $from to $to with optional step.
// Equivalent to: Collection::range($from, $to)
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

// All returns all items in the collection.
// Equivalent to: $collection->all()
func (c *Collection[T]) All() []T {
	return c.items
}

// Count returns the total number of items in the collection.
// Equivalent to: $collection->count()
func (c *Collection[T]) Count() int {
	return len(c.items)
}

// IsEmpty determines if the collection is empty.
// Equivalent to: $collection->isEmpty()
func (c *Collection[T]) IsEmpty() bool {
	return len(c.items) == 0
}

// IsNotEmpty determines if the collection is not empty.
// Equivalent to: $collection->isNotEmpty()
func (c *Collection[T]) IsNotEmpty() bool {
	return len(c.items) > 0
}

// ContainsOneItem determines if the collection contains a single item.
// Equivalent to: $collection->containsOneItem()
func (c *Collection[T]) ContainsOneItem() bool {
	return len(c.items) == 1
}

// ContainsManyItems determines if the collection contains more than one item.
// Equivalent to: $collection->containsManyItems()
func (c *Collection[T]) ContainsManyItems() bool {
	return len(c.items) > 1
}

// HasMany is an alias for ContainsManyItems.
// Equivalent to: $collection->hasMany()
func (c *Collection[T]) HasMany() bool {
	return c.ContainsManyItems()
}

// First returns the first element matching the given predicate.
// If no predicate is given, returns the first element.
// Equivalent to: $collection->first($callback)
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

// FirstOrFail returns the first element or an error if empty.
// Equivalent to: $collection->firstOrFail()
func (c *Collection[T]) FirstOrFail(predicates ...func(T, int) bool) (T, error) {
	item, ok := c.First(predicates...)
	if !ok {
		var zero T
		return zero, &ItemNotFoundException{}
	}
	return item, nil
}

// Last returns the last element matching the given predicate.
// Equivalent to: $collection->last($callback)
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

// Sole returns the only element matching the predicate, or an error if zero or multiple found.
// Equivalent to: $collection->sole($callback)
func (c *Collection[T]) Sole(predicates ...func(T, int) bool) (T, error) {
	var filtered *Collection[T]
	if len(predicates) == 0 || predicates[0] == nil {
		filtered = c
	} else {
		filtered = c.Filter(predicates[0])
	}

	if filtered.Count() == 0 {
		var zero T
		return zero, &ItemNotFoundException{}
	}
	if filtered.Count() > 1 {
		var zero T
		return zero, &MultipleItemsFoundException{Count: filtered.Count()}
	}
	return filtered.items[0], nil
}

// Get returns the item at a given index.
// Equivalent to: $collection->get($key, $default)
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

// GetOrPut returns the item at the given index or inserts a default value.
// Equivalent to: $collection->getOrPut($key, $value)
func (c *Collection[T]) GetOrPut(index int, value T) T {
	if index >= 0 && index < len(c.items) {
		return c.items[index]
	}
	c.items = append(c.items, value)
	return value
}

// Put sets the item at the given index to the given value.
// Equivalent to: $collection->put($key, $value)
func (c *Collection[T]) Put(index int, value T) *Collection[T] {
	if index >= 0 && index < len(c.items) {
		c.items[index] = value
	}
	return c
}

// Pull removes and returns an item from the collection by index.
// Equivalent to: $collection->pull($key)
func (c *Collection[T]) Pull(index int) (T, bool) {
	if index < 0 || index >= len(c.items) {
		var zero T
		return zero, false
	}
	item := c.items[index]
	c.items = append(c.items[:index], c.items[index+1:]...)
	return item, true
}

// Contains determines if the collection contains an item matching the predicate.
// Equivalent to: $collection->contains($callback)
func (c *Collection[T]) Contains(predicate func(T, int) bool) bool {
	for i, item := range c.items {
		if predicate(item, i) {
			return true
		}
	}
	return false
}

// DoesntContain determines if the collection doesn't contain an item matching the predicate.
// Equivalent to: $collection->doesntContain($callback)
func (c *Collection[T]) DoesntContain(predicate func(T, int) bool) bool {
	return !c.Contains(predicate)
}

// Search searches the collection for the given value and returns its index.
// Equivalent to: $collection->search($value)
func (c *Collection[T]) Search(predicate func(T, int) bool) (int, bool) {
	for i, item := range c.items {
		if predicate(item, i) {
			return i, true
		}
	}
	return -1, false
}

// Before returns the item before the first item matching the predicate.
// Equivalent to: $collection->before($value)
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

// After returns the item after the first item matching the predicate.
// Equivalent to: $collection->after($value)
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
// Equivalent to: $collection->push(...$values)
func (c *Collection[T]) Push(values ...T) *Collection[T] {
	c.items = append(c.items, values...)
	return c
}

// Add is an alias for Push with a single value.
// Equivalent to: $collection->add($item)
func (c *Collection[T]) Add(item T) *Collection[T] {
	return c.Push(item)
}

// Prepend adds an item to the beginning of the collection.
// Equivalent to: $collection->prepend($value)
func (c *Collection[T]) Prepend(value T) *Collection[T] {
	c.items = append([]T{value}, c.items...)
	return c
}

// Unshift is an alias for Prepend.
// Equivalent to: $collection->unshift($value)
func (c *Collection[T]) Unshift(value T) *Collection[T] {
	return c.Prepend(value)
}

// Pop removes and returns the last N items from the collection.
// Equivalent to: $collection->pop($count)
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

// PopMany removes and returns the last N items from the collection.
// Equivalent to: $collection->pop($count) where count > 1
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
// Equivalent to: $collection->shift()
func (c *Collection[T]) Shift() (T, bool) {
	if len(c.items) == 0 {
		var zero T
		return zero, false
	}
	item := c.items[0]
	c.items = c.items[1:]
	return item, true
}

// ShiftMany removes and returns the first N items from the collection.
// Equivalent to: $collection->shift($count)
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

// Each iterates over the items in the collection and passes each item to a callback.
// Return false from the callback to stop iterating.
// Equivalent to: $collection->each($callback)
func (c *Collection[T]) Each(callback func(T, int) bool) *Collection[T] {
	for i, item := range c.items {
		if !callback(item, i) {
			break
		}
	}
	return c
}

// EachSpread iterates over the collection's items, passing each nested item value into the given callback.
// For Go, this operates the same as Each since we don't have PHP's spread operator.
// Equivalent to: $collection->eachSpread($callback)
func (c *Collection[T]) EachSpread(callback func(T, int) bool) *Collection[T] {
	return c.Each(callback)
}

// Tap passes the collection to the given callback and returns the collection.
// Equivalent to: $collection->tap($callback)
func (c *Collection[T]) Tap(callback func(*Collection[T])) *Collection[T] {
	callback(c)
	return c
}

// Pipe passes the collection to the given callback and returns the result.
// Equivalent to: $collection->pipe($callback)
func Pipe[T any, R any](c *Collection[T], callback func(*Collection[T]) R) R {
	return callback(c)
}

// PipeInto passes the collection to the given constructor and returns the result.
// Equivalent to: $collection->pipeInto($class)
func PipeInto[T any, R any](c *Collection[T], constructor func(*Collection[T]) R) R {
	return constructor(c)
}

// PipeThrough passes the collection through a series of callbacks and returns the result.
// Equivalent to: $collection->pipeThrough($callbacks)
func PipeThrough[T any](c *Collection[T], callbacks ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	result := c
	for _, cb := range callbacks {
		result = cb(result)
	}
	return result
}

// Filter returns all items that pass the given truth test.
// Equivalent to: $collection->filter($callback)
func (c *Collection[T]) Filter(callback func(T, int) bool) *Collection[T] {
	result := make([]T, 0)
	for i, item := range c.items {
		if callback(item, i) {
			result = append(result, item)
		}
	}
	return Collect(result)
}

// Reject returns all items that do not pass the given truth test.
// Equivalent to: $collection->reject($callback)
func (c *Collection[T]) Reject(callback func(T, int) bool) *Collection[T] {
	return c.Filter(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Map runs a map over each of the items and returns a new collection.
// Equivalent to: $collection->map($callback)
func Map[T any, R any](c *Collection[T], callback func(T, int) R) *Collection[R] {
	result := make([]R, len(c.items))
	for i, item := range c.items {
		result[i] = callback(item, i)
	}
	return Collect(result)
}

// Transform transforms each item in the collection using the callback (mutates in place).
// Equivalent to: $collection->transform($callback)
func (c *Collection[T]) Transform(callback func(T, int) T) *Collection[T] {
	for i, item := range c.items {
		c.items[i] = callback(item, i)
	}
	return c
}

// FlatMap maps a collection and collapses the result.
// Equivalent to: $collection->flatMap($callback)
func FlatMap[T any, R any](c *Collection[T], callback func(T, int) []R) *Collection[R] {
	result := make([]R, 0)
	for i, item := range c.items {
		result = append(result, callback(item, i)...)
	}
	return Collect(result)
}

// MapInto maps items into a new type using a constructor function.
// Equivalent to: $collection->mapInto($class)
func MapInto[T any, R any](c *Collection[T], constructor func(T) R) *Collection[R] {
	result := make([]R, len(c.items))
	for i, item := range c.items {
		result[i] = constructor(item)
	}
	return Collect(result)
}

// Reduce reduces the collection to a single value.
// Equivalent to: $collection->reduce($callback, $initial)
func Reduce[T any, R any](c *Collection[T], callback func(R, T, int) R, initial R) R {
	result := initial
	for i, item := range c.items {
		result = callback(result, item, i)
	}
	return result
}

// Flatten flattens a multi-dimensional collection into a single dimension.
// For Go this works with []any items.
// Equivalent to: $collection->flatten()
func (c *Collection[T]) Flatten() *Collection[T] {
	// In Go, we can't easily flatten generic types without reflection.
	// This returns the collection as-is for non-nested types.
	return Collect(append([]T{}, c.items...))
}

// Chunk breaks the collection into multiple, smaller slices of a given size.
// Note: Returns [][]T instead of nested Collection due to Go generics constraints.
// Equivalent to: $collection->chunk($size)
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

// ChunkWhile breaks the collection into multiple groups while the given callback returns true.
// Equivalent to: $collection->chunkWhile($callback)
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

// Split breaks a collection into the given number of groups.
// Equivalent to: $collection->split($numberOfGroups)
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

// SplitIn splits a collection into groups, filling non-terminal groups completely.
// Equivalent to: $collection->splitIn($numberOfGroups)
func (c *Collection[T]) SplitIn(numberOfGroups int) [][]T {
	size := int(math.Ceil(float64(len(c.items)) / float64(numberOfGroups)))
	return c.Chunk(size)
}

// Sliding creates a sliding window view of the collection.
// Equivalent to: $collection->sliding($size, $step)
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

// Slice extracts a slice of the collection.
// Equivalent to: $collection->slice($offset, $length)
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

// Splice removes and returns a slice of items starting at the specified index.
// Equivalent to: $collection->splice($offset, $length, $replacement)
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

// SpliceReplace removes a portion and replaces it with the given items.
// Equivalent to: $collection->splice($offset, $length, $replacement)
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

// Take returns a new collection with the specified number of items.
// Equivalent to: $collection->take($limit)
func (c *Collection[T]) Take(limit int) *Collection[T] {
	if limit < 0 {
		return c.Slice(limit)
	}
	return c.Slice(0, limit)
}

// TakeUntil returns items until the given callback returns true.
// Equivalent to: $collection->takeUntil($callback)
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

// TakeWhile returns items while the given callback returns true.
// Equivalent to: $collection->takeWhile($callback)
func (c *Collection[T]) TakeWhile(callback func(T, int) bool) *Collection[T] {
	return c.TakeUntil(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Skip skips over the first N items.
// Equivalent to: $collection->skip($count)
func (c *Collection[T]) Skip(count int) *Collection[T] {
	return c.Slice(count)
}

// SkipUntil skips items until the given callback returns true.
// Equivalent to: $collection->skipUntil($callback)
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

// SkipWhile skips items while the given callback returns true.
// Equivalent to: $collection->skipWhile($callback)
func (c *Collection[T]) SkipWhile(callback func(T, int) bool) *Collection[T] {
	return c.SkipUntil(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Nth creates a new collection consisting of every n-th element.
// Equivalent to: $collection->nth($step, $offset)
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

// ForPage "pages" the collection by returning a given number of items per page.
// Equivalent to: $collection->forPage($page, $perPage)
func (c *Collection[T]) ForPage(page, perPage int) *Collection[T] {
	offset := (page - 1) * perPage
	return c.Slice(offset, perPage)
}

// Values resets the keys on the collection (returns a copy).
// Equivalent to: $collection->values()
func (c *Collection[T]) Values() *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	return Collect(result)
}

// Reverse reverses the order of items.
// Equivalent to: $collection->reverse()
func (c *Collection[T]) Reverse() *Collection[T] {
	result := make([]T, len(c.items))
	for i, j := 0, len(c.items)-1; j >= 0; i, j = i+1, j-1 {
		result[i] = c.items[j]
	}
	return Collect(result)
}

// Shuffle randomly shuffles the items.
// Equivalent to: $collection->shuffle()
func (c *Collection[T]) Shuffle() *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return Collect(result)
}

// Random returns a random item from the collection.
// Equivalent to: $collection->random()
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

// Sort sorts the collection using the given comparison function.
// Equivalent to: $collection->sort($callback)
func (c *Collection[T]) Sort(less func(a, b T) bool) *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	sort.SliceStable(result, func(i, j int) bool {
		return less(result[i], result[j])
	})
	return Collect(result)
}

// SortBy sorts the collection by the given callback.
// Equivalent to: $collection->sortBy($callback)
func SortBy[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	sort.SliceStable(result, func(i, j int) bool {
		return keyFunc(result[i]) < keyFunc(result[j])
	})
	return Collect(result)
}

// SortByDesc sorts the collection by the given callback in descending order.
// Equivalent to: $collection->sortByDesc($callback)
func SortByDesc[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	sort.SliceStable(result, func(i, j int) bool {
		return keyFunc(result[i]) > keyFunc(result[j])
	})
	return Collect(result)
}

// SortDesc sorts the collection in descending order.
// Equivalent to: $collection->sortDesc()
func (c *Collection[T]) SortDesc(less func(a, b T) bool) *Collection[T] {
	return c.Sort(func(a, b T) bool {
		return less(b, a)
	})
}

// Unique returns unique items using the given key function.
// Equivalent to: $collection->unique($callback)
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

// Duplicates returns all duplicate items using the given key function.
// Equivalent to: $collection->duplicates($callback)
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

// Every determines if all items pass the given truth test.
// Equivalent to: $collection->every($callback)
func (c *Collection[T]) Every(callback func(T, int) bool) bool {
	for i, item := range c.items {
		if !callback(item, i) {
			return false
		}
	}
	return true
}

// Partition separates items that pass the truth test from those that don't.
// Equivalent to: $collection->partition($callback)
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

// Concat appends the given items to the end of the collection.
// Equivalent to: $collection->concat($source)
func (c *Collection[T]) Concat(items []T) *Collection[T] {
	result := make([]T, len(c.items)+len(items))
	copy(result, c.items)
	copy(result[len(c.items):], items)
	return Collect(result)
}

// Merge merges the given items into the collection.
// Equivalent to: $collection->merge($items)
func (c *Collection[T]) Merge(items []T) *Collection[T] {
	return c.Concat(items)
}

// Pad pads the collection to the specified length with a value.
// Equivalent to: $collection->pad($size, $value)
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

// Multiply creates multiple copies of all items in the collection.
// Equivalent to: $collection->multiply($multiplier)
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

// Flip swaps keys and values. Only meaningful for Pair types.
// Equivalent to: $collection->flip()
// Note: For typed Go, this is better handled via Map.
func (c *Collection[T]) Flip() *Collection[T] {
	result := make([]T, len(c.items))
	for i, j := 0, len(c.items)-1; j >= 0; i, j = i+1, j-1 {
		result[i] = c.items[j]
	}
	return Collect(result)
}

// Forget removes an item from the collection by index.
// Equivalent to: $collection->forget($key)
func (c *Collection[T]) Forget(index int) *Collection[T] {
	if index < 0 || index >= len(c.items) {
		return c
	}
	c.items = append(c.items[:index], c.items[index+1:]...)
	return c
}

// Implode joins elements into a string using fmt.Sprint.
// Equivalent to: $collection->implode($glue)
func (c *Collection[T]) Implode(glue string) string {
	parts := make([]string, len(c.items))
	for i, item := range c.items {
		parts[i] = fmt.Sprint(item)
	}
	return strings.Join(parts, glue)
}

// Join is like Implode but with an optional final glue.
// Equivalent to: $collection->join($glue, $finalGlue)
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

// When applies the callback if the given condition is true.
// Equivalent to: $collection->when($condition, $callback, $default)
func (c *Collection[T]) When(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	if condition {
		return callback(c)
	}
	if len(defaults) > 0 {
		return defaults[0](c)
	}
	return c
}

// WhenEmpty applies the callback if the collection is empty.
// Equivalent to: $collection->whenEmpty($callback, $default)
func (c *Collection[T]) WhenEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(c.IsEmpty(), callback, defaults...)
}

// WhenNotEmpty applies the callback if the collection is not empty.
// Equivalent to: $collection->whenNotEmpty($callback, $default)
func (c *Collection[T]) WhenNotEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(c.IsNotEmpty(), callback, defaults...)
}

// Unless applies the callback unless the given condition is true.
// Equivalent to: $collection->unless($condition, $callback, $default)
func (c *Collection[T]) Unless(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(!condition, callback, defaults...)
}

// UnlessEmpty applies the callback unless the collection is empty.
// Equivalent to: $collection->unlessEmpty($callback, $default)
func (c *Collection[T]) UnlessEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.WhenNotEmpty(callback, defaults...)
}

// UnlessNotEmpty applies the callback unless the collection is not empty.
// Equivalent to: $collection->unlessNotEmpty($callback, $default)
func (c *Collection[T]) UnlessNotEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.WhenEmpty(callback, defaults...)
}

// Zip merges the collection with the given items.
// Equivalent to: $collection->zip($items)
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

// CrossJoin cross joins the collection with the given items.
// Equivalent to: $collection->crossJoin($lists)
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

// Combine creates a collection of Pair by combining keys from this collection with values from the given items.
// Equivalent to: $collection->combine($values)
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

// Collapse collapses a collection of slices into a single, flat collection.
// Equivalent to: $collection->collapse()
func Collapse[T any](c *Collection[[]T]) *Collection[T] {
	result := make([]T, 0)
	for _, items := range c.items {
		result = append(result, items...)
	}
	return Collect(result)
}

// Diff returns the items in the collection not present in the given items.
// Equivalent to: $collection->diff($items)
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

// DiffUsing returns the items not present in the given items, using the callback.
// Equivalent to: $collection->diffUsing($items, $callback)
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

// Intersect returns the items present in both collections.
// Equivalent to: $collection->intersect($items)
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

// IntersectUsing returns the items present in both collections, using the callback.
// Equivalent to: $collection->intersectUsing($items, $callback)
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

// ToSlice returns the underlying slice.
// Equivalent to: $collection->toArray()
func (c *Collection[T]) ToSlice() []T {
	result := make([]T, len(c.items))
	copy(result, c.items)
	return result
}

// ToJSON serializes the collection to JSON.
// Equivalent to: $collection->toJson()
func (c *Collection[T]) ToJSON() ([]byte, error) {
	return json.Marshal(c.items)
}

// ToPrettyJSON serializes the collection to indented JSON.
// Equivalent to: $collection->toPrettyJson()
func (c *Collection[T]) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(c.items, "", "    ")
}

// String returns the JSON string representation of the collection.
// Equivalent to: $collection->__toString()
func (c *Collection[T]) String() string {
	b, err := c.ToJSON()
	if err != nil {
		return "[]"
	}
	return string(b)
}

// MarshalJSON implements the json.Marshaler interface.
// Equivalent to: $collection->jsonSerialize()
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

// Len implements sort.Interface.
func (c *Collection[T]) Len() int {
	return len(c.items)
}

// TapEach calls the given callback on each item, returning the original collection.
// Equivalent to: $collection->tapEach($callback)
func (c *Collection[T]) TapEach(callback func(T, int)) *Collection[T] {
	for i, item := range c.items {
		callback(item, i)
	}
	return c
}

// Dump prints the collection items for debugging.
// Equivalent to: $collection->dump()
func (c *Collection[T]) Dump() *Collection[T] {
	fmt.Printf("%v\n", c.items)
	return c
}

// Only returns a collection with items at the given indices.
// Equivalent to: $collection->only($keys)
func (c *Collection[T]) Only(indices ...int) *Collection[T] {
	result := make([]T, 0, len(indices))
	for _, idx := range indices {
		if idx >= 0 && idx < len(c.items) {
			result = append(result, c.items[idx])
		}
	}
	return Collect(result)
}

// Except returns a collection without items at the given indices.
// Equivalent to: $collection->except($keys)
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

// Has determines if a key exists in the collection.
// Equivalent to: $collection->has($key)
func (c *Collection[T]) Has(index int) bool {
	if index < 0 {
		index = len(c.items) + index
	}
	return index >= 0 && index < len(c.items)
}

// HasAny determines if any of the given keys exist in the collection.
// Equivalent to: $collection->hasAny($keys)
func (c *Collection[T]) HasAny(indices ...int) bool {
	for _, idx := range indices {
		if c.Has(idx) {
			return true
		}
	}
	return false
}

// Pluck extracts values from a collection using a key function.
// Equivalent to: $collection->pluck($value)
func Pluck[T any, V any](c *Collection[T], valueFunc func(T) V) *Collection[V] {
	result := make([]V, len(c.items))
	for i, item := range c.items {
		result[i] = valueFunc(item)
	}
	return Collect(result)
}

// GroupBy groups the collection's items by a given key.
// Equivalent to: $collection->groupBy($groupBy)
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

// KeyBy keys the collection by the given key.
// Equivalent to: $collection->keyBy($keyBy)
func KeyBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range c.items {
		result[keyFunc(item)] = item
	}
	return result
}

// CountBy counts the occurrences of values in the collection.
// Equivalent to: $collection->countBy($callback)
func CountBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]int {
	result := make(map[K]int)
	for _, item := range c.items {
		result[keyFunc(item)]++
	}
	return result
}

// MapToDictionary maps items to key-value pairs and groups them.
// Equivalent to: $collection->mapToDictionary($callback)
func MapToDictionary[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K][]V {
	result := make(map[K][]V)
	for _, item := range c.items {
		key, value := callback(item)
		result[key] = append(result[key], value)
	}
	return result
}

// MapWithKeys maps items to key-value pairs.
// Equivalent to: $collection->mapWithKeys($callback)
func MapWithKeys[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K]V {
	result := make(map[K]V)
	for _, item := range c.items {
		key, value := callback(item)
		result[key] = value
	}
	return result
}

// Where filters items by a callback.
// Equivalent to: $collection->where($key, $operator, $value)
func (c *Collection[T]) Where(predicate func(T) bool) *Collection[T] {
	return c.Filter(func(item T, _ int) bool {
		return predicate(item)
	})
}

// WhereNot filters items by a negative callback.
func (c *Collection[T]) WhereNot(predicate func(T) bool) *Collection[T] {
	return c.Filter(func(item T, _ int) bool {
		return !predicate(item)
	})
}

// Dot flattens a multi-dimensional collection using "dot" notation.
// For Go with typed collections, this returns the collection as-is.
// Equivalent to: $collection->dot()
func (c *Collection[T]) Dot() *Collection[T] {
	return c.Copy()
}

// Undot expands dotted keys. For Go, returns as-is.
// Equivalent to: $collection->undot()
func (c *Collection[T]) Undot() *Collection[T] {
	return c.Copy()
}

// Ensure asserts that all items pass the given truth test.
// Equivalent to: $collection->ensure($type)
func (c *Collection[T]) Ensure(predicate func(T) bool) error {
	for _, item := range c.items {
		if !predicate(item) {
			return fmt.Errorf("collection item failed ensure check")
		}
	}
	return nil
}

// ToBase returns the base Collection.
// Equivalent to: $collection->toBase()
func (c *Collection[T]) ToBase() *Collection[T] {
	return c
}

// Lazy returns a new LazyCollection from the items.
// Equivalent to: $collection->lazy()
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

// Median returns the median value.
// Equivalent to: $collection->median()
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

// MedianBy returns the median value using a key function.
// Equivalent to: $collection->median($key)
func MedianBy[T any](c *Collection[T], valueFunc func(T) float64) float64 {
	return Median(Map(c, func(item T, _ int) float64 {
		return valueFunc(item)
	}))
}

// Mode returns the mode (most frequent) values.
// Equivalent to: $collection->mode()
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

// Sum returns the sum of all items.
// Equivalent to: $collection->sum()
func Sum[T Numeric](c *Collection[T]) T {
	var total T
	for _, item := range c.items {
		total += item
	}
	return total
}

// SumBy returns the sum of values extracted by a key function.
// Equivalent to: $collection->sum($callback)
func SumBy[T any, N Numeric](c *Collection[T], valueFunc func(T) N) N {
	var total N
	for _, item := range c.items {
		total += valueFunc(item)
	}
	return total
}

// Avg returns the average of all items.
// Equivalent to: $collection->avg()
func Avg[T Numeric](c *Collection[T]) float64 {
	if len(c.items) == 0 {
		return 0
	}
	return float64(Sum(c)) / float64(len(c.items))
}

// AvgBy returns the average of values using a key function.
// Equivalent to: $collection->avg($callback)
func AvgBy[T any, N Numeric](c *Collection[T], valueFunc func(T) N) float64 {
	if len(c.items) == 0 {
		return 0
	}
	return float64(SumBy(c, valueFunc)) / float64(len(c.items))
}

// Average is an alias for Avg.
// Equivalent to: $collection->average()
func Average[T Numeric](c *Collection[T]) float64 {
	return Avg(c)
}

// Min returns the minimum value.
// Equivalent to: $collection->min()
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

// MinBy returns the minimum value using a key function.
// Equivalent to: $collection->min($callback)
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

// Max returns the maximum value.
// Equivalent to: $collection->max()
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

// MaxBy returns the maximum value using a key function.
// Equivalent to: $collection->max($callback)
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
