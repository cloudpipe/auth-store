package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// Context provides shared state among route handlers.
type Context struct {
	Settings

	Storage Storage
}

// Settings contains configuration options loaded from the environment.
type Settings struct {
	Port      int
	LogLevel  string
	LogColors bool
	MongoURL  string
	CACert    string
	Cert      string
	Key       string
}

// Load reads configuration settings from the environment and validates them.
func (c *Context) Load() error {
	if err := envconfig.Process("AUTH", &c.Settings); err != nil {
		return err
	}

	if c.Port == 0 {
		c.Port = 8000
	}

	if c.LogLevel == "" {
		c.LogLevel = "info"
	}

	if c.MongoURL == "" {
		c.MongoURL = "mongo"
	}

	if c.CACert == "" {
		c.CACert = "/certificates/ca.pem"
	}

	if c.Cert == "" {
		c.Cert = "/certificates/auth-store-cert.pem"
	}

	if c.Key == "" {
		c.Key = "/certificates/auth-store-key.pem"
	}

	if _, err := log.ParseLevel(c.LogLevel); err != nil {
		return err
	}

	return nil
}

// NewContext loads configuration from the environment and applies immediate, global settings.
func NewContext() (*Context, error) {
	c := &Context{}

	if err := c.Load(); err != nil {
		return c, err
	}

	// Configure the logging level and formatter.

	level, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		return c, err
	}
	log.SetLevel(level)

	log.SetFormatter(&log.TextFormatter{
		ForceColors: c.LogColors,
	})

	// Summarize the loaded settings.

	log.WithFields(log.Fields{
		"port":           c.Port,
		"logging level":  c.LogLevel,
		"log with color": c.LogColors,
		"mongo URL":      c.MongoURL,
		"CA cert":        c.CACert,
		"cert":           c.Cert,
		"key":            c.Key,
	}).Info("Initializing with loaded settings.")

	// Connect to MongoDB

	c.Storage, err = NewMongoStorage(c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// ListenAddr generates an address to bind the net/http server to based on the current settings.
func (c *Context) ListenAddr() string {
	return fmt.Sprintf(":%d", c.Port)
}
