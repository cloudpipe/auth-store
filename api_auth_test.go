package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthTestStorage struct {
	NullStorage

	Created *Account
}

func (storage *AuthTestStorage) CreateAccount(account *Account) error {
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
