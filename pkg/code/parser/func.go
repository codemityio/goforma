package parser

import "github.com/codemityio/goforma/pkg/code/doc"

// WithRootPath configuration option.
func WithRootPath(rootPath string) Option {
	return func(dcp *DefaultParser) {
		dcp.rootPath = rootPath
	}
}

// WithDocParser configuration option.
func WithDocParser(docParser doc.Parser) Option {
	return func(dcp *DefaultParser) {
		dcp.docParser = docParser
	}
}
