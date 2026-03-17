package collection

import "fmt"

// ItemNotFoundException is returned when an item lookup fails.
type ItemNotFoundException struct {
	Message string
}

func (e *ItemNotFoundException) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "item not found"
}

// MultipleItemsFoundException is returned when a single item is expected but multiple are found.
type MultipleItemsFoundException struct {
	Count   int
	Message string
}

func (e *MultipleItemsFoundException) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("multiple items found: %d items", e.Count)
}
