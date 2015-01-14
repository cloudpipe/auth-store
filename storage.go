package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Storage provides high-level interactions with an underlying storage mechanism.
type Storage interface {
	CreateAccount(account *Account) error
	FindAccount(name string) (*Account, error)
	AddKeyToAccount(name, key string) error
	RevokeKeyFromAccount(name, key string) error
	AccountHasKey(name, key string) (bool, error)
}

// MongoStorage is a Storage implementation that connects to a real MongoDB cluster.
type MongoStorage struct {
	Database *mgo.Database
}

// NewMongoStorage establishes a connection to a MongoDB cluster.
func NewMongoStorage(c *Context) (*MongoStorage, error) {
	session, err := mgo.Dial(c.MongoURL)
	if err != nil {
		return nil, err
	}

	return &MongoStorage{Database: session.DB("auth")}, nil
}

func (storage *MongoStorage) accounts() *mgo.Collection {
	return storage.Database.C("accounts")
}

// CreateAccount persists an Account model into Mongo as it's currently populated.
func (storage *MongoStorage) CreateAccount(account *Account) error {
	return storage.accounts().Insert(account)
}

// FindAccount queries for an existing account with a specified name. If no such account exists,
// nil is returned.
func (storage *MongoStorage) FindAccount(name string) (*Account, error) {
	var account Account
	err := storage.accounts().FindId(name).One(&account)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	return &account, err
}

// AddKeyToAccount appends a newly generated API key to an existing account.
func (storage *MongoStorage) AddKeyToAccount(name, key string) error {
	return storage.accounts().UpdateId(name, bson.M{
		"$push": bson.M{"api_keys": key},
	})
}

// RevokeKeyFromAccount removes an API key from an account.
func (storage *MongoStorage) RevokeKeyFromAccount(name, key string) error {
	return storage.accounts().UpdateId(name, bson.M{
		"$pull": bson.M{"api_keys": key},
	})
}

// AccountHasKey returns true if the named account has an associated API key that matches the
// provided one, or false if it does not.
func (storage *MongoStorage) AccountHasKey(name, key string) (bool, error) {
	n, err := storage.accounts().Find(bson.M{
		"_id":      name,
		"api_keys": key,
	}).Count()

	return n == 1, err
}

// NullStorage provides no-op implementations of Storage methods. It's useful for selective
// overriding in unit tests.
type NullStorage struct{}

// CreateAccount is a no-op.
func (storage NullStorage) CreateAccount(*Account) error {
	return nil
}

// FindAccount always fails to find an account.
func (storage NullStorage) FindAccount(name string) (*Account, error) {
	return nil, nil
}

// AddKeyToAccount is a no-op.
func (storage NullStorage) AddKeyToAccount(name, key string) error {
	return nil
}

// RevokeKeyFromAccount is a no-op.
func (storage NullStorage) RevokeKeyFromAccount(name, key string) error {
	return nil
}

// AccountHasKey always returns false.
func (storage NullStorage) AccountHasKey(name, key string) (bool, error) {
	return false, nil
}

// Ensure that NullStorage obeys the Storage interface.
var _ Storage = NullStorage{}
