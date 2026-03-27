package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/codemityio/goforma/pkg/code/doc"
	"golang.org/x/tools/go/packages"
)

// DefaultParser default parser service.
type DefaultParser struct {
	rootPath   string
	docParser  doc.Parser
	interfaces map[string]*Interface
	// types used to perform type lookup
	types map[string]*Type
}

// Parse function to scan and parse a provided path.
func (p *DefaultParser) Parse(path string) (*CodeMap[*Var, *Type, *Func, *Const], error) {
	output := CodeMap[*Var, *Type, *Func, *Const]{
		Var:   []*Var{},
		Type:  []*Type{},
		Func:  []*Func{},
		Const: []*Const{},
	}

	pkgs, err := packages.Load(
		&packages.Config{ //nolint:exhaustruct // not required to be exhaustive...
			Mode: packages.NeedName |
				packages.NeedFiles |
				packages.NeedImports |
				packages.NeedDeps |
				packages.NeedTypes |
				packages.NeedSyntax |
				packages.NeedTypesInfo |
				packages.NeedModule,
			Dir: p.rootPath,
		},
		path,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPkgLoad, err)
	}

	errs := make(chan error)

	go func() {
		p.interfaces = p.collectAllInterfaces(pkgs)

		for _, pkg := range pkgs {
			p.parse(&output, pkg, errs)
		}

		for key := range p.types {
			p.enrichType(p.types[key])

			// dedupe
			p.types[key].Fields = p.dedupeSliceOfFields(p.types[key].Fields)
			p.types[key].Interfaces = p.dedupeSliceOfInterfaces(p.types[key].Interfaces)
			p.types[key].Methods = p.dedupeSliceOfMethods(p.types[key].Methods)
			p.types[key].Embedded = p.dedupeSliceOfTypeDesc(p.types[key].Embedded)
		}

		errs <- nil
	}()

	err = <-errs
	if err != nil {
		return nil, err
	}

	return &output, nil
}

// Helper method to collect all interfaces.
func (p *DefaultParser) collectAllInterfaces(pkgs []*packages.Package) map[string]*Interface {
	interfaces := map[string]*Interface{}

	for _, pkg := range pkgs {
		for _, name := range pkg.Types.Scope().Names() {
			obj := pkg.Types.Scope().Lookup(name)

			typeName, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}

			if iface, ok := typeName.Type().Underlying().(*types.Interface); ok {
				item := Interface{
					Name:        name,
					PackageName: pkg.Name,
					PackagePath: pkg.PkgPath,
					iface:       iface,
				}

				interfaces[pkg.PkgPath+"."+name] = &item
			}
		}
	}

	return interfaces
}

func (p *DefaultParser) parse(
	output *CodeMap[*Var, *Type, *Func, *Const], pkg *packages.Package, errs chan error,
) {
	// collect all methods defined with receivers for the current package
	packageMethods := p.collectAllPackageMethods(pkg, errs)

	for _, syn := range pkg.Syntax {
		for _, decl := range syn.Decls {
			p.inspect(output, pkg, decl, packageMethods, errs)
		}
	}
}

func (p *DefaultParser) inspect(
	output *CodeMap[*Var, *Type, *Func, *Const],
	pkg *packages.Package, decl ast.Decl,
	packageMethods map[string][]*Method,
	errs chan error,
) {
	ast.Inspect(decl, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		start := pkg.Fset.Position(node.Pos())
		end := pkg.Fset.Position(node.End())

		rootPath, err := p.normalizeRootPath(start.Filename)
		if err != nil {
			errs <- err

			return true
		}

		relativePath, err := filepath.Rel(rootPath, start.Filename)
		if err != nil {
			errs <- fmt.Errorf(
				"%w: %w: with root path `%s` and file name `%s`",
				ErrGetRelPath, err, rootPath, start.Filename,
			)

			return true
		}

		// functions defined within a package not having any receiver attached...
		if funcDecl, ok := node.(*ast.FuncDecl); ok && funcDecl.Recv == nil {
			output.Func = append(
				output.Func,
				p.describeFuncDecl(relativePath, pkg, start, end, funcDecl),
			)
		}

		p.inspectSwitch(output, relativePath, pkg, start, end, node, packageMethods, errs)

		return true
	})
}

func (p *DefaultParser) inspectSwitch(
	output *CodeMap[*Var, *Type, *Func, *Const],
	path string, pkg *packages.Package, start, end token.Position, node ast.Node,
	packageMethods map[string][]*Method,
	errs chan error,
) {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return
	}

	if !p.isPackageScope(start, pkg) {
		return
	}

	switch genDecl.Tok { //nolint:exhaustive // not required...
	case token.VAR:
		for _, spec := range genDecl.Specs {
			if varSpec, ok := spec.(*ast.ValueSpec); ok {
				output.Var = append(
					output.Var,
					p.describeVarDecl(path, pkg, start, end, varSpec)...,
				)
			}
		}
	case token.CONST:
		for _, spec := range genDecl.Specs {
			if constSpec, ok := spec.(*ast.ValueSpec); ok {
				output.Const = append(
					output.Const,
					p.describeConstDecl(path, pkg, start, end, constSpec)...,
				)
			}
		}
	default:
		for _, spec := range genDecl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				tp := p.describeTypeDecl(path, pkg, start, end, typeSpec, packageMethods, errs)

				p.types[tp.PackagePath+"."+tp.Name] = tp

				output.Type = append(output.Type, tp)
			}
		}
	}
}

// Helper method to check if the checked item is in the global scope using token positions.
func (p *DefaultParser) isPackageScope(pos token.Position, pkg *packages.Package) bool {
	isGlobal := true

	for _, file := range pkg.Syntax {
		fset := pkg.Fset

		if fset.Position(file.Pos()).Filename != pos.Filename {
			continue
		}

		// traverse the entire AST for this file to check for function scopes
		ast.Inspect(file, func(node ast.Node) bool {
			switch nt := node.(type) {
			case *ast.FuncDecl:
				// check if pos falls within a named function’s range
				if pos.Line > fset.Position(nt.Pos()).Line &&
					pos.Line < fset.Position(nt.End()).Line {
					isGlobal = false

					return false // stop further inspection within this function
				}
			case *ast.FuncLit:
				// Check if pos falls within an anonymous function (closure) range
				if pos.Line > fset.Position(nt.Pos()).Line &&
					pos.Line < fset.Position(nt.End()).Line {
					isGlobal = false

					return false // stop further inspection within this function literal
				}
			}

			return true // continue inspection
		})

		if !isGlobal {
			break // Stop once we confirm it's not global
		}
	}

	// if the variable is not within any function's position range, it’s in package scope
	return isGlobal
}

// Helper method to describe functions declaration.
func (p *DefaultParser) describeFuncDecl(
	path string, pkg *packages.Package, start, end token.Position, funcDecl *ast.FuncDecl,
) *Func {
	item := Func{
		FilePath: path,
		Start: &Position{
			Offset: start.Offset,
			Line:   start.Line,
			Column: start.Column,
		},
		End: &Position{
			Offset: end.Offset,
			Line:   end.Line,
			Column: end.Column,
		},
		PackageName: pkg.Name,
		PackagePath: pkg.PkgPath,
		Name:        funcDecl.Name.Name,
		Signature:   p.getFuncSignature(funcDecl.Type, pkg.TypesInfo),
		Params:      nil,
		IsExported:  ast.IsExported(funcDecl.Name.Name),
		Doc:         "",
	}

	item.Params = p.getFuncTypeParams(funcDecl, pkg.TypesInfo)

	if funcDecl.Doc != nil {
		item.Doc = p.getDoc(funcDecl.Doc)
	}

	return &item
}

// Helper method to describe variables declaration.
func (p *DefaultParser) describeVarDecl(
	path string, pkg *packages.Package, start, end token.Position, varSpec *ast.ValueSpec,
) []*Var {
	items := make([]*Var, 0, len(varSpec.Names))

	for _, name := range varSpec.Names {
		item := Var{
			FilePath: path,
			Start: &Position{
				Offset: start.Offset,
				Line:   start.Line,
				Column: start.Column,
			},
			End: &Position{
				Offset: end.Offset,
				Line:   end.Line,
				Column: end.Column,
			},
			PackageName: pkg.Name,
			PackagePath: pkg.PkgPath,
			Name:        name.Name,
			Type:        nil,
			TypeInfo:    "",
			Params:      nil,
			IsExported:  ast.IsExported(name.Name),
			Doc:         "",
		}

		item.Type, item.TypeInfo, item.Doc = p.resolveTypeInfo(
			varSpec.Type,
			varSpec.Doc,
			name,
			pkg.TypesInfo,
		)

		item.Params = p.getVarParamsWithConstraints(varSpec, pkg.TypesInfo)

		items = append(items, &item)
	}

	return items
}

// Helper method to describe constants declaration.
func (p *DefaultParser) describeConstDecl(
	path string, pkg *packages.Package, start, end token.Position, constSpec *ast.ValueSpec,
) []*Const {
	items := make([]*Const, 0, len(constSpec.Names))

	for _, name := range constSpec.Names {
		item := Const{
			FilePath: path,
			Start: &Position{
				Offset: start.Offset,
				Line:   start.Line,
				Column: start.Column,
			},
			End: &Position{
				Offset: end.Offset,
				Line:   end.Line,
				Column: end.Column,
			},
			PackageName: pkg.Name,
			PackagePath: pkg.PkgPath,
			Name:        name.Name,
			Type:        nil,
			TypeInfo:    "",
			IsExported:  ast.IsExported(name.Name),
			Doc:         "",
		}

		item.Type, item.TypeInfo, item.Doc = p.resolveTypeInfo(
			constSpec.Type,
			constSpec.Doc,
			name,
			pkg.TypesInfo,
		)

		items = append(items, &item)
	}

	return items
}

func (p *DefaultParser) resolveTypeInfo( //nolint:nonamedreturns
	specType ast.Expr,
	doc *ast.CommentGroup,
	name *ast.Ident,
	typesInfo *types.Info,
) (typ *TypeDesc, typeInfo string, docStr string) {
	if specType != nil {
		typ = p.describeType(specType, typesInfo)
	}

	if obj, ok := typesInfo.Defs[name]; ok && obj != nil {
		typeInfo = obj.Type().String()
	}

	if doc != nil {
		docStr = p.getDoc(doc)
	}

	return
}

// Helper method to describe types declaration.
func (p *DefaultParser) describeTypeDecl(
	path string, pkg *packages.Package, start, end token.Position, typeSpec *ast.TypeSpec,
	packageMethods map[string][]*Method,
	errs chan error,
) *Type {
	// handle types (structs, interfaces, function types, aliases)
	item := Type{
		FilePath: path,
		Start: &Position{
			Offset: start.Offset,
			Line:   start.Line,
			Column: start.Column,
		},
		End: &Position{
			Offset: end.Offset,
			Line:   end.Line,
			Column: end.Column,
		},
		PackageName: pkg.Name,
		PackagePath: pkg.PkgPath,
		Name:        typeSpec.Name.Name,
		Type:        p.describeType(typeSpec.Type, pkg.TypesInfo),
		Interfaces:  make([]*Interface, 0),
		TypeInfo:    "",
		IsTypeAlias: false,
		Embedded:    make([]*TypeDesc, 0),
		Params:      nil,
		Fields:      make([]*Field, 0),
		Methods:     make([]*Method, 0),
		IsExported:  ast.IsExported(typeSpec.Name.Name),
		Doc:         "",
	}

	item.IsTypeAlias = typeSpec.Assign != token.NoPos

	// load the type information for the package
	if obj, ok := pkg.TypesInfo.Defs[typeSpec.Name]; ok && obj != nil {
		item.TypeInfo = obj.Type().String()

		if namedType, ok := obj.Type().(*types.Named); ok {
			item.Interfaces = p.getInterfaces(namedType)
		}
	}

	item.Embedded = p.getEmbeddedTypes(typeSpec, pkg.TypesInfo)

	if typeSpec.Doc != nil {
		item.Doc = p.getDoc(typeSpec.Doc)
	}

	item.Params = p.getTypeParamsWithConstraints(typeSpec, pkg.TypesInfo)

	// if the type is a struct, get its fields
	if structType, ok := typeSpec.Type.(*ast.StructType); ok {
		item.Fields = p.getStructFields(structType.Fields, pkg.TypesInfo)
	}

	// assign type methods
	if _, ok := packageMethods[typeSpec.Name.Name]; ok {
		item.Methods = packageMethods[typeSpec.Name.Name]
	}

	// get interface methods
	item.Methods = append(item.Methods, p.getInterfaceMethods(typeSpec, pkg, errs)...)

	return &item
}

// Helper method to get interfaces of a named types.
func (p *DefaultParser) getInterfaces(namedType *types.Named) []*Interface {
	interfaces := make([]*Interface, 0)

	if _, ok := namedType.Underlying().(*types.Interface); !ok {
		typeMethods := p.collectMethodSet(namedType)
		pointerMethods := p.collectMethodSet(types.NewPointer(namedType))

		for _, iface := range p.interfaces {
			counter := 0

			for ifaceMethod := range iface.iface.Methods() {
				// check if the method exists in either the type's method set or pointer method set
				for _, method := range append(typeMethods, pointerMethods...) {
					if method.Name() == ifaceMethod.Name() &&
						types.Identical(method.Type(), ifaceMethod.Type()) {
						counter++
					}
				}
			}

			// if counter is 0 it means no methods matched so the interface should not be considered at all
			if counter > 0 && counter == iface.iface.NumMethods() {
				interfaces = append(interfaces, &Interface{
					Name:        iface.Name,
					PackageName: iface.PackageName,
					PackagePath: iface.PackagePath,
					iface:       nil,
				})
			}
		}
	}

	// added to ensure deterministic output
	sort.Slice(interfaces, func(i, j int) bool {
		return interfaces[i].Name < interfaces[j].Name
	})

	return interfaces
}

// Helper method to collect a method set of a specific type.
func (p *DefaultParser) collectMethodSet(t types.Type) []*types.Func {
	methods := make([]*types.Func, 0)

	methodSet := types.NewMethodSet(t)

	for method := range methodSet.Methods() {
		if mthd, ok := method.Obj().(*types.Func); ok {
			methods = append(methods, mthd)
		}
	}

	return methods
}

// Helper method to get a human-readable description of the underlying type (instead of ast.* types).
func (p *DefaultParser) describeType( //nolint:funlen,gocyclo,cyclop // not an issue at this stage...
	expr ast.Expr,
	typesInfo *types.Info,
) *TypeDesc {
	tpe := TypeDesc{
		Label: "unknown",
		Links: make([]*Link, 0),
	}

	if expr == nil {
		return &tpe
	}

	switch tp := expr.(type) {
	// An *ast.Ident represents an identifier, which is a name used for variables, types, functions, constants,
	// and package imports in Go source code. Essentially, *ast.Ident is a fundamental node in the abstract syntax
	// tree (AST) that represents any named entity in a Go program.
	case *ast.Ident:
		tpe.Label = tp.Name

		pkgName, pkgPath, pkgAlias, isPackage := p.describePkg(tp, typesInfo)

		// for package description, keep the original package name, alias should be kept as an additional value
		if isPackage {
			tpe.Label = pkgName
		}

		tpe.Links = append(tpe.Links, p.prepareLink(
			tp.Name, p.getUnderlying(expr, typesInfo), pkgName, pkgPath, pkgAlias,
		))
	// An ast.SelectorExpr represents an expression of the form X.Y, where X is an expression (often an identifier
	// representing a package or type) and Y is a field, method, or type that is being selected from X.
	case *ast.SelectorExpr:
		lbl := p.describeType(tp.X, typesInfo).Label + "." + tp.Sel.Name

		tpe.Label = lbl

		ident, ok := tp.X.(*ast.Ident)

		if !ok {
			break
		}

		pkgName, pkgPath, pkgAlias, _ := p.describePkg(ident, typesInfo)

		tpe.Links = append(tpe.Links, p.prepareLink(
			tp.Sel.Name, p.getUnderlying(expr, typesInfo), pkgName, pkgPath, pkgAlias,
		))
	case *ast.StarExpr:
		ti := p.describeType(tp.X, typesInfo)

		tpe.Label = "*" + ti.Label
		tpe.Links = append(tpe.Links, ti.Links...)
	case *ast.ArrayType:
		ti := p.describeType(tp.Elt, typesInfo)

		tpe.Label = "[]" + ti.Label
		tpe.Links = append(tpe.Links, ti.Links...)
	case *ast.MapType:
		kti := p.describeType(tp.Key, typesInfo)
		vti := p.describeType(tp.Value, typesInfo)

		tpe.Label = "map[" + kti.Label + "]" + vti.Label
		tpe.Links = append(tpe.Links, kti.Links...)
		tpe.Links = append(tpe.Links, vti.Links...)
	case *ast.FuncType:
		sig := p.getFuncSignature(tp, typesInfo)

		tpe.Label = "func" + sig.Label
		tpe.Links = append(tpe.Links, sig.Links...)
	case *ast.StructType:
		tpe.Label = "struct"
		tpe.Links = append(
			tpe.Links,
			p.prepareLink("struct", "struct", "", "", ""),
		)
	case *ast.InterfaceType:
		tpe.Label = "interface"
		tpe.Links = append(
			tpe.Links,
			p.prepareLink("interface", "interface", "", "", ""),
		)
	case *ast.Ellipsis:
		ti := p.describeType(tp.Elt, typesInfo)

		tpe.Label = "..." + ti.Label
		tpe.Links = append(tpe.Links, ti.Links...)
	case *ast.IndexExpr:
		xti := p.describeType(tp.X, typesInfo)
		iti := p.describeType(tp.Index, typesInfo)

		tpe.Label = xti.Label + "[" + iti.Label + "]"
		tpe.Links = append(tpe.Links, xti.Links...)
		tpe.Links = append(tpe.Links, iti.Links...)
	case *ast.IndexListExpr:
		var typeArgs []string

		for _, arg := range tp.Indices {
			td := p.describeType(arg, typesInfo)

			typeArgs = append(typeArgs, td.Label)

			tpe.Links = append(tpe.Links, td.Links...)
		}

		dtx := p.describeType(tp.X, typesInfo)

		tpe.Links = append(tpe.Links, dtx.Links...)
		tpe.Label = dtx.Label + "[" + strings.Join(typeArgs, ", ") + "]"
	default:
		tpe.Label = fmt.Sprintf("%T", expr)
	}

	tpe.Links = p.dedupeSliceOfLinks(tpe.Links)

	return &tpe
}

func (p *DefaultParser) enrichType(item *Type) {
	if !p.isSelector(item.Type.Label) && !p.isValidIdentifier(item.Type.Label) {
		return
	}

	selector := p.findSelector(item.Type.Label, item.Type.Links)

	if _, ok := p.types[selector]; ok {
		underlying := p.types[selector]

		// Type Alias:
		//
		// - No explicit conversion is needed; they are treated as the same type.
		// - All methods of the original type are automatically available.

		// Type Definition:
		//
		// - Creates a new type.
		// - Requires explicit conversion to assign or use values interchangeably with the base type.
		// - Methods must be explicitly defined for the new type.

		item.Fields = append(item.Fields, underlying.Fields...)
		item.Embedded = append(item.Embedded, underlying.Embedded...)

		if !item.IsTypeAlias {
			return
		}

		for _, method := range underlying.Methods {
			if method.IsExported {
				item.Methods = append(item.Methods, method)
			}
		}

		for _, iface := range underlying.Interfaces {
			ifaceSel := iface.PackagePath + "." + iface.Name
			if _, ok := p.types[ifaceSel]; ok && p.types[ifaceSel].IsExported {
				item.Interfaces = append(item.Interfaces, iface)
			}
		}
	}
}

func (p *DefaultParser) findSelector(label string, links []*Link) string {
	split := strings.Split(label, ".")

	for _, link := range links {
		if link.Name == split[len(split)-1] {
			return link.PackagePath + "." + link.Name
		}
	}

	return ""
}

// Helper function to get underlying type.
func (p *DefaultParser) getUnderlying(expr ast.Expr, typesInfo *types.Info) string {
	if tav, ok := typesInfo.Types[expr]; ok {
		return tav.Type.Underlying().String()
	}

	return ""
}

// Helper function to prepare link.
func (p *DefaultParser) prepareLink(name, underlying, pkgName, pkgPath, pkgAlias string) *Link {
	return &Link{
		Name:         name,
		Underlying:   underlying,
		PackageName:  pkgName,
		PackagePath:  pkgPath,
		PackageAlias: pkgAlias,
	}
}

// Helper method to describe a current object package.
// Returns values in the following order: PkgName, PkgPath, PkgAlias, IsPackage.
func (p *DefaultParser) describePkg(
	ident *ast.Ident,
	typesInfo *types.Info,
) (string, string, string, bool) {
	obj := typesInfo.ObjectOf(ident)

	if obj != nil && obj.Pkg() != nil {
		originalPkgName := obj.Pkg().Name()
		pkgPath := obj.Pkg().Path()
		alias := ""
		isPackage := false

		if pkgName, ok := typesInfo.Uses[ident].(*types.PkgName); ok {
			originalPkgName = pkgName.Imported().Name()
			pkgPath = pkgName.Imported().Path()

			if pkgName.Imported().Name() != pkgName.Name() {
				alias = pkgName.Name() // the alias used in the import
			}

			isPackage = true
		}

		return originalPkgName, pkgPath, alias, isPackage
	}

	return "", "", "", false
}

// Helper method to collect all package methods with a receiver.
func (p *DefaultParser) collectAllPackageMethods(
	pkg *packages.Package,
	errs chan error,
) map[string][]*Method {
	methods := make(map[string][]*Method)

	for _, syn := range pkg.Syntax {
		for _, decl := range syn.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)

			if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
				continue
			}

			typeName, method := p.describeMethodDecl(pkg, funcDecl, pkg.TypesInfo, errs)

			methods[typeName] = append(methods[typeName], method)
		}
	}

	return methods
}

// Helper method to describe method declaration.
func (p *DefaultParser) describeMethodDecl(
	pkg *packages.Package, decl *ast.FuncDecl, typesInfo *types.Info, errs chan error,
) (string, *Method) {
	var (
		typeName          string
		isPointerReceiver bool
	)

	receiver := decl.Recv.List[0].Type

	switch recv := receiver.(type) {
	case *ast.StarExpr:
		// pointer receiver (*Type)
		if ident, ok := recv.X.(*ast.Ident); ok {
			typeName = ident.Name
			isPointerReceiver = true
		}
	case *ast.Ident:
		// value receiver (Type)
		typeName = recv.Name
	}

	if typeName == "" {
		return "", nil
	}

	start := pkg.Fset.Position(decl.Pos())
	end := pkg.Fset.Position(decl.End())

	rootPath, err := p.normalizeRootPath(start.Filename)
	if err != nil {
		errs <- err

		return "", nil
	}

	relativePath, err := filepath.Rel(rootPath, start.Filename)
	if err != nil {
		errs <- fmt.Errorf(
			"%w: %w: with root path `%s` and file name `%s`",
			ErrGetRelPath, err, rootPath, start.Filename,
		)

		return "", nil
	}

	method := Method{
		FilePath: relativePath,
		Start: &Position{
			Offset: start.Offset,
			Line:   start.Line,
			Column: start.Column,
		},
		End: &Position{
			Offset: end.Offset,
			Line:   end.Line,
			Column: end.Column,
		},
		Name:              decl.Name.Name,
		Signature:         p.getFuncSignature(decl.Type, typesInfo),
		IsPointerReceiver: isPointerReceiver,
		IsExported:        ast.IsExported(decl.Name.Name),
		Doc:               "",
	}

	if decl.Doc != nil {
		method.Doc = p.getDoc(decl.Doc)
	}

	return typeName, &method
}

// Helper method to describe function signatures.
func (p *DefaultParser) getFuncSignature(funcType *ast.FuncType, typesInfo *types.Info) *Signature {
	signature := Signature{
		Label:     "",
		Alt:       "",
		Arguments: []*Argument{},
		Returns:   []*Return{},
		Links:     make([]*Link, 0),
	}

	arguments := make([]string, 0)
	argumentsAlt := make([]string, 0)

	for _, param := range funcType.Params.List {
		args, prm, alts, lns := p.getFuncSignatureArguments(param, typesInfo)

		signature.Arguments = append(signature.Arguments, args...)
		arguments = append(arguments, prm...)
		argumentsAlt = append(argumentsAlt, alts...)
		signature.Links = append(signature.Links, lns...)
	}

	results := make([]string, 0)

	if funcType.Results != nil {
		for _, result := range funcType.Results.List {
			rtns, rsls, lns := p.getFuncSignatureReturns(result, typesInfo)

			signature.Returns = append(signature.Returns, rtns...)
			results = append(results, rsls...)
			signature.Links = append(signature.Links, lns...)
		}
	}

	signature.Links = p.dedupeSliceOfLinks(signature.Links)

	if len(results) > 0 {
		signature.Label = fmt.Sprintf(
			"(%s) (%s)",
			strings.Join(arguments, ", "), strings.Join(results, ", "),
		)

		signature.Alt = fmt.Sprintf(
			"(%s) (%s)",
			strings.Join(argumentsAlt, ", "), strings.Join(results, ", "),
		)

		return &signature
	}

	signature.Label = fmt.Sprintf("(%s)", strings.Join(arguments, ", "))
	signature.Alt = fmt.Sprintf("(%s)", strings.Join(argumentsAlt, ", "))

	return &signature
}

// Helper method to describe function signature arguments.
func (p *DefaultParser) getFuncSignatureArguments(
	field *ast.Field, typesInfo *types.Info,
) ([]*Argument, []string, []string, []*Link) {
	td := p.describeType(field.Type, typesInfo)

	paramType := td.Label

	links := make([]*Link, 0, len(td.Links))
	params := make([]string, 0, len(field.Names))
	alts := make([]string, 0, len(field.Names))
	arguments := make([]*Argument, 0, len(field.Names))

	links = append(links, td.Links...)

	for _, fieldName := range field.Names {
		argument := Argument{
			Name: fieldName.Name,
			Type: td,
		}

		params = append(params, fmt.Sprintf("%s %s", fieldName.Name, paramType))
		alts = append(alts, paramType)

		arguments = append(arguments, &argument)
	}

	if len(field.Names) > 0 {
		return arguments, params, alts, links
	}

	params = append(params, paramType)

	argument := Argument{
		Name: "",
		Type: td,
	}

	arguments = append(arguments, &argument)

	return arguments, params, alts, links
}

// Helper method to describe function signature returns.
func (p *DefaultParser) getFuncSignatureReturns(
	field *ast.Field,
	typesInfo *types.Info,
) ([]*Return, []string, []*Link) {
	td := p.describeType(field.Type, typesInfo)

	links := make([]*Link, 0, len(td.Links))
	results := make([]string, 0, len(field.Names))
	returns := make([]*Return, 0, len(field.Names))

	if len(field.Names) > 0 {
		for _, resultName := range field.Names {
			retrn := Return{
				Name: resultName.Name,
				Type: td,
			}

			results = append(results, fmt.Sprintf("%s %s", resultName.Name, td.Label))
			returns = append(returns, &retrn)
		}
	} else {
		retrn := Return{
			Name: "",
			Type: td,
		}

		results = append(results, td.Label)
		returns = append(returns, &retrn)
	}

	links = append(links, td.Links...)

	return returns, results, links
}

// Helper method to get type parameters for generic types, including their constraints.
func (p *DefaultParser) getVarParamsWithConstraints(
	valueSpec *ast.ValueSpec,
	typesInfo *types.Info,
) *Params {
	params := Params{
		Label: "",
		List:  []*Param{},
	}

	pl := make([]*Param, 0)
	sl := make([]string, 0)

	// check if the type of the variable is an IndexExpr (for generics)
	switch expr := valueSpec.Type.(type) {
	case *ast.IndexExpr:
		// single generic argument (for older Go versions with single generic type)
		if ident, ok := expr.Index.(*ast.Ident); ok {
			dt := p.describeType(ident, typesInfo)

			sl = append(sl, fmt.Sprintf("%s %s", ident.Name, dt.Label))
			pl = append(pl, &Param{
				Name:      ident.Name,
				Constrain: dt,
			})
		}
	case *ast.IndexListExpr:
		// multiple generic arguments (Go 1.18+)
		for _, arg := range expr.Indices {
			if ident, ok := arg.(*ast.Ident); ok {
				dt := p.describeType(ident, typesInfo)

				sl = append(sl, fmt.Sprintf("%s %s", ident.Name, dt.Label))
				pl = append(pl, &Param{
					Name:      ident.Name,
					Constrain: dt,
				})
			}
		}
	}

	if len(pl) == 0 {
		return nil
	}

	params.Label = fmt.Sprintf("[%s]", strings.Join(sl, ", "))
	params.List = pl

	return &params
}

// Helper method to get type parameters for generic types, including their constraints.
func (p *DefaultParser) getTypeParamsWithConstraints(
	typeSpec *ast.TypeSpec,
	typesInfo *types.Info,
) *Params {
	if typeSpec.TypeParams == nil {
		return nil
	}

	params := Params{
		Label: "",
		List:  []*Param{},
	}

	pl := make([]*Param, 0)
	sl := make([]string, 0)

	for _, param := range typeSpec.TypeParams.List {
		for _, paramName := range param.Names {
			dt := p.describeType(param.Type, typesInfo)

			sl = append(sl, fmt.Sprintf("%s %s", paramName.Name, dt.Label))
			pl = append(pl, &Param{
				Name:      paramName.Name,
				Constrain: dt,
			})
		}
	}

	params.Label = fmt.Sprintf("[%s]", strings.Join(sl, ", "))
	params.List = pl

	return &params
}

// Helper method to get function type params.
func (p *DefaultParser) getFuncTypeParams(funcDecl *ast.FuncDecl, typesInfo *types.Info) *Params {
	if funcDecl.Type.TypeParams == nil {
		return nil
	}

	params := Params{
		Label: "",
		List:  nil,
	}

	pl := make([]*Param, 0)
	sl := make([]string, 0)

	for _, param := range funcDecl.Type.TypeParams.List {
		for _, paramName := range param.Names {
			dt := p.describeType(param.Type, typesInfo)

			sl = append(sl, fmt.Sprintf("%s %s", paramName.Name, dt.Label))
			pl = append(pl, &Param{
				Name:      paramName.Name,
				Constrain: dt,
			})
		}
	}

	params.Label = fmt.Sprintf("[%s]", strings.Join(sl, ", "))
	params.List = pl

	return &params
}

// Helper method to get struct fields and their types.
func (p *DefaultParser) getStructFields(fields *ast.FieldList, typesInfo *types.Info) []*Field {
	list := make([]*Field, 0)

	for _, field := range fields.List {
		fieldType := p.describeType(field.Type, typesInfo)
		for _, fieldName := range field.Names {
			item := Field{
				Name:       fieldName.Name,
				Type:       fieldType,
				IsExported: ast.IsExported(fieldName.Name),
				Tags:       nil,
				Doc:        "",
			}

			if field.Tag != nil {
				item.Tags = p.parseTags(field.Tag)
			}

			if field.Doc != nil {
				item.Doc = p.getDoc(field.Doc)
			}

			list = append(list, &item)
		}
	}

	return list
}

// Helper method to parses a tag string and returns a structured list of tags.
func (p *DefaultParser) parseTags(basicList *ast.BasicLit) *Tags {
	tags := Tags{
		Label: regexp.MustCompile(`\s+`).ReplaceAllString(basicList.Value, " "),
		Tags:  make([]*Tag, 0),
	}

	// split the tag string by spaces to handle multiple tags
	keyValuePairs := strings.Split(strings.Trim(tags.Label, "`"), " ")

	const tagPartsSize = 2

	// process each key-value pair
	for _, keyValue := range keyValuePairs {
		// split by colon to separate the tag key and value
		parts := strings.SplitN(keyValue, ":", tagPartsSize)

		if len(parts) != tagPartsSize {
			continue // skip malformed tags e.g. in case multiple space separation is used
		}

		key := parts[0]
		value := strings.Trim(parts[1], `"`)

		// split value by comma in case of multiple options (e.g., `json:"name,omitempty"`)
		valueParts := strings.SplitSeq(value, ",")

		for part := range valueParts {
			tags.Tags = append(tags.Tags, &Tag{
				Name:  key,
				Value: part,
			})
		}
	}

	return &tags
}

// Helper method to get all embedded types.
func (p *DefaultParser) getEmbeddedTypes(
	typeSpec *ast.TypeSpec,
	typesInfo *types.Info,
) []*TypeDesc {
	desc := make([]*TypeDesc, 0)

	switch tp := typeSpec.Type.(type) {
	case *ast.InterfaceType:
		for _, method := range tp.Methods.List {
			if len(method.Names) == 0 {
				desc = append(desc, p.describeType(method.Type, typesInfo))
			}
		}
	case *ast.StructType:
		for _, field := range tp.Fields.List {
			if len(field.Names) == 0 {
				desc = append(desc, p.describeType(field.Type, typesInfo))
			}
		}
	}

	return desc
}

// Helper method to collect all methods for a given interface type.
func (p *DefaultParser) getInterfaceMethods(
	typeSpec *ast.TypeSpec, pkg *packages.Package, errs chan error,
) []*Method {
	var methods []*Method

	if ifaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		for _, field := range ifaceType.Methods.List {
			// each field could be a method or embedded interface
			for _, name := range field.Names {
				funcType, ok := field.Type.(*ast.FuncType)
				if !ok {
					// not a valid function type, skip this field
					return nil
				}

				start := pkg.Fset.Position(funcType.Pos())
				end := pkg.Fset.Position(funcType.End())

				rootPath, err := p.normalizeRootPath(start.Filename)
				if err != nil {
					errs <- err

					return nil
				}

				relativePath, err := filepath.Rel(rootPath, start.Filename)
				if err != nil {
					errs <- fmt.Errorf(
						"%w: %w: with root path `%s` and file name `%s`",
						ErrGetRelPath, err, rootPath, start.Filename,
					)

					return nil
				}

				method := Method{
					FilePath: relativePath,
					Start: &Position{
						Offset: start.Offset,
						Line:   start.Line,
						Column: start.Column,
					},
					End: &Position{
						Offset: end.Offset,
						Line:   end.Line,
						Column: end.Column,
					},
					Name:              name.Name,
					Signature:         p.getFuncSignature(funcType, pkg.TypesInfo),
					IsPointerReceiver: false,
					IsExported:        ast.IsExported(name.Name),
					Doc:               "",
				}

				if field.Doc != nil {
					method.Doc = p.getDoc(field.Doc)
				}

				methods = append(methods, &method)
			}
		}
	}

	return methods
}

func (p *DefaultParser) isSelector(input string) bool {
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
	return p.isValidIdentifier(parts[0]) && p.isValidIdentifier(parts[1])
}

func (p *DefaultParser) isValidIdentifier(input string) bool {
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

// Helper method to deduplicate slice of link.
func (p *DefaultParser) dedupeSliceOfLinks(input []*Link) []*Link {
	unique := make(map[string]struct{})
	result := make([]*Link, 0)

	for _, val := range input {
		if _, exists := unique[val.Name]; !exists {
			unique[val.Name] = struct{}{}

			result = append(result, val)
		}
	}

	return result
}

// Helper method to deduplicate slice of fields.
func (p *DefaultParser) dedupeSliceOfFields(input []*Field) []*Field {
	unique := make(map[string]struct{})
	result := make([]*Field, 0)

	for _, val := range input {
		if _, exists := unique[val.Name]; !exists {
			unique[val.Name] = struct{}{}

			result = append(result, val)
		}
	}

	return result
}

// Helper method to deduplicate slice of interfaces.
func (p *DefaultParser) dedupeSliceOfInterfaces(input []*Interface) []*Interface {
	unique := make(map[string]struct{})
	result := make([]*Interface, 0)

	for _, val := range input {
		if _, exists := unique[val.Name]; !exists {
			unique[val.Name] = struct{}{}

			result = append(result, val)
		}
	}

	return result
}

// Helper method to deduplicate slice of methods.
func (p *DefaultParser) dedupeSliceOfMethods(input []*Method) []*Method {
	unique := make(map[string]struct{})
	result := make([]*Method, 0)

	for _, val := range input {
		if _, exists := unique[val.Name]; !exists {
			unique[val.Name] = struct{}{}

			result = append(result, val)
		}
	}

	return result
}

// Helper method to deduplicate slice of type desc.
func (p *DefaultParser) dedupeSliceOfTypeDesc(input []*TypeDesc) []*TypeDesc {
	unique := make(map[string]struct{})
	result := make([]*TypeDesc, 0)

	for _, val := range input {
		if _, exists := unique[val.Label]; !exists {
			unique[val.Label] = struct{}{}

			result = append(result, val)
		}
	}

	return result
}

// Helper method to get code doc.
func (p *DefaultParser) getDoc(doc *ast.CommentGroup) string {
	lines := make([]string, len(doc.List))

	for i, comment := range doc.List {
		lines[i] = comment.Text
	}

	if p.docParser != nil {
		return p.docParser.Parse(lines)
	}

	return strings.Join(lines, "\n")
}

func (p *DefaultParser) normalizeRootPath(filename string) (string, error) {
	if filepath.IsAbs(filename) && !filepath.IsAbs(p.rootPath) {
		path, err := filepath.Abs(p.rootPath)
		if err != nil {
			return "", fmt.Errorf("%w: `%s`", ErrRootPathNormalise, path)
		}

		return path, nil
	}

	return p.rootPath, nil
}
