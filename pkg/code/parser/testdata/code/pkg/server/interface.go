package server

import (
	"code/pkg/integration"
	intg "code/pkg/integration"
)

type Configurator interface {
	integration.BytesLoader
	intg.EnvLoader
}
