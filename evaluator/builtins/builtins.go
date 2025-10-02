package builtins

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"
)

type Builtin struct {
	Name string
	Fn   func(args ...interface{}) interface{}
}

var builtins = map[string]*Builtin{
	"print": Print,
	"fetch": Fetch,
}

func Get(name string) (*Builtin, bool) {
	builtin, ok := builtins[name]
	return builtin, ok
}

// Print is a built-in function that prints values to stdout
//
// Syntax: print(arg1, arg2, ...)
//
// Examples:
//
//	print("Hello, World!")           → prints: Hello, World!
//	print("x =", 42)                  → prints: x = 42
var Print = &Builtin{
	Name: "print",
	Fn: func(args ...interface{}) interface{} {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(toString(arg))
		}
		fmt.Println()
		return nil
	},
}

// Fetch is a built-in function that makes HTTP requests
//
// Syntax: fetch(url, options)
//
// Options object (all optional):
//   - method: string ("GET", "POST", "PUT", "DELETE", "PATCH", etc.) - default: "GET"
//   - headers: object with string key-value pairs
//   - body: string request body (for POST, PUT, PATCH)
//
// Examples:
//
//	// Simple GET request
//	let response = fetch("https://api.example.com/data")
//	print(response.body)
//
//	// POST request with JSON body
//	let response = fetch("https://api.example.com/users", {
//	    method: "POST",
//	    headers: {
//	        "Content-Type": "application/json",
//	        "Authorization": "Bearer token123"
//	    },
//	    body: '{"name": "Alice", "age": 30}'
//	})
var Fetch = &Builtin{
	Name: "fetch",
	Fn: func(args ...interface{}) interface{} {
		if len(args) < 1 || len(args) > 2 {
			return map[string]interface{}{
				"error": "fetch requires 1 or 2 arguments (url, options?)",
			}
		}

		url := toString(args[0])

		method := "GET"
		headers := make(map[string]string)
		var bodyStr string

		// Parse options if provided
		if len(args) == 2 {
			// Use reflection to handle any map type
			v := reflect.ValueOf(args[1])
			if v.Kind() != reflect.Map {
				return map[string]interface{}{
					"error": fmt.Sprintf("second argument must be an options object, got %T", args[1]),
				}
			}

			if methodVal := v.MapIndex(reflect.ValueOf("method")); methodVal.IsValid() {
				method = toString(methodVal.Interface())
			}

			if headersVal := v.MapIndex(reflect.ValueOf("headers")); headersVal.IsValid() {
				headersInterface := headersVal.Interface()
				hv := reflect.ValueOf(headersInterface)
				if hv.Kind() == reflect.Map {
					for _, key := range hv.MapKeys() {
						keyStr := key.String()
						val := hv.MapIndex(key)
						if val.IsValid() {
							headers[keyStr] = toString(val.Interface())
						}
					}
				}
			}

			if bodyVal := v.MapIndex(reflect.ValueOf("body")); bodyVal.IsValid() {
				bodyStr = toString(bodyVal.Interface())
			}
		}

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		var bodyReader io.Reader
		if bodyStr != "" {
			bodyReader = bytes.NewBufferString(bodyStr)
		}

		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return map[string]interface{}{
				"error": err.Error(),
			}
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			return map[string]interface{}{
				"error": err.Error(),
			}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return map[string]interface{}{
				"error": err.Error(),
			}
		}

		responseHeaders := make(map[string]interface{})
		for key, values := range resp.Header {
			if len(values) == 1 {
				responseHeaders[key] = values[0]
			} else {
				responseHeaders[key] = values
			}
		}

		return map[string]interface{}{
			"status":     float64(resp.StatusCode),
			"statusText": resp.Status,
			"body":       string(body),
			"headers":    responseHeaders,
			"ok":         resp.StatusCode >= 200 && resp.StatusCode < 300,
		}
	},
}

func toString(val interface{}) string {
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
			result += k + ": " + toString(val)
			first = false
		}
		result += "}"
		return result
	default:
		// Handle other types (functions, etc.)
		return fmt.Sprintf("%v", v)
	}
}
