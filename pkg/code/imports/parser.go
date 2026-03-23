package imports

import (
	"fmt"
	"go/token"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

type DefaultParser struct {
	rootPath                                        string
	depth                                           int
	packagesMap                                     map[string]*packages.Package
	owned                                           []string
	excludePaths                                    []string
	excludeStandard, excludeVendor, excludeInternal bool
}

func (p *DefaultParser) Parse(path string) (*Packages, error) {
	cfg := &packages.Config{ //nolint:exhaustruct // not required to be exhaustive...
		Mode: packages.NeedImports | packages.NeedDeps | packages.NeedName | packages.NeedFiles,
		Dir:  p.rootPath,
		Fset: token.NewFileSet(),
	}

	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPkgLoad, err)
	}

	// create a map to hold package information
	p.packagesMap = make(map[string]*packages.Package)

	// collect all packages, including dependencies
	p.collectPackages(pkgs, 0, p.depth)

	// build the hierarchy tree
	return p.build(), nil
}

// Collect packages and their dependencies recursively.
func (p *DefaultParser) collectPackages(pkgs []*packages.Package, currentDepth, maxDepth int) {
	if currentDepth > maxDepth {
		return
	}

	for _, pkg := range pkgs {
		if len(pkg.GoFiles) == 0 {
			continue // skip packages with no Go files
		}

		if p.isExclude(p.sanitiseVendorPath(pkg.PkgPath)) {
			continue
		}

		if _, exists := p.packagesMap[pkg.PkgPath]; !exists {
			p.packagesMap[pkg.PkgPath] = pkg

			// convert pkg.Imports to a slice and recursively collect imports
			for _, imp := range pkg.Imports {
				p.collectPackages([]*packages.Package{imp}, currentDepth+1, maxDepth)
			}
		}
	}
}

// Helper function to build the package import hierarchy.
func (p *DefaultParser) build() *Packages {
	list := Packages{
		Paths: []string{},
		List:  []*Package{},
	}

	for _, pkg := range p.packagesMap {
		pkgPath := p.sanitiseVendorPath(pkg.PkgPath)

		if p.isExclude(pkgPath) {
			continue
		}

		list.Paths = append(list.Paths, pkgPath)

		item := Package{
			Name:       pkg.Name,
			Path:       pkgPath,
			Label:      p.sanitisePath(pkg.PkgPath),
			IsLocal:    p.isLocal(filepath.Dir(pkg.GoFiles[0])),
			IsOwned:    p.isOwned(pkg.PkgPath),
			IsExternal: p.isExternal(pkg.PkgPath),
			IsStandard: p.isStandard(pkg.PkgPath),
			IsInternal: p.isInternal(pkg.PkgPath),
			IsVendor:   p.isVendor(pkg.PkgPath),
			DepPaths:   nil,
		}

		for importPath := range pkg.Imports {
			if p.isExclude(importPath) {
				continue
			}

			item.DepPaths = append(item.DepPaths, importPath)
		}

		item.DepPaths = p.dedupeSliceOfString(item.DepPaths)

		sort.Strings(item.DepPaths)

		list.List = append(list.List, &item)
	}

	list.Paths = p.dedupeSliceOfString(list.Paths)

	sort.Strings(list.Paths)

	sort.Slice(list.List, func(i, j int) bool {
		return list.List[i].Path < list.List[j].Path
	})

	return &list
}

// Helper function to check if a package is from the current project library.
func (p *DefaultParser) isLocal(path string) bool {
	cleanRootPath := filepath.Clean(p.rootPath)
	cleanPath := filepath.Clean(path)

	return strings.HasPrefix(cleanPath, cleanRootPath)
}

// Helper function to check if a package is owned.
func (p *DefaultParser) isOwned(path string) bool {
	owned := false

	for _, o := range p.owned {
		if strings.HasPrefix(path, o) {
			owned = true

			break
		}
	}

	return owned
}

// Helper function to check if a package is owned.
func (p *DefaultParser) isExclude(path string) bool {
	for _, o := range p.excludePaths {
		if strings.HasPrefix(path, o) {
			return true
		}
	}

	if p.excludeStandard && p.isStandard(path) {
		return true
	}

	if p.excludeVendor && p.isVendor(path) {
		return true
	}

	if p.excludeInternal && p.isInternal(path) {
		return true
	}

	return false
}

// Helper function to check if a package is from the external library (3rd party).
func (p *DefaultParser) isExternal(path string) bool {
	return !p.isStandard(path) && !p.isLocal(path)
}

// Helper function to check if a package is from the standard library.
func (p *DefaultParser) isStandard(path string) bool {
	return !strings.Contains(path, ".")
}

// Helper function to check if a package is from the internal directory.
func (p *DefaultParser) isInternal(path string) bool {
	re := regexp.MustCompile(`(^|/)internal(/|$)`)

	return re.MatchString(path)
}

// Helper function to check if a package is from the vendor directory.
func (p *DefaultParser) isVendor(path string) bool {
	re := regexp.MustCompile(`^vendor(/|$)`)

	return re.MatchString(path)
}

// Helper function to sanitise path (remove vendor part).
func (p *DefaultParser) sanitiseVendorPath(path string) string {
	if p.isLocal(path) {
		return path
	}

	reVendor := regexp.MustCompile(`^vendor(/|$)`)

	return strings.Trim(filepath.Clean(reVendor.ReplaceAllString(path, "/")), `/`)
}

// Helper function to sanitise path (remove vendor and internal parts for more readability).
func (p *DefaultParser) sanitisePath(path string) string {
	if p.isLocal(path) {
		return path
	}

	reInternal := regexp.MustCompile(`(^|/)internal(/|$)`)

	reVendor := regexp.MustCompile(`^vendor(/|$)`)

	sanitizedPath := reInternal.ReplaceAllString(path, "/")       // replace "internal" with "/"
	sanitizedPath = reVendor.ReplaceAllString(sanitizedPath, "/") // replace "vendor" with "/"

	return strings.Trim(filepath.Clean(sanitizedPath), `/`)
}

// Helper function to deduplicate slice of string.
func (p *DefaultParser) dedupeSliceOfString(input []string) []string {
	unique := make(map[string]struct{})

	result := make([]string, 0)

	for _, val := range input {
		if _, exists := unique[val]; !exists {
			unique[val] = struct{}{}

			result = append(result, val)
		}
	}

	return result
}
