// Package config contains Client configuration struct.
package config

// Transport contains client transport configuration.
type Transport struct {
	MaxConnsPerHost     int64 `env:"MAX_CONNS_PER_HOST"      json:"maxConnsPerHost"`
	MaxIdleConns        int64 `env:"MAX_IDLE_CONNS"          json:"maxIdleConns"`
	MaxIdleConnsPerHost int64 `env:"MAX_IDLE_CONNS_PER_HOST" json:"maxIdleConnsPerHost"`
	TimeoutInSeconds    int64 `env:"TIMEOUT_IN_SECONDS"      json:"timeoutInSeconds"`
}
