package doc

import (
	"regexp"
	"strings"
)

// DefaultParser default code doc block parser.
type DefaultParser struct{}

// Parse method to parse code doc block.
func (p *DefaultParser) Parse(lines []string) string {
	var builder strings.Builder

	lineCommentPattern := regexp.MustCompile(`^\s*//\s?`) // pattern to strip off leading "//"

	indent := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			// Strip the leading "//" but preserve the text, adding one-space indentation
			line = lineCommentPattern.ReplaceAllString(line, indent)

			builder.WriteString(line)
			builder.WriteString("\n")

			continue
		}

		builder.WriteString(line)
	}

	return builder.String()
}
