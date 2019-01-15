package compiler

// Symbol represents a symbol in the symbol table.
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}
