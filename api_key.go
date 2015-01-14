package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// KeyHandler dispatches requests made to the /keys resource to relevant subhandlers
func KeyHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		KeyGenerationHandler(c, w, r)
	case "DELETE":
		KeyRevocationHandler(c, w, r)
	default:
		APIError{
			Message: fmt.Sprintf("Unsupported method %s. Only POST is accepted for this resource.",
				r.Method),
		}.Log("").Report(w, http.StatusMethodNotAllowed)
	}
}

// KeyGenerationHandler generates a new API key for a provided user account. It persists the new
// key in storage and returns it as a plaintext string.
func KeyGenerationHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	// Validate the credentials provided as query parameters.
	accountName, password, ok := ExtractPasswordCredentials(w, r, "Key generation")
	if !ok {
		return
	}

	rejectAuth := func() {
		APIError{
			UserMessage: "Incorrect account name or password.",
			LogMessage:  "Authentication failure for account.",
		}.Log(accountName).Report(w, http.StatusUnauthorized)
	}

	account, err := c.Storage.FindAccount(accountName)
	if err != nil {
		APIError{
			UserMessage: "Internal storage error. Please try again later.",
			LogMessage:  fmt.Sprintf("Error finding account: %v", err),
		}.Log(accountName).Report(w, http.StatusInternalServerError)
		return
	}
	if account == nil {
		// Account does not exist. Treat this exactly like a failed password attempt.

		// Thwart timing attacks by doing a fake bcrypt comparison.
		(&Account{}).HasPassword(password)

		rejectAuth()
		return
	}

	if !account.HasPassword(password) {
		// BZZZZZZZT
		rejectAuth()
		return
	}

	// Success. Generate the new key, put it in Mongo, and return it as a plaintext response.
	key, err := account.GenerateAPIKey()
	if err != nil {
		APIError{
			UserMessage: "Unable to generate your API key. Please try again later.",
			LogMessage:  fmt.Sprintf("Unable to generate API key: %v", err),
		}.Log(accountName).Report(w, http.StatusInternalServerError)
		return
	}

	if err := c.Storage.AddKeyToAccount(accountName, key); err != nil {
		APIError{
			UserMessage: "Unable to generate your API key. Please try again later.",
			LogMessage:  fmt.Sprintf("Unable to store API key in MongoDB: %v", err),
		}.Log(accountName).Report(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))

	log.WithFields(log.Fields{
		"account": accountName,
		"key":     key,
	}).Info("A new API key has been generated.")
}

// KeyRevocationHandler marks an API key as invalid for a specific account.
func KeyRevocationHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	accountName, apiKey, ok := ExtractKeyCredentials(w, r, "Key revocation")
	if !ok {
		return
	}

	if err := c.Storage.RevokeKeyFromAccount(accountName, apiKey); err != nil {
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
	}

	// Success!
	w.WriteHeader(http.StatusNoContent)

	log.WithFields(log.Fields{
		"account": accountName,
		"key":     apiKey,
	}).Info("An existing API key has revoked.")
}
