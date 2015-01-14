package main

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
)

// ValidateHandler determines whether or not an API key is valid for a specific account.
func ValidateHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	if !MethodOk(w, r, "GET") {
		return
	}

	if err := r.ParseForm(); err != nil {
		APIError{
			Message: fmt.Sprintf("Unable to parse URL parameters: %v", err),
		}.Log("").Report(w, http.StatusBadRequest)
		return
	}

	accountName, apiKey := r.FormValue("accountName"), r.FormValue("apiKey")
	if accountName == "" || apiKey == "" {
		APIError{
			UserMessage: `Missing required query parameters "accountName" and "apiKey".`,
			LogMessage:  "Key validation request missing required query parameters.",
		}.Log("").Report(w, http.StatusBadRequest)
		return
	}

	ok, err := c.Storage.AccountHasKey(accountName, apiKey)
	if err != nil {
		if err == mgo.ErrNotFound {
			APIError{
				Message: "Unrecognized account or API key.",
			}.Log(accountName).Report(w, http.StatusUnauthorized)
			return
		}
		APIError{
			UserMessage: "Internal storage error encountered. Please try again later.",
			LogMessage:  fmt.Sprintf("Storage error: %v", err),
		}.Log(accountName).Report(w, http.StatusInternalServerError)
		return
	}

	if ok {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
