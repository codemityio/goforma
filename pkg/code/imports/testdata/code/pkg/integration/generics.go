package integration

import "context"

type ExternalCustomFunction[I ExternalInput, O ExternalOutput] func(context.Context, I) (O, error)
