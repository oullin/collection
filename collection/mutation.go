package collection

import "fmt"

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

	newItems := make([]T, len(c.items)-1)
	copy(newItems, c.items[:len(c.items)-1])
	c.items = newItems

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

	items := make([]T, count)
	copy(items, c.items[idx:])
	popped := Collect(items)

	remaining := make([]T, idx)
	copy(remaining, c.items[:idx])
	c.items = remaining

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

	newItems := make([]T, len(c.items)-1)
	copy(newItems, c.items[1:])
	c.items = newItems

	return item, true
}

// ShiftMany removes and returns the first n items from the collection.
func (c *Collection[T]) ShiftMany(count int) *Collection[T] {
	if count >= len(c.items) {
		shifted := Collect(c.items)
		c.items = make([]T, 0)

		return shifted
	}

	items := make([]T, count)
	copy(items, c.items[:count])
	shifted := Collect(items)

	remaining := make([]T, len(c.items)-count)
	copy(remaining, c.items[count:])
	c.items = remaining

	return shifted
}

// Forget removes an item from the collection by index, mutating the collection.
func (c *Collection[T]) Forget(index int) *Collection[T] {
	if index < 0 || index >= len(c.items) {
		return c
	}

	c.items = append(c.items[:index], c.items[index+1:]...)

	return c
}

// Transform applies the callback to each item in place, mutating the collection.
func (c *Collection[T]) Transform(callback func(T, int) T) *Collection[T] {
	for i, item := range c.items {
		c.items[i] = callback(item, i)
	}

	return c
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
