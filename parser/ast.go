package parser

import (
	"bytes"
	"strings"
)

const (
	nullRep = "<null>"
)

// Node represents a node in the AST.
type Node interface {
	// Pos returns the position of first character belonging to the node.
	Pos() Pos
	// End returns the position of first character immediately after the node.
	End() Pos
	// String returns a string representation of the node.
	String() string
}

// IdentList represents a list of identifiers.
type IdentList struct {
	LParen  Pos
	VarArgs bool
	List    []*Ident
	RParen  Pos
}

// Pos returns the position of first character belonging to the node.
func (n *IdentList) Pos() Pos {
	if n.LParen.IsValid() {
		return n.LParen
	}
	if len(n.List) > 0 {
		return n.List[0].Pos()
	}
	return NoPos
}

// End returns the position of first character immediately after the node.
func (n *IdentList) End() Pos {
	if n.RParen.IsValid() {
		return n.RParen + 1
	}
	if l := len(n.List); l > 0 {
		return n.List[l-1].End()
	}
	return NoPos
}

// NumFields returns the number of fields.
func (n *IdentList) NumFields() int {
	if n == nil {
		return 0
	}
	return len(n.List)
}

func (n *IdentList) String() string {
	var list []string
	for i, e := range n.List {
		if n.VarArgs && i == len(n.List)-1 {
			list = append(list, "..."+e.String())
		} else {
			list = append(list, e.String())
		}
	}
	return "(" + strings.Join(list, ", ") + ")"
}

// ValuedIdentList represents a list of identifier with value pairs.
type ValuedIdentList struct {
	LParen  Pos
	VarArgs bool
	Names   []*Ident
	Values  []Expr
	RParen  Pos
}

// Pos returns the position of first character belonging to the node.
func (n *ValuedIdentList) Pos() Pos {
	if n.LParen.IsValid() {
		return n.LParen
	}
	if len(n.Names) > 0 {
		return n.Names[0].Pos()
	}
	return NoPos
}

// End returns the position of first character immediately after the node.
func (n *ValuedIdentList) End() Pos {
	if n.RParen.IsValid() {
		return n.RParen + 1
	}
	if l := len(n.Names); l > 0 {
		if n.VarArgs {
			return n.Names[l-1].End()
		}
		return n.Values[l-1].End()
	}
	return NoPos
}

// NumFields returns the number of fields.
func (n *ValuedIdentList) NumFields() int {
	if n == nil {
		return 0
	}
	return len(n.Names)
}

func (n *ValuedIdentList) String() string {
	var list []string
	for i, e := range n.Names {
		if n.VarArgs && i == len(n.Names)-1 {
			list = append(list, "..."+e.String())
		} else {
			list = append(list, e.String()+" = "+n.Values[i].String())
		}
	}
	return "(" + strings.Join(list, ", ") + ")"
}

// FuncParams represents a function paramsw.
type FuncParams struct {
	LParen Pos
	Args   *IdentList
	Kwargs *ValuedIdentList
	RParen Pos
}

// Pos returns the position of first character belonging to the node.
func (n *FuncParams) Pos() Pos {
	if n.LParen.IsValid() {
		return n.LParen
	}
	if n.Args != nil && len(n.Args.List) > 0 {
		return n.Args.List[0].Pos()
	}
	if n.Kwargs != nil && len(n.Kwargs.Names) > 0 {
		return n.Kwargs.Names[0].Pos()
	}
	return NoPos
}

// End returns the position of first character immediately after the node.
func (n *FuncParams) End() Pos {
	if n.RParen.IsValid() {
		return n.RParen + 1
	}
	if n.Kwargs != nil && len(n.Kwargs.Names) > 0 {
		return n.Kwargs.End()
	}
	if n.Args != nil && len(n.Args.List) > 0 {
		return n.Args.End()
	}
	return NoPos
}

func (n *FuncParams) String() string {
	buf := bytes.NewBufferString("(")
	if n.Args != nil && len(n.Args.List) > 0 {
		v := n.Args.String()
		buf.WriteString(v[1 : len(v)-1])
	}
	if n.Kwargs != nil && len(n.Kwargs.Names) > 0 {
		buf.WriteString("; ")
		v := n.Kwargs.String()
		buf.WriteString(v[1 : len(v)-1])
	}
	buf.WriteString(")")
	return buf.String()
}
