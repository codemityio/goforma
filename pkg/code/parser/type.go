package parser

import "go/types"

// Option function.
type Option func(p *DefaultParser)

type Position struct {
	Offset int
	Line   int
	Column int
}

type Var struct {
	FilePath    string
	Start       *Position
	End         *Position
	PackageName string
	PackagePath string
	Name        string
	Type        *TypeDesc
	TypeInfo    string
	Params      *Params
	IsExported  bool
	Doc         string
}

type Const struct {
	FilePath    string
	Start       *Position
	End         *Position
	PackageName string
	PackagePath string
	Name        string
	Type        *TypeDesc
	TypeInfo    string
	IsExported  bool
	Doc         string
}

type Func struct {
	FilePath    string
	Start       *Position
	End         *Position
	PackageName string
	PackagePath string
	Name        string
	Signature   *Signature
	Params      *Params
	IsExported  bool
	Doc         string
}

type Argument struct {
	Name string
	Type *TypeDesc
}

type Return struct {
	Name string
	Type *TypeDesc
}

type Signature struct {
	Label     string
	Alt       string
	Arguments []*Argument
	Returns   []*Return
	Links     []*Link
}

type Type struct {
	FilePath    string
	Start       *Position
	End         *Position
	PackageName string
	PackagePath string
	Name        string
	Type        *TypeDesc
	TypeInfo    string
	IsTypeAlias bool
	Interfaces  []*Interface
	Embedded    []*TypeDesc
	Params      *Params
	Fields      []*Field
	Methods     []*Method
	IsExported  bool
	Doc         string
}

type Params struct {
	Label string
	List  []*Param
}

type Param struct {
	Name      string
	Constrain *TypeDesc
}

type Field struct {
	Name       string
	Type       *TypeDesc
	IsExported bool
	Tags       *Tags
	Doc        string
}

type Tags struct {
	Label string
	Tags  []*Tag
}

type Tag struct {
	Name  string
	Value string
}

type Method struct {
	FilePath          string
	Start             *Position
	End               *Position
	Name              string
	Signature         *Signature
	IsPointerReceiver bool
	IsExported        bool
	Doc               string
}

type Interface struct {
	Name        string
	PackageName string
	PackagePath string
	iface       *types.Interface
}

type TypeDesc struct {
	Label string
	Links []*Link
}

type Link struct {
	Name         string
	Underlying   string
	PackageName  string
	PackagePath  string
	PackageAlias string
}

type CodeMap[V *Var, T *Type, F *Func, C *Const] struct {
	Var   []V
	Type  []T
	Func  []F
	Const []C
}
