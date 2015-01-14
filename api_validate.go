package main

import "net/http"

// ValidateHandler determines whether or not an API key is valid for a specific account.
func ValidateHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	if !MethodOk(w, r, "GET") {
		return
	}
}
