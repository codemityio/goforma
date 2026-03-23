package imports

// WithRootPath configuration option.
func WithRootPath(rootPath string) Option {
	return func(dcp *DefaultParser) {
		dcp.rootPath = rootPath
	}
}

// WithDepth configuration option.
func WithDepth(depth int) Option {
	return func(dcp *DefaultParser) {
		dcp.depth = depth
	}
}

// WithOwned configuration option.
func WithOwned(owned []string) Option {
	return func(dcp *DefaultParser) {
		dcp.owned = owned
	}
}

// WithExcludePaths configuration option.
func WithExcludePaths(paths []string) Option {
	return func(dcp *DefaultParser) {
		dcp.excludePaths = paths
	}
}

// WithExcludeStandard configuration option.
func WithExcludeStandard(exclude bool) Option {
	return func(dcp *DefaultParser) {
		dcp.excludeStandard = exclude
	}
}

// WithExcludeVendor configuration option.
func WithExcludeVendor(exclude bool) Option {
	return func(dcp *DefaultParser) {
		dcp.excludeVendor = exclude
	}
}

// WithExcludeInternal configuration option.
func WithExcludeInternal(exclude bool) Option {
	return func(dcp *DefaultParser) {
		dcp.excludeInternal = exclude
	}
}
