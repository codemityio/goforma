package gen

import "github.com/codemityio/goforma/pkg/code/parser"

// WithUMLGraphGeneratorConfig configuration option.
func WithUMLGraphGeneratorConfig(config *UMLGraphGeneratorConfig) UMLGraphGeneratorOption {
	return func(g *DefaultUMLGraphGenerator) {
		g.config = config
	}
}

// WithUMLGraphGeneratorCodeMap configuration option.
func WithUMLGraphGeneratorCodeMap(
	codeMap *parser.CodeMap[*parser.Var, *parser.Type, *parser.Func, *parser.Const],
) UMLGraphGeneratorOption {
	return func(gen *DefaultUMLGraphGenerator) {
		gen.codeMaps = map[string]*parser.CodeMap[*parser.Var, *parser.Type, *parser.Func, *parser.Const]{}
		gen.types = map[string]struct{}{}
		gen.typesCache = map[string]struct{}{}
		gen.linksCache = map[string]struct{}{}

		// categorise all elements by package path
		for _, val := range codeMap.Var {
			if _, ok := gen.codeMaps[val.PackagePath]; !ok {
				gen.codeMaps[val.PackagePath] = gen.initiateCodeMap()
			}

			gen.codeMaps[val.PackagePath].Var = append(gen.codeMaps[val.PackagePath].Var, val)
		}

		for _, val := range codeMap.Const {
			if _, ok := gen.codeMaps[val.PackagePath]; !ok {
				gen.codeMaps[val.PackagePath] = gen.initiateCodeMap()
			}

			gen.codeMaps[val.PackagePath].Const = append(gen.codeMaps[val.PackagePath].Const, val)
		}

		for _, val := range codeMap.Func {
			if _, ok := gen.codeMaps[val.PackagePath]; !ok {
				gen.codeMaps[val.PackagePath] = gen.initiateCodeMap()
			}

			gen.codeMaps[val.PackagePath].Func = append(gen.codeMaps[val.PackagePath].Func, val)
		}

		for _, val := range codeMap.Type {
			if _, ok := gen.codeMaps[val.PackagePath]; !ok {
				gen.codeMaps[val.PackagePath] = gen.initiateCodeMap()
			}

			gen.codeMaps[val.PackagePath].Type = append(gen.codeMaps[val.PackagePath].Type, val)
			gen.types[val.PackagePath+"."+val.Name] = struct{}{}
		}
	}
}
