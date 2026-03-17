package collection

import (
	"testing"
)

func TestHead(t *testing.T) {
	v, ok := Head([]int{1, 2, 3})
	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	_, ok = Head([]int{})
	if ok {
		t.Error("expected not found for empty slice")
	}
}

func TestLast_Helper(t *testing.T) {
	v, ok := Last([]int{1, 2, 3})
	if !ok || v != 3 {
		t.Errorf("expected 3, got %d", v)
	}
}

func TestWhenValue(t *testing.T) {
	v := WhenValue(true, "yes", "no")
	if v != "yes" {
		t.Errorf("expected 'yes', got '%s'", v)
	}

	v = WhenValue(false, "yes", "no")
	if v != "no" {
		t.Errorf("expected 'no', got '%s'", v)
	}
}

func TestWhenFunc(t *testing.T) {
	v := WhenFunc(true, func() int { return 42 }, func() int { return 0 })
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestDataGet(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
			"age":  25,
		},
	}

	v := DataGet(data, "user.name")
	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}

	v = DataGet(data, "user.email", "default@test.com")
	if v != "default@test.com" {
		t.Errorf("expected default, got %v", v)
	}

	v = DataGet(data, "missing.key")
	if v != nil {
		t.Errorf("expected nil, got %v", v)
	}
}

func TestDataSet(t *testing.T) {
	data := make(map[string]any)
	DataSet(data, "user.name", "Alice")

	v := DataGet(data, "user.name")
	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}
}

func TestDataHas(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
		},
	}

	if !DataHas(data, "user.name") {
		t.Error("expected true")
	}
	if DataHas(data, "user.email") {
		t.Error("expected false")
	}
}

func TestDataFill(t *testing.T) {
	data := map[string]any{
		"name": "Alice",
	}

	DataFill(data, "name", "Bob") // Should not overwrite
	if data["name"] != "Alice" {
		t.Errorf("expected Alice (not overwritten), got %v", data["name"])
	}

	DataFill(data, "email", "alice@test.com") // Should fill
	if data["email"] != "alice@test.com" {
		t.Errorf("expected email to be filled, got %v", data["email"])
	}
}

func TestDataForget(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name":  "Alice",
			"email": "alice@test.com",
		},
	}

	DataForget(data, "user.email")
	if DataHas(data, "user.email") {
		t.Error("expected email to be removed")
	}
	if !DataHas(data, "user.name") {
		t.Error("expected name to still exist")
	}
}

func TestValue(t *testing.T) {
	v := Value(42)
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestValueFunc(t *testing.T) {
	v := ValueFunc(func() int { return 42 })
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}
