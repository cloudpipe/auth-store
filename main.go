package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func main() {
	c, err := NewContext()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Unable to load application context.")
		return
	}

	log.WithFields(log.Fields{
		"address": c.ListenAddr(),
	}).Info("Auth API listening.")

	// v1 routes
	http.HandleFunc("/v1/style", BindContext(c, StyleHandler))
	http.HandleFunc("/v1/validate", BindContext(c, ValidateHandler))

	http.HandleFunc("/v1/accounts", BindContext(c, AccountHandler))
	http.HandleFunc("/v1/keys", BindContext(c, KeyHandler))

	err = http.ListenAndServeTLS(c.ListenAddr(), c.Cert, c.Key, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to launch auth API.")
	}
}

// ContextHandler is an HTTP HandlerFunc that accepts an additional parameter containing the
// server context.
type ContextHandler func(c *Context, w http.ResponseWriter, r *http.Request)

// BindContext returns an http.HandlerFunc that binds a ContextHandler to a specific Context.
func BindContext(c *Context, handler ContextHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { handler(c, w, r) }
}

// APIError consistently renders error conditions as a JSON payload.
type APIError struct {
	// If nonzero, this message will be displayed to the user in the generated response payload.
	UserMessage string `json:"message"`

	// If nonzero, this message will be displayed to operators in the process log.
	LogMessage string `json:"-"`

	// Used as both UserMessage and LogMessage if either are missing.
	Message string `json:"-"`
}

// Log emits a log message for an error.
func (err APIError) Log(username string) APIError {
	if err.LogMessage == "" {
		err.LogMessage = err.Message
	}

	f := log.Fields{}
	if username != "" {
		f["username"] = username
	}
	log.WithFields(f).Error(err.LogMessage)
	return err
}

// Report renders an error as an HTTP response with the correct content-type and HTTP status code.
func (err APIError) Report(w http.ResponseWriter, status int) APIError {
	if err.UserMessage == "" {
		err.UserMessage = err.Message
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	encodeErr := json.NewEncoder(w).Encode(err)
	if encodeErr != nil {
		fmt.Fprintf(w, `{"message":"Unable to encode error: %v"}`, encodeErr)
	}
	return err
}

// MethodOk tests the HTTP request method. If the method is correct, it does nothing and
// returns true. If it's incorrect, it generates a JSON error and returns false.
func MethodOk(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}

	APIError{
		Message: fmt.Sprintf("Unsupported method %s. Only %s is accepted for this resource.",
			r.Method, method),
	}.Log("").Report(w, http.StatusMethodNotAllowed)

	return false
}

// ExtractCredentials attempts to read an account name and API key from the request.
func ExtractCredentials(w http.ResponseWriter, r *http.Request) (accountName, apiKey string, ok bool) {
	if err := r.ParseForm(); err != nil {
		APIError{
			Message: fmt.Sprintf("Unable to parse URL parameters: %v", err),
		}.Log("").Report(w, http.StatusBadRequest)
		return "", "", false
	}

	accountName, apiKey = r.FormValue("accountName"), r.FormValue("apiKey")
	if accountName == "" || apiKey == "" {
		APIError{
			UserMessage: `Missing required query parameters "accountName" and "apiKey".`,
			LogMessage:  "Key validation request missing required query parameters.",
		}.Log("").Report(w, http.StatusBadRequest)
		return "", "", false
	}
	return accountName, apiKey, true
}
