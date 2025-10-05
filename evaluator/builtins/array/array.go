package array

import (
	"go-script/internal"
)

// ArrayReference wraps an array to make it mutable
// This allows array methods like push to modify the array in place
type ArrayReference struct {
	Elements *internal.Array
}

func NewArrayReference(elements internal.Array) *ArrayReference {
	return &ArrayReference{Elements: &elements}
}

func (ar *ArrayReference) Get(index int) internal.Value {
	if index < 0 || index >= len(*ar.Elements) {
		return nil
	}
	return (*ar.Elements)[index]
}

func (ar *ArrayReference) Set(index int, value internal.Value) {
	if index >= 0 && index < len(*ar.Elements) {
		(*ar.Elements)[index] = value
	}
}

func (ar *ArrayReference) Push(values ...internal.Value) float64 {
	*ar.Elements = append(*ar.Elements, values...)
	return float64(len(*ar.Elements))
}

func (ar *ArrayReference) Length() float64 {
	return float64(len(*ar.Elements))
}

func (ar *ArrayReference) GetElements() internal.Array {
	return *ar.Elements
}
