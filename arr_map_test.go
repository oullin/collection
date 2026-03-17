package collection

import (
	"reflect"
	"strings"
	"testing"
)

func TestArrGet(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
		},
	}
	v := ArrGet(data, "user.name")
	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}
}

func TestArrSet(t *testing.T) {
	data := make(map[string]any)
	ArrSet(data, "user.name", "Alice")
	v := ArrGet(data, "user.name")
	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}
}

func TestArrAdd(t *testing.T) {
	data := map[string]any{"name": "Alice"}
	ArrAdd(data, "name", "Bob") // Should not overwrite
	if data["name"] != "Alice" {
		t.Errorf("expected Alice, got %v", data["name"])
	}
	ArrAdd(data, "email", "alice@test.com")
	if data["email"] != "alice@test.com" {
		t.Errorf("expected email set, got %v", data["email"])
	}
}

func TestArrPull_Map(t *testing.T) {
	data := map[string]any{"name": "Alice", "age": 25}
	v := ArrPull(data, "name")
	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}
	if _, ok := data["name"]; ok {
		t.Error("expected name to be removed")
	}
}

func TestArrForget_Map(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name":  "Alice",
			"email": "alice@test.com",
		},
	}
	ArrForget(data, "user.email")
	if DataHas(data, "user.email") {
		t.Error("expected email to be removed")
	}
}

func TestArrHasMap(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
		},
	}
	if !ArrHasMap(data, "user.name") {
		t.Error("expected true")
	}
	if ArrHasMap(data, "user.email") {
		t.Error("expected false")
	}
}

func TestArrHasAnyMap(t *testing.T) {
	data := map[string]any{
		"name": "Alice",
	}
	if !ArrHasAnyMap(data, "email", "name") {
		t.Error("expected true")
	}
	if ArrHasAnyMap(data, "email", "phone") {
		t.Error("expected false")
	}
}

func TestArrDot(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name":  "Alice",
			"email": "alice@test.com",
		},
		"status": "active",
	}
	dotted := ArrDot(data)
	if dotted["user.name"] != "Alice" {
		t.Errorf("expected Alice, got %v", dotted["user.name"])
	}
	if dotted["status"] != "active" {
		t.Errorf("expected active, got %v", dotted["status"])
	}
}

func TestArrUndot(t *testing.T) {
	dotted := map[string]any{
		"user.name":  "Alice",
		"user.email": "alice@test.com",
	}
	undotted := ArrUndot(dotted)
	v := DataGet(undotted, "user.name")
	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}
}

func TestArrOnlyMap(t *testing.T) {
	data := map[string]any{"a": 1, "b": 2, "c": 3}
	result := ArrOnlyMap(data, "a", "c")
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
	if result["a"] != 1 {
		t.Error("expected a=1")
	}
}

func TestArrExceptMap(t *testing.T) {
	data := map[string]any{"a": 1, "b": 2, "c": 3}
	result := ArrExceptMap(data, "b")
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
	if _, ok := result["b"]; ok {
		t.Error("expected b to be excluded")
	}
}

func TestArrQuery(t *testing.T) {
	data := map[string]any{"name": "Alice", "age": 25}
	result := ArrQuery(data)
	if !strings.Contains(result, "name=Alice") {
		t.Errorf("expected name=Alice in query, got %s", result)
	}
}

func TestArrToCssClasses(t *testing.T) {
	classes := map[string]bool{
		"active":   true,
		"disabled": false,
		"visible":  true,
	}
	result := ArrToCssClasses(classes)
	if !strings.Contains(result, "active") {
		t.Error("expected 'active' in result")
	}
	if strings.Contains(result, "disabled") {
		t.Error("expected 'disabled' to be excluded")
	}
}

func TestArrToCssStyles(t *testing.T) {
	styles := map[string]bool{
		"color: red":       true,
		"display: none":    false,
		"font-weight: bold": true,
	}
	result := ArrToCssStyles(styles)
	if !strings.Contains(result, "color: red;") {
		t.Errorf("expected 'color: red;' in result, got '%s'", result)
	}
}

func TestArrPrependKeysWith(t *testing.T) {
	data := map[string]int{"name": 1, "age": 2}
	result := ArrPrependKeysWith(data, "user_")
	if _, ok := result["user_name"]; !ok {
		t.Error("expected 'user_name' key")
	}
}

func TestArrSortMap(t *testing.T) {
	data := map[string]any{"c": 3, "a": 1, "b": 2}
	result := ArrSortMap(data)
	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d", len(result))
	}
}

func TestArrMapMap(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	result := ArrMapMap(data, func(v int, k string) int { return v * 10 })
	if result["a"] != 10 {
		t.Errorf("expected 10, got %d", result["a"])
	}
}

func TestArrWhereMap(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	result := ArrWhereMap(data, func(v int, k string) bool { return v > 1 })
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestArrMapReplace(t *testing.T) {
	data := map[string]any{"a": 1, "b": 2}
	result := ArrMapReplace(data, map[string]any{"b": 20, "c": 3})
	if result["b"] != 20 {
		t.Errorf("expected 20, got %v", result["b"])
	}
	if result["c"] != 3 {
		t.Errorf("expected 3, got %v", result["c"])
	}
}

func TestArrIsAssoc(t *testing.T) {
	if !ArrIsAssoc(map[string]any{"a": 1}) {
		t.Error("expected true")
	}
	if ArrIsAssoc(map[string]any{}) {
		t.Error("expected false for empty map")
	}
}

func TestArrSortRecursiveMap(t *testing.T) {
	data := map[string]any{
		"b": 2,
		"a": map[string]any{
			"z": 26,
			"y": 25,
		},
	}
	result := ArrSortRecursiveMap(data)
	if result == nil {
		t.Error("expected non-nil result")
	}
	nested, ok := result["a"].(map[string]any)
	if !ok {
		t.Error("expected nested map")
	}
	_ = nested
}

func TestArrDot_WithPrefix(t *testing.T) {
	data := map[string]any{
		"name": "Alice",
	}
	result := ArrDot(data, "user")
	if result["user.name"] != "Alice" {
		t.Errorf("expected Alice at 'user.name', got %v", result["user.name"])
	}
}

func TestArrAccessible(t *testing.T) {
	if !ArrAccessible([]int{1, 2, 3}) {
		t.Error("expected true")
	}
	if ArrAccessible(nil) {
		t.Error("expected false for nil")
	}
}

func TestArrMapWithKeysMap(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	data := map[string]User{
		"a": {1, "Alice"},
		"b": {2, "Bob"},
	}
	result := ArrMapMap(data, func(u User, _ string) string { return u.Name })
	expected := map[string]string{"a": "Alice", "b": "Bob"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
