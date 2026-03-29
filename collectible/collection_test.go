package collectible

import (
	"reflect"
	"testing"

	"github.com/gocanto/collection/support"
)

func TestNew(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})

	if m.Count() != 3 {
		t.Errorf("expected 3, got %d", m.Count())
	}
}

func TestFromPairs(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	v, ok := m.Get("a")

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
}

func TestMapGet(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	v, ok := m.Get("a")

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	_, ok = m.Get("z")

	if ok {
		t.Error("expected not found")
	}

	v, _ = m.Get("z", 99)

	if v != 99 {
		t.Errorf("expected default 99, got %d", v)
	}
}

func TestMapGetOrPut(t *testing.T) {
	m := New(map[string]int{"a": 1})

	v := m.GetOrPut("a", 99)

	if v != 1 {
		t.Errorf("expected existing 1, got %d", v)
	}

	v = m.GetOrPut("b", 99)

	if v != 99 {
		t.Errorf("expected default 99, got %d", v)
	}

	if m.Count() != 2 {
		t.Errorf("expected 2 items, got %d", m.Count())
	}
}

func TestMapHas(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	if !m.Has("a") {
		t.Error("expected to have 'a'")
	}

	if m.Has("z") {
		t.Error("expected not to have 'z'")
	}
}

func TestMapHasAny(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	if !m.HasAny("z", "a") {
		t.Error("expected HasAny to return true")
	}

	if m.HasAny("x", "y", "z") {
		t.Error("expected HasAny to return false")
	}
}

func TestMapPut(t *testing.T) {
	m := New(map[string]int{})
	m.Put("a", 1)
	v, ok := m.Get("a")

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
}

func TestMapPull(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	v, ok := m.Pull("a")

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	if m.Has("a") {
		t.Error("expected 'a' to be removed")
	}
}

func TestMapForget(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	m.Forget("a", "c")

	if m.Count() != 1 {
		t.Errorf("expected 1, got %d", m.Count())
	}
}

func TestMapOnly(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	result := m.Only("a", "c")

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}

	if !result.Has("a") || !result.Has("c") {
		t.Error("expected to have 'a' and 'c'")
	}
}

func TestMapExcept(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	result := m.Except("b")

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}

	if result.Has("b") {
		t.Error("expected 'b' to be excluded")
	}
}

func TestMapKeys(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	keys := m.Keys()

	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestMapValues(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	values := m.Values()

	if len(values) != 2 {
		t.Errorf("expected 2 values, got %d", len(values))
	}
}

func TestMapContains(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	if !m.Contains(func(v int, k string) bool { return v == 2 }) {
		t.Error("expected to contain value 2")
	}

	if m.Contains(func(v int, k string) bool { return v == 99 }) {
		t.Error("expected not to contain 99")
	}
}

func TestMapFirst(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	v, ok := m.First()

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
}

func TestMapLast(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	v, ok := m.Last()

	if !ok || v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
}

func TestMapFilter(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	filtered := m.Filter(func(v int, k string) bool { return v > 1 })

	if filtered.Count() != 2 {
		t.Errorf("expected 2, got %d", filtered.Count())
	}
}

func TestMapReject(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	rejected := m.Reject(func(v int, k string) bool { return v > 1 })

	if rejected.Count() != 1 {
		t.Errorf("expected 1, got %d", rejected.Count())
	}
}

func TestMapValues_Transform(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	result := MapValues(m, func(v int, k string) int { return v * 10 })
	v, _ := result.Get("a")

	if v != 10 {
		t.Errorf("expected 10, got %d", v)
	}
}

func TestMapEvery(t *testing.T) {
	m := New(map[string]int{"a": 2, "b": 4})

	if !m.Every(func(v int, k string) bool { return v%2 == 0 }) {
		t.Error("expected all even")
	}
}

func TestMapPartition(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	pass, fail := m.Partition(func(v int, k string) bool { return v > 1 })

	if pass.Count() != 2 {
		t.Errorf("expected 2 passing, got %d", pass.Count())
	}

	if fail.Count() != 1 {
		t.Errorf("expected 1 failing, got %d", fail.Count())
	}
}

func TestMapMerge(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	result := m.Merge(map[string]int{"b": 20, "c": 3})
	v, _ := result.Get("b")

	if v != 20 {
		t.Errorf("expected 20 (overwritten), got %d", v)
	}

	if result.Count() != 3 {
		t.Errorf("expected 3, got %d", result.Count())
	}
}

func TestMapUnion(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	result := m.Union(map[string]int{"b": 20, "c": 3})
	v, _ := result.Get("b")

	if v != 2 {
		t.Errorf("expected 2 (not overwritten), got %d", v)
	}

	if result.Count() != 3 {
		t.Errorf("expected 3, got %d", result.Count())
	}
}

func TestMapDiffKeys(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	result := m.DiffKeys(map[string]int{"b": 99})

	if result.Has("b") {
		t.Error("expected 'b' to be excluded")
	}

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}
}

func TestMapIntersectByKeys(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	result := m.IntersectByKeys(map[string]int{"a": 99, "c": 99})

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}
}

func TestFlip(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	flipped := Flip(m)
	v, ok := flipped.Get(1)

	if !ok || v != "a" {
		t.Errorf("expected 'a', got '%s'", v)
	}
}

func TestSortKeys(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "c", Value: 3},
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	sorted := SortKeys(m)
	keys := sorted.Keys()
	expected := []string{"a", "b", "c"}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("expected %v, got %v", expected, keys)
	}
}

func TestSortKeysDesc(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "c", Value: 3},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	sorted := SortKeysDesc(m)
	keys := sorted.Keys()
	expected := []string{"c", "b", "a"}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("expected %v, got %v", expected, keys)
	}
}

func TestMapImplode(t *testing.T) {
	m := FromPairs(
		support.Pair[string, string]{Key: "a", Value: "hello"},
		support.Pair[string, string]{Key: "b", Value: "world"},
	)
	result := m.Implode(", ")

	if result != "hello, world" {
		t.Errorf("expected 'hello, world', got '%s'", result)
	}
}

func TestMapToJSON(t *testing.T) {
	m := New(map[string]int{"a": 1})
	b, err := m.ToJSON()

	if err != nil {
		t.Fatal(err)
	}

	if string(b) != `{"a":1}` {
		t.Errorf("expected {\"a\":1}, got %s", string(b))
	}
}

func TestMapCopy(t *testing.T) {
	m := New(map[string]int{"a": 1})
	m2 := m.Copy()
	m.Put("b", 2)

	if m2.Has("b") {
		t.Error("copy should not be affected")
	}
}

func TestMapToPairs(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	pairs := m.ToPairs()

	if len(pairs) != 2 {
		t.Errorf("expected 2 pairs, got %d", len(pairs))
	}
}

func TestMapWhen(t *testing.T) {
	m := New(map[string]int{"a": 1})
	result := m.When(true, func(mc *Collection[string, int]) *Collection[string, int] {
		mc.Put("b", 2)

		return mc
	})

	if !result.Has("b") {
		t.Error("expected 'b' to be added")
	}
}

func TestMapSearch(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	key, ok := m.Search(func(v int, k string) bool { return v == 2 })

	if !ok || key != "b" {
		t.Errorf("expected 'b', got '%s'", key)
	}
}

func TestMapIsEmpty(t *testing.T) {
	m := New(map[string]int{})

	if !m.IsEmpty() {
		t.Error("expected empty")
	}

	if m.IsNotEmpty() {
		t.Error("expected empty")
	}
}

func TestMapContainsOneItem(t *testing.T) {
	m := New(map[string]int{"a": 1})

	if !m.ContainsOneItem() {
		t.Error("expected true")
	}
}

func TestMapContainsManyItems(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	if !m.ContainsManyItems() {
		t.Error("expected true")
	}
}

func TestMapEach(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	sum := 0
	m.Each(func(v int, k string) bool {
		sum += v

		return true
	})

	if sum != 3 {
		t.Errorf("expected 3, got %d", sum)
	}
}

func TestMapDoesntContain(t *testing.T) {
	m := New(map[string]int{"a": 1})

	if !m.DoesntContain(func(v int, k string) bool { return v == 99 }) {
		t.Error("expected true")
	}
}

func TestMapSortKeysUsing(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "banana", Value: 1},
		support.Pair[string, int]{Key: "apple", Value: 2},
		support.Pair[string, int]{Key: "cherry", Value: 3},
	)
	sorted := m.SortKeysUsing(func(a, b string) bool { return a < b })
	keys := sorted.Keys()
	expected := []string{"apple", "banana", "cherry"}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("expected %v, got %v", expected, keys)
	}
}

func TestMapReplace(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	result := m.Replace(map[string]int{"b": 20, "c": 3})
	v, _ := result.Get("b")

	if v != 20 {
		t.Errorf("expected 20, got %d", v)
	}
}
