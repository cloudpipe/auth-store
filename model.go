package main

// Account is a user account.
type Account struct {
	Name           string `json:"name" bson:"_id"`
	HashedPassword []byte `json:"-" bson:"password"`
	Administrator  bool   `json:"admin" bson:"admin"`

	CreatedAt uint64 `json:"-" bson:"created_at"`
	UpdatedAt uint64 `json:"-" bson:"updated_at"`
}
