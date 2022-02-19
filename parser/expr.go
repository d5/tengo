package parser

import (
	"bytes"
	"strings"

	"github.com/d5/tengo/v2/token"
)

// Expr represents an expression node in the AST.
type Expr interface {
	Node
	exprNode()
}

// ArrayLit represents an array literal.
type ArrayLit struct {
	Elements []Expr
	LBrack   Pos
	RBrack   Pos
}

func (e *ArrayLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *ArrayLit) Pos() Pos {
	return e.LBrack
}

// End returns the position of first character immediately after the node.
func (e *ArrayLit) End() Pos {
	return e.RBrack + 1
}

func (e *ArrayLit) String() string {
	var elements []string
	for _, m := range e.Elements {
		elements = append(elements, m.String())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

// BadExpr represents a bad expression.
type BadExpr struct {
	From Pos
	To   Pos
}

func (e *BadExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *BadExpr) Pos() Pos {
	return e.From
}

// End returns the position of first character immediately after the node.
func (e *BadExpr) End() Pos {
	return e.To
}

func (e *BadExpr) String() string {
	return "<bad expression>"
}

// BinaryExpr represents a binary operator expression.
type BinaryExpr struct {
	LHS      Expr
	RHS      Expr
	Token    token.Token
	TokenPos Pos
}

func (e *BinaryExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *BinaryExpr) Pos() Pos {
	return e.LHS.Pos()
}

// End returns the position of first character immediately after the node.
func (e *BinaryExpr) End() Pos {
	return e.RHS.End()
}

func (e *BinaryExpr) String() string {
	return "(" + e.LHS.String() + " " + e.Token.String() +
		" " + e.RHS.String() + ")"
}

// BoolLit represents a boolean literal.
type BoolLit struct {
	Value    bool
	ValuePos Pos
	Literal  string
}

func (e *BoolLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *BoolLit) Pos() Pos {
	return e.ValuePos
}

// End returns the position of first character immediately after the node.
func (e *BoolLit) End() Pos {
	return Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *BoolLit) String() string {
	return e.Literal
}

// CallExprArgs represents a call expression arguments.
type CallExprArgs struct {
	Values   []Expr
	Ellipsis Pos
}

// CallExprKwargs represents a call expression keyword arguments.
type CallExprKwargs struct {
	Names    []*Ident
	Values   []Expr
	Ellipsis Pos
}

// CallExpr represents a function call expression.
type CallExpr struct {
	Func   Expr
	LParen Pos
	Args   CallExprArgs
	Kwargs CallExprKwargs
	RParen Pos
}

func (e *CallExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *CallExpr) Pos() Pos {
	return e.Func.Pos()
}

// End returns the position of first character immediately after the node.
func (e *CallExpr) End() Pos {
	return e.RParen + 1
}

func (e *CallExpr) String() string {
	var buf = bytes.NewBufferString(e.Func.String())
	buf.WriteString("(")
	if l := len(e.Args.Values); l > 0 {
		for i, e := range e.Args.Values {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(e.String())
		}
		if e.Args.Ellipsis.IsValid() {
			buf.WriteString("...")
		}
	}
	if l := len(e.Kwargs.Values); l > 0 {
		if len(e.Args.Values) == 0 {
			buf.WriteString(";")
		} else {
			buf.WriteString(", ")
		}

		qnt := l
		if e.Kwargs.Ellipsis.IsValid() {
			qnt--
		}
		for i, n := range e.Kwargs.Names[0:qnt] {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(n.String())
			buf.WriteString(" = ")
			buf.WriteString(e.Kwargs.Values[i].String())
		}
		if e.Kwargs.Ellipsis.IsValid() {
			if qnt > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(e.Kwargs.Values[qnt].String())
			buf.WriteString("...")
		}
	}
	buf.WriteString(")")
	return buf.String()
}

// CharLit represents a character literal.
type CharLit struct {
	Value    rune
	ValuePos Pos
	Literal  string
}

func (e *CharLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *CharLit) Pos() Pos {
	return e.ValuePos
}

// End returns the position of first character immediately after the node.
func (e *CharLit) End() Pos {
	return Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *CharLit) String() string {
	return e.Literal
}

// CondExpr represents a ternary conditional expression.
type CondExpr struct {
	Cond        Expr
	True        Expr
	False       Expr
	QuestionPos Pos
	ColonPos    Pos
}

func (e *CondExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *CondExpr) Pos() Pos {
	return e.Cond.Pos()
}

// End returns the position of first character immediately after the node.
func (e *CondExpr) End() Pos {
	return e.False.End()
}

func (e *CondExpr) String() string {
	return "(" + e.Cond.String() + " ? " + e.True.String() +
		" : " + e.False.String() + ")"
}

// ErrorExpr represents an error expression
type ErrorExpr struct {
	Expr     Expr
	ErrorPos Pos
	LParen   Pos
	RParen   Pos
}

func (e *ErrorExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *ErrorExpr) Pos() Pos {
	return e.ErrorPos
}

// End returns the position of first character immediately after the node.
func (e *ErrorExpr) End() Pos {
	return e.RParen
}

func (e *ErrorExpr) String() string {
	return "error(" + e.Expr.String() + ")"
}

// FloatLit represents a floating point literal.
type FloatLit struct {
	Value    float64
	ValuePos Pos
	Literal  string
}

func (e *FloatLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *FloatLit) Pos() Pos {
	return e.ValuePos
}

// End returns the position of first character immediately after the node.
func (e *FloatLit) End() Pos {
	return Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *FloatLit) String() string {
	return e.Literal
}

// FuncLit represents a function literal.
type FuncLit struct {
	Type *FuncType
	Body *BlockStmt
}

func (e *FuncLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *FuncLit) Pos() Pos {
	return e.Type.Pos()
}

// End returns the position of first character immediately after the node.
func (e *FuncLit) End() Pos {
	return e.Body.End()
}

func (e *FuncLit) String() string {
	return "func" + e.Type.Params.String() + " " + e.Body.String()
}

// FuncType represents a function type definition.
type FuncType struct {
	FuncPos Pos
	Params  *FuncParams
}

func (e *FuncType) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *FuncType) Pos() Pos {
	return e.FuncPos
}

// End returns the position of first character immediately after the node.
func (e *FuncType) End() Pos {
	return e.Params.End()
}

func (e *FuncType) String() string {
	return "func" + e.Params.String()
}

// Ident represents an identifier.
type Ident struct {
	Name    string
	NamePos Pos
}

func (e *Ident) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *Ident) Pos() Pos {
	return e.NamePos
}

// End returns the position of first character immediately after the node.
func (e *Ident) End() Pos {
	return Pos(int(e.NamePos) + len(e.Name))
}

func (e *Ident) String() string {
	if e != nil {
		return e.Name
	}
	return nullRep
}

// ImmutableExpr represents an immutable expression
type ImmutableExpr struct {
	Expr     Expr
	ErrorPos Pos
	LParen   Pos
	RParen   Pos
}

func (e *ImmutableExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *ImmutableExpr) Pos() Pos {
	return e.ErrorPos
}

// End returns the position of first character immediately after the node.
func (e *ImmutableExpr) End() Pos {
	return e.RParen
}

func (e *ImmutableExpr) String() string {
	return "immutable(" + e.Expr.String() + ")"
}

// ImportExpr represents an import expression
type ImportExpr struct {
	ModuleName string
	Token      token.Token
	TokenPos   Pos
}

func (e *ImportExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *ImportExpr) Pos() Pos {
	return e.TokenPos
}

// End returns the position of first character immediately after the node.
func (e *ImportExpr) End() Pos {
	// import("moduleName")
	return Pos(int(e.TokenPos) + 10 + len(e.ModuleName))
}

func (e *ImportExpr) String() string {
	return `import("` + e.ModuleName + `")"`
}

// IndexExpr represents an index expression.
type IndexExpr struct {
	Expr   Expr
	LBrack Pos
	Index  Expr
	RBrack Pos
}

func (e *IndexExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *IndexExpr) Pos() Pos {
	return e.Expr.Pos()
}

// End returns the position of first character immediately after the node.
func (e *IndexExpr) End() Pos {
	return e.RBrack + 1
}

func (e *IndexExpr) String() string {
	var index string
	if e.Index != nil {
		index = e.Index.String()
	}
	return e.Expr.String() + "[" + index + "]"
}

// IntLit represents an integer literal.
type IntLit struct {
	Value    int64
	ValuePos Pos
	Literal  string
}

func (e *IntLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *IntLit) Pos() Pos {
	return e.ValuePos
}

// End returns the position of first character immediately after the node.
func (e *IntLit) End() Pos {
	return Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *IntLit) String() string {
	return e.Literal
}

// MapElementLit represents a map element.
type MapElementLit struct {
	Key      string
	KeyPos   Pos
	ColonPos Pos
	Value    Expr
}

func (e *MapElementLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *MapElementLit) Pos() Pos {
	return e.KeyPos
}

// End returns the position of first character immediately after the node.
func (e *MapElementLit) End() Pos {
	return e.Value.End()
}

func (e *MapElementLit) String() string {
	return e.Key + ": " + e.Value.String()
}

// MapLit represents a map literal.
type MapLit struct {
	LBrace   Pos
	Elements []*MapElementLit
	RBrace   Pos
}

func (e *MapLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *MapLit) Pos() Pos {
	return e.LBrace
}

// End returns the position of first character immediately after the node.
func (e *MapLit) End() Pos {
	return e.RBrace + 1
}

func (e *MapLit) String() string {
	var elements []string
	for _, m := range e.Elements {
		elements = append(elements, m.String())
	}
	return "{" + strings.Join(elements, ", ") + "}"
}

// ParenExpr represents a parenthesis wrapped expression.
type ParenExpr struct {
	Expr   Expr
	LParen Pos
	RParen Pos
}

func (e *ParenExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *ParenExpr) Pos() Pos {
	return e.LParen
}

// End returns the position of first character immediately after the node.
func (e *ParenExpr) End() Pos {
	return e.RParen + 1
}

func (e *ParenExpr) String() string {
	return "(" + e.Expr.String() + ")"
}

// SelectorExpr represents a selector expression.
type SelectorExpr struct {
	Expr Expr
	Sel  Expr
}

func (e *SelectorExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *SelectorExpr) Pos() Pos {
	return e.Expr.Pos()
}

// End returns the position of first character immediately after the node.
func (e *SelectorExpr) End() Pos {
	return e.Sel.End()
}

func (e *SelectorExpr) String() string {
	return e.Expr.String() + "." + e.Sel.String()
}

// SliceExpr represents a slice expression.
type SliceExpr struct {
	Expr   Expr
	LBrack Pos
	Low    Expr
	High   Expr
	RBrack Pos
}

func (e *SliceExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *SliceExpr) Pos() Pos {
	return e.Expr.Pos()
}

// End returns the position of first character immediately after the node.
func (e *SliceExpr) End() Pos {
	return e.RBrack + 1
}

func (e *SliceExpr) String() string {
	var low, high string
	if e.Low != nil {
		low = e.Low.String()
	}
	if e.High != nil {
		high = e.High.String()
	}
	return e.Expr.String() + "[" + low + ":" + high + "]"
}

// StringLit represents a string literal.
type StringLit struct {
	Value    string
	ValuePos Pos
	Literal  string
}

func (e *StringLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *StringLit) Pos() Pos {
	return e.ValuePos
}

// End returns the position of first character immediately after the node.
func (e *StringLit) End() Pos {
	return Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *StringLit) String() string {
	return e.Literal
}

// UnaryExpr represents an unary operator expression.
type UnaryExpr struct {
	Expr     Expr
	Token    token.Token
	TokenPos Pos
}

func (e *UnaryExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *UnaryExpr) Pos() Pos {
	return e.Expr.Pos()
}

// End returns the position of first character immediately after the node.
func (e *UnaryExpr) End() Pos {
	return e.Expr.End()
}

func (e *UnaryExpr) String() string {
	return "(" + e.Token.String() + e.Expr.String() + ")"
}

// UndefinedLit represents an undefined literal.
type UndefinedLit struct {
	TokenPos Pos
}

func (e *UndefinedLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *UndefinedLit) Pos() Pos {
	return e.TokenPos
}

// End returns the position of first character immediately after the node.
func (e *UndefinedLit) End() Pos {
	return e.TokenPos + 9 // len(undefined) == 9
}

func (e *UndefinedLit) String() string {
	return "undefined"
}

// UndefinedKwargLit represents an undefined kwarg literal.
type UndefinedKwargLit struct {
	TokenPos Pos
}

func (e *UndefinedKwargLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *UndefinedKwargLit) Pos() Pos {
	return e.TokenPos
}

// End returns the position of first character immediately after the node.
func (e *UndefinedKwargLit) End() Pos {
	return e.TokenPos + 15 // len(kwarg_undefined) == 15
}

func (e *UndefinedKwargLit) String() string {
	return "undefined_kwarg"
}

// CalleeLit represents an literal callee function/closure flag.
type CalleeLit struct {
	TokenPos Pos
}

func (e *CalleeLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *CalleeLit) Pos() Pos {
	return e.TokenPos
}

// End returns the position of first character immediately after the node.
func (e *CalleeLit) End() Pos {
	return e.TokenPos + 6 // len(callee) == 15
}

func (e *CalleeLit) String() string {
	return "callee"
}

// CalledArgsLit represents an literal called arguments.
type CalledArgsLit struct {
	TokenPos Pos
}

func (e *CalledArgsLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *CalledArgsLit) Pos() Pos {
	return e.TokenPos
}

// End returns the position of first character immediately after the node.
func (e *CalledArgsLit) End() Pos {
	return e.TokenPos + 11 // len("called_args") == 11
}

func (e *CalledArgsLit) String() string {
	return "called_args"
}

// CalledArgsLit represents an literal called arguments.
type CalledKwargsLit struct {
	TokenPos Pos
}

func (e *CalledKwargsLit) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *CalledKwargsLit) Pos() Pos {
	return e.TokenPos
}

// End returns the position of first character immediately after the node.
func (e *CalledKwargsLit) End() Pos {
	return e.TokenPos + 13 // len("called_kwargs") == 13
}

func (e *CalledKwargsLit) String() string {
	return "called_kwargs"
}
