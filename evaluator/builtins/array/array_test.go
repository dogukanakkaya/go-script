package array

import (
	"go-script/internal"
	"testing"
)

func TestNewArrayReference(t *testing.T) {
	elements := internal.Array{1.0, 2.0, 3.0}
	arr := NewArrayReference(elements)

	if arr == nil {
		t.Fatal("NewArrayReference returned nil")
	}

	if len(*arr.Elements) != 3 {
		t.Errorf("Expected array length 3, got %d", len(*arr.Elements))
	}
}

func TestArrayReferenceGet(t *testing.T) {
	elements := internal.Array{"hello", 42.0, true}
	arr := NewArrayReference(elements)

	tests := []struct {
		index    int
		expected internal.Value
	}{
		{0, "hello"},
		{1, 42.0},
		{2, true},
		{-1, nil},
		{3, nil},
		{100, nil},
	}

	for _, tt := range tests {
		result := arr.Get(tt.index)
		if result != tt.expected {
			t.Errorf("Get(%d) = %v, expected %v", tt.index, result, tt.expected)
		}
	}
}

func TestArrayReferenceSet(t *testing.T) {
	elements := internal.Array{1.0, 2.0, 3.0}
	arr := NewArrayReference(elements)

	arr.Set(1, 99.0)
	if (*arr.Elements)[1] != 99.0 {
		t.Errorf("Set(1, 99.0) failed, got %v", (*arr.Elements)[1])
	}

	arr.Set(0, "first")
	if (*arr.Elements)[0] != "first" {
		t.Errorf("Set(0, 'first') failed, got %v", (*arr.Elements)[0])
	}

	arr.Set(2, "last")
	if (*arr.Elements)[2] != "last" {
		t.Errorf("Set(2, 'last') failed, got %v", (*arr.Elements)[2])
	}

	// Invalid sets (should not panic or modify array)
	originalLen := len(*arr.Elements)
	arr.Set(-1, "negative")
	arr.Set(3, "out of bounds")
	arr.Set(100, "way out")

	if len(*arr.Elements) != originalLen {
		t.Errorf("Invalid Set operations modified array length")
	}
}

func TestArrayReferencePush(t *testing.T) {
	elements := internal.Array{1.0, 2.0, 3.0}
	arr := NewArrayReference(elements)

	newLen := arr.Push(4.0)
	if newLen != 4.0 {
		t.Errorf("Push(4.0) returned %v, expected 4.0", newLen)
	}
	if len(*arr.Elements) != 4 {
		t.Errorf("Array length = %d, expected 4", len(*arr.Elements))
	}
	if (*arr.Elements)[3] != 4.0 {
		t.Errorf("Last element = %v, expected 4.0", (*arr.Elements)[3])
	}

	newLen = arr.Push(5.0, 6.0, 7.0)
	if newLen != 7.0 {
		t.Errorf("Push(5.0, 6.0, 7.0) returned %v, expected 7.0", newLen)
	}
	if len(*arr.Elements) != 7 {
		t.Errorf("Array length = %d, expected 7", len(*arr.Elements))
	}

	expected := []internal.Value{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}
	for i, exp := range expected {
		if (*arr.Elements)[i] != exp {
			t.Errorf("Element[%d] = %v, expected %v", i, (*arr.Elements)[i], exp)
		}
	}

	arr.Push("string", true, nil)
	if len(*arr.Elements) != 10 {
		t.Errorf("Array length = %d, expected 10 after pushing mixed types", len(*arr.Elements))
	}
}

func TestArrayReferenceLength(t *testing.T) {
	tests := []struct {
		elements internal.Array
		expected float64
	}{
		{internal.Array{}, 0.0},
		{internal.Array{1.0}, 1.0},
		{internal.Array{1.0, 2.0, 3.0}, 3.0},
		{internal.Array{"a", "b", "c", "d", "e"}, 5.0},
	}

	for _, tt := range tests {
		arr := NewArrayReference(tt.elements)
		length := arr.Length()
		if length != tt.expected {
			t.Errorf("Length() = %v, expected %v", length, tt.expected)
		}
	}
}

func TestArrayReferenceGetElements(t *testing.T) {
	elements := internal.Array{"hello", 42.0, true}
	arr := NewArrayReference(elements)

	result := arr.GetElements()

	if len(result) != 3 {
		t.Errorf("GetElements() length = %d, expected 3", len(result))
	}

	if result[0] != "hello" {
		t.Errorf("GetElements()[0] = %v, expected 'hello'", result[0])
	}
	if result[1] != 42.0 {
		t.Errorf("GetElements()[1] = %v, expected 42.0", result[1])
	}
	if result[2] != true {
		t.Errorf("GetElements()[2] = %v, expected true", result[2])
	}
}

// Test that modifications through ArrayReference affect the underlying array
func TestArrayReferenceMutability(t *testing.T) {
	elements := internal.Array{1.0, 2.0, 3.0}
	arr := NewArrayReference(elements)

	arr.Set(0, 10.0)

	if (*arr.Elements)[0] != 10.0 {
		t.Errorf("Modification not visible, got %v", (*arr.Elements)[0])
	}

	arr.Push(4.0)

	if len(*arr.Elements) != 4 {
		t.Errorf("Push didn't modify array, length = %d", len(*arr.Elements))
	}
}

func TestArrayReferenceEmptyArray(t *testing.T) {
	elements := internal.Array{}
	arr := NewArrayReference(elements)

	if arr.Length() != 0.0 {
		t.Errorf("Empty array length = %v, expected 0.0", arr.Length())
	}

	if arr.Get(0) != nil {
		t.Errorf("Get(0) on empty array = %v, expected nil", arr.Get(0))
	}

	newLen := arr.Push("first")
	if newLen != 1.0 {
		t.Errorf("Push to empty array returned %v, expected 1.0", newLen)
	}
	if (*arr.Elements)[0] != "first" {
		t.Errorf("First element = %v, expected 'first'", (*arr.Elements)[0])
	}
}

func TestArrayReferenceNestedArrays(t *testing.T) {
	inner1 := internal.Array{1.0, 2.0}
	inner2 := internal.Array{3.0, 4.0}
	elements := internal.Array{inner1, inner2}

	arr := NewArrayReference(elements)

	if arr.Length() != 2.0 {
		t.Errorf("Nested array length = %v, expected 2.0", arr.Length())
	}

	nested := arr.Get(0)
	if innerArr, ok := nested.(internal.Array); ok {
		if len(innerArr) != 2 {
			t.Errorf("Nested array length = %d, expected 2", len(innerArr))
		}
		if innerArr[0] != 1.0 {
			t.Errorf("Nested array[0] = %v, expected 1.0", innerArr[0])
		}
	} else {
		t.Error("Get(0) didn't return an array")
	}
}
