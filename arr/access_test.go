package arr_test

import (
	"reflect"
	"testing"

	"github.com/gocanto/collection/arr"
)

func TestFirst(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	v, ok := arr.First(items)

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	v, ok = arr.First(items, func(item int, _ int) bool { return item > 3 })

	if !ok || v != 4 {
		t.Errorf("expected 4, got %d", v)
	}

	_, ok = arr.First([]int{})

	if ok {
		t.Error("expected false for empty slice")
	}
}

func TestLast(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	v, ok := arr.Last(items)

	if !ok || v != 5 {
		t.Errorf("expected 5, got %d", v)
	}

	v, ok = arr.Last(items, func(item int, _ int) bool { return item < 4 })

	if !ok || v != 3 {
		t.Errorf("expected 3, got %d", v)
	}
}

func TestTake(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	result := arr.Take(items, 3)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	result = arr.Take(items, -2)
	expected = []int{4, 5}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
