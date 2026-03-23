// Package config contains JWT configuration struct.
package config

// JWT a JWT token signer configuration.
type JWT struct {
	ID      string `env:"ID"       json:"id,omitempty"`
	KeyPath string `env:"KEY_PATH" json:"keyPath,omitempty"`
}
