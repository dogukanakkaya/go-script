package fetch

import (
	"bytes"
	"fmt"
	"go-script/internal"
	"io"
	"net/http"
	"reflect"
	"time"
)

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
var Fetch = &internal.Builtin{
	Name: "fetch",
	Fn: func(args ...interface{}) interface{} {
		if len(args) < 1 || len(args) > 2 {
			return map[string]interface{}{
				"error": "fetch requires 1 or 2 arguments (url, options?)",
			}
		}

		url := internal.ToString(args[0])

		// Default options
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

			// Parse method
			if methodVal := v.MapIndex(reflect.ValueOf("method")); methodVal.IsValid() {
				method = internal.ToString(methodVal.Interface())
			}

			// Parse headers
			if headersVal := v.MapIndex(reflect.ValueOf("headers")); headersVal.IsValid() {
				headersInterface := headersVal.Interface()
				hv := reflect.ValueOf(headersInterface)
				if hv.Kind() == reflect.Map {
					for _, key := range hv.MapKeys() {
						keyStr := key.String()
						val := hv.MapIndex(key)
						if val.IsValid() {
							headers[keyStr] = internal.ToString(val.Interface())
						}
					}
				}
			}

			// Parse body
			if bodyVal := v.MapIndex(reflect.ValueOf("body")); bodyVal.IsValid() {
				bodyStr = internal.ToString(bodyVal.Interface())
			}
		}

		// Create HTTP client with timeout
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		// Create request
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

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		// Make request
		resp, err := client.Do(req)
		if err != nil {
			return map[string]interface{}{
				"error": err.Error(),
			}
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return map[string]interface{}{
				"error": err.Error(),
			}
		}

		// Parse response headers
		responseHeaders := make(map[string]interface{})
		for key, values := range resp.Header {
			if len(values) == 1 {
				responseHeaders[key] = values[0]
			} else {
				responseHeaders[key] = values
			}
		}

		// Return response object
		return map[string]interface{}{
			"status":     float64(resp.StatusCode),
			"statusText": resp.Status,
			"body":       string(body),
			"headers":    responseHeaders,
			"ok":         resp.StatusCode >= 200 && resp.StatusCode < 300,
		}
	},
}