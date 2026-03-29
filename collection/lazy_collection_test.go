package collection

import (
	"reflect"
	"testing"
)

func TestLazyFrom(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
	items := lc.All()
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected %v, got %v", expected, items)
	}
}

func TestLazyEmpty(t *testing.T) {
	lc := LazyEmpty[int]()

	if !lc.IsEmpty() {
		t.Error("expected empty")
	}

	if lc.IsNotEmpty() {
		t.Error("expected empty")
	}
}

func TestLazyRange(t *testing.T) {
	lc := LazyRange(1, 5)
	items := lc.All()
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected %v, got %v", expected, items)
	}

	lc2 := LazyRange(5, 1)
	items2 := lc2.All()
	expected2 := []int{5, 4, 3, 2, 1}

	if !reflect.DeepEqual(items2, expected2) {
		t.Errorf("expected %v, got %v", expected2, items2)
	}
}

func TestLazyTimes(t *testing.T) {
	lc := LazyTimes(3, func(i int) int { return i * 10 })
	expected := []int{10, 20, 30}

	if !reflect.DeepEqual(lc.All(), expected) {
		t.Errorf("expected %v, got %v", expected, lc.All())
	}
}

func TestLazyCount(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})

	if lc.Count() != 5 {
		t.Errorf("expected 5, got %d", lc.Count())
	}
}

func TestLazyEager(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
	c := lc.Eager()

	if c.Count() != 3 {
		t.Errorf("expected 3, got %d", c.Count())
	}
}

func TestLazyFirst(t *testing.T) {
	lc := LazyFrom([]int{10, 20, 30})
	v, ok := lc.First()

	if !ok || v != 10 {
		t.Errorf("expected 10, got %d", v)
	}

	v, ok = lc.First(func(item int, _ int) bool { return item > 15 })

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}
}

func TestLazyLast(t *testing.T) {
	lc := LazyFrom([]int{10, 20, 30})
	v, ok := lc.Last()

	if !ok || v != 30 {
		t.Errorf("expected 30, got %d", v)
	}

	v, ok = lc.Last(func(item int, _ int) bool { return item < 25 })

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}
}

func TestLazyGet(t *testing.T) {
	lc := LazyFrom([]int{10, 20, 30})
	v, ok := lc.Get(1)

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}

	_, ok = lc.Get(10)

	if ok {
		t.Error("expected not found")
	}
}

func TestLazyContains(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})

	if !lc.Contains(func(item int, _ int) bool { return item == 3 }) {
		t.Error("expected to contain 3")
	}

	if lc.Contains(func(item int, _ int) bool { return item == 10 }) {
		t.Error("expected not to contain 10")
	}
}

func TestLazySearch(t *testing.T) {
	lc := LazyFrom([]int{10, 20, 30})
	idx, ok := lc.Search(func(item int, _ int) bool { return item == 20 })

	if !ok || idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}
}

func TestLazyBefore(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	v, ok := lc.Before(func(item int, _ int) bool { return item == 3 })

	if !ok || v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
}

func TestLazyAfter(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	v, ok := lc.After(func(item int, _ int) bool { return item == 3 })

	if !ok || v != 4 {
		t.Errorf("expected 4, got %d", v)
	}
}

func TestLazyFilter(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	filtered := lc.Filter(func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{2, 4}

	if !reflect.DeepEqual(filtered.All(), expected) {
		t.Errorf("expected %v, got %v", expected, filtered.All())
	}
}

func TestLazyReject(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	rejected := lc.Reject(func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{1, 3, 5}

	if !reflect.DeepEqual(rejected.All(), expected) {
		t.Errorf("expected %v, got %v", expected, rejected.All())
	}
}

func TestLazyMap(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
	mapped := LazyMap(lc, func(item int, _ int) int { return item * 2 })
	expected := []int{2, 4, 6}

	if !reflect.DeepEqual(mapped.All(), expected) {
		t.Errorf("expected %v, got %v", expected, mapped.All())
	}
}

func TestLazyFlatMap(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
	result := LazyFlatMap(lc, func(item int, _ int) []int { return []int{item, item * 10} })
	expected := []int{1, 10, 2, 20, 3, 30}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazyTake(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
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

func TestLazyTakeUntil(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	result := lc.TakeUntil(func(item int, _ int) bool { return item == 4 })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazyTakeWhile(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	result := lc.TakeWhile(func(item int, _ int) bool { return item < 4 })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazySkip(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	result := lc.Skip(2)
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazySkipUntil(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	result := lc.SkipUntil(func(item int, _ int) bool { return item == 3 })
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazySkipWhile(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	result := lc.SkipWhile(func(item int, _ int) bool { return item < 3 })
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazySlice(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	result := lc.Slice(1, 3)
	expected := []int{2, 3, 4}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazyChunk(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
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

func TestLazyNth(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5, 6, 7, 8})
	result := lc.Nth(3)
	expected := []int{1, 4, 7}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazyConcat(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
	result := lc.Concat([]int{4, 5})
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazyPad(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
	result := lc.Pad(5, 0)
	expected := []int{1, 2, 3, 0, 0}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazyEvery(t *testing.T) {
	lc := LazyFrom([]int{2, 4, 6, 8})

	if !lc.Every(func(item int, _ int) bool { return item%2 == 0 }) {
		t.Error("expected all even")
	}
}

func TestLazyImplode(t *testing.T) {
	lc := LazyFrom([]string{"a", "b", "c"})
	result := lc.Implode(", ")

	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}
}

func TestLazyJoin(t *testing.T) {
	lc := LazyFrom([]string{"a", "b", "c"})
	result := lc.Join(", ", " and ")

	if result != "a, b and c" {
		t.Errorf("expected 'a, b and c', got '%s'", result)
	}
}

func TestLazyReduce(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	sum := LazyReduce(lc, func(carry int, item int, _ int) int { return carry + item }, 0)

	if sum != 15 {
		t.Errorf("expected 15, got %d", sum)
	}
}

func TestLazyUnique(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 2, 3, 3, 3})
	result := LazyUnique(lc, func(item int) int { return item })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazyPluck(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	lc := LazyFrom([]User{{"Alice", 25}, {"Bob", 30}})
	names := LazyPluck(lc, func(u User) string { return u.Name })
	expected := []string{"Alice", "Bob"}

	if !reflect.DeepEqual(names.All(), expected) {
		t.Errorf("expected %v, got %v", expected, names.All())
	}
}

func TestLazyGroupBy(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5, 6})
	groups := LazyGroupBy(lc, func(item int) string {
		if item%2 == 0 {
			return "even"
		}

		return "odd"
	})

	if len(groups["even"].All()) != 3 {
		t.Errorf("expected 3 evens, got %d", len(groups["even"].All()))
	}
}

func TestLazyKeyBy(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	lc := LazyFrom([]User{{1, "Alice"}, {2, "Bob"}})
	keyed := LazyKeyBy(lc, func(u User) int { return u.ID })

	if keyed[1].Name != "Alice" {
		t.Error("expected Alice")
	}
}

func TestLazyCountBy(t *testing.T) {
	lc := LazyFrom([]string{"apple", "banana", "apple"})
	counts := LazyCountBy(lc, func(item string) string { return item })

	if counts["apple"] != 2 {
		t.Errorf("expected 2, got %d", counts["apple"])
	}
}

func TestLazyRemember(t *testing.T) {
	callCount := 0
	lc := NewLazy(func(yield func(int) bool) {
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

func TestLazyContainsOneItem(t *testing.T) {
	if !LazyFrom([]int{1}).ContainsOneItem() {
		t.Error("expected true")
	}

	if LazyFrom([]int{1, 2}).ContainsOneItem() {
		t.Error("expected false")
	}
}

func TestLazyWhen(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
	result := lc.When(true, func(lc *LazyCollection[int]) *LazyCollection[int] {
		return lc.Take(2)
	})

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}
}

func TestLazyTapEach(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
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

func TestCollectionToLazy(t *testing.T) {
	c := New(1, 2, 3)
	lc := c.Lazy()

	if lc.Count() != 3 {
		t.Errorf("expected 3, got %d", lc.Count())
	}
}

func TestLazyEach(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3, 4, 5})
	sum := 0
	lc.Each(func(item int, _ int) bool {
		sum += item

		return item < 3
	})

	if sum != 6 { // 1+2+3
		t.Errorf("expected 6, got %d", sum)
	}
}

func TestLazySole(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
	v, err := lc.Sole(func(item int, _ int) bool { return item == 2 })

	if err != nil || v != 2 {
		t.Errorf("expected 2, got %d, err: %v", v, err)
	}
}

func TestLazyFirstOrFail(t *testing.T) {
	lc := LazyFrom([]int{1, 2, 3})
	v, err := lc.FirstOrFail()

	if err != nil || v != 1 {
		t.Errorf("expected 1, got %d, err: %v", v, err)
	}

	empty := LazyEmpty[int]()
	_, err = empty.FirstOrFail()

	if err == nil {
		t.Error("expected error")
	}
}
