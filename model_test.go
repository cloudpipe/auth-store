package main

import "testing"

func TestCreateAccount(t *testing.T) {
	account, err := NewAccount("sample", "secret")
	if err != nil {
		t.Fatalf("Unable to create account: %v", err)
	}

	if account.Name != "sample" {
		t.Errorf("Account had unexpected name: %s", account.Name)
	}

	if account.Administrator {
		t.Error("Account unexpectedly had administrator rights.")
	}

	if account.CreatedAt == 0 {
		t.Error("Account did not have creation time populated")
	}

	if account.CreatedAt != account.UpdatedAt {
		t.Errorf("Account creation and update time differ: %d == %d",
			account.CreatedAt, account.UpdatedAt)
	}
}

func TestHasPassword(t *testing.T) {
	account, err := NewAccount("sample", "secret")
	if err != nil {
		t.Fatalf("Unable to create account: %v", err)
	}

	if !account.HasPassword("secret") {
		t.Error("Correct password not accepted")
	}

	if account.HasPassword("wrong") {
		t.Error("Incorrect password accepted")
	}
}
