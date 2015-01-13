package main

import "github.com/kelseyhightower/envconfig"

// Context provides shared state among route handlers.
type Context struct {
	Settings
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

	return nil
}

// NewContext loads configuration from the environment and applies immediate, global settings.
func NewContext() (*Context, error) {
	c := &Context{}

	return c, nil
}
