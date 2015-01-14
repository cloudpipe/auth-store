package main

import (
	"fmt"
	"net/http"
)

// KeyHandler dispatches requests made to the /keys resource to relevant subhandlers
func KeyHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		KeyGenerationHandler(c, w, r)
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
	// Validate the credentials provided in basic auth.
	accountName, password, ok := r.BasicAuth()
	if !ok {
		APIError{
			UserMessage: "Please use HTTP basic authentication to provide an account name and password.",
			LogMessage:  "Key generation request failed due to missing credentials.",
		}.Log("").Report(w, http.StatusUnauthorized)
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
}
