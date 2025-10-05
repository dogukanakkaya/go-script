package environment

import (
	"testing"
)

func TestEnvironmentGet(t *testing.T) {
	env := New(nil)
	env.Set("x", 42.0)

	val, ok := env.Get("x")
	if !ok {
		t.Error("Variable 'x' should exist")
	}
	if num, ok := val.(float64); !ok || num != 42.0 {
		t.Errorf("Expected 42.0, got %v", val)
	}

	_, ok = env.Get("y")
	if ok {
		t.Error("Variable 'y' should not exist")
	}
}

func TestEnvironmentSet(t *testing.T) {
	env := New(nil)

	env.Set("x", 10.0)
	val, ok := env.Get("x")
	if !ok || val.(float64) != 10.0 {
		t.Error("Setting variable failed")
	}

	env.Set("x", 20.0)
	val, ok = env.Get("x")
	if !ok || val.(float64) != 20.0 {
		t.Error("Updating variable failed")
	}
}

func TestEnvironmentScoping(t *testing.T) {
	outer := New(nil)
	outer.Set("x", 10.0)

	inner := New(outer)
	inner.Set("y", 20.0)

	// Inner scope can access outer variable
	val, ok := inner.Get("x")
	if !ok || val.(float64) != 10.0 {
		t.Error("Inner scope should access outer variable")
	}

	// Inner scope has its own variable
	val, ok = inner.Get("y")
	if !ok || val.(float64) != 20.0 {
		t.Error("Inner scope should have its own variable")
	}

	// Outer scope cannot access inner variable
	_, ok = outer.Get("y")
	if ok {
		t.Error("Outer scope should not access inner variable")
	}
}

func TestEnvironmentUpdate(t *testing.T) {
	outer := New(nil)
	outer.Set("x", 10.0)

	inner := New(outer)
	inner.Update("x", 20.0)

	// Update should modify outer scope
	val, ok := outer.Get("x")
	if !ok || val.(float64) != 20.0 {
		t.Error("Update should modify variable in outer scope")
	}
}
