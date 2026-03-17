package collection

import (
	"fmt"
	"strings"
	"time"
)

// LazyCollection provides lazy evaluation for collection operations.
// Items are only computed when needed, making it memory-efficient for large datasets.
// Equivalent to: Illuminate\Support\LazyCollection
type LazyCollection[T any] struct {
	source func(yield func(T) bool)
}

// NewLazy creates a new LazyCollection from a generator function.
// The generator receives a yield function; call yield for each item and
// stop generating if yield returns false.
// Equivalent to: new LazyCollection($source)
func NewLazy[T any](source func(yield func(T) bool)) *LazyCollection[T] {
	return &LazyCollection[T]{source: source}
}

// LazyFrom creates a LazyCollection from a slice.
// Equivalent to: LazyCollection::make($items)
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
// Equivalent to: LazyCollection::empty()
func LazyEmpty[T any]() *LazyCollection[T] {
	return NewLazy[T](func(yield func(T) bool) {})
}

// LazyRange creates a lazy range of integers.
// Equivalent to: LazyCollection::range($from, $to)
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

// LazyTimes creates a lazy collection by invoking the callback N times.
// Equivalent to: LazyCollection::times($number, $callback)
func LazyTimes[T any](number int, callback func(int) T) *LazyCollection[T] {
	return NewLazy(func(yield func(T) bool) {
		for i := 1; i <= number; i++ {
			if !yield(callback(i)) {
				return
			}
		}
	})
}

// All eagerly evaluates the lazy collection and returns all items.
// Equivalent to: $lazy->all()
func (lc *LazyCollection[T]) All() []T {
	result := make([]T, 0)
	lc.source(func(item T) bool {
		result = append(result, item)
		return true
	})
	return result
}

// Eager eagerly evaluates and returns a Collection.
// Equivalent to: $lazy->eager()
func (lc *LazyCollection[T]) Eager() *Collection[T] {
	return Collect(lc.All())
}

// Collect converts the lazy collection to an eager collection.
// Equivalent to: $lazy->collect()
func (lc *LazyCollection[T]) Collect() *Collection[T] {
	return lc.Eager()
}

// Count returns the total number of items (eagerly evaluates).
// Equivalent to: $lazy->count()
func (lc *LazyCollection[T]) Count() int {
	count := 0
	lc.source(func(_ T) bool {
		count++
		return true
	})
	return count
}

// IsEmpty determines if the lazy collection is empty.
// Equivalent to: $lazy->isEmpty()
func (lc *LazyCollection[T]) IsEmpty() bool {
	empty := true
	lc.source(func(_ T) bool {
		empty = false
		return false
	})
	return empty
}

// IsNotEmpty determines if the lazy collection is not empty.
// Equivalent to: $lazy->isNotEmpty()
func (lc *LazyCollection[T]) IsNotEmpty() bool {
	return !lc.IsEmpty()
}

// ContainsOneItem determines if the lazy collection contains exactly one item.
// Equivalent to: $lazy->containsOneItem()
func (lc *LazyCollection[T]) ContainsOneItem() bool {
	return lc.Count() == 1
}

// ContainsManyItems determines if the lazy collection contains more than one item.
// Equivalent to: $lazy->containsManyItems()
func (lc *LazyCollection[T]) ContainsManyItems() bool {
	return lc.Count() > 1
}

// First returns the first element matching the given predicate.
// Equivalent to: $lazy->first($callback)
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

// FirstOrFail returns the first element or an error.
// Equivalent to: $lazy->firstOrFail()
func (lc *LazyCollection[T]) FirstOrFail(predicates ...func(T, int) bool) (T, error) {
	item, ok := lc.First(predicates...)
	if !ok {
		var zero T
		return zero, &ItemNotFoundException{}
	}
	return item, nil
}

// Last returns the last element matching the given predicate.
// Equivalent to: $lazy->last($callback)
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

// Sole returns the only element matching the predicate.
// Equivalent to: $lazy->sole($callback)
func (lc *LazyCollection[T]) Sole(predicates ...func(T, int) bool) (T, error) {
	return lc.Eager().Sole(predicates...)
}

// Get returns the item at a given index.
// Equivalent to: $lazy->get($key)
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

// Contains determines if the lazy collection contains an item.
// Equivalent to: $lazy->contains($callback)
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

// DoesntContain determines if the lazy collection doesn't contain an item.
// Equivalent to: $lazy->doesntContain($callback)
func (lc *LazyCollection[T]) DoesntContain(predicate func(T, int) bool) bool {
	return !lc.Contains(predicate)
}

// Search searches the lazy collection for the given value.
// Equivalent to: $lazy->search($value)
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

// Before returns the item before the first matching item.
// Equivalent to: $lazy->before($value)
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

// After returns the item after the first matching item.
// Equivalent to: $lazy->after($value)
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

// Each iterates over the items.
// Equivalent to: $lazy->each($callback)
func (lc *LazyCollection[T]) Each(callback func(T, int) bool) *LazyCollection[T] {
	idx := 0
	lc.source(func(item T) bool {
		cont := callback(item, idx)
		idx++
		return cont
	})
	return lc
}

// Tap passes the lazy collection to the given callback.
// Equivalent to: $lazy->tap($callback)
func (lc *LazyCollection[T]) Tap(callback func(*LazyCollection[T])) *LazyCollection[T] {
	callback(lc)
	return lc
}

// Filter returns a new lazy collection with items passing the predicate.
// Equivalent to: $lazy->filter($callback)
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

// Reject returns items that don't pass the predicate.
// Equivalent to: $lazy->reject($callback)
func (lc *LazyCollection[T]) Reject(callback func(T, int) bool) *LazyCollection[T] {
	return lc.Filter(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// LazyMap maps items lazily.
// Equivalent to: $lazy->map($callback)
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

// LazyFlatMap flat maps items lazily.
// Equivalent to: $lazy->flatMap($callback)
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

// Take returns the first N items lazily.
// Equivalent to: $lazy->take($limit)
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

// TakeUntil returns items until the callback returns true.
// Equivalent to: $lazy->takeUntil($callback)
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

// TakeWhile returns items while the callback returns true.
// Equivalent to: $lazy->takeWhile($callback)
func (lc *LazyCollection[T]) TakeWhile(callback func(T, int) bool) *LazyCollection[T] {
	return lc.TakeUntil(func(item T, idx int) bool {
		return !callback(item, idx)
	})
}

// TakeUntilTimeout returns items until the given timeout.
// Equivalent to: $lazy->takeUntilTimeout($timeout)
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

// Skip skips the first N items.
// Equivalent to: $lazy->skip($count)
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

// SkipUntil skips items until the callback returns true.
// Equivalent to: $lazy->skipUntil($callback)
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

// SkipWhile skips items while the callback returns true.
// Equivalent to: $lazy->skipWhile($callback)
func (lc *LazyCollection[T]) SkipWhile(callback func(T, int) bool) *LazyCollection[T] {
	return lc.SkipUntil(func(item T, idx int) bool {
		return !callback(item, idx)
	})
}

// Slice extracts a lazy slice.
// Equivalent to: $lazy->slice($offset, $length)
func (lc *LazyCollection[T]) Slice(offset int, lengths ...int) *LazyCollection[T] {
	result := lc.Skip(offset)
	if len(lengths) > 0 {
		result = result.Take(lengths[0])
	}
	return result
}

// Chunk breaks the lazy collection into chunks and returns them eagerly as a slice of slices.
// Note: Go does not allow recursive generic instantiation, so this returns [][]T directly.
// Equivalent to: $lazy->chunk($size)
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

// ChunkWhile breaks the lazy collection into chunks while the callback returns true.
// Equivalent to: $lazy->chunkWhile($callback)
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

// Nth returns every n-th element.
// Equivalent to: $lazy->nth($step, $offset)
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

// Concat appends items lazily.
// Equivalent to: $lazy->concat($source)
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

// Pad pads the lazy collection to the given length.
// Equivalent to: $lazy->pad($size, $value)
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

// TapEach calls the given callback on each item without affecting the collection.
// Equivalent to: $lazy->tapEach($callback)
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

// Throttle adds a delay between each item.
// Equivalent to: $lazy->throttle($seconds)
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

// Remember caches yielded items so they're not recomputed on subsequent iterations.
// Equivalent to: $lazy->remember()
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

// Every determines if all items pass the given test.
// Equivalent to: $lazy->every($callback)
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

// Has determines if a key exists.
// Equivalent to: $lazy->has($key)
func (lc *LazyCollection[T]) Has(index int) bool {
	_, ok := lc.Get(index)
	return ok
}

// HasAny determines if any of the given keys exist.
// Equivalent to: $lazy->hasAny($keys)
func (lc *LazyCollection[T]) HasAny(indices ...int) bool {
	for _, idx := range indices {
		if lc.Has(idx) {
			return true
		}
	}
	return false
}

// Implode joins elements lazily into a string.
// Equivalent to: $lazy->implode($glue)
func (lc *LazyCollection[T]) Implode(glue string) string {
	parts := make([]string, 0)
	lc.source(func(item T) bool {
		parts = append(parts, fmt.Sprint(item))
		return true
	})
	return strings.Join(parts, glue)
}

// Join is like Implode but with an optional final glue.
// Equivalent to: $lazy->join($glue, $finalGlue)
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

// When applies the callback if the given condition is true.
// Equivalent to: $lazy->when($condition, $callback, $default)
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
// Equivalent to: $lazy->whenEmpty($callback, $default)
func (lc *LazyCollection[T]) WhenEmpty(callback func(*LazyCollection[T]) *LazyCollection[T], defaults ...func(*LazyCollection[T]) *LazyCollection[T]) *LazyCollection[T] {
	return lc.When(lc.IsEmpty(), callback, defaults...)
}

// WhenNotEmpty applies the callback if the lazy collection is not empty.
// Equivalent to: $lazy->whenNotEmpty($callback, $default)
func (lc *LazyCollection[T]) WhenNotEmpty(callback func(*LazyCollection[T]) *LazyCollection[T], defaults ...func(*LazyCollection[T]) *LazyCollection[T]) *LazyCollection[T] {
	return lc.When(lc.IsNotEmpty(), callback, defaults...)
}

// Unless applies the callback unless the condition is true.
// Equivalent to: $lazy->unless($condition, $callback, $default)
func (lc *LazyCollection[T]) Unless(condition bool, callback func(*LazyCollection[T]) *LazyCollection[T], defaults ...func(*LazyCollection[T]) *LazyCollection[T]) *LazyCollection[T] {
	return lc.When(!condition, callback, defaults...)
}

// LazyReduce reduces the lazy collection to a single value.
// Equivalent to: $lazy->reduce($callback, $initial)
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

// LazyUnique returns unique items lazily.
// Equivalent to: $lazy->unique($callback)
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

// LazyPluck extracts values lazily.
// Equivalent to: $lazy->pluck($value)
func LazyPluck[T any, V any](lc *LazyCollection[T], valueFunc func(T) V) *LazyCollection[V] {
	return LazyMap(lc, func(item T, _ int) V {
		return valueFunc(item)
	})
}

// LazyGroupBy groups items lazily (evaluates eagerly for grouping).
// Equivalent to: $lazy->groupBy($groupBy)
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

// LazyKeyBy keys items by the given function (evaluates eagerly).
// Equivalent to: $lazy->keyBy($keyBy)
func LazyKeyBy[T any, K comparable](lc *LazyCollection[T], keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	lc.source(func(item T) bool {
		result[keyFunc(item)] = item
		return true
	})
	return result
}

// LazyCountBy counts occurrences lazily.
// Equivalent to: $lazy->countBy($callback)
func LazyCountBy[T any, K comparable](lc *LazyCollection[T], keyFunc func(T) K) map[K]int {
	result := make(map[K]int)
	lc.source(func(item T) bool {
		result[keyFunc(item)]++
		return true
	})
	return result
}

// Dump prints the lazy collection items for debugging.
// Equivalent to: $lazy->dump()
func (lc *LazyCollection[T]) Dump() *LazyCollection[T] {
	items := lc.All()
	fmt.Printf("%v\n", items)
	return LazyFrom(items)
}
