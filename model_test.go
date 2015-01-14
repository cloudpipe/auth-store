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

func TestGenerateAPIKey(t *testing.T) {
	account, err := NewAccount("sample", "secret")
	if err != nil {
		t.Fatalf("Unable to create account: %v", err)
	}

	if len(account.APIKeys) != 1 {
		t.Errorf("Expected newly created account to have 1 API key, but had %d", len(account.APIKeys))
	}

	key, err := account.GenerateAPIKey()
	if err != nil {
		t.Errorf("Unexpected error generating an API key: %v", err)
	}

	if len(account.APIKeys) != 2 {
		t.Errorf("Expected account to have two API keys, but had %d", len(account.APIKeys))
	}

	if account.APIKeys[1] != key {
		t.Errorf("Expected the generated key [%s] to match the account [%s]", key, account.APIKeys[1])
	}
}
