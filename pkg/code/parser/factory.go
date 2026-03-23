package parser

// New factory function.
func New(options ...Option) *DefaultParser {
	dp := &DefaultParser{
		rootPath:   ".",
		docParser:  nil,
		interfaces: map[string]*Interface{},
		types:      map[string]*Type{},
	}

	for _, option := range options {
		option(dp)
	}

	return dp
}
