package gen

// UMLGraphGeneratorOption option function.
type UMLGraphGeneratorOption func(p *DefaultUMLGraphGenerator)

// UMLGraphGeneratorConfig configuration.
type UMLGraphGeneratorConfig struct {
	Legend, Primitive, Var, Const, Func, NotExported, Doc bool
}

// TypeInfoParts to represent parsed type info.
type TypeInfoParts struct {
	Name, PackagePath string
}
