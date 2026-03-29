package lazy

import (
	"reflect"
	"testing"
)

func TestFrom(t *testing.T) {
	lc := From([]int{1, 2, 3})
	items := lc.All()
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected %v, got %v", expected, items)
	}
}

func TestEmpty(t *testing.T) {
	lc := Empty[int]()

	if !lc.IsEmpty() {
		t.Error("expected empty")
	}

	if lc.IsNotEmpty() {
		t.Error("expected empty")
	}
}

func TestRange(t *testing.T) {
	lc := Range(1, 5)
	items := lc.All()
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected %v, got %v", expected, items)
	}

	lc2 := Range(5, 1)
	items2 := lc2.All()
	expected2 := []int{5, 4, 3, 2, 1}

	if !reflect.DeepEqual(items2, expected2) {
		t.Errorf("expected %v, got %v", expected2, items2)
	}
}

func TestTimes(t *testing.T) {
	lc := Times(3, func(i int) int { return i * 10 })
	expected := []int{10, 20, 30}

	if !reflect.DeepEqual(lc.All(), expected) {
		t.Errorf("expected %v, got %v", expected, lc.All())
	}
}

func TestCount(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})

	if lc.Count() != 5 {
		t.Errorf("expected 5, got %d", lc.Count())
	}
}

func TestEager(t *testing.T) {
	lc := From([]int{1, 2, 3})
	items := lc.Eager()

	if len(items) != 3 {
		t.Errorf("expected 3, got %d", len(items))
	}
}

func TestFirst(t *testing.T) {
	lc := From([]int{10, 20, 30})
	v, ok := lc.First()

	if !ok || v != 10 {
		t.Errorf("expected 10, got %d", v)
	}

	v, ok = lc.First(func(item int, _ int) bool { return item > 15 })

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}
}

func TestLast(t *testing.T) {
	lc := From([]int{10, 20, 30})
	v, ok := lc.Last()

	if !ok || v != 30 {
		t.Errorf("expected 30, got %d", v)
	}

	v, ok = lc.Last(func(item int, _ int) bool { return item < 25 })

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}
}

func TestGet(t *testing.T) {
	lc := From([]int{10, 20, 30})
	v, ok := lc.Get(1)

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}

	_, ok = lc.Get(10)

	if ok {
		t.Error("expected not found")
	}
}

func TestContains(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})

	if !lc.Contains(func(item int, _ int) bool { return item == 3 }) {
		t.Error("expected to contain 3")
	}

	if lc.Contains(func(item int, _ int) bool { return item == 10 }) {
		t.Error("expected not to contain 10")
	}
}

func TestSearch(t *testing.T) {
	lc := From([]int{10, 20, 30})
	idx, ok := lc.Search(func(item int, _ int) bool { return item == 20 })

	if !ok || idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}
}

func TestBefore(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	v, ok := lc.Before(func(item int, _ int) bool { return item == 3 })

	if !ok || v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
}

func TestAfter(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	v, ok := lc.After(func(item int, _ int) bool { return item == 3 })

	if !ok || v != 4 {
		t.Errorf("expected 4, got %d", v)
	}
}

func TestFilter(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	filtered := lc.Filter(func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{2, 4}

	if !reflect.DeepEqual(filtered.All(), expected) {
		t.Errorf("expected %v, got %v", expected, filtered.All())
	}
}

func TestReject(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	rejected := lc.Reject(func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{1, 3, 5}

	if !reflect.DeepEqual(rejected.All(), expected) {
		t.Errorf("expected %v, got %v", expected, rejected.All())
	}
}

func TestMap(t *testing.T) {
	lc := From([]int{1, 2, 3})
	mapped := Map(lc, func(item int, _ int) int { return item * 2 })
	expected := []int{2, 4, 6}

	if !reflect.DeepEqual(mapped.All(), expected) {
		t.Errorf("expected %v, got %v", expected, mapped.All())
	}
}

func TestFlatMap(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := FlatMap(lc, func(item int, _ int) []int { return []int{item, item * 10} })
	expected := []int{1, 10, 2, 20, 3, 30}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestTake(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	taken := lc.Take(3)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(taken.All(), expected) {
		t.Errorf("expected %v, got %v", expected, taken.All())
	}

	// Negative take
	taken2 := lc.Take(-2)
	expected2 := []int{4, 5}

	if !reflect.DeepEqual(taken2.All(), expected2) {
		t.Errorf("expected %v, got %v", expected2, taken2.All())
	}
}

func TestTakeUntil(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.TakeUntil(func(item int, _ int) bool { return item == 4 })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestTakeWhile(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.TakeWhile(func(item int, _ int) bool { return item < 4 })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSkip(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.Skip(2)
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSkipUntil(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.SkipUntil(func(item int, _ int) bool { return item == 3 })
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSkipWhile(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.SkipWhile(func(item int, _ int) bool { return item < 3 })
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSlice(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.Slice(1, 3)
	expected := []int{2, 3, 4}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestChunk(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	chunks := lc.Chunk(2)

	if len(chunks) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(chunks))
	}

	if !reflect.DeepEqual(chunks[0], []int{1, 2}) {
		t.Errorf("expected [1 2], got %v", chunks[0])
	}

	if !reflect.DeepEqual(chunks[2], []int{5}) {
		t.Errorf("expected [5], got %v", chunks[2])
	}
}

func TestNth(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5, 6, 7, 8})
	result := lc.Nth(3)
	expected := []int{1, 4, 7}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestConcat(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.Concat([]int{4, 5})
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestPad(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.Pad(5, 0)
	expected := []int{1, 2, 3, 0, 0}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestEvery(t *testing.T) {
	lc := From([]int{2, 4, 6, 8})

	if !lc.Every(func(item int, _ int) bool { return item%2 == 0 }) {
		t.Error("expected all even")
	}
}

func TestImplode(t *testing.T) {
	lc := From([]string{"a", "b", "c"})
	result := lc.Implode(", ")

	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}
}

func TestJoin(t *testing.T) {
	lc := From([]string{"a", "b", "c"})
	result := lc.Join(", ", " and ")

	if result != "a, b and c" {
		t.Errorf("expected 'a, b and c', got '%s'", result)
	}
}

func TestReduce(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	sum := Reduce(lc, func(carry int, item int, _ int) int { return carry + item }, 0)

	if sum != 15 {
		t.Errorf("expected 15, got %d", sum)
	}
}

func TestUnique(t *testing.T) {
	lc := From([]int{1, 2, 2, 3, 3, 3})
	result := Unique(lc, func(item int) int { return item })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestPluck(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	lc := From([]User{{"Alice", 25}, {"Bob", 30}})
	names := Pluck(lc, func(u User) string { return u.Name })
	expected := []string{"Alice", "Bob"}

	if !reflect.DeepEqual(names.All(), expected) {
		t.Errorf("expected %v, got %v", expected, names.All())
	}
}

func TestGroupBy(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5, 6})
	groups := GroupBy(lc, func(item int) string {
		if item%2 == 0 {
			return "even"
		}

		return "odd"
	})

	if len(groups["even"].All()) != 3 {
		t.Errorf("expected 3 evens, got %d", len(groups["even"].All()))
	}
}

func TestKeyBy(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	lc := From([]User{{1, "Alice"}, {2, "Bob"}})
	keyed := KeyBy(lc, func(u User) int { return u.ID })

	if keyed[1].Name != "Alice" {
		t.Error("expected Alice")
	}
}

func TestCountBy(t *testing.T) {
	lc := From([]string{"apple", "banana", "apple"})
	counts := CountBy(lc, func(item string) string { return item })

	if counts["apple"] != 2 {
		t.Errorf("expected 2, got %d", counts["apple"])
	}
}

func TestRemember(t *testing.T) {
	callCount := 0
	lc := New(func(yield func(int) bool) {
		callCount++

		for i := 1; i <= 3; i++ {
			if !yield(i) {
				return
			}
		}
	})

	remembered := lc.Remember()
	// First evaluation
	remembered.All()
	// Second evaluation - should use cache
	remembered.All()

	if callCount != 1 {
		t.Errorf("expected source called once, got %d", callCount)
	}
}

func TestContainsOneItem(t *testing.T) {
	if !From([]int{1}).ContainsOneItem() {
		t.Error("expected true")
	}

	if From([]int{1, 2}).ContainsOneItem() {
		t.Error("expected false")
	}
}

func TestWhen(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.When(true, func(lc *Collection[int]) *Collection[int] {
		return lc.Take(2)
	})

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}
}

func TestTapEach(t *testing.T) {
	lc := From([]int{1, 2, 3})
	sum := 0
	result := lc.TapEach(func(item int, _ int) {
		sum += item
	})
	// TapEach is lazy, need to evaluate
	result.All()

	if sum != 6 {
		t.Errorf("expected 6, got %d", sum)
	}
}

func TestFromSlice(t *testing.T) {
	items := []int{1, 2, 3}
	lc := From(items)

	if lc.Count() != 3 {
		t.Errorf("expected 3, got %d", lc.Count())
	}
}

func TestEach(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	sum := 0
	lc.Each(func(item int, _ int) bool {
		sum += item

		return item < 3
	})

	if sum != 6 { // 1+2+3
		t.Errorf("expected 6, got %d", sum)
	}
}

func TestSole(t *testing.T) {
	lc := From([]int{1, 2, 3})
	v, err := lc.Sole(func(item int, _ int) bool { return item == 2 })

	if err != nil || v != 2 {
		t.Errorf("expected 2, got %d, err: %v", v, err)
	}
}

func TestFirstOrFail(t *testing.T) {
	lc := From([]int{1, 2, 3})
	v, err := lc.FirstOrFail()

	if err != nil || v != 1 {
		t.Errorf("expected 1, got %d, err: %v", v, err)
	}

	empty := Empty[int]()
	_, err = empty.FirstOrFail()

	if err == nil {
		t.Error("expected error")
	}
}
