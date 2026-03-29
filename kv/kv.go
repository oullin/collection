package kv

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// --- Dot-notation helpers for nested map[string]any ---

// Get retrieves a value from a nested map using dot notation.
// If the key is not found, the first default value is returned, or nil.
func Get(target map[string]any, key string, defaults ...any) any {
	if key == "" {
		return target
	}

	segments := strings.Split(key, ".")

	var current any = target

	for _, segment := range segments {
		switch v := current.(type) {
		case map[string]any:
			val, ok := v[segment]

			if !ok {
				if len(defaults) > 0 {
					return defaults[0]
				}

				return nil
			}

			current = val
		default:
			if len(defaults) > 0 {
				return defaults[0]
			}

			return nil
		}
	}

	return current
}

// Set sets a value in a nested map using dot notation.
// By default existing values are overwritten; pass false to preserve them.
func Set(target map[string]any, key string, value any, overwrite ...bool) map[string]any {
	shouldOverwrite := true

	if len(overwrite) > 0 {
		shouldOverwrite = overwrite[0]
	}

	segments := strings.Split(key, ".")
	current := target

	for i, segment := range segments {
		if i == len(segments)-1 {
			if shouldOverwrite {
				current[segment] = value
			} else {
				if _, exists := current[segment]; !exists {
					current[segment] = value
				}
			}
		} else {
			if _, exists := current[segment]; !exists {
				current[segment] = make(map[string]any)
			}

			if next, ok := current[segment].(map[string]any); ok {
				current = next
			} else {
				newMap := make(map[string]any)
				current[segment] = newMap
				current = newMap
			}
		}
	}

	return target
}

// Has reports whether the given dot-notated key exists in the nested map.
func Has(target map[string]any, key string) bool {
	if key == "" {
		return false
	}

	segments := strings.Split(key, ".")

	var current any = target

	for _, segment := range segments {
		switch v := current.(type) {
		case map[string]any:
			val, ok := v[segment]

			if !ok {
				return false
			}

			current = val
		default:
			return false
		}
	}

	return true
}

// Fill sets the value at the given dot-notated key only if it does not already exist.
func Fill(target map[string]any, key string, value any) map[string]any {
	return Set(target, key, value, false)
}

// Add is an alias for [Fill]. It sets a value only if the key does not already exist.
func Add(target map[string]any, key string, value any) map[string]any {
	return Fill(target, key, value)
}

// Forget removes the value at the given dot-notated key from the nested map.
func Forget(target map[string]any, key string) map[string]any {
	segments := strings.Split(key, ".")

	if len(segments) == 1 {
		delete(target, key)

		return target
	}

	current := target

	for i := 0; i < len(segments)-1; i++ {
		val, ok := current[segments[i]]

		if !ok {
			return target
		}

		next, ok := val.(map[string]any)

		if !ok {
			return target
		}

		current = next
	}

	delete(current, segments[len(segments)-1])

	return target
}

// ForgetMany removes one or more dot-notated keys from the map.
func ForgetMany(items map[string]any, keys ...string) map[string]any {
	for _, key := range keys {
		Forget(items, key)
	}

	return items
}

// Pull removes a key from the map and returns its value.
// If the key does not exist, the first default value is returned.
func Pull(items map[string]any, key string, defaults ...any) any {
	value := Get(items, key, defaults...)
	Forget(items, key)

	return value
}

// HasAll reports whether all of the given dot-notated keys exist in the map.
func HasAll(items map[string]any, keys ...string) bool {
	for _, key := range keys {
		if !Has(items, key) {
			return false
		}
	}

	return true
}

// HasAny reports whether at least one of the given dot-notated keys exists in the map.
func HasAny(items map[string]any, keys ...string) bool {
	for _, key := range keys {
		if Has(items, key) {
			return true
		}
	}

	return false
}

// --- Dot / Undot ---

// Dot flattens a nested map into a single-level map with dot-notated keys.
// An optional prefix is prepended to each key.
func Dot(items map[string]any, prepend ...string) map[string]any {
	prefix := ""

	if len(prepend) > 0 {
		prefix = prepend[0]
	}

	result := make(map[string]any)
	dotRecursive(items, prefix, result)

	return result
}

func dotRecursive(items map[string]any, prefix string, result map[string]any) {
	for key, value := range items {
		fullKey := key

		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]any); ok {
			dotRecursive(nested, fullKey, result)
		} else {
			result[fullKey] = value
		}
	}
}

// Undot expands a flat map with dot-notated keys into a nested map.
func Undot(items map[string]any) map[string]any {
	result := make(map[string]any)

	for key, value := range items {
		Set(result, key, value)
	}

	return result
}

// --- Map filtering and selection ---

// Only returns a new map containing only the specified keys.
func Only(items map[string]any, keys ...string) map[string]any {
	result := make(map[string]any)

	for _, key := range keys {
		if v, ok := items[key]; ok {
			result[key] = v
		}
	}

	return result
}

// Except returns a new map with the specified keys removed.
func Except(items map[string]any, keys ...string) map[string]any {
	result := make(map[string]any, len(items))

	for k, v := range items {
		result[k] = v
	}

	for _, key := range keys {
		delete(result, key)
	}

	return result
}

// IsAssoc reports whether the map is non-empty.
func IsAssoc(items map[string]any) bool {
	return len(items) > 0
}

// --- Encoding ---

// Query encodes the map as a URL query string.
func Query(items map[string]any) string {
	params := url.Values{}

	for k, v := range items {
		params.Set(k, fmt.Sprint(v))
	}

	return params.Encode()
}

// ToCssClasses returns a space-separated string of CSS class names
// whose corresponding boolean values are true. Classes are sorted alphabetically.
func ToCssClasses(classes map[string]bool) string {
	result := make([]string, 0)

	for class, condition := range classes {
		if condition {
			result = append(result, class)
		}
	}

	sort.Strings(result)

	return strings.Join(result, " ")
}

// ToCssStyles returns a space-separated string of CSS style declarations
// whose corresponding boolean values are true. A trailing semicolon is appended
// to each style if not already present. Styles are sorted alphabetically.
func ToCssStyles(styles map[string]bool) string {
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

// --- Sorting ---

// Sort returns a new map with the same entries, iterated by sorted keys.
// Note: Go maps do not guarantee iteration order; the returned map
// is sorted at construction time.
func Sort(items map[string]any) map[string]any {
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

// SortRecursive returns a new map with entries sorted by key,
// recursing into any nested maps.
func SortRecursive(items map[string]any) map[string]any {
	result := make(map[string]any, len(items))
	keys := make([]string, 0, len(items))

	for k := range items {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := items[k]

		if nested, ok := v.(map[string]any); ok {
			result[k] = SortRecursive(nested)
		} else {
			result[k] = v
		}
	}

	return result
}

// --- Map transformation ---

// Map applies a callback to each value in the map and returns a new map
// with the transformed values.
func Map[V any, R any](items map[string]V, callback func(V, string) R) map[string]R {
	result := make(map[string]R, len(items))

	for k, v := range items {
		result[k] = callback(v, k)
	}

	return result
}

// Where returns a new map containing only the entries for which
// the callback returns true.
func Where[V any](items map[string]V, callback func(V, string) bool) map[string]V {
	result := make(map[string]V)

	for k, v := range items {
		if callback(v, k) {
			result[k] = v
		}
	}

	return result
}

// PrependKeysWith returns a new map with the given prefix prepended to every key.
func PrependKeysWith[V any](items map[string]V, prefix string) map[string]V {
	result := make(map[string]V, len(items))

	for k, v := range items {
		result[prefix+k] = v
	}

	return result
}

// Replace returns a new map with entries from items, overwritten by
// any matching keys from the replacement maps.
func Replace(items map[string]any, replacements ...map[string]any) map[string]any {
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
