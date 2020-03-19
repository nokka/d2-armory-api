package env

import (
	"strconv"
	"syscall"
)

// DefaultClient is the default client backed by syscall.Getenv
var DefaultClient = &Client{syscall.Getenv}

// Client for retrieving environment variables.
type Client struct {
	Getenv func(string) (string, bool)
}

// Bool returns a bool from the environment, error if unable to parse to bool, or fallback if not set.
func (c *Client) Bool(key string, fallback bool) (bool, error) {
	s, ok := c.Getenv(key)
	if !ok {
		return fallback, nil
	}

	return strconv.ParseBool(s)
}

// Int returns an int from the environment, error if unable to parse to int, or fallback if not set.
func (c *Client) Int(key string, fallback int) (int, error) {
	s, ok := c.Getenv(key)
	if !ok {
		return fallback, nil
	}

	return strconv.Atoi(s)
}

// String returns a string from the environment, or fallback if not set.
func (c *Client) String(key, fallback string) string {
	if value, ok := c.Getenv(key); ok {
		return value
	}

	return fallback
}

// Bool returns a bool from the environment, error if unable to parse to bool, or fallback if not set.
func Bool(key string, fallback bool) (bool, error) {
	return DefaultClient.Bool(key, fallback)
}

// Int returns an int from the environment, error if unable to parse to int, or fallback if not set.
func Int(key string, fallback int) (int, error) {
	return DefaultClient.Int(key, fallback)
}

// String returns a string from the environment, or fallback if not set.
func String(key, fallback string) string {
	return DefaultClient.String(key, fallback)
}
