package internal

import "fmt"

type Builtin struct {
	Name string
	Fn   func(args ...interface{}) interface{}
}

func ToString(val interface{}) string {
	if val == nil {
		return "nil"
	}

	switch v := val.(type) {
	case string:
		return v
	case float64:
		// Format integers without decimal point
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case map[string]interface{}:
		// Format object as {key: value, ...}
		result := "{"
		first := true
		for k, val := range v {
			if !first {
				result += ", "
			}
			result += k + ": " + ToString(val)
			first = false
		}
		result += "}"
		return result
	default:
		// Handle other types (functions, etc.)
		return fmt.Sprintf("%v", v)
	}
}
