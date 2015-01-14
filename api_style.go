package main

import (
	"fmt"
	"net/http"
)

// StyleHandler reports this auth backend's "style" attribute. This is reported through cloudpipe
// to API consumers to provide them with a hint about other auth interactions that are possible at
// this endpoint.
func StyleHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	if !MethodOk(w, r, "GET") {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "auth-store")
}
