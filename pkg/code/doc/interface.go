package doc

// Parser code doc block parser.
type Parser interface {
	Parse(lines []string) string
}
