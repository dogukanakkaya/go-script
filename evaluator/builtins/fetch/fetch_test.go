package fetch

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchSuccess(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	// Call fetch with test server URL
	result := Fetch.Fn(server.URL)

	// Check result is a map
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	// Check status
	status, ok := response["status"].(float64)
	if !ok || status != 200 {
		t.Errorf("Expected status 200, got %v", response["status"])
	}

	// Check body
	body, ok := response["body"].(string)
	if !ok || body != `{"message": "success"}` {
		t.Errorf("Expected body %q, got %q", `{"message": "success"}`, body)
	}

	// Check ok field
	okField, ok := response["ok"].(bool)
	if !ok || !okField {
		t.Errorf("Expected ok=true, got %v", response["ok"])
	}
}

func TestFetchNotFound(t *testing.T) {
	// Create a test HTTP server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Not Found`))
	}))
	defer server.Close()

	// Call fetch
	result := Fetch.Fn(server.URL)

	// Check result
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	// Check status
	status, ok := response["status"].(float64)
	if !ok || status != 404 {
		t.Errorf("Expected status 404, got %v", response["status"])
	}

	// Check ok field (should be false for 404)
	okField, ok := response["ok"].(bool)
	if !ok || okField {
		t.Errorf("Expected ok=false, got %v", response["ok"])
	}
}

func TestFetchInvalidURL(t *testing.T) {
	// Call fetch with invalid URL
	result := Fetch.Fn("not-a-valid-url")

	// Check result contains error
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	// Check error field exists
	if _, hasError := response["error"]; !hasError {
		t.Errorf("Expected error field in response, got %v", response)
	}
}

func TestFetchNoArgs(t *testing.T) {
	// Call fetch with no arguments
	result := Fetch.Fn()

	// Check result contains error
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	// Check error message
	errorMsg, ok := response["error"].(string)
	if !ok || errorMsg != "fetch requires 1 or 2 arguments (url, options?)" {
		t.Errorf("Expected specific error message, got %q", errorMsg)
	}
}

func TestFetchTooManyArgs(t *testing.T) {
	// Call fetch with too many arguments
	result := Fetch.Fn("url1", "url2", "url3")

	// Check result contains error
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	// Check error message
	errorMsg, ok := response["error"].(string)
	if !ok || errorMsg != "fetch requires 1 or 2 arguments (url, options?)" {
		t.Errorf("Expected specific error message, got %q", errorMsg)
	}
}

func TestFetchStatusCodes(t *testing.T) {
	tests := []struct {
		statusCode int
		expectedOk bool
	}{
		{200, true},
		{201, true},
		{204, true},
		{299, true},
		{300, false},
		{400, false},
		{404, false},
		{500, false},
	}

	for _, tt := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.statusCode)
			w.Write([]byte("test"))
		}))

		result := Fetch.Fn(server.URL)
		response, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", result)
		}

		status := response["status"].(float64)
		if int(status) != tt.statusCode {
			t.Errorf("Expected status %d, got %v", tt.statusCode, status)
		}

		okField := response["ok"].(bool)
		if okField != tt.expectedOk {
			t.Errorf("For status %d, expected ok=%v, got %v", tt.statusCode, tt.expectedOk, okField)
		}

		server.Close()
	}
}

func TestFetchWithCustomMethod(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Echo back the method
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(r.Method))
		}))

		options := map[string]interface{}{
			"method": method,
		}

		result := Fetch.Fn(server.URL, options)
		response, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", result)
		}

		body := response["body"].(string)
		if body != method {
			t.Errorf("Expected method %s, got %s", method, body)
		}

		server.Close()
	}
}

func TestFetchWithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check custom headers
		authHeader := r.Header.Get("Authorization")
		contentType := r.Header.Get("Content-Type")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(authHeader + "|" + contentType))
	}))
	defer server.Close()

	options := map[string]interface{}{
		"headers": map[string]interface{}{
			"Authorization": "Bearer token123",
			"Content-Type":  "application/json",
		},
	}

	result := Fetch.Fn(server.URL, options)
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	body := response["body"].(string)
	expected := "Bearer token123|application/json"
	if body != expected {
		t.Errorf("Expected body %q, got %q", expected, body)
	}
}

func TestFetchWithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and echo back the request body
		bodyBytes, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
	}))
	defer server.Close()

	requestBody := `{"name": "Alice", "age": 30}`
	options := map[string]interface{}{
		"method": "POST",
		"body":   requestBody,
	}

	result := Fetch.Fn(server.URL, options)
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	body := response["body"].(string)
	if body != requestBody {
		t.Errorf("Expected body %q, got %q", requestBody, body)
	}
}

func TestFetchResponseHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "test-value")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	result := Fetch.Fn(server.URL)
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	// Check headers exist
	headers, ok := response["headers"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected headers to be map[string]interface{}, got %T", response["headers"])
	}

	// Check specific header
	customHeader, exists := headers["X-Custom-Header"]
	if !exists {
		t.Errorf("Expected X-Custom-Header to exist in response headers")
	}
	if customHeader != "test-value" {
		t.Errorf("Expected X-Custom-Header=test-value, got %v", customHeader)
	}

	// Check Content-Type header
	contentType, exists := headers["Content-Type"]
	if !exists {
		t.Errorf("Expected Content-Type to exist in response headers")
	}
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type=application/json, got %v", contentType)
	}
}

func TestFetchComplexRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header")
		}
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("Expected Authorization header")
		}

		// Read body
		bodyBytes, _ := io.ReadAll(r.Body)
		bodyStr := string(bodyBytes)

		// Send response
		w.Header().Set("X-Response-Id", "123")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": true, "echo": "` + bodyStr + `"}`))
	}))
	defer server.Close()

	options := map[string]interface{}{
		"method": "POST",
		"headers": map[string]interface{}{
			"Content-Type":  "application/json",
			"Authorization": "Bearer secret",
		},
		"body": `{"action": "create", "data": "test"}`,
	}

	result := Fetch.Fn(server.URL, options)
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	// Check status
	status := response["status"].(float64)
	if status != 201 {
		t.Errorf("Expected status 201, got %v", status)
	}

	// Check ok
	okField := response["ok"].(bool)
	if !okField {
		t.Errorf("Expected ok=true")
	}

	// Check body contains expected data
	body := response["body"].(string)
	if !bytes.Contains([]byte(body), []byte("success")) {
		t.Errorf("Expected body to contain 'success', got %q", body)
	}

	// Check response headers
	headers := response["headers"].(map[string]interface{})
	if headers["X-Response-Id"] != "123" {
		t.Errorf("Expected X-Response-Id=123, got %v", headers["X-Response-Id"])
	}
}

func TestFetchInvalidOptions(t *testing.T) {
	// Call fetch with invalid options type
	result := Fetch.Fn("http://example.com", "not-an-object")

	// Check result contains error
	response, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	// Check error message
	errorMsg, ok := response["error"].(string)
	if !ok {
		t.Errorf("Expected error field to be string, got %T", response["error"])
	}
	// Error message now includes type information
	if !ok || len(errorMsg) == 0 {
		t.Errorf("Expected error message, got %q", errorMsg)
	}
}
