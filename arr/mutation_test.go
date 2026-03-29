package arr_test

import (
	"reflect"
	"testing"

	"github.com/gocanto/collection/arr"
)

func TestPrepend(t *testing.T) {
	items := []int{2, 3}
	result := arr.Prepend(items, 1)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestPush(t *testing.T) {
	items := []int{1, 2}
	result := arr.Push(items, 3, 4)
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
