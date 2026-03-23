package parser

// Parser code packages parser.
type Parser interface {
	Parse(pkgPath string) (map[string][]any, error)
}
