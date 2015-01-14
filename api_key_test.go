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
	Revoked     *string
}

func (storage *KeyTestStorage) consumeError() error {
	err := storage.NextError
	storage.NextError = nil
	return err
}

func (storage *KeyTestStorage) FindAccount(name string) (*Account, error) {
	if err := storage.consumeError(); err != nil {
		return nil, err
	}

	return storage.FoundAccount, nil
}

func (storage *KeyTestStorage) AddKeyToAccount(name, key string) error {
	if err := storage.consumeError(); err != nil {
		return err
	}

	storage.AccountName = &name
	storage.Appended = &key
	return nil
}

func (storage *KeyTestStorage) RevokeKeyFromAccount(name, key string) error {
	if err := storage.consumeError(); err != nil {
		return err
	}

	storage.AccountName = &name
	storage.Revoked = &key
	return nil
}

func TestKeyGenerationSuccess(t *testing.T) {
	r := HTTPRequest(t, "POST", "https://localhost/v1/keys",
		`accountName=someone%40gmail.com&password=secret`)
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

func TestKeyGenerationBadPassword(t *testing.T) {
	r := HTTPRequest(t, "POST", "https://localhost/v1/keys",
		`accountName=someone%40gmail.com&password=wrongwrongwrong`)
	w := httptest.NewRecorder()
	a, err := NewAccount("someone@gmail.com", "correct")
	if err != nil {
		t.Fatalf("Unable to create account: %v", err)
	}
	s := &KeyTestStorage{FoundAccount: a}
	c := &Context{Storage: s}

	KeyHandler(c, w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected response code %d, but was %d", http.StatusUnauthorized, w.Code)
	}
}

func TestKeyGenerationBadAccountName(t *testing.T) {
	r := HTTPRequest(t, "POST", "https://localhost/v1/keys",
		`accountName=unknown%40gmail.com&password=vacuouslytrue`)
	w := httptest.NewRecorder()
	c := &Context{Storage: &KeyTestStorage{}}

	KeyHandler(c, w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected response code %d, but was %d", http.StatusUnauthorized, w.Code)
	}
}

func TestKeyRevocationSuccess(t *testing.T) {
	r := HTTPRequest(t, "DELETE", "https://localhost/v1/keys?accountName=someone&apiKey=123abc", "")
	w := httptest.NewRecorder()
	a, err := NewAccount("someone", "secret")
	if err != nil {
		t.Fatalf("Unable to create account: %v", err)
	}
	s := &KeyTestStorage{FoundAccount: a}
	c := &Context{Storage: s}

	KeyRevocationHandler(c, w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected response code %d, but was %d", http.StatusNoContent, w.Code)
	}

	if s.AccountName == nil || s.Revoked == nil {
		t.Fatal("Expected storage to process key revocation")
	}

	if *s.AccountName != "someone" {
		t.Errorf("Unexpected account name [%s]", *s.AccountName)
	}

	if *s.Revoked != "123abc" {
		t.Errorf("Unexpected revoked key [%s]", *s.Revoked)
	}
}
