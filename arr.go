package collection

import (
	"math/rand/v2"
	"sort"
	"strings"
)

// Arr provides static helper functions for working with slices and maps,
// mirroring Laravel's Illuminate\Support\Arr class.

// ArrAccessible checks if the given value is array accessible (always true for slices/maps in Go).
// Equivalent to: Arr::accessible($value)
func ArrAccessible(value any) bool {
	if value == nil {
		return false
	}
	return true
}

// ArrIsList determines if the given slice is a list (sequential integer keys starting from 0).
// In Go, all slices are lists.
// Equivalent to: Arr::isList($array)
func ArrIsList[T any](items []T) bool {
	return true
}

// ArrFirst returns the first element matching the callback.
// Equivalent to: Arr::first($array, $callback, $default)
func ArrFirst[T any](items []T, callbacks ...func(T, int) bool) (T, bool) {
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

// ArrLast returns the last element matching the callback.
// Equivalent to: Arr::last($array, $callback, $default)
func ArrLast[T any](items []T, callbacks ...func(T, int) bool) (T, bool) {
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

// ArrTake returns the first N items from a slice.
// Equivalent to: Arr::take($array, $limit)
func ArrTake[T any](items []T, limit int) []T {
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

// ArrOnly returns a subset of items from the slice at the given indices.
// Equivalent to: Arr::only($array, $keys)
func ArrOnly[T any](items []T, indices []int) []T {
	result := make([]T, 0, len(indices))
	for _, idx := range indices {
		if idx >= 0 && idx < len(items) {
			result = append(result, items[idx])
		}
	}
	return result
}

// ArrExcept returns all items except those at the given indices.
// Equivalent to: Arr::except($array, $keys)
func ArrExcept[T any](items []T, indices []int) []T {
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

// ArrFlatten flattens a multi-dimensional slice.
// Equivalent to: Arr::flatten($array)
func ArrFlatten[T any](items [][]T) []T {
	result := make([]T, 0)
	for _, inner := range items {
		result = append(result, inner...)
	}
	return result
}

// ArrCollapse collapses an array of arrays into a single array.
// Equivalent to: Arr::collapse($array)
func ArrCollapse[T any](items [][]T) []T {
	return ArrFlatten(items)
}

// ArrWrap wraps the given value in a slice if it is not already a slice.
// Equivalent to: Arr::wrap($value)
func ArrWrap[T any](value T) []T {
	return []T{value}
}

// ArrWrapSlice wraps an existing slice (identity).
// Equivalent to: Arr::wrap($value) when value is array
func ArrWrapSlice[T any](value []T) []T {
	return value
}

// ArrPrepend adds an item to the beginning of a slice.
// Equivalent to: Arr::prepend($array, $value)
func ArrPrepend[T any](items []T, value T) []T {
	return append([]T{value}, items...)
}

// ArrPush appends one or more items to a slice.
// Equivalent to: Arr::push($array, ...$values)
func ArrPush[T any](items []T, values ...T) []T {
	return append(items, values...)
}

// ArrShuffle shuffles the items in a slice.
// Equivalent to: Arr::shuffle($array)
func ArrShuffle[T any](items []T) []T {
	result := make([]T, len(items))
	copy(result, items)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}

// ArrRandom returns random items from a slice.
// Equivalent to: Arr::random($array, $number)
func ArrRandom[T any](items []T, counts ...int) []T {
	count := 1
	if len(counts) > 0 {
		count = counts[0]
	}
	shuffled := ArrShuffle(items)
	if count >= len(shuffled) {
		return shuffled
	}
	return shuffled[:count]
}

// ArrSort sorts a slice.
// Equivalent to: Arr::sort($array)
func ArrSort[T any](items []T, less func(a, b T) bool) []T {
	result := make([]T, len(items))
	copy(result, items)
	sort.SliceStable(result, func(i, j int) bool {
		return less(result[i], result[j])
	})
	return result
}

// ArrSortDesc sorts a slice in descending order.
// Equivalent to: Arr::sortDesc($array)
func ArrSortDesc[T any](items []T, less func(a, b T) bool) []T {
	return ArrSort(items, func(a, b T) bool {
		return less(b, a)
	})
}

// ArrSortRecursive recursively sorts a slice.
// Equivalent to: Arr::sortRecursive($array)
func ArrSortRecursive[T any](items []T, less func(a, b T) bool) []T {
	return ArrSort(items, less)
}

// ArrWhere filters a slice by the given callback.
// Equivalent to: Arr::where($array, $callback)
func ArrWhere[T any](items []T, callback func(T, int) bool) []T {
	result := make([]T, 0)
	for i, item := range items {
		if callback(item, i) {
			result = append(result, item)
		}
	}
	return result
}

// ArrWhereNotNull filters out nil/zero values.
// Equivalent to: Arr::whereNotNull($array)
func ArrWhereNotNull[T comparable](items []T) []T {
	var zero T
	result := make([]T, 0)
	for _, item := range items {
		if item != zero {
			result = append(result, item)
		}
	}
	return result
}

// ArrReject returns items that don't pass the given callback.
// Equivalent to: Arr::reject($array, $callback)
func ArrReject[T any](items []T, callback func(T, int) bool) []T {
	return ArrWhere(items, func(item T, index int) bool {
		return !callback(item, index)
	})
}

// ArrPartition separates items that pass the test from those that don't.
// Equivalent to: Arr::partition($array, $callback)
func ArrPartition[T any](items []T, callback func(T, int) bool) ([]T, []T) {
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

// ArrEvery determines if all items pass the given test.
// Equivalent to: Arr::every($array, $callback)
func ArrEvery[T any](items []T, callback func(T, int) bool) bool {
	for i, item := range items {
		if !callback(item, i) {
			return false
		}
	}
	return true
}

// ArrSome determines if any items pass the given test.
// Equivalent to: Arr::some($array, $callback)
func ArrSome[T any](items []T, callback func(T, int) bool) bool {
	for i, item := range items {
		if callback(item, i) {
			return true
		}
	}
	return false
}

// ArrExists checks if a key exists in the slice.
// Equivalent to: Arr::exists($array, $key)
func ArrExists[T any](items []T, index int) bool {
	return index >= 0 && index < len(items)
}

// ArrHas determines if given keys exist in the slice.
// Equivalent to: Arr::has($array, $keys)
func ArrHas[T any](items []T, indices ...int) bool {
	for _, idx := range indices {
		if idx < 0 || idx >= len(items) {
			return false
		}
	}
	return true
}

// ArrHasAny determines if any of the given keys exist.
// Equivalent to: Arr::hasAny($array, $keys)
func ArrHasAny[T any](items []T, indices ...int) bool {
	for _, idx := range indices {
		if idx >= 0 && idx < len(items) {
			return true
		}
	}
	return false
}

// ArrJoin joins slice items into a string with a glue.
// Equivalent to: Arr::join($array, $glue, $finalGlue)
func ArrJoin(items []string, glue string, finalGlues ...string) string {
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

// ArrCrossJoin cross joins multiple slices.
// Equivalent to: Arr::crossJoin(...$arrays)
func ArrCrossJoin[T any](lists ...[]T) [][]T {
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

// ArrDivide returns two slices: one with keys (indices), one with values.
// Equivalent to: Arr::divide($array)
func ArrDivide[T any](items []T) ([]int, []T) {
	keys := make([]int, len(items))
	values := make([]T, len(items))
	for i, item := range items {
		keys[i] = i
		values[i] = item
	}
	return keys, values
}

// ArrMap applies a callback to each item.
// Equivalent to: Arr::map($array, $callback)
func ArrMap[T any, R any](items []T, callback func(T, int) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = callback(item, i)
	}
	return result
}

// ArrMapWithKeys maps items to key-value pairs.
// Equivalent to: Arr::mapWithKeys($array, $callback)
func ArrMapWithKeys[T any, K comparable, V any](items []T, callback func(T) (K, V)) map[K]V {
	result := make(map[K]V)
	for _, item := range items {
		key, value := callback(item)
		result[key] = value
	}
	return result
}

// ArrMapSpread maps items using spread arguments.
// For Go, this is equivalent to ArrMap.
// Equivalent to: Arr::mapSpread($array, $callback)
func ArrMapSpread[T any, R any](items []T, callback func(T, int) R) []R {
	return ArrMap(items, callback)
}

// ArrKeyBy keys a slice by the given key function.
// Equivalent to: Arr::keyBy($array, $keyBy)
func ArrKeyBy[T any, K comparable](items []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range items {
		result[keyFunc(item)] = item
	}
	return result
}

// ArrPluck extracts values using a key function.
// Equivalent to: Arr::pluck($array, $value)
func ArrPluck[T any, V any](items []T, valueFunc func(T) V) []V {
	result := make([]V, len(items))
	for i, item := range items {
		result[i] = valueFunc(item)
	}
	return result
}
