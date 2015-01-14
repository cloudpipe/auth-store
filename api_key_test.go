package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type KeyTestStorage struct {
	NullStorage

	NextError error

	FoundAccount *Account

	AccountName *string
	Appended    *string
}

func (storage *KeyTestStorage) FindAccount(name string) (*Account, error) {
	if err := storage.NextError; err != nil {
		storage.NextError = nil
		return nil, err
	}

	return storage.FoundAccount, nil
}

func (storage *KeyTestStorage) AddKeyToAccount(name, key string) error {
	if err := storage.NextError; err != nil {
		storage.NextError = nil
		return err
	}

	storage.AccountName = &name
	storage.Appended = &key
	return nil
}

func TestKeyGenerationSuccess(t *testing.T) {
	r := HTTPRequest(t, "POST", "https://localhost/v1/keys", "")
	r.SetBasicAuth("someone@gmail.com", "secret")
	w := httptest.NewRecorder()
	a, err := NewAccount("someone@gmail.com", "secret")
	if err != nil {
		t.Fatalf("Unable to create account: %v", err)
	}
	s := &KeyTestStorage{FoundAccount: a}
	c := &Context{Storage: s}

	KeyHandler(c, w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected response code %d, but was %d", http.StatusOK, w.Code)
	}

	if ctype := w.HeaderMap.Get("Content-Type"); ctype != "text/plain" {
		t.Errorf("Expected content type of [text/plain], but got [%s]", ctype)
	}

	key := w.Body.String()

	if key == "" {
		t.Error("Expected response to contain an API key")
	}

	t.Logf("Generated API key: [%s]", key)

	if s.AccountName == nil || s.Appended == nil {
		t.Fatal("Expected generated key to be appended to storage")
	}

	if *s.AccountName != "someone@gmail.com" {
		t.Errorf("Expected account [someone@gmail.com] to be modified, but was [%s]", *s.AccountName)
	}

	if *s.Appended != key {
		t.Errorf("Expected API key [%s] to be appended to account, but was [%s]", key, *s.Appended)
	}
}
