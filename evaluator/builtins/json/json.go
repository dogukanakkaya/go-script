package json

import (
	encodingjson "encoding/json"
	"fmt"
)

// Builtin represents a built-in function
type Builtin struct {
	Name string
	Fn   func(args ...interface{}) interface{}
}

// JSON is a namespace object containing stringify and parse methods
var JSON = map[string]*Builtin{
	"stringify": Stringify,
	"parse":     Parse,
}

// Stringify converts a JavaScript value to a JSON string
//
// Syntax: JSON.stringify(value)
//
// Examples:
//
//	let obj = { name: "Alice", age: 30 }
//	let jsonStr = JSON.stringify(obj)
//	print(jsonStr)  → {"age":30,"name":"Alice"}
var Stringify = &Builtin{
	Name: "stringify",
	Fn: func(args ...interface{}) interface{} {
		if len(args) != 1 {
			return map[string]interface{}{
				"error": "JSON.stringify requires exactly 1 argument",
			}
		}

		// Convert to JSON
		jsonBytes, err := encodingjson.Marshal(args[0])
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("JSON.stringify error: %v", err),
			}
		}

		return string(jsonBytes)
	},
}

// Parse converts a JSON string to a JavaScript value
//
// Syntax: JSON.parse(jsonString)
//
// Examples:
//
//	let jsonStr = '{"name":"Alice","age":30}'
//	let obj = JSON.parse(jsonStr)
//	print(obj.name)  → Alice
var Parse = &Builtin{
	Name: "parse",
	Fn: func(args ...interface{}) interface{} {
		if len(args) != 1 {
			return map[string]interface{}{
				"error": "JSON.parse requires exactly 1 argument",
			}
		}

		// Get the JSON string
		jsonStr, ok := args[0].(string)
		if !ok {
			return map[string]interface{}{
				"error": fmt.Sprintf("JSON.parse requires a string argument, got %T", args[0]),
			}
		}

		// Parse JSON into generic structure
		var result interface{}
		err := encodingjson.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("JSON.parse error: %v", err),
			}
		}

		// Convert numbers to float64 and return
		return convertJSONTypes(result)
	},
}

// convertJSONTypes converts JSON types to JavaScript-compatible types
// JSON numbers come as float64, which is what we want
// JSON objects come as map[string]interface{}, which works
// JSON arrays come as []interface{}, which works
func convertJSONTypes(val interface{}) interface{} {
	switch v := val.(type) {
	case map[string]interface{}:
		// Recursively convert nested objects
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = convertJSONTypes(value)
		}
		return result
	case []interface{}:
		// Recursively convert array elements
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = convertJSONTypes(value)
		}
		return result
	default:
		// Primitives (string, float64, bool, nil) are already correct
		return v
	}
}
