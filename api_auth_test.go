package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/mgo.v2"
)

type AuthTestStorage struct {
	NullStorage

	NextError error
	Created   *Account
}

func (storage *AuthTestStorage) CreateAccount(account *Account) error {
	if err := storage.NextError; err != nil {
		storage.NextError = nil
		return err
	}

	storage.Created = account
	return nil
}

func TestCreateHandlerSuccess(t *testing.T) {
	r := HTTPRequest(t, "POST", "https://localhost/v1/accounts", `
	{
		"name": "someone@gmail.com",
		"password": "secret"
	}
	`)
	w := httptest.NewRecorder()
	s := &AuthTestStorage{}
	c := &Context{Storage: s}

	CreateHandler(c, w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected response code %d, but was %d", http.StatusCreated, w.Code)
	}

	if w.Body.Len() != 0 {
		t.Errorf("Expected empty body, but got:<<<\n%s>>>", w.Body.String())
	}

	acct := s.Created

	if acct == nil {
		t.Fatal("Account not created")
	}

	if acct.Name != "someone@gmail.com" {
		t.Errorf("Account had unexpected name: [%s]", acct.Name)
	}
}

func TestCreateHandlerInvalidJSON(t *testing.T) {
	r := HTTPRequest(t, "POST", "https://localhost/v1/accounts", `{ "wat"`)
	w := httptest.NewRecorder()
	c := &Context{}

	CreateHandler(c, w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected response code %d, but was %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateHandlerDuplicateAccount(t *testing.T) {
	r := HTTPRequest(t, "POST", "https://localhost/v1/accounts", `
	{
		"name": "duplicate@gmail.com",
		"password": "wat"
	}
	`)
	w := httptest.NewRecorder()
	s := &AuthTestStorage{
		// See https://github.com/go-mgo/mgo/blob/445c05a1261a0941bc48d898c8eb3ee18ab398c3/session.go#L2116
		NextError: &mgo.QueryError{Code: 11000},
	}
	c := &Context{Storage: s}

	CreateHandler(c, w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected response code %d, but was %d", http.StatusConflict, w.Code)
	}
}

func TestCreateHandlerStorageFailure(t *testing.T) {
	r := HTTPRequest(t, "POST", "https://localhost/v1/accounts", `
	{
		"name": "duplicate@gmail.com",
		"password": "wat"
	}
	`)
	w := httptest.NewRecorder()
	s := &AuthTestStorage{NextError: errors.New("WTF")}
	c := &Context{Storage: s}

	CreateHandler(c, w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected response code %d, but was %d", http.StatusInternalServerError, w.Code)
	}
}
