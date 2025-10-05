package evaluator

import "go-script/internal"

// ArrayReference wraps an array to make it mutable
// This allows array methods like push to modify the array in place
type ArrayReference struct {
	Elements *internal.Array
}

func NewArrayReference(elements internal.Array) *ArrayReference {
	return &ArrayReference{Elements: &elements}
}

func (ar *ArrayReference) Get(index int) Value {
	if index < 0 || index >= len(*ar.Elements) {
		return nil
	}
	return (*ar.Elements)[index]
}

func (ar *ArrayReference) Set(index int, value Value) {
	if index >= 0 && index < len(*ar.Elements) {
		(*ar.Elements)[index] = value
	}
}

func (ar *ArrayReference) Push(values ...Value) float64 {
	*ar.Elements = append(*ar.Elements, values...)
	return float64(len(*ar.Elements))
}

func (ar *ArrayReference) Length() float64 {
	return float64(len(*ar.Elements))
}

func (ar *ArrayReference) GetElements() internal.Array {
	return *ar.Elements
}

func GetArrayProperty(arr *ArrayReference, property string) Value {
	switch property {
	case "length":
		return arr.Length()
	case "push":
		return createPushMethod(arr)
	case "map":
		return createMapMethod(arr)
	case "filter":
		return createFilterMethod(arr)
	default:
		return nil
	}
}

func createPushMethod(arr *ArrayReference) Value {
	return &internal.Builtin{
		Fn: func(args ...interface{}) interface{} {
			values := make([]Value, len(args))
			for i, arg := range args {
				values[i] = arg
			}
			return arr.Push(values...)
		},
	}
}

func createMapMethod(arr *ArrayReference) Value {
	return &internal.Builtin{
		Fn: func(args ...interface{}) interface{} {
			if len(args) == 0 {
				return nil
			}

			fn, ok := args[0].(*Function)
			if !ok {
				return nil
			}

			result := make(Array, 0, len(*arr.Elements))
			for i, elem := range *arr.Elements {
				callbackEnv := New(fn.Env)

				if len(fn.Parameters) > 0 {
					callbackEnv.Set(fn.Parameters[0], elem)
				}
				if len(fn.Parameters) > 1 {
					callbackEnv.Set(fn.Parameters[1], float64(i))
				}
				if len(fn.Parameters) > 2 {
					callbackEnv.Set(fn.Parameters[2], arr)
				}

				callbackResult := Eval(fn.Body, callbackEnv)

				if returnVal, ok := callbackResult.(*ReturnValue); ok {
					callbackResult = returnVal.Value
				}

				result = append(result, callbackResult)
			}

			return NewArrayReference(result)
		},
	}
}

func createFilterMethod(arr *ArrayReference) Value {
	return &internal.Builtin{
		Fn: func(args ...interface{}) interface{} {
			if len(args) == 0 {
				return nil
			}

			fn, ok := args[0].(*Function)
			if !ok {
				return nil
			}

			result := make(Array, 0)
			for i, elem := range *arr.Elements {
				callbackEnv := New(fn.Env)

				if len(fn.Parameters) > 0 {
					callbackEnv.Set(fn.Parameters[0], elem)
				}
				if len(fn.Parameters) > 1 {
					callbackEnv.Set(fn.Parameters[1], float64(i))
				}
				if len(fn.Parameters) > 2 {
					callbackEnv.Set(fn.Parameters[2], arr)
				}

				callbackResult := Eval(fn.Body, callbackEnv)

				if returnVal, ok := callbackResult.(*ReturnValue); ok {
					callbackResult = returnVal.Value
				}

				if isTruthy(callbackResult) {
					result = append(result, elem)
				}
			}

			return NewArrayReference(result)
		},
	}
}
