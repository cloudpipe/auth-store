package main

import (
	"os"
	"testing"
)

func TestLoadFromEnvironment(t *testing.T) {
	c := &Context{}

	os.Setenv("AUTH_INTERNALPORT", "1111")
	os.Setenv("AUTH_EXTERNALPORT", "2222")
	os.Setenv("AUTH_LOGLEVEL", "debug")
	os.Setenv("AUTH_LOGCOLORS", "true")
	os.Setenv("AUTH_MONGOURL", "server.example.com")
	os.Setenv("AUTH_INTERNALCACERT", "/lockbox/internal-ca.pem")
	os.Setenv("AUTH_INTERNALCERT", "/lockbox/internal-cert.pem")
	os.Setenv("AUTH_INTERNALKEY", "/lockbox/internal-key.pem")
	os.Setenv("AUTH_EXTERNALCERT", "/lockbox/external-cert.pem")
	os.Setenv("AUTH_EXTERNALKEY", "/lockbox/external-key.pem")

	if err := c.Load(); err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}

	if c.InternalPort != 1111 {
		t.Errorf("Unexpected internal port: [%d]", c.InternalPort)
	}

	if c.ExternalPort != 2222 {
		t.Errorf("Unexpected external port: [%d]", c.ExternalPort)
	}

	if c.LogLevel != "debug" {
		t.Errorf("Unexpected log level: [%s]", c.LogLevel)
	}

	if !c.LogColors {
		t.Error("Expected log coloring to be enabled")
	}

	if c.MongoURL != "server.example.com" {
		t.Errorf("Unexpected MongoDB URL: [%s]", c.MongoURL)
	}

	if c.InternalCACert != "/lockbox/internal-ca.pem" {
		t.Errorf("Unexpected internal CA certificate path: [%s]", c.InternalCACert)
	}

	if c.InternalCert != "/lockbox/internal-cert.pem" {
		t.Errorf("Unexpected internal certificate path: [%s]", c.InternalCert)
	}

	if c.InternalKey != "/lockbox/internal-key.pem" {
		t.Errorf("Unexpected internal private key path: [%s]", c.InternalKey)
	}

	if c.ExternalCert != "/lockbox/external-cert.pem" {
		t.Errorf("Unexpected external certificate path: [%s]", c.ExternalCert)
	}

	if c.ExternalKey != "/lockbox/external-key.pem" {
		t.Errorf("Unexpected external private key path: [%s]", c.ExternalKey)
	}
}

func TestDefaultValues(t *testing.T) {
	c := &Context{}

	os.Setenv("AUTH_INTERNALPORT", "")
	os.Setenv("AUTH_EXTERNALPORT", "")
	os.Setenv("AUTH_LOGLEVEL", "")
	os.Setenv("AUTH_LOGCOLORS", "")
	os.Setenv("AUTH_MONGOURL", "")
	os.Setenv("AUTH_INTERNALCACERT", "")
	os.Setenv("AUTH_INTERNALCERT", "")
	os.Setenv("AUTH_INTERNALKEY", "")
	os.Setenv("AUTH_EXTERNALCERT", "")
	os.Setenv("AUTH_EXTERNALKEY", "")

	if err := c.Load(); err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}

	if c.InternalPort != 8001 {
		t.Errorf("Unexpected internal port: [%d]", c.InternalPort)
	}

	if c.ExternalPort != 8000 {
		t.Errorf("Unexpected external port: [%d]", c.ExternalPort)
	}

	if c.LogLevel != "info" {
		t.Errorf("Unexpected log level: [%s]", c.LogLevel)
	}

	if c.LogColors {
		t.Error("Expected log coloring to be disabled by default")
	}

	if c.MongoURL != "mongo" {
		t.Errorf("Unexpected MongoDB URL: [%s]", c.MongoURL)
	}

	if c.InternalCACert != "/certificates/ca.pem" {
		t.Errorf("Unexpected internal CA certificate path: [%s]", c.InternalCACert)
	}

	if c.InternalCert != "/certificates/auth-store-cert.pem" {
		t.Errorf("Unexpected internal certificate path: [%s]", c.InternalCert)
	}

	if c.InternalKey != "/certificates/auth-store-key.pem" {
		t.Errorf("Unexpected internal private key path: [%s]", c.InternalKey)
	}

	if c.ExternalCert != "/certificates/external-cert.pem" {
		t.Errorf("Unexpected external certificate path: [%s]", c.ExternalCert)
	}

	if c.ExternalKey != "/certificates/external-key.pem" {
		t.Errorf("Unexpected external private key: [%s]", c.ExternalKey)
	}
}
