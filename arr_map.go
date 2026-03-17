package collection

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// ArrGet retrieves a value from a nested map using "dot" notation.
// Equivalent to: Arr::get($array, $key, $default)
func ArrGet(items map[string]any, key string, defaults ...any) any {
	return DataGet(items, key, defaults...)
}

// ArrSet sets a value in a nested map using "dot" notation.
// Equivalent to: Arr::set($array, $key, $value)
func ArrSet(items map[string]any, key string, value any) map[string]any {
	return DataSet(items, key, value)
}

// ArrAdd adds an element to a map if it doesn't exist.
// Equivalent to: Arr::add($array, $key, $value)
func ArrAdd(items map[string]any, key string, value any) map[string]any {
	return DataFill(items, key, value)
}

// ArrPull removes and returns a value from a map.
// Equivalent to: Arr::pull($array, $key, $default)
func ArrPull(items map[string]any, key string, defaults ...any) any {
	value := ArrGet(items, key, defaults...)
	DataForget(items, key)
	return value
}

// ArrForget removes one or more keys from a map.
// Equivalent to: Arr::forget($array, $keys)
func ArrForget(items map[string]any, keys ...string) map[string]any {
	for _, key := range keys {
		DataForget(items, key)
	}
	return items
}

// ArrHasMap determines if keys exist in a map.
// Equivalent to: Arr::has($array, $keys)
func ArrHasMap(items map[string]any, keys ...string) bool {
	for _, key := range keys {
		if !DataHas(items, key) {
			return false
		}
	}
	return true
}

// ArrHasAnyMap determines if any of the given keys exist.
// Equivalent to: Arr::hasAny($array, $keys)
func ArrHasAnyMap(items map[string]any, keys ...string) bool {
	for _, key := range keys {
		if DataHas(items, key) {
			return true
		}
	}
	return false
}

// ArrDot flattens a multi-dimensional map into a single level using "dot" notation.
// Equivalent to: Arr::dot($array, $prepend)
func ArrDot(items map[string]any, prepend ...string) map[string]any {
	prefix := ""
	if len(prepend) > 0 {
		prefix = prepend[0]
	}
	result := make(map[string]any)
	arrDotRecursive(items, prefix, result)
	return result
}

func arrDotRecursive(items map[string]any, prefix string, result map[string]any) {
	for key, value := range items {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		if nested, ok := value.(map[string]any); ok {
			arrDotRecursive(nested, fullKey, result)
		} else {
			result[fullKey] = value
		}
	}
}

// ArrUndot expands a "dot" notated map into a multi-dimensional map.
// Equivalent to: Arr::undot($array)
func ArrUndot(items map[string]any) map[string]any {
	result := make(map[string]any)
	for key, value := range items {
		DataSet(result, key, value)
	}
	return result
}

// ArrOnlyMap returns only the specified keys from a map.
// Equivalent to: Arr::only($array, $keys)
func ArrOnlyMap(items map[string]any, keys ...string) map[string]any {
	result := make(map[string]any)
	for _, key := range keys {
		if v, ok := items[key]; ok {
			result[key] = v
		}
	}
	return result
}

// ArrExceptMap returns all items except those with specified keys.
// Equivalent to: Arr::except($array, $keys)
func ArrExceptMap(items map[string]any, keys ...string) map[string]any {
	result := make(map[string]any, len(items))
	for k, v := range items {
		result[k] = v
	}
	for _, key := range keys {
		delete(result, key)
	}
	return result
}

// ArrIsAssoc determines if a map is an associative array (always true for maps).
// Equivalent to: Arr::isAssoc($array)
func ArrIsAssoc(items map[string]any) bool {
	return len(items) > 0
}

// ArrQuery builds a query string from a map.
// Equivalent to: Arr::query($array)
func ArrQuery(items map[string]any) string {
	params := url.Values{}
	for k, v := range items {
		params.Set(k, fmt.Sprint(v))
	}
	return params.Encode()
}

// ArrToCssClasses builds a CSS class string from a map of conditions.
// Equivalent to: Arr::toCssClasses($array)
func ArrToCssClasses(classes map[string]bool) string {
	result := make([]string, 0)
	for class, condition := range classes {
		if condition {
			result = append(result, class)
		}
	}
	sort.Strings(result)
	return strings.Join(result, " ")
}

// ArrToCssStyles builds a CSS style string from a map of conditions.
// Equivalent to: Arr::toCssStyles($styles)
func ArrToCssStyles(styles map[string]bool) string {
	result := make([]string, 0)
	for style, condition := range styles {
		if condition {
			if !strings.HasSuffix(style, ";") {
				style += ";"
			}
			result = append(result, style)
		}
	}
	sort.Strings(result)
	return strings.Join(result, " ")
}

// ArrSortMap sorts a map by keys.
// Equivalent to: Arr::sort($array)
func ArrSortMap(items map[string]any) map[string]any {
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result := make(map[string]any, len(items))
	for _, k := range keys {
		result[k] = items[k]
	}
	return result
}

// ArrSortRecursiveMap recursively sorts a map by keys.
// Equivalent to: Arr::sortRecursive($array)
func ArrSortRecursiveMap(items map[string]any) map[string]any {
	result := make(map[string]any, len(items))
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := items[k]
		if nested, ok := v.(map[string]any); ok {
			result[k] = ArrSortRecursiveMap(nested)
		} else {
			result[k] = v
		}
	}
	return result
}

// ArrMapMap applies a callback to each value in a map.
// Equivalent to: Arr::map($array, $callback)
func ArrMapMap[V any, R any](items map[string]V, callback func(V, string) R) map[string]R {
	result := make(map[string]R, len(items))
	for k, v := range items {
		result[k] = callback(v, k)
	}
	return result
}

// ArrWhereMap filters a map by callback.
// Equivalent to: Arr::where($array, $callback)
func ArrWhereMap[V any](items map[string]V, callback func(V, string) bool) map[string]V {
	result := make(map[string]V)
	for k, v := range items {
		if callback(v, k) {
			result[k] = v
		}
	}
	return result
}

// ArrPrependKeysWith prepends all keys with the given string.
// Equivalent to: Arr::prependKeysWith($array, $prependWith)
func ArrPrependKeysWith[V any](items map[string]V, prefix string) map[string]V {
	result := make(map[string]V, len(items))
	for k, v := range items {
		result[prefix+k] = v
	}
	return result
}

// ArrMapReplace replaces items in a map.
// Equivalent to: Arr::replace($array, ...$replacements)
func ArrMapReplace(items map[string]any, replacements ...map[string]any) map[string]any {
	result := make(map[string]any, len(items))
	for k, v := range items {
		result[k] = v
	}
	for _, replacement := range replacements {
		for k, v := range replacement {
			result[k] = v
		}
	}
	return result
}
