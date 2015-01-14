package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type ValidateTestStorage struct {
	NullStorage

	Accept bool
}

func (storage *ValidateTestStorage) AccountHasKey(name, key string) (bool, error) {
	return storage.Accept, nil
}

func TestValidateHandlerSuccess(t *testing.T) {
	r := HTTPRequest(t, "GET", "https://localhost/v1/validate?accountName=someone&apiKey=ff01ab", "")
	w := httptest.NewRecorder()
	s := &ValidateTestStorage{Accept: true}
	c := &Context{Storage: s}

	ValidateHandler(c, w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected response code %d, but was %d", http.StatusNoContent, w.Code)
	}
}

func TestValidateHandlerReject(t *testing.T) {
	r := HTTPRequest(t, "GET", "https://localhost/v1/validate?accountName=someone&apiKey=ff01ab", "")
	w := httptest.NewRecorder()
	s := &ValidateTestStorage{Accept: false}
	c := &Context{Storage: s}

	ValidateHandler(c, w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected response code %d, but was %d", http.StatusNotFound, w.Code)
	}
}
