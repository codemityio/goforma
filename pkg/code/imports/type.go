package imports

// Option function.
type Option func(p *DefaultParser)

// Package item structure.
type Package struct {
	Name, Path, Label string
	// IsLocal indicates the package is part of the scanned project.
	IsLocal, IsOwned                             bool
	IsExternal, IsStandard, IsInternal, IsVendor bool
	DepPaths                                     []string
}

// Packages holding all information about the used packages.
type Packages struct {
	Paths []string
	List  []*Package
}
