package collection

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	c := New(1, 2, 3)

	if c.Count() != 3 {
		t.Errorf("expected 3 items, got %d", c.Count())
	}

	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(c.All(), expected) {
		t.Errorf("expected %v, got %v", expected, c.All())
	}
}

func TestCollect(t *testing.T) {
	items := []string{"a", "b", "c"}
	c := Collect(items)

	if c.Count() != 3 {
		t.Errorf("expected 3, got %d", c.Count())
	}
}

func TestEmpty(t *testing.T) {
	c := Empty[int]()

	if !c.IsEmpty() {
		t.Error("expected empty collection")
	}

	if c.IsNotEmpty() {
		t.Error("expected empty collection")
	}
}

func TestTimes(t *testing.T) {
	c := Times(5, func(i int) int { return i * 2 })
	expected := []int{2, 4, 6, 8, 10}

	if !reflect.DeepEqual(c.All(), expected) {
		t.Errorf("expected %v, got %v", expected, c.All())
	}
}

func TestRange(t *testing.T) {
	c := Range(1, 5)
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(c.All(), expected) {
		t.Errorf("expected %v, got %v", expected, c.All())
	}

	c2 := Range(5, 1)
	expected2 := []int{5, 4, 3, 2, 1}

	if !reflect.DeepEqual(c2.All(), expected2) {
		t.Errorf("expected %v, got %v", expected2, c2.All())
	}
}

func TestContainsOneItem(t *testing.T) {
	if !New(1).ContainsOneItem() {
		t.Error("expected true")
	}

	if New(1, 2).ContainsOneItem() {
		t.Error("expected false")
	}
}

func TestContainsManyItems(t *testing.T) {
	if !New(1, 2).ContainsManyItems() {
		t.Error("expected true")
	}

	if New(1).ContainsManyItems() {
		t.Error("expected false")
	}
}

func TestFirst(t *testing.T) {
	c := New(1, 2, 3, 4, 5)

	v, ok := c.First()

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	v, ok = c.First(func(item int, _ int) bool { return item > 3 })

	if !ok || v != 4 {
		t.Errorf("expected 4, got %d", v)
	}

	_, ok = c.First(func(item int, _ int) bool { return item > 10 })

	if ok {
		t.Error("expected not found")
	}
}

func TestFirstOrFail(t *testing.T) {
	c := New(1, 2, 3)
	v, err := c.FirstOrFail()

	if err != nil || v != 1 {
		t.Errorf("expected 1, got %d, err: %v", v, err)
	}

	empty := Empty[int]()
	_, err = empty.FirstOrFail()

	if err == nil {
		t.Error("expected error")
	}
}

func TestLast(t *testing.T) {
	c := New(1, 2, 3, 4, 5)

	v, ok := c.Last()

	if !ok || v != 5 {
		t.Errorf("expected 5, got %d", v)
	}

	v, ok = c.Last(func(item int, _ int) bool { return item < 3 })

	if !ok || v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
}

func TestSole(t *testing.T) {
	c := New(1, 2, 3)
	v, err := c.Sole(func(item int, _ int) bool { return item == 2 })

	if err != nil || v != 2 {
		t.Errorf("expected 2, got %d, err: %v", v, err)
	}

	_, err = c.Sole(func(item int, _ int) bool { return item > 1 })

	if err == nil {
		t.Error("expected MultipleItemsFoundError")
	}

	_, err = c.Sole(func(item int, _ int) bool { return item > 10 })

	if err == nil {
		t.Error("expected ItemNotFoundError")
	}
}

func TestGet(t *testing.T) {
	c := New(10, 20, 30)

	v, ok := c.Get(1)

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}

	v, ok = c.Get(-1)

	if !ok || v != 30 {
		t.Errorf("expected 30 for negative index, got %d", v)
	}

	_, ok = c.Get(10)

	if ok {
		t.Error("expected not found")
	}

	v, _ = c.Get(10, 99)

	if v != 99 {
		t.Errorf("expected default 99, got %d", v)
	}
}

func TestGetOrPut(t *testing.T) {
	c := New(1, 2, 3)
	v := c.GetOrPut(1, 99)

	if v != 2 {
		t.Errorf("expected existing value 2, got %d", v)
	}

	v = c.GetOrPut(10, 99)

	if v != 99 {
		t.Errorf("expected default 99, got %d", v)
	}

	if c.Count() != 4 {
		t.Errorf("expected 4 items after put, got %d", c.Count())
	}
}

func TestPut(t *testing.T) {
	c := New(1, 2, 3)
	c.Put(1, 99)
	v, _ := c.Get(1)

	if v != 99 {
		t.Errorf("expected 99, got %d", v)
	}
}

func TestPull(t *testing.T) {
	c := New(10, 20, 30)
	v, ok := c.Pull(1)

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}

	if c.Count() != 2 {
		t.Errorf("expected 2 items, got %d", c.Count())
	}
}

func TestContains(t *testing.T) {
	c := New(1, 2, 3, 4, 5)

	if !c.Contains(func(item int, _ int) bool { return item == 3 }) {
		t.Error("expected to contain 3")
	}

	if c.Contains(func(item int, _ int) bool { return item == 10 }) {
		t.Error("expected not to contain 10")
	}
}

func TestDoesntContain(t *testing.T) {
	c := New(1, 2, 3)

	if !c.DoesntContain(func(item int, _ int) bool { return item == 10 }) {
		t.Error("expected not to contain 10")
	}
}

func TestSearch(t *testing.T) {
	c := New(10, 20, 30)
	idx, ok := c.Search(func(item int, _ int) bool { return item == 20 })

	if !ok || idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}

	_, ok = c.Search(func(item int, _ int) bool { return item == 99 })

	if ok {
		t.Error("expected not found")
	}
}

func TestBefore(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	v, ok := c.Before(func(item int, _ int) bool { return item == 3 })

	if !ok || v != 2 {
		t.Errorf("expected 2, got %d", v)
	}

	_, ok = c.Before(func(item int, _ int) bool { return item == 1 })

	if ok {
		t.Error("expected not found for first element")
	}
}

func TestAfter(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	v, ok := c.After(func(item int, _ int) bool { return item == 3 })

	if !ok || v != 4 {
		t.Errorf("expected 4, got %d", v)
	}

	_, ok = c.After(func(item int, _ int) bool { return item == 5 })

	if ok {
		t.Error("expected not found for last element")
	}
}

func TestPush(t *testing.T) {
	c := New(1, 2)
	c.Push(3, 4)

	if c.Count() != 4 {
		t.Errorf("expected 4, got %d", c.Count())
	}
}

func TestPrepend(t *testing.T) {
	c := New(2, 3)
	c.Prepend(1)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(c.All(), expected) {
		t.Errorf("expected %v, got %v", expected, c.All())
	}
}

func TestPop(t *testing.T) {
	c := New(1, 2, 3)
	v, ok := c.Pop()

	if !ok || v != 3 {
		t.Errorf("expected 3, got %d", v)
	}

	if c.Count() != 2 {
		t.Errorf("expected 2 items, got %d", c.Count())
	}
}

func TestPopMany(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	popped := c.PopMany(2)

	if !reflect.DeepEqual(popped.All(), []int{4, 5}) {
		t.Errorf("expected [4 5], got %v", popped.All())
	}

	if c.Count() != 3 {
		t.Errorf("expected 3 items, got %d", c.Count())
	}
}

func TestShift(t *testing.T) {
	c := New(1, 2, 3)
	v, ok := c.Shift()

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	if c.Count() != 2 {
		t.Errorf("expected 2, got %d", c.Count())
	}
}

func TestShiftMany(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	shifted := c.ShiftMany(2)

	if !reflect.DeepEqual(shifted.All(), []int{1, 2}) {
		t.Errorf("expected [1 2], got %v", shifted.All())
	}

	if c.Count() != 3 {
		t.Errorf("expected 3 items, got %d", c.Count())
	}
}

func TestEach(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	sum := 0
	c.Each(func(item int, _ int) bool {
		sum += item

		return true
	})

	if sum != 15 {
		t.Errorf("expected sum 15, got %d", sum)
	}

	// Test early break
	count := 0
	c.Each(func(item int, _ int) bool {
		count++

		return item < 3
	})

	if count != 3 {
		t.Errorf("expected 3 iterations, got %d", count)
	}
}

func TestFilter(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	evens := c.Filter(func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{2, 4}

	if !reflect.DeepEqual(evens.All(), expected) {
		t.Errorf("expected %v, got %v", expected, evens.All())
	}
}

func TestReject(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	odds := c.Reject(func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{1, 3, 5}

	if !reflect.DeepEqual(odds.All(), expected) {
		t.Errorf("expected %v, got %v", expected, odds.All())
	}
}

func TestMap(t *testing.T) {
	c := New(1, 2, 3)
	doubled := Map(c, func(item int, _ int) int { return item * 2 })
	expected := []int{2, 4, 6}

	if !reflect.DeepEqual(doubled.All(), expected) {
		t.Errorf("expected %v, got %v", expected, doubled.All())
	}
}

func TestMapDifferentTypes(t *testing.T) {
	c := New(1, 2, 3)
	strs := Map(c, func(item int, _ int) string {
		return string(rune('a' + item - 1))
	})
	expected := []string{"a", "b", "c"}

	if !reflect.DeepEqual(strs.All(), expected) {
		t.Errorf("expected %v, got %v", expected, strs.All())
	}
}

func TestTransform(t *testing.T) {
	c := New(1, 2, 3)
	c.Transform(func(item int, _ int) int { return item * 10 })
	expected := []int{10, 20, 30}

	if !reflect.DeepEqual(c.All(), expected) {
		t.Errorf("expected %v, got %v", expected, c.All())
	}
}

func TestFlatMap(t *testing.T) {
	c := New(1, 2, 3)
	result := FlatMap(c, func(item int, _ int) []int { return []int{item, item * 10} })
	expected := []int{1, 10, 2, 20, 3, 30}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestReduce(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	sum := Reduce(c, func(carry int, item int, _ int) int { return carry + item }, 0)

	if sum != 15 {
		t.Errorf("expected 15, got %d", sum)
	}
}

func TestChunk(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	chunks := c.Chunk(2)

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

func TestChunkWhile(t *testing.T) {
	c := New(1, 1, 2, 2, 3)
	chunks := c.ChunkWhile(func(item int, _ int, current []int) bool {
		return item == current[len(current)-1]
	})

	if len(chunks) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(chunks))
	}
}

func TestSplit(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	groups := c.Split(3)

	if len(groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(groups))
	}
}

func TestSliding(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	windows := c.Sliding(3)

	if len(windows) != 3 {
		t.Errorf("expected 3 windows, got %d", len(windows))
	}

	if !reflect.DeepEqual(windows[0], []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", windows[0])
	}

	// With step
	windows2 := c.Sliding(2, 2)

	if len(windows2) != 2 {
		t.Errorf("expected 2 windows, got %d", len(windows2))
	}
}

func TestSlice(t *testing.T) {
	c := New(1, 2, 3, 4, 5)

	s := c.Slice(2)
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(s.All(), expected) {
		t.Errorf("expected %v, got %v", expected, s.All())
	}

	s2 := c.Slice(1, 2)
	expected2 := []int{2, 3}

	if !reflect.DeepEqual(s2.All(), expected2) {
		t.Errorf("expected %v, got %v", expected2, s2.All())
	}

	// Negative offset
	s3 := c.Slice(-2)
	expected3 := []int{4, 5}

	if !reflect.DeepEqual(s3.All(), expected3) {
		t.Errorf("expected %v, got %v", expected3, s3.All())
	}
}

func TestTake(t *testing.T) {
	c := New(1, 2, 3, 4, 5)

	taken := c.Take(3)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(taken.All(), expected) {
		t.Errorf("expected %v, got %v", expected, taken.All())
	}

	// Negative take
	taken2 := c.Take(-2)
	expected2 := []int{4, 5}

	if !reflect.DeepEqual(taken2.All(), expected2) {
		t.Errorf("expected %v, got %v", expected2, taken2.All())
	}
}

func TestTakeUntil(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := c.TakeUntil(func(item int, _ int) bool { return item == 4 })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestTakeWhile(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := c.TakeWhile(func(item int, _ int) bool { return item < 4 })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSkip(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := c.Skip(2)
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSkipUntil(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := c.SkipUntil(func(item int, _ int) bool { return item == 3 })
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSkipWhile(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := c.SkipWhile(func(item int, _ int) bool { return item < 3 })
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestNth(t *testing.T) {
	c := New(1, 2, 3, 4, 5, 6, 7, 8)
	result := c.Nth(3)
	expected := []int{1, 4, 7}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}

	result2 := c.Nth(3, 1)
	expected2 := []int{2, 5, 8}

	if !reflect.DeepEqual(result2.All(), expected2) {
		t.Errorf("expected %v, got %v", expected2, result2.All())
	}
}

func TestForPage(t *testing.T) {
	c := New(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	result := c.ForPage(2, 3)
	expected := []int{4, 5, 6}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestReverse(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := c.Reverse()
	expected := []int{5, 4, 3, 2, 1}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestShuffle(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	shuffled := c.Shuffle()

	if shuffled.Count() != 5 {
		t.Errorf("expected 5 items, got %d", shuffled.Count())
	}
}

func TestSort(t *testing.T) {
	c := New(3, 1, 4, 1, 5, 9, 2, 6)
	sorted := c.Sort(func(a, b int) bool { return a < b })
	expected := []int{1, 1, 2, 3, 4, 5, 6, 9}

	if !reflect.DeepEqual(sorted.All(), expected) {
		t.Errorf("expected %v, got %v", expected, sorted.All())
	}
}

func TestSortBy(t *testing.T) {
	type Item struct {
		Name string
		Age  int
	}

	c := New(Item{"Bob", 30}, Item{"Alice", 25}, Item{"Charlie", 35})
	sorted := SortBy(c, func(item Item) int { return item.Age })

	if sorted.All()[0].Name != "Alice" {
		t.Errorf("expected Alice first, got %s", sorted.All()[0].Name)
	}
}

func TestSortByDesc(t *testing.T) {
	type Item struct {
		Name string
		Age  int
	}

	c := New(Item{"Bob", 30}, Item{"Alice", 25}, Item{"Charlie", 35})
	sorted := SortByDesc(c, func(item Item) int { return item.Age })

	if sorted.All()[0].Name != "Charlie" {
		t.Errorf("expected Charlie first, got %s", sorted.All()[0].Name)
	}
}

func TestUnique(t *testing.T) {
	c := New(1, 2, 2, 3, 3, 3, 4)
	unique := Unique(c, func(item int) int { return item })
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(unique.All(), expected) {
		t.Errorf("expected %v, got %v", expected, unique.All())
	}
}

func TestDuplicates(t *testing.T) {
	c := New(1, 2, 2, 3, 3, 3)
	dups := Duplicates(c, func(item int) int { return item })
	expected := []int{2, 3, 3}

	if !reflect.DeepEqual(dups.All(), expected) {
		t.Errorf("expected %v, got %v", expected, dups.All())
	}
}

func TestEvery(t *testing.T) {
	c := New(2, 4, 6, 8)

	if !c.Every(func(item int, _ int) bool { return item%2 == 0 }) {
		t.Error("expected all even")
	}

	if c.Every(func(item int, _ int) bool { return item > 5 }) {
		t.Error("expected not all > 5")
	}
}

func TestPartition(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	pass, fail := c.Partition(func(item int, _ int) bool { return item > 3 })

	if !reflect.DeepEqual(pass.All(), []int{4, 5}) {
		t.Errorf("expected [4 5], got %v", pass.All())
	}

	if !reflect.DeepEqual(fail.All(), []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", fail.All())
	}
}

func TestConcat(t *testing.T) {
	c := New(1, 2, 3)
	result := c.Concat([]int{4, 5})
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestPad(t *testing.T) {
	c := New(1, 2, 3)

	padded := c.Pad(5, 0)
	expected := []int{1, 2, 3, 0, 0}

	if !reflect.DeepEqual(padded.All(), expected) {
		t.Errorf("expected %v, got %v", expected, padded.All())
	}

	paddedLeft := c.Pad(-5, 0)
	expected2 := []int{0, 0, 1, 2, 3}

	if !reflect.DeepEqual(paddedLeft.All(), expected2) {
		t.Errorf("expected %v, got %v", expected2, paddedLeft.All())
	}
}

func TestMultiply(t *testing.T) {
	c := New(1, 2, 3)
	result := c.Multiply(3)
	expected := []int{1, 2, 3, 1, 2, 3, 1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestForget(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	c.Forget(2)
	expected := []int{1, 2, 4, 5}

	if !reflect.DeepEqual(c.All(), expected) {
		t.Errorf("expected %v, got %v", expected, c.All())
	}
}

func TestImplode(t *testing.T) {
	c := New("a", "b", "c")
	result := c.Implode(", ")

	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}
}

func TestJoin(t *testing.T) {
	c := New("a", "b", "c")

	result := c.Join(", ")

	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}

	result2 := c.Join(", ", " and ")

	if result2 != "a, b and c" {
		t.Errorf("expected 'a, b and c', got '%s'", result2)
	}
}

func TestWhen(t *testing.T) {
	c := New(1, 2, 3)

	result := c.When(true, func(c *Collection[int]) *Collection[int] {
		return c.Push(4)
	})

	if result.Count() != 4 {
		t.Errorf("expected 4 items, got %d", result.Count())
	}

	c2 := New(1, 2, 3)
	result2 := c2.When(false, func(c *Collection[int]) *Collection[int] {
		return c.Push(4)
	})

	if result2.Count() != 3 {
		t.Errorf("expected 3 items, got %d", result2.Count())
	}
}

func TestUnless(t *testing.T) {
	c := New(1, 2, 3)
	result := c.Unless(false, func(c *Collection[int]) *Collection[int] {
		return c.Push(4)
	})

	if result.Count() != 4 {
		t.Errorf("expected 4, got %d", result.Count())
	}
}

func TestDiff(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := Diff(c, []int{2, 4})
	expected := []int{1, 3, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestIntersect(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := Intersect(c, []int{2, 4, 6})
	expected := []int{2, 4}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestDiffUsing(t *testing.T) {
	c := New(1, 2, 3)
	result := c.DiffUsing([]int{2, 3}, func(a, b int) bool { return a == b })

	if !reflect.DeepEqual(result.All(), []int{1}) {
		t.Errorf("expected [1], got %v", result.All())
	}
}

func TestIntersectUsing(t *testing.T) {
	c := New(1, 2, 3)
	result := c.IntersectUsing([]int{2, 3, 4}, func(a, b int) bool { return a == b })

	if !reflect.DeepEqual(result.All(), []int{2, 3}) {
		t.Errorf("expected [2 3], got %v", result.All())
	}
}

func TestZip(t *testing.T) {
	c := New(1, 2, 3)
	result := Zip(c, []int{10, 20, 30})

	if len(result.All()) != 3 {
		t.Errorf("expected 3 pairs, got %d", len(result.All()))
	}

	if !reflect.DeepEqual(result.All()[0], []int{1, 10}) {
		t.Errorf("expected [1 10], got %v", result.All()[0])
	}
}

func TestCrossJoin(t *testing.T) {
	c := New(1, 2)
	result := CrossJoin(c, []int{10, 20})

	if len(result.All()) != 4 {
		t.Errorf("expected 4 combinations, got %d", len(result.All()))
	}
}

func TestCombine(t *testing.T) {
	keys := New("name", "age")
	result := Combine(keys, []string{"John", "30"})

	if result.Count() != 2 {
		t.Errorf("expected 2 pairs, got %d", result.Count())
	}

	if result.All()[0].Key != "name" || result.All()[0].Value != "John" {
		t.Error("unexpected pair values")
	}
}

func TestCollapse(t *testing.T) {
	c := New([]int{1, 2}, []int{3, 4}, []int{5})
	result := Collapse(c)
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestPluck(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	c := New(User{"Alice", 25}, User{"Bob", 30})
	names := Pluck(c, func(u User) string { return u.Name })
	expected := []string{"Alice", "Bob"}

	if !reflect.DeepEqual(names.All(), expected) {
		t.Errorf("expected %v, got %v", expected, names.All())
	}
}

func TestGroupBy(t *testing.T) {
	c := New(1, 2, 3, 4, 5, 6)
	groups := GroupBy(c, func(item int) string {
		if item%2 == 0 {
			return "even"
		}

		return "odd"
	})

	if groups["even"].Count() != 3 {
		t.Errorf("expected 3 evens, got %d", groups["even"].Count())
	}

	if groups["odd"].Count() != 3 {
		t.Errorf("expected 3 odds, got %d", groups["odd"].Count())
	}
}

func TestKeyBy(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	c := New(User{1, "Alice"}, User{2, "Bob"})
	keyed := KeyBy(c, func(u User) int { return u.ID })

	if keyed[1].Name != "Alice" {
		t.Error("expected Alice at key 1")
	}
}

func TestCountBy(t *testing.T) {
	c := New("apple", "banana", "apple", "cherry", "banana", "apple")
	counts := CountBy(c, func(item string) string { return item })

	if counts["apple"] != 3 {
		t.Errorf("expected 3 apples, got %d", counts["apple"])
	}
}

func TestMapToDictionary(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	dict := MapToDictionary(c, func(item int) (string, int) {
		if item%2 == 0 {
			return "even", item
		}

		return "odd", item
	})

	if len(dict["even"]) != 2 {
		t.Errorf("expected 2 evens, got %d", len(dict["even"]))
	}
}

func TestMapWithKeys(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	c := New(User{1, "Alice"}, User{2, "Bob"})
	result := MapWithKeys(c, func(u User) (int, string) { return u.ID, u.Name })

	if result[1] != "Alice" {
		t.Error("expected Alice at key 1")
	}
}

func TestOnly(t *testing.T) {
	c := New(10, 20, 30, 40, 50)
	result := c.Only(1, 3)
	expected := []int{20, 40}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestExcept(t *testing.T) {
	c := New(10, 20, 30, 40, 50)
	result := c.Except(1, 3)
	expected := []int{10, 30, 50}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestHas(t *testing.T) {
	c := New(1, 2, 3)

	if !c.Has(0) {
		t.Error("expected index 0 to exist")
	}

	if !c.Has(-1) {
		t.Error("expected negative index to work")
	}

	if c.Has(10) {
		t.Error("expected index 10 to not exist")
	}
}

func TestSum(t *testing.T) {
	c := New(1, 2, 3, 4, 5)

	if Sum(c) != 15 {
		t.Errorf("expected 15, got %d", Sum(c))
	}
}

func TestSumBy(t *testing.T) {
	type Item struct {
		Price float64
	}

	c := New(Item{10.5}, Item{20.5}, Item{30.0})
	result := SumBy(c, func(i Item) float64 { return i.Price })

	if result != 61.0 {
		t.Errorf("expected 61.0, got %f", result)
	}
}

func TestAvg(t *testing.T) {
	c := New(1.0, 2.0, 3.0, 4.0, 5.0)
	result := Avg(c)

	if result != 3.0 {
		t.Errorf("expected 3.0, got %f", result)
	}
}

func TestMin(t *testing.T) {
	c := New(3, 1, 4, 1, 5, 9)
	v, ok := Min(c)

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
}

func TestMax(t *testing.T) {
	c := New(3, 1, 4, 1, 5, 9)
	v, ok := Max(c)

	if !ok || v != 9 {
		t.Errorf("expected 9, got %d", v)
	}
}

func TestMedian(t *testing.T) {
	c := New(1.0, 2.0, 3.0, 4.0, 5.0)
	result := Median(c)

	if result != 3.0 {
		t.Errorf("expected 3.0, got %f", result)
	}

	c2 := New(1.0, 2.0, 3.0, 4.0)
	result2 := Median(c2)

	if result2 != 2.5 {
		t.Errorf("expected 2.5, got %f", result2)
	}
}

func TestMode(t *testing.T) {
	c := New(1, 2, 2, 3, 3, 3)
	result := Mode(c)

	if len(result) != 1 || result[0] != 3 {
		t.Errorf("expected [3], got %v", result)
	}
}

func TestToJSON(t *testing.T) {
	c := New(1, 2, 3)
	b, err := c.ToJSON()

	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "[1,2,3]" {
		t.Errorf("expected [1,2,3], got %s", string(b))
	}
}

func TestMarshalJSON(t *testing.T) {
	c := New("a", "b", "c")
	b, err := json.Marshal(c)

	if err != nil {
		t.Fatal(err)
	}

	expected := `["a","b","c"]`

	if string(b) != expected {
		t.Errorf("expected %s, got %s", expected, string(b))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	c := Empty[int]()
	err := json.Unmarshal([]byte("[1,2,3]"), c)

	if err != nil {
		t.Fatal(err)
	}

	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(c.All(), expected) {
		t.Errorf("expected %v, got %v", expected, c.All())
	}
}

func TestCopy(t *testing.T) {
	c := New(1, 2, 3)
	c2 := c.Copy()
	c.Push(4)

	if c2.Count() != 3 {
		t.Error("copy should not be affected by changes to original")
	}
}

func TestSplice(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	removed := c.Splice(1, 2)

	if !reflect.DeepEqual(removed.All(), []int{2, 3}) {
		t.Errorf("expected [2 3], got %v", removed.All())
	}

	if !reflect.DeepEqual(c.All(), []int{1, 4, 5}) {
		t.Errorf("expected [1 4 5], got %v", c.All())
	}
}

func TestSpliceReplace(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	removed := c.SpliceReplace(1, 2, []int{20, 30, 40})

	if !reflect.DeepEqual(removed.All(), []int{2, 3}) {
		t.Errorf("expected [2 3] removed, got %v", removed.All())
	}

	if !reflect.DeepEqual(c.All(), []int{1, 20, 30, 40, 4, 5}) {
		t.Errorf("expected [1 20 30 40 4 5], got %v", c.All())
	}
}

func TestWhere(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := c.Where(func(item int) bool { return item > 3 })
	expected := []int{4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestRandom(t *testing.T) {
	c := New(1, 2, 3, 4, 5)
	result := c.Random(2)

	if result.Count() != 2 {
		t.Errorf("expected 2 items, got %d", result.Count())
	}
}

func TestTap(t *testing.T) {
	c := New(1, 2, 3)
	tapped := false
	result := c.Tap(func(c *Collection[int]) {
		tapped = true
	})

	if !tapped {
		t.Error("expected tap callback to be called")
	}

	if result != c {
		t.Error("expected same collection returned")
	}
}

func TestPipe(t *testing.T) {
	c := New(1, 2, 3)
	result := Pipe(c, func(c *Collection[int]) int {
		return Sum(c)
	})

	if result != 6 {
		t.Errorf("expected 6, got %d", result)
	}
}

func TestPipeThrough(t *testing.T) {
	c := New(1, 2, 3)
	result := PipeThrough(c,
		func(c *Collection[int]) *Collection[int] {
			return c.Filter(func(item int, _ int) bool { return item > 1 })
		},
		func(c *Collection[int]) *Collection[int] {
			return c.Push(10)
		},
	)
	expected := []int{2, 3, 10}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestMinBy(t *testing.T) {
	type Item struct {
		Name string
		Age  int
	}

	c := New(Item{"Bob", 30}, Item{"Alice", 25}, Item{"Charlie", 35})
	v, ok := MinBy(c, func(item Item) int { return item.Age })

	if !ok || v.Name != "Alice" {
		t.Errorf("expected Alice, got %s", v.Name)
	}
}

func TestMaxBy(t *testing.T) {
	type Item struct {
		Name string
		Age  int
	}

	c := New(Item{"Bob", 30}, Item{"Alice", 25}, Item{"Charlie", 35})
	v, ok := MaxBy(c, func(item Item) int { return item.Age })

	if !ok || v.Name != "Charlie" {
		t.Errorf("expected Charlie, got %s", v.Name)
	}
}

func TestMapInto(t *testing.T) {
	c := New(1, 2, 3)
	result := MapInto(c, func(v int) string {
		return string(rune('A' + v - 1))
	})
	expected := []string{"A", "B", "C"}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestEnsure(t *testing.T) {
	c := New(2, 4, 6)
	err := c.Ensure(func(v int) bool { return v%2 == 0 })

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = c.Ensure(func(v int) bool { return v > 10 })

	if err == nil {
		t.Error("expected error")
	}
}

func TestString(t *testing.T) {
	c := New(1, 2, 3)
	s := c.String()

	if s != "[1,2,3]" {
		t.Errorf("expected [1,2,3], got %s", s)
	}
}

func TestDot(t *testing.T) {
	c := New(1, 2, 3)
	d := c.Dot()

	if !reflect.DeepEqual(d.All(), c.All()) {
		t.Error("expected copy")
	}
}

func TestSortDesc(t *testing.T) {
	c := New(1, 3, 2, 5, 4)
	sorted := c.SortDesc(func(a, b int) bool { return a < b })
	expected := []int{5, 4, 3, 2, 1}

	if !reflect.DeepEqual(sorted.All(), expected) {
		t.Errorf("expected %v, got %v", expected, sorted.All())
	}
}

func TestValues(t *testing.T) {
	c := New(10, 20, 30)
	v := c.Values()

	if !reflect.DeepEqual(v.All(), c.All()) {
		t.Error("expected same values")
	}
	// Ensure it's a copy
	v.Push(40)

	if c.Count() != 3 {
		t.Error("values should return a copy")
	}
}

func TestToSlice(t *testing.T) {
	c := New(1, 2, 3)
	s := c.ToSlice()
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("expected %v, got %v", expected, s)
	}
}

func TestTapEach(t *testing.T) {
	c := New(1, 2, 3)
	sum := 0
	result := c.TapEach(func(item int, _ int) {
		sum += item
	})

	if sum != 6 {
		t.Errorf("expected 6, got %d", sum)
	}

	if result != c {
		t.Error("expected same collection")
	}
}

func TestSplitIn(t *testing.T) {
	c := New(1, 2, 3, 4, 5, 6, 7)
	groups := c.SplitIn(3)
	total := 0

	for _, g := range groups {
		total += len(g)
	}

	if total != 7 {
		t.Errorf("expected 7 total items, got %d", total)
	}
}

// Test that Mode handles ties
func TestModeTie(t *testing.T) {
	c := New(1, 1, 2, 2, 3)
	result := Mode(c)
	sort.Ints(result)
	expected := []int{1, 2}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestWhenEmpty(t *testing.T) {
	empty := Empty[int]()
	result := empty.WhenEmpty(func(c *Collection[int]) *Collection[int] {
		return c.Push(1, 2, 3)
	})

	if result.Count() != 3 {
		t.Errorf("expected 3, got %d", result.Count())
	}
}

func TestWhenNotEmpty(t *testing.T) {
	c := New(1, 2, 3)
	called := false
	c.WhenNotEmpty(func(c *Collection[int]) *Collection[int] {
		called = true

		return c
	})

	if !called {
		t.Error("expected callback to be called")
	}
}
