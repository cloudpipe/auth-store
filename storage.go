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

func (storage *MongoStorage) CreateAccount(account *Account) error {
	return storage.accounts().Insert(account)
}

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
		"$push": bson.M{"apiKeys": key},
	})
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

// Ensure that NullStorage obeys the Storage interface.
var _ Storage = NullStorage{}
