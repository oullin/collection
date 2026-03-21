package arr_test

import (
	"reflect"
	"testing"

	"github.com/gocanto/collection/arr"
)

func TestAccessible(t *testing.T) {
	if !arr.Accessible([]int{1, 2, 3}) {
		t.Error("expected true for non-nil slice")
	}
	if arr.Accessible(nil) {
		t.Error("expected false for nil")
	}
}

func TestIsList(t *testing.T) {
	if !arr.IsList([]int{1, 2, 3}) {
		t.Error("expected true")
	}
}

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

func TestOnly(t *testing.T) {
	items := []int{10, 20, 30, 40, 50}
	result := arr.Only(items, []int{1, 3})
	expected := []int{20, 40}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestExcept(t *testing.T) {
	items := []int{10, 20, 30, 40, 50}
	result := arr.Except(items, []int{1, 3})
	expected := []int{10, 30, 50}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestFlatten(t *testing.T) {
	items := [][]int{{1, 2}, {3, 4}, {5}}
	result := arr.Flatten(items)
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestCollapse(t *testing.T) {
	items := [][]int{{1, 2}, {3, 4}}
	result := arr.Collapse(items)
	if len(result) != 4 {
		t.Errorf("expected 4 items, got %d", len(result))
	}
}

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

func TestShuffle(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := arr.Shuffle(items)
	if len(result) != 5 {
		t.Errorf("expected 5 items, got %d", len(result))
	}
}

func TestRandom(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := arr.Random(items, 2)
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}

func TestSort(t *testing.T) {
	items := []int{3, 1, 4, 1, 5}
	result := arr.Sort(items, func(a, b int) bool { return a < b })
	expected := []int{1, 1, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSortDesc(t *testing.T) {
	items := []int{3, 1, 4, 1, 5}
	result := arr.SortDesc(items, func(a, b int) bool { return a < b })
	expected := []int{5, 4, 3, 1, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestWhere(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := arr.Where(items, func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{2, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestWhereNotNull(t *testing.T) {
	items := []string{"a", "", "b", "", "c"}
	result := arr.WhereNotNull(items)
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestReject(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := arr.Reject(items, func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{1, 3, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestPartition(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	pass, fail := arr.Partition(items, func(item int, _ int) bool { return item > 3 })
	if !reflect.DeepEqual(pass, []int{4, 5}) {
		t.Errorf("expected [4 5], got %v", pass)
	}
	if !reflect.DeepEqual(fail, []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", fail)
	}
}

func TestEvery(t *testing.T) {
	items := []int{2, 4, 6}
	if !arr.Every(items, func(item int, _ int) bool { return item%2 == 0 }) {
		t.Error("expected all even")
	}
	if arr.Every([]int{2, 3, 4}, func(item int, _ int) bool { return item%2 == 0 }) {
		t.Error("expected not all even")
	}
}

func TestSome(t *testing.T) {
	items := []int{1, 2, 3}
	if !arr.Some(items, func(item int, _ int) bool { return item == 2 }) {
		t.Error("expected to find 2")
	}
	if arr.Some(items, func(item int, _ int) bool { return item == 10 }) {
		t.Error("expected not to find 10")
	}
}

func TestExists(t *testing.T) {
	items := []int{10, 20, 30}
	if !arr.Exists(items, 1) {
		t.Error("expected index 1 to exist")
	}
	if arr.Exists(items, 10) {
		t.Error("expected index 10 to not exist")
	}
}

func TestHas(t *testing.T) {
	items := []int{10, 20, 30}
	if !arr.Has(items, 0, 1, 2) {
		t.Error("expected all indices to exist")
	}
	if arr.Has(items, 0, 5) {
		t.Error("expected to fail for index 5")
	}
}

func TestHasAny(t *testing.T) {
	items := []int{10, 20, 30}
	if !arr.HasAny(items, 0, 99) {
		t.Error("expected at least one valid index")
	}
	if arr.HasAny(items, 5, 10) {
		t.Error("expected no valid indices")
	}
}

func TestJoin(t *testing.T) {
	result := arr.Join([]string{"a", "b", "c"}, ", ")
	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}

	result = arr.Join([]string{"a", "b", "c"}, ", ", " and ")
	if result != "a, b and c" {
		t.Errorf("expected 'a, b and c', got '%s'", result)
	}
}

func TestCrossJoin(t *testing.T) {
	result := arr.CrossJoin([]int{1, 2}, []int{10, 20})
	if len(result) != 4 {
		t.Errorf("expected 4 combinations, got %d", len(result))
	}
}

func TestDivide(t *testing.T) {
	items := []string{"a", "b", "c"}
	keys, values := arr.Divide(items)
	if !reflect.DeepEqual(keys, []int{0, 1, 2}) {
		t.Errorf("expected [0 1 2], got %v", keys)
	}
	if !reflect.DeepEqual(values, items) {
		t.Errorf("expected %v, got %v", items, values)
	}
}

func TestMap(t *testing.T) {
	items := []int{1, 2, 3}
	result := arr.Map(items, func(item int, _ int) int { return item * 2 })
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapWithKeys(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	users := []User{{1, "Alice"}, {2, "Bob"}}
	result := arr.MapWithKeys(users, func(u User) (int, string) { return u.ID, u.Name })
	if result[1] != "Alice" {
		t.Error("expected Alice at key 1")
	}
}

func TestKeyBy(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	users := []User{{1, "Alice"}, {2, "Bob"}}
	result := arr.KeyBy(users, func(u User) int { return u.ID })
	if result[1].Name != "Alice" {
		t.Error("expected Alice")
	}
}

func TestPluck(t *testing.T) {
	type User struct {
		Name string
	}
	users := []User{{"Alice"}, {"Bob"}}
	result := arr.Pluck(users, func(u User) string { return u.Name })
	expected := []string{"Alice", "Bob"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestWrap(t *testing.T) {
	result := arr.Wrap(42)
	if !reflect.DeepEqual(result, []int{42}) {
		t.Errorf("expected [42], got %v", result)
	}
}

func TestWrapSlice(t *testing.T) {
	input := []int{1, 2, 3}
	result := arr.WrapSlice(input)
	if !reflect.DeepEqual(result, input) {
		t.Errorf("expected %v, got %v", input, result)
	}
}

func TestMapSpread(t *testing.T) {
	items := []int{1, 2, 3}
	result := arr.MapSpread(items, func(item int, _ int) int { return item * 10 })
	expected := []int{10, 20, 30}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSortRecursive(t *testing.T) {
	items := []int{3, 1, 2}
	result := arr.SortRecursive(items, func(a, b int) bool { return a < b })
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
