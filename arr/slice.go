package arr

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
)

// Accessible reports whether the given value is non-nil.
func Accessible(value any) bool {
	return value != nil
}

// IsList reports whether the given slice is a list.
// In Go, all slices are sequential lists, so this always returns true.
func IsList[T any](items []T) bool {
	return true
}

// First returns the first element that matches the optional callback.
// If no callback is provided, the first element of the slice is returned.
// The second return value reports whether a match was found.
func First[T any](items []T, callbacks ...func(T, int) bool) (T, bool) {
	if len(callbacks) == 0 || callbacks[0] == nil {
		if len(items) > 0 {
			return items[0], true
		}

		var zero T

		return zero, false
	}

	callback := callbacks[0]

	for i, item := range items {
		if callback(item, i) {
			return item, true
		}
	}

	var zero T

	return zero, false
}

// Last returns the last element that matches the optional callback.
// If no callback is provided, the last element of the slice is returned.
// The second return value reports whether a match was found.
func Last[T any](items []T, callbacks ...func(T, int) bool) (T, bool) {
	if len(callbacks) == 0 || callbacks[0] == nil {
		if len(items) > 0 {
			return items[len(items)-1], true
		}

		var zero T

		return zero, false
	}

	callback := callbacks[0]

	for i := len(items) - 1; i >= 0; i-- {
		if callback(items[i], i) {
			return items[i], true
		}
	}

	var zero T

	return zero, false
}

// Take returns up to limit elements from the slice.
// A positive limit takes from the front; a negative limit takes from the end.
func Take[T any](items []T, limit int) []T {
	if limit < 0 {
		start := len(items) + limit

		if start < 0 {
			start = 0
		}

		result := make([]T, len(items)-start)
		copy(result, items[start:])

		return result
	}

	if limit > len(items) {
		limit = len(items)
	}

	result := make([]T, limit)
	copy(result, items[:limit])

	return result
}

// Only returns the elements at the given indices.
func Only[T any](items []T, indices []int) []T {
	result := make([]T, 0, len(indices))

	for _, idx := range indices {
		if idx >= 0 && idx < len(items) {
			result = append(result, items[idx])
		}
	}

	return result
}

// Except returns all elements except those at the given indices.
func Except[T any](items []T, indices []int) []T {
	excludeSet := make(map[int]bool, len(indices))

	for _, idx := range indices {
		excludeSet[idx] = true
	}

	result := make([]T, 0)

	for i, item := range items {
		if !excludeSet[i] {
			result = append(result, item)
		}
	}

	return result
}

// Flatten flattens a slice of slices into a single slice.
func Flatten[T any](items [][]T) []T {
	result := make([]T, 0)

	for _, inner := range items {
		result = append(result, inner...)
	}

	return result
}

// Collapse merges a slice of slices into a single slice.
// It is an alias for [Flatten].
func Collapse[T any](items [][]T) []T {
	return Flatten(items)
}

// Wrap wraps the given value in a single-element slice.
func Wrap[T any](value T) []T {
	return []T{value}
}

// WrapSlice returns the given slice unchanged.
// Use this to wrap a value that is already a slice without nesting it.
func WrapSlice[T any](value []T) []T {
	return value
}

// Prepend inserts an element at the beginning of the slice and returns the new slice.
func Prepend[T any](items []T, value T) []T {
	return append([]T{value}, items...)
}

// Push appends one or more elements to the end of the slice and returns the new slice.
func Push[T any](items []T, values ...T) []T {
	return append(items, values...)
}

// Shuffle returns a new slice with the elements in random order.
func Shuffle[T any](items []T) []T {
	result := make([]T, len(items))
	copy(result, items)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

// Random returns a new slice containing count random elements.
// If count is omitted it defaults to 1.
func Random[T any](items []T, counts ...int) []T {
	count := 1

	if len(counts) > 0 {
		count = counts[0]
	}

	shuffled := Shuffle(items)

	if count >= len(shuffled) {
		return shuffled
	}

	return shuffled[:count]
}

// Sort returns a new slice sorted using the provided less function.
// The sort is stable.
func Sort[T any](items []T, less func(a, b T) bool) []T {
	result := make([]T, len(items))
	copy(result, items)
	sort.SliceStable(result, func(i, j int) bool {
		return less(result[i], result[j])
	})

	return result
}

// SortDesc returns a new slice sorted in descending order using the provided less function.
func SortDesc[T any](items []T, less func(a, b T) bool) []T {
	return Sort(items, func(a, b T) bool {
		return less(b, a)
	})
}

// SortRecursive sorts a slice using the provided less function.
// For flat slices this behaves identically to [Sort].
func SortRecursive[T any](items []T, less func(a, b T) bool) []T {
	return Sort(items, less)
}

// Where returns the elements for which the callback returns true.
func Where[T any](items []T, callback func(T, int) bool) []T {
	result := make([]T, 0)

	for i, item := range items {
		if callback(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// WhereNotNull returns all elements that are not the zero value of their type.
func WhereNotNull[T comparable](items []T) []T {
	var zero T
	result := make([]T, 0)

	for _, item := range items {
		if item != zero {
			result = append(result, item)
		}
	}

	return result
}

// Reject returns the elements for which the callback returns false.
func Reject[T any](items []T, callback func(T, int) bool) []T {
	return Where(items, func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Partition splits elements into two slices: those that pass the callback
// and those that do not.
func Partition[T any](items []T, callback func(T, int) bool) ([]T, []T) {
	pass := make([]T, 0)
	fail := make([]T, 0)

	for i, item := range items {
		if callback(item, i) {
			pass = append(pass, item)
		} else {
			fail = append(fail, item)
		}
	}

	return pass, fail
}

// Every reports whether every element passes the callback test.
func Every[T any](items []T, callback func(T, int) bool) bool {
	for i, item := range items {
		if !callback(item, i) {
			return false
		}
	}

	return true
}

// Some reports whether any element passes the callback test.
func Some[T any](items []T, callback func(T, int) bool) bool {
	for i, item := range items {
		if callback(item, i) {
			return true
		}
	}

	return false
}

// Exists reports whether the given index is valid for the slice.
func Exists[T any](items []T, index int) bool {
	return index >= 0 && index < len(items)
}

// Has reports whether all of the given indices are valid for the slice.
func Has[T any](items []T, indices ...int) bool {
	for _, idx := range indices {
		if idx < 0 || idx >= len(items) {
			return false
		}
	}

	return true
}

// HasAny reports whether at least one of the given indices is valid for the slice.
func HasAny[T any](items []T, indices ...int) bool {
	for _, idx := range indices {
		if idx >= 0 && idx < len(items) {
			return true
		}
	}

	return false
}

// Join concatenates string slice elements with a glue string.
// An optional final glue is used between the last two elements.
func Join(items []string, glue string, finalGlues ...string) string {
	if len(items) == 0 {
		return ""
	}

	if len(items) == 1 {
		return items[0]
	}

	if len(finalGlues) > 0 && finalGlues[0] != "" {
		last := items[len(items)-1]
		rest := items[:len(items)-1]

		return strings.Join(rest, glue) + finalGlues[0] + last
	}

	return strings.Join(items, glue)
}

// CrossJoin returns the Cartesian product of the given slices.
func CrossJoin[T any](lists ...[]T) [][]T {
	results := [][]T{{}}

	for _, list := range lists {
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

	return results
}

// Divide returns two slices: one containing the indices and one containing the values.
func Divide[T any](items []T) ([]int, []T) {
	keys := make([]int, len(items))
	values := make([]T, len(items))

	for i, item := range items {
		keys[i] = i
		values[i] = item
	}

	return keys, values
}

// Map applies a callback to each element and returns a slice of the results.
func Map[T any, R any](items []T, callback func(T, int) R) []R {
	result := make([]R, len(items))

	for i, item := range items {
		result[i] = callback(item, i)
	}

	return result
}

// MapWithKeys applies a callback to each element, producing key-value pairs
// that are collected into a map.
func MapWithKeys[T any, K comparable, V any](items []T, callback func(T) (K, V)) map[K]V {
	result := make(map[K]V)

	for _, item := range items {
		key, value := callback(item)
		result[key] = value
	}

	return result
}

// MapSpread applies a callback to each element and returns a slice of the results.
// In Go this is identical to [Map].
func MapSpread[T any, R any](items []T, callback func(T, int) R) []R {
	return Map(items, callback)
}

// KeyBy indexes the slice elements by the key returned from keyFunc.
func KeyBy[T any, K comparable](items []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T)

	for _, item := range items {
		result[keyFunc(item)] = item
	}

	return result
}

// Pluck extracts a value from each element using valueFunc and returns
// the collected values as a slice.
func Pluck[T any, V any](items []T, valueFunc func(T) V) []V {
	result := make([]V, len(items))

	for i, item := range items {
		result[i] = valueFunc(item)
	}

	return result
}

// Get returns the element at the given index.
// If the index is out of bounds, the first default value is returned,
// or the zero value of T.
func Get[T any](items []T, index int, defaults ...T) T {
	if index >= 0 && index < len(items) {
		return items[index]
	}

	if len(defaults) > 0 {
		return defaults[0]
	}

	var zero T

	return zero
}

// Set returns a new slice with the element at the given index replaced.
// If the index is out of bounds, the original slice is returned unchanged.
func Set[T any](items []T, index int, value T) []T {
	if index < 0 || index >= len(items) {
		result := make([]T, len(items))
		copy(result, items)

		return result
	}

	result := make([]T, len(items))
	copy(result, items)
	result[index] = value

	return result
}

// Forget returns a new slice with the element at the given index removed.
// If the index is out of bounds, the original slice is returned unchanged.
func Forget[T any](items []T, index int) []T {
	if index < 0 || index >= len(items) {
		result := make([]T, len(items))
		copy(result, items)

		return result
	}

	result := make([]T, 0, len(items)-1)
	result = append(result, items[:index]...)
	result = append(result, items[index+1:]...)

	return result
}

// Pull removes the element at the given index from the slice and returns both
// the removed element and the remaining slice. If the index is out of bounds,
// the zero value and the original slice are returned.
func Pull[T any](items []T, index int) (T, []T) {
	if index < 0 || index >= len(items) {
		var zero T
		result := make([]T, len(items))
		copy(result, items)

		return zero, result
	}

	value := items[index]
	result := make([]T, 0, len(items)-1)
	result = append(result, items[:index]...)
	result = append(result, items[index+1:]...)

	return value, result
}

// Sole returns the only element matching the optional callback.
// It returns an error if zero or more than one element matches.
func Sole[T any](items []T, callbacks ...func(T, int) bool) (T, error) {
	if len(callbacks) == 0 || callbacks[0] == nil {
		if len(items) == 1 {
			return items[0], nil
		}

		if len(items) == 0 {
			var zero T

			return zero, fmt.Errorf("no items found")
		}

		var zero T

		return zero, fmt.Errorf("multiple items found: %d items", len(items))
	}

	callback := callbacks[0]

	var result T
	found := 0

	for i, item := range items {
		if callback(item, i) {
			result = item
			found++

			if found > 1 {
				var zero T

				return zero, fmt.Errorf("multiple items found: %d+ items", found)
			}
		}
	}

	if found == 0 {
		var zero T

		return zero, fmt.Errorf("no items found")
	}

	return result, nil
}

// SortRecursiveDesc sorts a slice in descending order using the provided less function.
func SortRecursiveDesc[T any](items []T, less func(a, b T) bool) []T {
	return SortDesc(items, less)
}
