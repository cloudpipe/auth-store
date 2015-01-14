package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Account is a user account.
type Account struct {
	Name           string `json:"name" bson:"_id"`
	HashedPassword []byte `json:"-" bson:"password"`
	Administrator  bool   `json:"admin" bson:"admin"`

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

	return account, nil
}

// HasPassword returns true if the supplied password is correct for the existing account.
func (account *Account) HasPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(account.HashedPassword, []byte(password)) == nil
}
