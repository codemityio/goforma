package integration

import (
	"code/pkg/config"
	cfg "code/pkg/config"
)

type Config struct {
	URL       config.URL
	Transport *config.Transport
	TLS       *config.TLS
	JWT       cfg.JWT
}
