package main

import (
	"os"
	"testing"
)

func TestLoadFromEnvironment(t *testing.T) {
	c := &Context{}

	os.Setenv("AUTH_PORT", "4321")
	os.Setenv("AUTH_LOGLEVEL", "debug")
	os.Setenv("AUTH_LOGCOLORS", "true")
	os.Setenv("AUTH_MONGOURL", "server.example.com")
	os.Setenv("AUTH_CACERT", "/lockbox/ca.pem")
	os.Setenv("AUTH_CERT", "/lockbox/cert.pem")
	os.Setenv("AUTH_KEY", "/lockbox/key.pem")

	if err := c.Load(); err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}

	if c.Port != 4321 {
		t.Errorf("Unexpected port: [%d]", c.Port)
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

	if c.CACert != "/lockbox/ca.pem" {
		t.Errorf("Unexpected CA certificate path: [%s]", c.CACert)
	}

	if c.Cert != "/lockbox/cert.pem" {
		t.Errorf("Unexpected certificate path: [%s]", c.Cert)
	}

	if c.Key != "/lockbox/key.pem" {
		t.Errorf("Unexpected private key path: [%s]", c.Key)
	}
}

func TestDefaultValues(t *testing.T) {
	c := &Context{}

	os.Setenv("AUTH_PORT", "")
	os.Setenv("AUTH_LOGLEVEL", "")
	os.Setenv("AUTH_LOGCOLORS", "")
	os.Setenv("AUTH_MONGOURL", "")
	os.Setenv("AUTH_CACERT", "")
	os.Setenv("AUTH_CERT", "")
	os.Setenv("AUTH_KEY", "")

	if err := c.Load(); err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}

	if c.Port != 8000 {
		t.Errorf("Unexpected port: [%d]", c.Port)
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

	if c.CACert != "/certificates/ca.pem" {
		t.Errorf("Unexpected CA certificate path: [%s]", c.CACert)
	}

	if c.Cert != "/certificates/auth-store-cert.pem" {
		t.Errorf("Unexpected certificate path: [%s]", c.Cert)
	}

	if c.Key != "/certificates/auth-store-key.pem" {
		t.Errorf("Unexpected private key path: [%s]", c.Key)
	}
}
