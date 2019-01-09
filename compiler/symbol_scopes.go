package compiler

type SymbolScope string

const (
	ScopeGlobal  SymbolScope = "GLOBAL"
	ScopeLocal               = "LOCAL"
	ScopeBuiltin             = "BUILTIN"
	ScopeFree                = "FREE"
)
