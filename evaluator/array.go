package evaluator

import (
	"go-script/environment"
	"go-script/evaluator/builtins/array"
	"go-script/internal"
)

func GetArrayProperty(arr *array.ArrayReference, property string) Value {
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

func createPushMethod(arr *array.ArrayReference) Value {
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

func createMapMethod(arr *array.ArrayReference) Value {
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
				callbackEnv := environment.New(fn.Env)

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

			return array.NewArrayReference(result)
		},
	}
}

func createFilterMethod(arr *array.ArrayReference) Value {
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
				callbackEnv := environment.New(fn.Env)

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

			return array.NewArrayReference(result)
		},
	}
}
