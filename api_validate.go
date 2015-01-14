package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// ValidateHandler determines whether or not an API key is valid for a specific account.
func ValidateHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	if !MethodOk(w, r, "GET") {
		return
	}

	accountName, apiKey, ok := ExtractKeyCredentials(w, r, "Key validation")
	if !ok {
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

	var message string
	if ok {
		w.WriteHeader(http.StatusNoContent)
		message = "API key successfully validated."
	} else {
		w.WriteHeader(http.StatusNotFound)
		message = "Invalid API key encountered."
	}

	log.WithFields(log.Fields{
		"account": accountName,
		"key":     apiKey,
	}).Info(message)
}
