package compiler

import (
	"fmt"
	"io"
	"reflect"

	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

// Compiler compiles the AST into a bytecode.
type Compiler struct {
	constants   []objects.Object
	symbolTable *SymbolTable
	scopes      []CompilationScope
	scopeIndex  int
	loops       []*Loop
	loopIndex   int
	trace       io.Writer
	indent      int
}

// NewCompiler creates a Compiler.
func NewCompiler(symbolTable *SymbolTable, trace io.Writer) *Compiler {
	mainScope := CompilationScope{
		instructions: make([]byte, 0),
	}

	if symbolTable == nil {
		symbolTable = NewSymbolTable()
	}

	for idx, fn := range objects.Builtins {
		symbolTable.DefineBuiltin(idx, fn.Name)
	}

	return &Compiler{
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
		loopIndex:   -1,
		trace:       trace,
	}
}

// Compile compiles the AST node.
func (c *Compiler) Compile(node ast.Node) error {
	if c.trace != nil {
		if node != nil {
			defer un(trace(c, fmt.Sprintf("%s (%s)", node.String(), reflect.TypeOf(node).Elem().Name())))
		} else {
			defer un(trace(c, "<nil>"))
		}
	}

	switch node := node.(type) {
	case *ast.File:
		for _, stmt := range node.Stmts {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}
	case *ast.ExprStmt:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
		c.emit(OpPop)
	case *ast.IncDecStmt:
		op := token.AddAssign
		if node.Token == token.Dec {
			op = token.SubAssign
		}

		return c.compileAssign([]ast.Expr{node.Expr}, []ast.Expr{&ast.IntLit{Value: 1}}, op)
	case *ast.ParenExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
	case *ast.BinaryExpr:
		if node.Token == token.LAnd || node.Token == token.LOr {
			if err := c.compileLogical(node); err != nil {
				return err
			}

			return nil
		}

		if node.Token == token.Less {
			if err := c.Compile(node.RHS); err != nil {
				return err
			}

			if err := c.Compile(node.LHS); err != nil {
				return err
			}

			c.emit(OpGreaterThan)

			return nil
		} else if node.Token == token.LessEq {
			if err := c.Compile(node.RHS); err != nil {
				return err
			}
			if err := c.Compile(node.LHS); err != nil {
				return err
			}

			c.emit(OpGreaterThanEqual)

			return nil
		}

		if err := c.Compile(node.LHS); err != nil {
			return err
		}
		if err := c.Compile(node.RHS); err != nil {
			return err
		}

		switch node.Token {
		case token.Add:
			c.emit(OpAdd)
		case token.Sub:
			c.emit(OpSub)
		case token.Mul:
			c.emit(OpMul)
		case token.Quo:
			c.emit(OpDiv)
		case token.Rem:
			c.emit(OpRem)
		case token.Greater:
			c.emit(OpGreaterThan)
		case token.GreaterEq:
			c.emit(OpGreaterThanEqual)
		case token.Equal:
			c.emit(OpEqual)
		case token.NotEqual:
			c.emit(OpNotEqual)
		case token.And:
			c.emit(OpBAnd)
		case token.Or:
			c.emit(OpBOr)
		case token.Xor:
			c.emit(OpBXor)
		case token.AndNot:
			c.emit(OpBAndNot)
		case token.Shl:
			c.emit(OpBShiftLeft)
		case token.Shr:
			c.emit(OpBShiftRight)
		default:
			return fmt.Errorf("unknown operator: %s", node.Token.String())
		}
	case *ast.IntLit:
		c.emit(OpConstant, c.addConstant(&objects.Int{Value: node.Value}))
	case *ast.FloatLit:
		c.emit(OpConstant, c.addConstant(&objects.Float{Value: node.Value}))
	case *ast.BoolLit:
		if node.Value {
			c.emit(OpTrue)
		} else {
			c.emit(OpFalse)
		}
	case *ast.StringLit:
		c.emit(OpConstant, c.addConstant(&objects.String{Value: node.Value}))
	case *ast.CharLit:
		c.emit(OpConstant, c.addConstant(&objects.Char{Value: node.Value}))
	case *ast.UndefinedLit:
		c.emit(OpNull)
	case *ast.UnaryExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}

		switch node.Token {
		case token.Not:
			c.emit(OpLNot)
		case token.Sub:
			c.emit(OpMinus)
		case token.Xor:
			c.emit(OpBComplement)
		case token.Add:
			// do nothing?
		default:
			return fmt.Errorf("unknown operator: %s", node.Token.String())
		}
	case *ast.IfStmt:
		// open new symbol table for the statement
		c.symbolTable = c.symbolTable.Fork(true)
		defer func() {
			c.symbolTable = c.symbolTable.Parent(false)
		}()

		if node.Init != nil {
			if err := c.Compile(node.Init); err != nil {
				return err
			}
		}

		if err := c.Compile(node.Cond); err != nil {
			return err
		}

		// first jump placeholder
		jumpPos1 := c.emit(OpJumpFalsy, 0)

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		if node.Else != nil {
			// second jump placeholder
			jumpPos2 := c.emit(OpJump, 0)

			// update first jump offset
			curPos := len(c.currentInstructions())
			c.changeOperand(jumpPos1, curPos)

			if err := c.Compile(node.Else); err != nil {
				return err
			}

			// update second jump offset
			curPos = len(c.currentInstructions())
			c.changeOperand(jumpPos2, curPos)
		} else {
			// update first jump offset
			curPos := len(c.currentInstructions())
			c.changeOperand(jumpPos1, curPos)
		}

	case *ast.ForStmt:
		return c.compileForStmt(node)
	case *ast.ForInStmt:
		return c.compileForInStmt(node)
	case *ast.BranchStmt:
		if node.Token == token.Break {
			curLoop := c.currentLoop()
			if curLoop == nil {
				return fmt.Errorf("break statement outside loop")
			}
			pos := c.emit(OpJump, 0)
			curLoop.Breaks = append(curLoop.Breaks, pos)
		} else if node.Token == token.Continue {
			curLoop := c.currentLoop()
			if curLoop == nil {
				return fmt.Errorf("continue statement outside loop")
			}
			pos := c.emit(OpJump, 0)
			curLoop.Continues = append(curLoop.Continues, pos)
		} else {
			return fmt.Errorf("unknown branch statement: %s", node.Token.String())
		}
	case *ast.BlockStmt:
		for _, stmt := range node.Stmts {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}
	case *ast.AssignStmt:

		if err := c.compileAssign(node.LHS, node.RHS, node.Token); err != nil {
			return err
		}
	case *ast.Ident:
		symbol, _, ok := c.symbolTable.Resolve(node.Name)
		if !ok {
			return fmt.Errorf("undefined variable: %s", node.Name)
		}

		switch symbol.Scope {
		case ScopeGlobal:
			c.emit(OpGetGlobal, symbol.Index)
		case ScopeLocal:
			c.emit(OpGetLocal, symbol.Index)
		case ScopeBuiltin:
			c.emit(OpGetBuiltin, symbol.Index)
		case ScopeFree:
			c.emit(OpGetFree, symbol.Index)
		}
	case *ast.ArrayLit:
		for _, elem := range node.Elements {
			if err := c.Compile(elem); err != nil {
				return err
			}
		}

		c.emit(OpArray, len(node.Elements))
	case *ast.MapLit:
		for _, elt := range node.Elements {
			// key
			c.emit(OpConstant, c.addConstant(&objects.String{Value: elt.Key}))

			// value
			if err := c.Compile(elt.Value); err != nil {
				return err
			}
		}

		c.emit(OpMap, len(node.Elements)*2)
	case *ast.SelectorExpr: // selector on RHS side
		if err := c.Compile(node.Expr); err != nil {
			return err
		}

		if err := c.Compile(node.Sel); err != nil {
			return err
		}

		c.emit(OpIndex)
	case *ast.IndexExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}

		if err := c.Compile(node.Index); err != nil {
			return err
		}

		c.emit(OpIndex)
	case *ast.SliceExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}

		if node.Low != nil {
			if err := c.Compile(node.Low); err != nil {
				return err
			}
		} else {
			c.emit(OpNull)
		}

		if node.High != nil {
			if err := c.Compile(node.High); err != nil {
				return err
			}
		} else {
			c.emit(OpNull)
		}

		c.emit(OpSliceIndex)
	case *ast.FuncLit:
		c.enterScope()

		for _, p := range node.Type.Params.List {
			c.symbolTable.Define(p.Name)
		}

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		// add OpReturn if function returns nothing
		if !c.lastInstructionIs(OpReturnValue) && !c.lastInstructionIs(OpReturn) {
			c.emit(OpReturn)
		}

		freeSymbols := c.symbolTable.FreeSymbols()
		numLocals := c.symbolTable.MaxSymbols()
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			switch s.Scope {
			case ScopeLocal:
				c.emit(OpGetLocal, s.Index)
			case ScopeFree:
				c.emit(OpGetFree, s.Index)
			}
		}

		compiledFunction := &objects.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Type.Params.List),
		}

		if len(freeSymbols) > 0 {
			c.emit(OpClosure, c.addConstant(compiledFunction), len(freeSymbols))
		} else {
			c.emit(OpConstant, c.addConstant(compiledFunction))
		}
	case *ast.ReturnStmt:
		if c.symbolTable.Parent(true) == nil {
			// outside the function
			return fmt.Errorf("return statement outside function")
		}

		switch len(node.Results) {
		case 0:
			c.emit(OpReturn)
		case 1:
			if err := c.Compile(node.Results[0]); err != nil {
				return err
			}

			c.emit(OpReturnValue, 1)
		default:
			return fmt.Errorf("multi-value return not implemented")
		}
	case *ast.CallExpr:
		if err := c.Compile(node.Func); err != nil {
			return err
		}

		for _, arg := range node.Args {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}

		c.emit(OpCall, len(node.Args))
	}

	return nil
}

// Bytecode returns a compiled bytecode.
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (c *Compiler) addConstant(o objects.Object) int {
	c.constants = append(c.constants, o)

	if c.trace != nil {
		c.printTrace(fmt.Sprintf("CONST %04d %s", len(c.constants)-1, o))
	}

	return len(c.constants) - 1
}

func (c *Compiler) addInstruction(b []byte) int {
	posNewIns := len(c.currentInstructions())

	c.scopes[c.scopeIndex].instructions = append(c.currentInstructions(), b...)

	return posNewIns
}

func (c *Compiler) setLastInstruction(op Opcode, pos int) {
	c.scopes[c.scopeIndex].lastInstructions[1] = c.scopes[c.scopeIndex].lastInstructions[0]

	c.scopes[c.scopeIndex].lastInstructions[0].Opcode = op
	c.scopes[c.scopeIndex].lastInstructions[0].Position = pos
}

func (c *Compiler) lastInstructionIs(op Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstructions[0].Opcode == op
}

func (c *Compiler) removeLastInstruction() {
	lastPos := c.scopes[c.scopeIndex].lastInstructions[0].Position

	if c.trace != nil {
		c.printTrace(fmt.Sprintf("DELET %s",
			FormatInstructions(c.scopes[c.scopeIndex].instructions[lastPos:], lastPos)[0]))
	}

	c.scopes[c.scopeIndex].instructions = c.currentInstructions()[:lastPos]
	c.scopes[c.scopeIndex].lastInstructions[0] = c.scopes[c.scopeIndex].lastInstructions[1]
}

func (c *Compiler) replaceInstruction(pos int, inst []byte) {
	copy(c.currentInstructions()[pos:], inst)

	if c.trace != nil {
		c.printTrace(fmt.Sprintf("REPLC %s",
			FormatInstructions(c.scopes[c.scopeIndex].instructions[pos:], pos)[0]))
	}
}

func (c *Compiler) changeOperand(opPos int, operand ...int) {
	op := Opcode(c.currentInstructions()[opPos])
	inst := MakeInstruction(op, operand...)

	c.replaceInstruction(opPos, inst)
}

func (c *Compiler) emit(opcode Opcode, operands ...int) int {
	inst := MakeInstruction(opcode, operands...)
	pos := c.addInstruction(inst)
	c.setLastInstruction(opcode, pos)

	if c.trace != nil {
		c.printTrace(fmt.Sprintf("EMIT  %s",
			FormatInstructions(c.scopes[c.scopeIndex].instructions[pos:], pos)[0]))
	}

	return pos
}

func (c *Compiler) printTrace(a ...interface{}) {
	const (
		dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
		n    = len(dots)
	)

	i := 2 * c.indent
	for i > n {
		_, _ = fmt.Fprint(c.trace, dots)
		i -= n
	}
	_, _ = fmt.Fprint(c.trace, dots[0:i])
	_, _ = fmt.Fprintln(c.trace, a...)
}

func trace(c *Compiler, msg string) *Compiler {
	c.printTrace(msg, "{")
	c.indent++

	return c
}

func un(c *Compiler) {
	c.indent--
	c.printTrace("}")
}
