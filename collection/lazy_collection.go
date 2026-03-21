package collection

import (
	"fmt"
	"iter"
	"strings"
	"time"
)

// LazyCollection provides lazy evaluation for collection operations.
// Items are only computed when needed, making it memory-efficient for large datasets.
type LazyCollection[T any] struct {
	source iter.Seq[T]
}

// NewLazy creates a new LazyCollection from an iterator function.
// The source receives a yield function; call yield for each item and
// stop generating if yield returns false.
func NewLazy[T any](source iter.Seq[T]) *LazyCollection[T] {
	return &LazyCollection[T]{source: source}
}

// LazyFrom creates a LazyCollection from a slice.
func LazyFrom[T any](items []T) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	})
}

// LazyEmpty creates an empty LazyCollection.
func LazyEmpty[T any]() *LazyCollection[T] {
	return NewLazy[T](func(yield func(T) bool) {})
}

// LazyRange creates a lazy collection of sequential integers from start to end (inclusive).
// If from > to, the sequence counts downward.
func LazyRange(from, to int) *LazyCollection[int] {
	return NewLazy(func(yield func(int) bool) {
		if from <= to {
			for i := from; i <= to; i++ {
				if !yield(i) {
					return
				}
			}
		} else {
			for i := from; i >= to; i-- {
				if !yield(i) {
					return
				}
			}
		}
	})
}

// LazyTimes creates a lazy collection by invoking the callback n times.
// The callback receives a 1-based index.
func LazyTimes[T any](number int, callback func(int) T) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		for i := 1; i <= number; i++ {
			if !yield(callback(i)) {
				return
			}
		}
	})
}

// Iter returns the underlying iterator for use with range loops.
func (lc *LazyCollection[T]) Iter() iter.Seq[T] {
	return lc.source
}

// All eagerly evaluates the lazy collection and returns all items as a slice.
func (lc *LazyCollection[T]) All() []T {
	result := make([]T, 0)
	lc.source(func(item T) bool {
		result = append(result, item)
		return true
	})
	return result
}

// Eager eagerly evaluates the lazy collection and returns a Collection.
func (lc *LazyCollection[T]) Eager() *Collection[T] {
	return Collect(lc.All())
}

// Collect converts the lazy collection to an eager Collection.
func (lc *LazyCollection[T]) Collect() *Collection[T] {
	return lc.Eager()
}

// Count returns the total number of items by eagerly evaluating the collection.
func (lc *LazyCollection[T]) Count() int {
	count := 0
	lc.source(func(_ T) bool {
		count++
		return true
	})
	return count
}

// IsEmpty reports whether the lazy collection contains no items.
func (lc *LazyCollection[T]) IsEmpty() bool {
	empty := true
	lc.source(func(_ T) bool {
		empty = false
		return false
	})
	return empty
}

// IsNotEmpty reports whether the lazy collection contains at least one item.
func (lc *LazyCollection[T]) IsNotEmpty() bool {
	return !lc.IsEmpty()
}

// ContainsOneItem reports whether the lazy collection contains exactly one item.
func (lc *LazyCollection[T]) ContainsOneItem() bool {
	return lc.Count() == 1
}

// ContainsManyItems reports whether the lazy collection contains more than one item.
func (lc *LazyCollection[T]) ContainsManyItems() bool {
	return lc.Count() > 1
}

// First returns the first element matching the optional predicate.
// If no predicate is given, it returns the first element.
// The second return value indicates whether a matching element was found.
func (lc *LazyCollection[T]) First(predicates ...func(T, int) bool) (T, bool) {
	var result T
	found := false
	idx := 0
	if len(predicates) == 0 || predicates[0] == nil {
		lc.source(func(item T) bool {
			result = item
			found = true
			return false
		})
	} else {
		predicate := predicates[0]
		lc.source(func(item T) bool {
			if predicate(item, idx) {
				result = item
				found = true
				return false
			}
			idx++
			return true
		})
	}
	return result, found
}

// FirstOrFail returns the first element matching the optional predicate,
// or an ItemNotFoundError if no element is found.
func (lc *LazyCollection[T]) FirstOrFail(predicates ...func(T, int) bool) (T, error) {
	item, ok := lc.First(predicates...)
	if !ok {
		var zero T
		return zero, &ItemNotFoundError{}
	}
	return item, nil
}

// Last returns the last element matching the optional predicate.
// The second return value indicates whether a matching element was found.
func (lc *LazyCollection[T]) Last(predicates ...func(T, int) bool) (T, bool) {
	var result T
	found := false
	idx := 0
	if len(predicates) == 0 || predicates[0] == nil {
		lc.source(func(item T) bool {
			result = item
			found = true
			return true
		})
	} else {
		predicate := predicates[0]
		lc.source(func(item T) bool {
			if predicate(item, idx) {
				result = item
				found = true
			}
			idx++
			return true
		})
	}
	return result, found
}

// Sole returns the only element matching the optional predicate.
// It returns an error if zero or more than one element matches.
func (lc *LazyCollection[T]) Sole(predicates ...func(T, int) bool) (T, error) {
	return lc.Eager().Sole(predicates...)
}

// Get returns the item at the given zero-based index.
// The second return value indicates whether the index exists.
func (lc *LazyCollection[T]) Get(index int) (T, bool) {
	var result T
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if idx == index {
			result = item
			found = true
			return false
		}
		idx++
		return true
	})
	return result, found
}

// Contains reports whether any item satisfies the predicate.
func (lc *LazyCollection[T]) Contains(predicate func(T, int) bool) bool {
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if predicate(item, idx) {
			found = true
			return false
		}
		idx++
		return true
	})
	return found
}

// Some is an alias for Contains.
func (lc *LazyCollection[T]) Some(predicate func(T, int) bool) bool {
	return lc.Contains(predicate)
}

// DoesntContain reports whether no item satisfies the predicate.
func (lc *LazyCollection[T]) DoesntContain(predicate func(T, int) bool) bool {
	return !lc.Contains(predicate)
}

// Search returns the index of the first item satisfying the predicate.
// The second return value indicates whether a match was found.
func (lc *LazyCollection[T]) Search(predicate func(T, int) bool) (int, bool) {
	resultIdx := -1
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if predicate(item, idx) {
			resultIdx = idx
			found = true
			return false
		}
		idx++
		return true
	})
	return resultIdx, found
}

// Before returns the item immediately before the first item matching the predicate.
// The second return value indicates whether such an item exists.
func (lc *LazyCollection[T]) Before(predicate func(T, int) bool) (T, bool) {
	var prev T
	hasPrev := false
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if predicate(item, idx) {
			found = true
			return false
		}
		prev = item
		hasPrev = true
		idx++
		return true
	})
	if found && hasPrev {
		return prev, true
	}
	var zero T
	return zero, false
}

// After returns the item immediately after the first item matching the predicate.
// The second return value indicates whether such an item exists.
func (lc *LazyCollection[T]) After(predicate func(T, int) bool) (T, bool) {
	matched := false
	var result T
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if matched {
			result = item
			found = true
			return false
		}
		if predicate(item, idx) {
			matched = true
		}
		idx++
		return true
	})
	return result, found
}

// Each iterates over items, calling the callback for each one.
// If the callback returns false, iteration stops early.
func (lc *LazyCollection[T]) Each(callback func(T, int) bool) *LazyCollection[T] {
	idx := 0
	lc.source(func(item T) bool {
		cont := callback(item, idx)
		idx++
		return cont
	})
	return lc
}

// Tap passes the lazy collection to the given callback and returns it unchanged.
func (lc *LazyCollection[T]) Tap(callback func(*LazyCollection[T])) *LazyCollection[T] {
	callback(lc)
	return lc
}

// Filter returns a new lazy collection containing only items for which the callback returns true.
func (lc *LazyCollection[T]) Filter(callback func(T, int) bool) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		idx := 0
		lc.source(func(item T) bool {
			if callback(item, idx) {
				if !yield(item) {
					return false
				}
			}
			idx++
			return true
		})
	})
}

// Reject returns a new lazy collection excluding items for which the callback returns true.
func (lc *LazyCollection[T]) Reject(callback func(T, int) bool) *LazyCollection[T] {
	return lc.Filter(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// LazyMap transforms each item using the callback, returning a new lazy collection.
func LazyMap[T any, R any](lc *LazyCollection[T], callback func(T, int) R) *LazyCollection[R] {
	return NewLazy(func(yield func(R) bool) {
		idx := 0
		lc.source(func(item T) bool {
			if !yield(callback(item, idx)) {
				return false
			}
			idx++
			return true
		})
	})
}

// LazyFlatMap transforms each item into a slice and flattens the results into a new lazy collection.
func LazyFlatMap[T any, R any](lc *LazyCollection[T], callback func(T, int) []R) *LazyCollection[R] {
	return NewLazy(func(yield func(R) bool) {
		idx := 0
		lc.source(func(item T) bool {
			for _, r := range callback(item, idx) {
				if !yield(r) {
					return false
				}
			}
			idx++
			return true
		})
	})
}

// Take returns a new lazy collection with at most limit items.
// A negative limit takes from the end (requires eager evaluation).
func (lc *LazyCollection[T]) Take(limit int) *LazyCollection[T] {
	if limit < 0 {
		// For negative take, we need to eagerly evaluate
		items := lc.All()
		start := len(items) + limit
		if start < 0 {
			start = 0
		}
		return LazyFrom(items[start:])
	}
	return NewLazy(func(yield func(T) bool) {
		count := 0
		lc.source(func(item T) bool {
			if count >= limit {
				return false
			}
			count++
			return yield(item)
		})
	})
}

// TakeUntil returns items from the beginning until the callback returns true.
func (lc *LazyCollection[T]) TakeUntil(callback func(T, int) bool) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		idx := 0
		lc.source(func(item T) bool {
			if callback(item, idx) {
				return false
			}
			idx++
			return yield(item)
		})
	})
}

// TakeWhile returns items from the beginning while the callback returns true.
func (lc *LazyCollection[T]) TakeWhile(callback func(T, int) bool) *LazyCollection[T] {
	return lc.TakeUntil(func(item T, idx int) bool {
		return !callback(item, idx)
	})
}

// TakeUntilTimeout returns items until the given duration has elapsed.
func (lc *LazyCollection[T]) TakeUntilTimeout(timeout time.Duration) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		deadline := time.Now().Add(timeout)
		lc.source(func(item T) bool {
			if time.Now().After(deadline) {
				return false
			}
			return yield(item)
		})
	})
}

// Skip returns a new lazy collection that skips the first count items.
func (lc *LazyCollection[T]) Skip(count int) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		skipped := 0
		lc.source(func(item T) bool {
			if skipped < count {
				skipped++
				return true
			}
			return yield(item)
		})
	})
}

// SkipUntil skips items until the callback returns true, then yields the rest.
func (lc *LazyCollection[T]) SkipUntil(callback func(T, int) bool) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		found := false
		idx := 0
		lc.source(func(item T) bool {
			if !found && callback(item, idx) {
				found = true
			}
			if found {
				return yield(item)
			}
			idx++
			return true
		})
	})
}

// SkipWhile skips items while the callback returns true, then yields the rest.
func (lc *LazyCollection[T]) SkipWhile(callback func(T, int) bool) *LazyCollection[T] {
	return lc.SkipUntil(func(item T, idx int) bool {
		return !callback(item, idx)
	})
}

// Slice returns a subset of the lazy collection starting at offset with an optional length.
func (lc *LazyCollection[T]) Slice(offset int, lengths ...int) *LazyCollection[T] {
	result := lc.Skip(offset)
	if len(lengths) > 0 {
		result = result.Take(lengths[0])
	}
	return result
}

// Chunk eagerly evaluates the lazy collection and splits items into groups of the given size.
func (lc *LazyCollection[T]) Chunk(size int) [][]T {
	var result [][]T
	chunk := make([]T, 0, size)
	lc.source(func(item T) bool {
		chunk = append(chunk, item)
		if len(chunk) >= size {
			result = append(result, chunk)
			chunk = make([]T, 0, size)
		}
		return true
	})
	if len(chunk) > 0 {
		result = append(result, chunk)
	}
	return result
}

// ChunkWhile splits the lazy collection into groups where consecutive items satisfy the callback.
// A new chunk is started when the callback returns false.
func (lc *LazyCollection[T]) ChunkWhile(callback func(T, int, []T) bool) [][]T {
	var result [][]T
	current := make([]T, 0)
	idx := 0
	lc.source(func(item T) bool {
		if len(current) == 0 || callback(item, idx, current) {
			current = append(current, item)
		} else {
			result = append(result, current)
			current = []T{item}
		}
		idx++
		return true
	})
	if len(current) > 0 {
		result = append(result, current)
	}
	return result
}

// Nth returns every step-th element, starting from an optional offset.
func (lc *LazyCollection[T]) Nth(step int, offsets ...int) *LazyCollection[T] {
	offset := 0
	if len(offsets) > 0 {
		offset = offsets[0]
	}
	return NewLazy(func(yield func(T) bool) {
		idx := 0
		lc.source(func(item T) bool {
			if idx >= offset && (idx-offset)%step == 0 {
				if !yield(item) {
					return false
				}
			}
			idx++
			return true
		})
	})
}

// Concat returns a new lazy collection with the given items appended.
func (lc *LazyCollection[T]) Concat(items []T) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		lc.source(func(item T) bool {
			return yield(item)
		})
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	})
}

// Pad returns a new lazy collection padded to the given size with the specified value.
// A negative size pads at the beginning; a positive size pads at the end.
func (lc *LazyCollection[T]) Pad(size int, value T) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		count := 0
		absSize := size
		if absSize < 0 {
			absSize = -absSize
		}
		if size < 0 {
			// Evaluate to know count
			items := lc.All()
			padCount := absSize - len(items)
			if padCount > 0 {
				for i := 0; i < padCount; i++ {
					if !yield(value) {
						return
					}
				}
			}
			for _, item := range items {
				if !yield(item) {
					return
				}
			}
		} else {
			lc.source(func(item T) bool {
				count++
				return yield(item)
			})
			for count < size {
				if !yield(value) {
					return
				}
				count++
			}
		}
	})
}

// TapEach returns a new lazy collection that calls the callback on each item as it passes through.
func (lc *LazyCollection[T]) TapEach(callback func(T, int)) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		idx := 0
		lc.source(func(item T) bool {
			callback(item, idx)
			idx++
			return yield(item)
		})
	})
}

// Throttle returns a new lazy collection that inserts a delay between each yielded item.
func (lc *LazyCollection[T]) Throttle(delay time.Duration) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		first := true
		lc.source(func(item T) bool {
			if !first {
				time.Sleep(delay)
			}
			first = false
			return yield(item)
		})
	})
}

// Remember returns a lazy collection that caches items on first iteration,
// so subsequent iterations reuse the cached values.
func (lc *LazyCollection[T]) Remember() *LazyCollection[T] {
	var cache []T
	cached := false
	return NewLazy(func(yield func(T) bool) {
		if cached {
			for _, item := range cache {
				if !yield(item) {
					return
				}
			}
			return
		}
		cache = make([]T, 0)
		lc.source(func(item T) bool {
			cache = append(cache, item)
			return yield(item)
		})
		cached = true
	})
}

// Every reports whether all items satisfy the callback.
func (lc *LazyCollection[T]) Every(callback func(T, int) bool) bool {
	result := true
	idx := 0
	lc.source(func(item T) bool {
		if !callback(item, idx) {
			result = false
			return false
		}
		idx++
		return true
	})
	return result
}

// Has reports whether the given zero-based index exists in the collection.
func (lc *LazyCollection[T]) Has(index int) bool {
	_, ok := lc.Get(index)
	return ok
}

// HasAny reports whether any of the given indices exist in the collection.
func (lc *LazyCollection[T]) HasAny(indices ...int) bool {
	for _, idx := range indices {
		if lc.Has(idx) {
			return true
		}
	}
	return false
}

// HasSole reports whether exactly one item matches the optional predicate.
// If no predicate is given, it checks whether the collection has exactly one item.
func (lc *LazyCollection[T]) HasSole(predicates ...func(T, int) bool) bool {
	count := 0
	idx := 0
	if len(predicates) == 0 || predicates[0] == nil {
		lc.source(func(_ T) bool {
			count++
			return count < 2
		})
	} else {
		predicate := predicates[0]
		lc.source(func(item T) bool {
			if predicate(item, idx) {
				count++
				if count > 1 {
					return false
				}
			}
			idx++
			return true
		})
	}
	return count == 1
}

// Implode joins all items into a single string separated by the given glue.
func (lc *LazyCollection[T]) Implode(glue string) string {
	parts := make([]string, 0)
	lc.source(func(item T) bool {
		parts = append(parts, fmt.Sprint(item))
		return true
	})
	return strings.Join(parts, glue)
}

// Join joins all items into a string with the given separator.
// An optional final separator can be provided for the last element.
func (lc *LazyCollection[T]) Join(glue string, finalGlues ...string) string {
	parts := make([]string, 0)
	lc.source(func(item T) bool {
		parts = append(parts, fmt.Sprint(item))
		return true
	})
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

// When applies the callback if the condition is true; otherwise applies the optional default.
func (lc *LazyCollection[T]) When(condition bool, callback func(*LazyCollection[T]) *LazyCollection[T], defaults ...func(*LazyCollection[T]) *LazyCollection[T]) *LazyCollection[T] {
	if condition {
		return callback(lc)
	}
	if len(defaults) > 0 {
		return defaults[0](lc)
	}
	return lc
}

// WhenEmpty applies the callback if the lazy collection is empty.
func (lc *LazyCollection[T]) WhenEmpty(callback func(*LazyCollection[T]) *LazyCollection[T], defaults ...func(*LazyCollection[T]) *LazyCollection[T]) *LazyCollection[T] {
	return lc.When(lc.IsEmpty(), callback, defaults...)
}

// WhenNotEmpty applies the callback if the lazy collection is not empty.
func (lc *LazyCollection[T]) WhenNotEmpty(callback func(*LazyCollection[T]) *LazyCollection[T], defaults ...func(*LazyCollection[T]) *LazyCollection[T]) *LazyCollection[T] {
	return lc.When(lc.IsNotEmpty(), callback, defaults...)
}

// Unless applies the callback unless the condition is true.
func (lc *LazyCollection[T]) Unless(condition bool, callback func(*LazyCollection[T]) *LazyCollection[T], defaults ...func(*LazyCollection[T]) *LazyCollection[T]) *LazyCollection[T] {
	return lc.When(!condition, callback, defaults...)
}

// LazyReduce reduces the lazy collection to a single value by applying the callback
// to an accumulator and each item in sequence.
func LazyReduce[T any, R any](lc *LazyCollection[T], callback func(R, T, int) R, initial R) R {
	result := initial
	idx := 0
	lc.source(func(item T) bool {
		result = callback(result, item, idx)
		idx++
		return true
	})
	return result
}

// LazyUnique returns a lazy collection containing only unique items as determined by the key function.
func LazyUnique[T any, K comparable](lc *LazyCollection[T], keyFunc func(T) K) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		seen := make(map[K]bool)
		lc.source(func(item T) bool {
			key := keyFunc(item)
			if !seen[key] {
				seen[key] = true
				return yield(item)
			}
			return true
		})
	})
}

// LazyPluck extracts a value from each item using the given function, returning a new lazy collection.
func LazyPluck[T any, V any](lc *LazyCollection[T], valueFunc func(T) V) *LazyCollection[V] {
	return LazyMap(lc, func(item T, _ int) V {
		return valueFunc(item)
	})
}

// LazyGroupBy groups items by the key returned by the given function.
// This requires eager evaluation to build the groups.
func LazyGroupBy[T any, K comparable](lc *LazyCollection[T], keyFunc func(T) K) map[K]*LazyCollection[T] {
	groups := make(map[K][]T)
	lc.source(func(item T) bool {
		key := keyFunc(item)
		groups[key] = append(groups[key], item)
		return true
	})
	result := make(map[K]*LazyCollection[T])
	for key, items := range groups {
		result[key] = LazyFrom(items)
	}
	return result
}

// LazyKeyBy indexes items by the key returned by the given function.
// Duplicate keys cause the later value to overwrite the earlier one.
func LazyKeyBy[T any, K comparable](lc *LazyCollection[T], keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	lc.source(func(item T) bool {
		result[keyFunc(item)] = item
		return true
	})
	return result
}

// LazyCountBy counts occurrences of each key returned by the given function.
func LazyCountBy[T any, K comparable](lc *LazyCollection[T], keyFunc func(T) K) map[K]int {
	result := make(map[K]int)
	lc.source(func(item T) bool {
		result[keyFunc(item)]++
		return true
	})
	return result
}

// Dump prints the lazy collection items to stdout for debugging and returns a new
// lazy collection backed by the evaluated items.
func (lc *LazyCollection[T]) Dump() *LazyCollection[T] {
	items := lc.All()
	fmt.Printf("%v\n", items)
	return LazyFrom(items)
}
