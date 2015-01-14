package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// AccountHandler dispatches requests to handlers that manage the /account resource based on
// request method.
func AccountHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		CreateHandler(c, w, r)
	default:
		APIError{
			Message: fmt.Sprintf("Unsupported method %s. Only POST is accepted for this resource.",
				r.Method),
		}.Log("").Report(w, http.StatusMethodNotAllowed)
	}
}

// CreateHandler creates and persists a new account based on a username and password. An error is
// returned if the username is not unique. Otherwise, an accepted status is returned.
func CreateHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	accountName, password, ok := ExtractPasswordCredentials(w, r, "Account creation")
	if !ok {
		return
	}

	account, err := NewAccount(accountName, password)
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
				accountName,
			),
		}.Log("").Report(w, http.StatusConflict)
		return
	}
	if err != nil {
		APIError{Message: "Internal storage error."}.Log(accountName).Report(w, http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"account": accountName,
	}).Info("Account created successfully.")

	w.WriteHeader(http.StatusCreated)
}
