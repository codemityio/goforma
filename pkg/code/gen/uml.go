package gen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/codemityio/goforma/pkg/code/parser"
)

//go:embed uml.puml.tpl
var umlPlantumlTpl string

// DefaultUMLGraphGenerator default generator implementation.
type DefaultUMLGraphGenerator struct {
	codeMaps       map[string]*parser.CodeMap[*parser.Var, *parser.Type, *parser.Func, *parser.Const]
	types          map[string]struct{}
	typesCache     map[string]struct{}
	linksCache     map[string]struct{}
	primitiveTypes map[string]bool
	config         *UMLGraphGeneratorConfig
}

// Generate generate UML diagram.
func (g *DefaultUMLGraphGenerator) Generate() (string, error) {
	// create a new template and parse the letter into it
	tmpl, err := template.New("plantuml").Funcs(template.FuncMap{
		"mapColour":         g.mapColour,
		"concat":            g.concat,
		"isFullPath":        g.isFullPath,
		"isExported":        g.isExported,
		"getFullPathParts":  g.getFullPathParts,
		"generateElementID": g.generateElementID,
		"typeExists":        g.typeExists,
		"typeInCache":       g.typeInCache,
		"linkInCache":       g.linkInCache,
		"isAny":             g.isAny,
		"isInterface":       g.isInterface,
		"isError":           g.isError,
		"isStruct":          g.isStruct,
		"isPrimitive":       g.isPrimitive,
		"isComposite":       g.isComposite,
		"isPointer":         g.isPointer,
		"isSelector":        g.isSelector,
		"getTypeInitial":    g.getTypeInitial,
		"mapTypeColour":     g.mapTypeColour,
		"skipLegend":        g.skipLegend,
		"skipPrimitive":     g.skipPrimitive,
		"skipVar":           g.skipVar,
		"skipConst":         g.skipConst,
		"skipFunc":          g.skipFunc,
		"skipNotExported":   g.skipNotExported,
		"skipDoc":           g.skipDoc,
		"ternary":           g.ternary,
	}).Parse(umlPlantumlTpl)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateParse, err)
	}

	var buf bytes.Buffer

	// execute the template, passing in the data structure
	err = tmpl.Execute(&buf, g.codeMaps)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateExecute, err)
	}

	return buf.String(), nil
}

func (g *DefaultUMLGraphGenerator) initiateCodeMap() *parser.CodeMap[
	*parser.Var, *parser.Type, *parser.Func, *parser.Const,
] {
	codeMap := parser.CodeMap[*parser.Var, *parser.Type, *parser.Func, *parser.Const]{
		Var:   make([]*parser.Var, 0),
		Const: make([]*parser.Const, 0),
		Func:  make([]*parser.Func, 0),
		Type:  make([]*parser.Type, 0),
	}

	return &codeMap
}

func (g *DefaultUMLGraphGenerator) generateElementID(input string) string {
	re := regexp.MustCompile(`[.,\- /()*\[\]]`)

	return re.ReplaceAllString(input, "_")
}

func (g *DefaultUMLGraphGenerator) concat(input ...string) string {
	return strings.Join(input, "")
}

func (g *DefaultUMLGraphGenerator) typeExists(input string) bool {
	key := strings.Trim(input, ".")

	if _, ok := g.types[key]; ok {
		return true
	}

	return false
}

func (g *DefaultUMLGraphGenerator) typeInCache(input string) bool {
	input = strings.ReplaceAll(input, " ", "")

	if _, ok := g.typesCache[input]; ok {
		return true
	}

	g.typesCache[input] = struct{}{}

	return false
}

func (g *DefaultUMLGraphGenerator) linkInCache(input string) bool {
	input = strings.ReplaceAll(input, " ", "")

	if _, ok := g.linksCache[input]; ok {
		return true
	}

	g.linksCache[input] = struct{}{}

	return false
}

func (g *DefaultUMLGraphGenerator) getTypeInitial(input string) string {
	if len(input) == 0 {
		return "?"
	}

	if g.isPrimitive(input) {
		return "P"
	}

	if g.isInterface(input) || g.isError(input) {
		return "I"
	}

	if g.isStruct(input) {
		return "S"
	}

	if g.isVar(input) {
		return "V"
	}

	if g.isConst(input) {
		return "C"
	}

	if g.isFunc(input) {
		return "F"
	}

	return "T"
}

func (g *DefaultUMLGraphGenerator) mapTypeColour(input string) string {
	typeCheckers := []struct {
		checker func(string) bool
		colour  string
	}{
		{g.isPrimitive, ColourPrimitive},
		{g.isComposite, ColourComposite},
		{g.isPointer, ColourPointer},
		{g.isInterface, ColourInterface},
		{g.isError, ColourInterface},
		{g.isStruct, ColourStruct},
		{g.isVar, ColourVar},
		{g.isConst, ColourConst},
		{g.isFunc, ColourFunc},
		{g.isExternal, ColourExternal},
	}

	for _, typeChecker := range typeCheckers {
		if typeChecker.checker(input) {
			return typeChecker.colour
		}
	}

	return ColourType
}

func (g *DefaultUMLGraphGenerator) mapColour(input string) string {
	switch input {
	case "name":
		return ColourName
	case "embedded":
		return ColourEmbedded
	case "fieldType":
		return ColourFieldType
	case "methodSignature":
		return ColourMethodSignature
	case "funcSignature":
		return ColourFuncSignature
	default:
		return ColourWhite
	}
}

func (g *DefaultUMLGraphGenerator) isPrimitive(input string) bool {
	return g.primitiveTypes[input]
}

func (g *DefaultUMLGraphGenerator) isComposite(input string) bool {
	if g.isPointer(input) {
		input = strings.TrimPrefix(input, "*")
	}

	// check if it's a primitive type
	if g.primitiveTypes[input] {
		return false
	}

	// check if it's a composite type by specific patterns
	switch {
	case strings.HasPrefix(input, "[]"): // slice type
		return true
	case strings.HasPrefix(input, "map["): // map type
		return true
	case strings.HasPrefix(input, "chan "): // channel type
		return true
	case strings.HasPrefix(input, "func("): // function type
		return true
	default:
		return false
	}
}

func (g *DefaultUMLGraphGenerator) isPointer(input string) bool {
	return strings.HasPrefix(input, "*")
}

func (g *DefaultUMLGraphGenerator) isAny(input string) bool {
	return input == "any"
}

func (g *DefaultUMLGraphGenerator) isInterface(input string) bool {
	return input == "interface" || strings.HasPrefix(input, "interface{")
}

func (g *DefaultUMLGraphGenerator) isError(input string) bool {
	return input == "error"
}

func (g *DefaultUMLGraphGenerator) isStruct(input string) bool {
	return input == "struct" || strings.HasPrefix(input, "struct{")
}

func (g *DefaultUMLGraphGenerator) isVar(input string) bool {
	return input == "var"
}

func (g *DefaultUMLGraphGenerator) isConst(input string) bool {
	return input == "const"
}

func (g *DefaultUMLGraphGenerator) isFunc(input string) bool {
	return input == "func" || strings.HasPrefix(input, "func(")
}

func (g *DefaultUMLGraphGenerator) isExternal(input string) bool {
	return input == "external"
}

func (g *DefaultUMLGraphGenerator) isSelector(input string) bool {
	// if there's a generic part, it must be at the end
	if i := strings.Index(input, "["); i != -1 {
		if !strings.HasSuffix(input, "]") {
			return false
		}

		input = input[:i]
	}

	parts := strings.Split(input, ".")

	// a selector should have exactly one dot, so it must split into two parts
	if len(parts) != selectorPartsLimit {
		return false
	}

	// both parts should be valid Go identifiers
	return g.isValidIdentifier(parts[0]) && g.isValidIdentifier(parts[1])
}

func (g *DefaultUMLGraphGenerator) isValidIdentifier(input string) bool {
	if input == "" {
		return false
	}

	// the first character must be a letter or underscore
	if !unicode.IsLetter(rune(input[0])) && input[0] != '_' {
		return false
	}

	// the rest can be letters, digits, or underscores
	for _, ch := range input[1:] {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false
		}
	}

	return true
}

func (g *DefaultUMLGraphGenerator) isFullPath(input string) bool {
	sel := regexp.MustCompile(regexFullPath)

	return sel.MatchString(input)
}

func (g *DefaultUMLGraphGenerator) isExported(input string) bool {
	return ast.IsExported(input)
}

func (g *DefaultUMLGraphGenerator) getFullPathParts(input string) *TypeInfoParts {
	sel := regexp.MustCompile(regexFullPath)

	if matches := sel.FindStringSubmatch(input); matches != nil {
		return &TypeInfoParts{
			PackagePath: strings.TrimRight(matches[1], "."),
			Name:        matches[2],
		}
	}

	return nil
}

func (g *DefaultUMLGraphGenerator) skipLegend() bool {
	if g.config == nil {
		return true
	}

	return !g.config.Legend
}

func (g *DefaultUMLGraphGenerator) skipPrimitive() bool {
	if g.config == nil {
		return true
	}

	return !g.config.Primitive
}

func (g *DefaultUMLGraphGenerator) skipVar() bool {
	if g.config == nil {
		return true
	}

	return !g.config.Var
}

func (g *DefaultUMLGraphGenerator) skipConst() bool {
	if g.config == nil {
		return true
	}

	return !g.config.Const
}

func (g *DefaultUMLGraphGenerator) skipFunc() bool {
	if g.config == nil {
		return true
	}

	return !g.config.Func
}

func (g *DefaultUMLGraphGenerator) skipNotExported() bool {
	if g.config == nil {
		return true
	}

	return !g.config.NotExported
}

func (g *DefaultUMLGraphGenerator) skipDoc() bool {
	if g.config == nil {
		return true
	}

	return !g.config.Doc
}

func (g *DefaultUMLGraphGenerator) ternary(cond bool, trueVal, falseVal string) string {
	if cond {
		return trueVal
	}

	return falseVal
}
