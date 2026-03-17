package collection

import (
	"strings"
)

// Value resolves a value or a callback.
// Equivalent to: value($value, ...$args)
func Value[T any](value T) T {
	return value
}

// ValueFunc resolves a callback to a value.
// Equivalent to: value($callback)
func ValueFunc[T any](callback func() T) T {
	return callback()
}

// Head returns the first element of a slice.
// Equivalent to: head($array)
func Head[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	return items[0], true
}

// Last returns the last element of a slice.
// Equivalent to: last($array)
func Last[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	return items[len(items)-1], true
}

// WhenValue returns the value if condition is true, otherwise the default.
// Equivalent to: when($condition, $value, $default)
func WhenValue[T any](condition bool, value T, defaults ...T) T {
	if condition {
		return value
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	var zero T
	return zero
}

// WhenFunc returns the callback result if condition is true, otherwise the default.
// Equivalent to: when($condition, $callback, $default)
func WhenFunc[T any](condition bool, callback func() T, defaults ...func() T) T {
	if condition {
		return callback()
	}
	if len(defaults) > 0 {
		return defaults[0]()
	}
	var zero T
	return zero
}

// DataGet retrieves a value from a nested map using "dot" notation.
// Equivalent to: data_get($target, $key, $default)
func DataGet(target map[string]any, key string, defaults ...any) any {
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

// DataSet sets a value in a nested map using "dot" notation.
// Equivalent to: data_set($target, $key, $value, $overwrite)
func DataSet(target map[string]any, key string, value any, overwrite ...bool) map[string]any {
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

// DataHas checks if a key exists in a nested map using "dot" notation.
// Equivalent to: data_has($target, $key)
func DataHas(target map[string]any, key string) bool {
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

// DataFill sets a value only if the key doesn't already exist.
// Equivalent to: data_fill($target, $key, $value)
func DataFill(target map[string]any, key string, value any) map[string]any {
	return DataSet(target, key, value, false)
}

// DataForget removes a key from a nested map using "dot" notation.
// Equivalent to: data_forget($target, $key)
func DataForget(target map[string]any, key string) map[string]any {
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
