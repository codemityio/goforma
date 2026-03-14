package imports

import "golang.org/x/tools/go/packages"

// New factory function.
func New(options ...Option) *DefaultParser {
	dcp := &DefaultParser{
		rootPath:        ".",
		depth:           DefaultDepth,
		packagesMap:     map[string]*packages.Package{},
		owned:           []string{},
		excludePaths:    []string{},
		excludeStandard: false,
		excludeVendor:   false,
		excludeInternal: false,
	}

	for _, option := range options {
		option(dcp)
	}

	return dcp
}
