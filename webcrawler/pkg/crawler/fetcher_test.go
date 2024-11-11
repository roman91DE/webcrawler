package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchAll(t *testing.T) {
	tests := []struct {
		name     string
		urls     []string
		validate func(results []FetchResult) bool
	}{
		{
			name: "All valid URLs",
			urls: []string{
				httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Hello World"))
				})).URL,
			},
			validate: func(results []FetchResult) bool {
				return len(results) == 1 && results[0].Status == http.StatusOK && results[0].Content == "Hello World"
			},
		},
		{
			name: "Invalid URL",
			urls: []string{
				"ht!tp://invalid-url",
			},
			validate: func(results []FetchResult) bool {
				return len(results) == 1 && results[0].Error != nil
			},
		},
		{
			name: "Timeout URL",
			urls: []string{
				httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(6 * time.Second)
				})).URL,
			},
			validate: func(results []FetchResult) bool {
				return len(results) == 1 && results[0].Error != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FetchAll(tt.urls)
			if !tt.validate(results) {
				t.Errorf("Validation failed for test %s", tt.name)
			}
		})
	}
}

func TestFetch(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		statusCode  int
		expectError bool
	}{
		{
			name:        "Valid URL",
			url:         httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })).URL,
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "Invalid URL Format",
			url:         "ht!tp://invalid-url",
			statusCode:  0,
			expectError: true,
		},
		{
			name: "Non-responsive server",
			url: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(20 * time.Second)
			})).URL,
			statusCode:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan FetchResult)
			go fetch(tt.url, ch)
			result := <-ch

			if tt.expectError && result.Error == nil {
				t.Errorf("Expected an error but didn't get one for test %s", tt.name)
			}
			if !tt.expectError && result.Error != nil {
				t.Errorf("Did not expect an error but got one for test %s: %v", tt.name, result.Error)
			}
			if result.Status != tt.statusCode && !tt.expectError {
				t.Errorf("Expected status code %d but got %d for test %s", tt.statusCode, result.Status, tt.name)
			}
		})
	}
}
