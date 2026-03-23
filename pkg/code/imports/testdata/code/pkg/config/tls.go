// Package config contains TLS configuration struct.
package config

// TLS contains server tls configuration.
type TLS struct {
	Cert   string `env:"CERT"    json:"cert"`
	Key    string `env:"KEY"     json:"key"`
	CACert string `env:"CA_CERT" json:"caCert"`
}
