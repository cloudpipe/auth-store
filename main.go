package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	go ServeInternal(c)
	ServeExternal(c)
}

// ServeInternal configures and launches the internal API.
func ServeInternal(c *Context) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("auth-store internal API alive and running."))
	})

	mux.HandleFunc("/v1/style", BindContext(c, StyleHandler))
	mux.HandleFunc("/v1/validate", BindContext(c, ValidateHandler))

	// Load TLS credentials used by the internal API.

	caCertPool := x509.NewCertPool()

	caCertPEM, err := ioutil.ReadFile(c.InternalCACert)
	if err != nil {
		log.Debug("Hint: if you're running in dev mode, try running script/genkeys first.")
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Unable to load CA certificate for internal API.")
	}
	caCertPool.AppendCertsFromPEM(caCertPEM)

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
	}

	server := &http.Server{
		Addr:      c.InternalListenAddr(),
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	log.WithFields(log.Fields{
		"address": c.InternalListenAddr(),
	}).Info("Internal auth API listening.")

	err = server.ListenAndServeTLS(c.InternalCert, c.InternalKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Unable to launch internal auth API.")
	}
}

// ServeExternal configures and launches the external API.
func ServeExternal(c *Context) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("auth-store external API alive and running."))
	})

	mux.HandleFunc("/v1/accounts", BindContext(c, AccountHandler))
	mux.HandleFunc("/v1/keys", BindContext(c, KeyHandler))

	server := &http.Server{
		Addr:    c.ExternalListenAddr(),
		Handler: mux,
	}

	log.WithFields(log.Fields{
		"address": c.ExternalListenAddr(),
	}).Info("External auth API listening.")

	err := server.ListenAndServeTLS(c.ExternalCert, c.ExternalKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Unable to launch external auth API.")
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

// ExtractKeyCredentials attempts to read an account name and API key from the request.
func ExtractKeyCredentials(w http.ResponseWriter, r *http.Request, requestName string) (accountName, apiKey string, ok bool) {
	return extractCredentials(w, r, requestName, "apiKey")
}

// ExtractPasswordCredentials attempts to read an account name and password from a request form.
func ExtractPasswordCredentials(w http.ResponseWriter, r *http.Request, requestName string) (accountName, password string, ok bool) {
	return extractCredentials(w, r, requestName, "password")
}

func extractCredentials(w http.ResponseWriter, r *http.Request, requestName, credentialName string) (accountName, credential string, ok bool) {
	if err := r.ParseForm(); err != nil {
		APIError{
			Message: fmt.Sprintf("Unable to parse URL parameters: %v", err),
		}.Log("").Report(w, http.StatusBadRequest)
		return "", "", false
	}

	accountName, credential = r.FormValue("accountName"), r.FormValue(credentialName)
	if accountName == "" || credential == "" {
		APIError{
			UserMessage: fmt.Sprintf(
				`Missing required parameters "accountName" and "%s".`,
				credentialName,
			),
			LogMessage: fmt.Sprintf(
				"%s request missing required query parameters.",
				requestName,
			),
		}.Log("").Report(w, http.StatusBadRequest)
		return "", "", false
	}
	return accountName, credential, true
}
