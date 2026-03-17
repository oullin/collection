package collection

import (
	"reflect"
	"testing"
)

func TestArrFirst(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	v, ok := ArrFirst(items)
	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	v, ok = ArrFirst(items, func(item int, _ int) bool { return item > 3 })
	if !ok || v != 4 {
		t.Errorf("expected 4, got %d", v)
	}
}

func TestArrLast(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	v, ok := ArrLast(items)
	if !ok || v != 5 {
		t.Errorf("expected 5, got %d", v)
	}

	v, ok = ArrLast(items, func(item int, _ int) bool { return item < 4 })
	if !ok || v != 3 {
		t.Errorf("expected 3, got %d", v)
	}
}

func TestArrTake(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	result := ArrTake(items, 3)
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	result2 := ArrTake(items, -2)
	expected2 := []int{4, 5}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("expected %v, got %v", expected2, result2)
	}
}

func TestArrOnly(t *testing.T) {
	items := []int{10, 20, 30, 40, 50}
	result := ArrOnly(items, []int{1, 3})
	expected := []int{20, 40}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrExcept(t *testing.T) {
	items := []int{10, 20, 30, 40, 50}
	result := ArrExcept(items, []int{1, 3})
	expected := []int{10, 30, 50}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrFlatten(t *testing.T) {
	items := [][]int{{1, 2}, {3, 4}, {5}}
	result := ArrFlatten(items)
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrPrepend(t *testing.T) {
	items := []int{2, 3}
	result := ArrPrepend(items, 1)
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrPush(t *testing.T) {
	items := []int{1, 2}
	result := ArrPush(items, 3, 4)
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrShuffle(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := ArrShuffle(items)
	if len(result) != 5 {
		t.Errorf("expected 5 items, got %d", len(result))
	}
}

func TestArrSort(t *testing.T) {
	items := []int{3, 1, 4, 1, 5}
	result := ArrSort(items, func(a, b int) bool { return a < b })
	expected := []int{1, 1, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrSortDesc(t *testing.T) {
	items := []int{3, 1, 4, 1, 5}
	result := ArrSortDesc(items, func(a, b int) bool { return a < b })
	expected := []int{5, 4, 3, 1, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrWhere(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := ArrWhere(items, func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{2, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrWhereNotNull(t *testing.T) {
	items := []string{"a", "", "b", "", "c"}
	result := ArrWhereNotNull(items)
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrReject(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := ArrReject(items, func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{1, 3, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrPartition(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	pass, fail := ArrPartition(items, func(item int, _ int) bool { return item > 3 })
	if !reflect.DeepEqual(pass, []int{4, 5}) {
		t.Errorf("expected [4 5], got %v", pass)
	}
	if !reflect.DeepEqual(fail, []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", fail)
	}
}

func TestArrEvery(t *testing.T) {
	items := []int{2, 4, 6}
	if !ArrEvery(items, func(item int, _ int) bool { return item%2 == 0 }) {
		t.Error("expected all even")
	}
}

func TestArrSome(t *testing.T) {
	items := []int{1, 2, 3}
	if !ArrSome(items, func(item int, _ int) bool { return item == 2 }) {
		t.Error("expected to find 2")
	}
}

func TestArrExists(t *testing.T) {
	items := []int{10, 20, 30}
	if !ArrExists(items, 1) {
		t.Error("expected index 1 to exist")
	}
	if ArrExists(items, 10) {
		t.Error("expected index 10 to not exist")
	}
}

func TestArrHas(t *testing.T) {
	items := []int{10, 20, 30}
	if !ArrHas(items, 0, 1, 2) {
		t.Error("expected all indices to exist")
	}
	if ArrHas(items, 0, 5) {
		t.Error("expected to fail for index 5")
	}
}

func TestArrJoin(t *testing.T) {
	result := ArrJoin([]string{"a", "b", "c"}, ", ")
	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}

	result2 := ArrJoin([]string{"a", "b", "c"}, ", ", " and ")
	if result2 != "a, b and c" {
		t.Errorf("expected 'a, b and c', got '%s'", result2)
	}
}

func TestArrCrossJoin(t *testing.T) {
	result := ArrCrossJoin([]int{1, 2}, []int{10, 20})
	if len(result) != 4 {
		t.Errorf("expected 4 combinations, got %d", len(result))
	}
}

func TestArrDivide(t *testing.T) {
	items := []string{"a", "b", "c"}
	keys, values := ArrDivide(items)
	if !reflect.DeepEqual(keys, []int{0, 1, 2}) {
		t.Errorf("expected [0 1 2], got %v", keys)
	}
	if !reflect.DeepEqual(values, items) {
		t.Errorf("expected %v, got %v", items, values)
	}
}

func TestArrMap(t *testing.T) {
	items := []int{1, 2, 3}
	result := ArrMap(items, func(item int, _ int) int { return item * 2 })
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrMapWithKeys(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	users := []User{{1, "Alice"}, {2, "Bob"}}
	result := ArrMapWithKeys(users, func(u User) (int, string) { return u.ID, u.Name })
	if result[1] != "Alice" {
		t.Error("expected Alice at key 1")
	}
}

func TestArrKeyBy(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	users := []User{{1, "Alice"}, {2, "Bob"}}
	result := ArrKeyBy(users, func(u User) int { return u.ID })
	if result[1].Name != "Alice" {
		t.Error("expected Alice")
	}
}

func TestArrPluck(t *testing.T) {
	type User struct {
		Name string
	}
	users := []User{{"Alice"}, {"Bob"}}
	result := ArrPluck(users, func(u User) string { return u.Name })
	expected := []string{"Alice", "Bob"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrIsList(t *testing.T) {
	if !ArrIsList([]int{1, 2, 3}) {
		t.Error("expected true")
	}
}

func TestArrRandom(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := ArrRandom(items, 2)
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}
