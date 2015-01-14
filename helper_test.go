package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func HTTPRequest(t *testing.T, method, url, body string) *http.Request {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	r, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("Unable to create a request: %v", err)
	}
	if body != "" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return r
}
