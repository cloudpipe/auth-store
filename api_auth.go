package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
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

// CreateHandler creates and persists a new account based on a username and password. An error is
// returned if the username is not unique. Otherwise, an accepted status is returned.
func CreateHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		APIError{
			Message: fmt.Sprintf("Unable to parse JSON from your request: %v", err),
		}.Log("").Report(w, http.StatusBadRequest)
		return
	}

	account, err := NewAccount(req.Name, req.Password)
	if err != nil {
		APIError{
			Message: fmt.Sprintf("Unable to create account: %v", err),
		}.Log("").Report(w, http.StatusInternalServerError)
		return
	}

	err = c.Storage.CreateAccount(account)
	if mgo.IsDup(err) {
		APIError{
			Message: fmt.Sprintf(
				`The account name "%s" has already been taken. Please choose another.`,
				req.Name,
			),
		}.Log("").Report(w, http.StatusConflict)
		return
	}
	if err != nil {
		APIError{Message: "Internal storage error."}.Log(req.Name).Report(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}