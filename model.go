package main

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// APIKeyLength determines how large generated API keys are.
const APIKeyLength = 64

// Account is a user account.
type Account struct {
	Name           string `json:"name" bson:"_id"`
	HashedPassword []byte `json:"-" bson:"password"`
	Administrator  bool   `json:"admin" bson:"admin"`

	APIKeys []string `json:"-" bson:"api_keys"`

	CreatedAt int64 `json:"-" bson:"created_at"`
	UpdatedAt int64 `json:"-" bson:"updated_at"`
}

// NewAccount initializes a new Account given a username and password.
func NewAccount(name, password string) (*Account, error) {
	now := time.Now().UnixNano()
	account := &Account{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	account.HashedPassword = hashed

	if _, err = account.GenerateAPIKey(); err != nil {
		return account, err
	}

	return account, nil
}

// HasPassword returns true if the supplied password is correct for the existing account.
func (account *Account) HasPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(account.HashedPassword, []byte(password)) == nil
}

// GenerateAPIKey securely creates an API key and attaches it to the associated account.
func (account *Account) GenerateAPIKey() (string, error) {
	b := make([]byte, APIKeyLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	key := hex.EncodeToString(b)

	account.APIKeys = append(account.APIKeys, key)

	return key, nil
}
