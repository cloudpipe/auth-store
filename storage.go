package main

import "gopkg.in/mgo.v2"

// Storage provides high-level interactions with an underlying storage mechanism.
type Storage interface{}

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

// NullStorage provides no-op implementations of Storage methods. It's useful for selective
// overriding in unit tests.
type NullStorage struct{}

// Ensure that NullStorage obeys the Storage interface.
var _ Storage = NullStorage{}
