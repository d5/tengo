package tengo

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/d5/tengo/internal"
	"github.com/d5/tengo/internal/token"
)

// Compiler compiles the AST into a bytecode.
type Compiler struct {
	file            *internal.SourceFile
	parent          *Compiler
	modulePath      string
	constants       []Object
	symbolTable     *internal.SymbolTable
	scopes          []internal.CompilationScope
	scopeIndex      int
	modules         *ModuleMap
	compiledModules map[string]*CompiledFunction
	allowFileImport bool
	loops           []*internal.Loop
	loopIndex       int
	trace           io.Writer
	indent          int
}

// NewCompiler creates a Compiler.
func NewCompiler(
	file *internal.SourceFile,
	symbolTable *internal.SymbolTable,
	constants []Object,
	modules *ModuleMap,
	trace io.Writer,
) *Compiler {
	mainScope := internal.CompilationScope{
		SymbolInit: make(map[string]bool),
		SourceMap:  make(map[int]internal.Pos),
	}

	// symbol table
	if symbolTable == nil {
		symbolTable = internal.NewSymbolTable()
	}

	// add builtin functions to the symbol table
	for idx, fn := range builtinFuncs {
		symbolTable.DefineBuiltin(idx, fn.Name)
	}

	// builtin modules
	if modules == nil {
		modules = NewModuleMap()
	}

	return &Compiler{
		file:            file,
		symbolTable:     symbolTable,
		constants:       constants,
		scopes:          []internal.CompilationScope{mainScope},
		scopeIndex:      0,
		loopIndex:       -1,
		trace:           trace,
		modules:         modules,
		compiledModules: make(map[string]*CompiledFunction),
	}
}

// Compile compiles the AST node.
func (c *Compiler) Compile(node internal.Node) error {
	if c.trace != nil {
		if node != nil {
			defer untracec(tracec(c, fmt.Sprintf("%s (%s)",
				node.String(), reflect.TypeOf(node).Elem().Name())))
		} else {
			defer untracec(tracec(c, "<nil>"))
		}
	}

	switch node := node.(type) {
	case *internal.File:
		for _, stmt := range node.Stmts {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}
	case *internal.ExprStmt:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
		c.emit(node, internal.OpPop)
	case *internal.IncDecStmt:
		op := token.AddAssign
		if node.Token == token.Dec {
			op = token.SubAssign
		}
		return c.compileAssign(node, []internal.Expr{node.Expr},
			[]internal.Expr{&internal.IntLit{Value: 1}}, op)
	case *internal.ParenExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
	case *internal.BinaryExpr:
		if node.Token == token.LAnd || node.Token == token.LOr {
			return c.compileLogical(node)
		}
		if node.Token == token.Less {
			if err := c.Compile(node.RHS); err != nil {
				return err
			}
			if err := c.Compile(node.LHS); err != nil {
				return err
			}
			c.emit(node, internal.OpBinaryOp, int(token.Greater))
			return nil
		} else if node.Token == token.LessEq {
			if err := c.Compile(node.RHS); err != nil {
				return err
			}
			if err := c.Compile(node.LHS); err != nil {
				return err
			}
			c.emit(node, internal.OpBinaryOp, int(token.GreaterEq))
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
			c.emit(node, internal.OpBinaryOp, int(token.Add))
		case token.Sub:
			c.emit(node, internal.OpBinaryOp, int(token.Sub))
		case token.Mul:
			c.emit(node, internal.OpBinaryOp, int(token.Mul))
		case token.Quo:
			c.emit(node, internal.OpBinaryOp, int(token.Quo))
		case token.Rem:
			c.emit(node, internal.OpBinaryOp, int(token.Rem))
		case token.Greater:
			c.emit(node, internal.OpBinaryOp, int(token.Greater))
		case token.GreaterEq:
			c.emit(node, internal.OpBinaryOp, int(token.GreaterEq))
		case token.Equal:
			c.emit(node, internal.OpEqual)
		case token.NotEqual:
			c.emit(node, internal.OpNotEqual)
		case token.And:
			c.emit(node, internal.OpBinaryOp, int(token.And))
		case token.Or:
			c.emit(node, internal.OpBinaryOp, int(token.Or))
		case token.Xor:
			c.emit(node, internal.OpBinaryOp, int(token.Xor))
		case token.AndNot:
			c.emit(node, internal.OpBinaryOp, int(token.AndNot))
		case token.Shl:
			c.emit(node, internal.OpBinaryOp, int(token.Shl))
		case token.Shr:
			c.emit(node, internal.OpBinaryOp, int(token.Shr))
		default:
			return c.errorf(node, "invalid binary operator: %s",
				node.Token.String())
		}
	case *internal.IntLit:
		c.emit(node, internal.OpConstant,
			c.addConstant(&Int{Value: node.Value}))
	case *internal.FloatLit:
		c.emit(node, internal.OpConstant,
			c.addConstant(&Float{Value: node.Value}))
	case *internal.BoolLit:
		if node.Value {
			c.emit(node, internal.OpTrue)
		} else {
			c.emit(node, internal.OpFalse)
		}
	case *internal.StringLit:
		if len(node.Value) > MaxStringLen {
			return c.error(node, ErrStringLimit)
		}
		c.emit(node, internal.OpConstant,
			c.addConstant(&String{Value: node.Value}))
	case *internal.CharLit:
		c.emit(node, internal.OpConstant,
			c.addConstant(&Char{Value: node.Value}))
	case *internal.UndefinedLit:
		c.emit(node, internal.OpNull)
	case *internal.UnaryExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}

		switch node.Token {
		case token.Not:
			c.emit(node, internal.OpLNot)
		case token.Sub:
			c.emit(node, internal.OpMinus)
		case token.Xor:
			c.emit(node, internal.OpBComplement)
		case token.Add:
			// do nothing?
		default:
			return c.errorf(node,
				"invalid unary operator: %s", node.Token.String())
		}
	case *internal.IfStmt:
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
		jumpPos1 := c.emit(node, internal.OpJumpFalsy, 0)
		if err := c.Compile(node.Body); err != nil {
			return err
		}
		if node.Else != nil {
			// second jump placeholder
			jumpPos2 := c.emit(node, internal.OpJump, 0)

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
	case *internal.ForStmt:
		return c.compileForStmt(node)
	case *internal.ForInStmt:
		return c.compileForInStmt(node)
	case *internal.BranchStmt:
		if node.Token == token.Break {
			curLoop := c.currentLoop()
			if curLoop == nil {
				return c.errorf(node, "break not allowed outside loop")
			}
			pos := c.emit(node, internal.OpJump, 0)
			curLoop.Breaks = append(curLoop.Breaks, pos)
		} else if node.Token == token.Continue {
			curLoop := c.currentLoop()
			if curLoop == nil {
				return c.errorf(node, "continue not allowed outside loop")
			}
			pos := c.emit(node, internal.OpJump, 0)
			curLoop.Continues = append(curLoop.Continues, pos)
		} else {
			panic(fmt.Errorf("invalid branch statement: %s",
				node.Token.String()))
		}
	case *internal.BlockStmt:
		if len(node.Stmts) == 0 {
			return nil
		}

		c.symbolTable = c.symbolTable.Fork(true)
		defer func() {
			c.symbolTable = c.symbolTable.Parent(false)
		}()

		for _, stmt := range node.Stmts {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}
	case *internal.AssignStmt:
		err := c.compileAssign(node, node.LHS, node.RHS, node.Token)
		if err != nil {
			return err
		}
	case *internal.Ident:
		symbol, _, ok := c.symbolTable.Resolve(node.Name)
		if !ok {
			return c.errorf(node, "unresolved reference '%s'", node.Name)
		}

		switch symbol.Scope {
		case internal.ScopeGlobal:
			c.emit(node, internal.OpGetGlobal, symbol.Index)
		case internal.ScopeLocal:
			c.emit(node, internal.OpGetLocal, symbol.Index)
		case internal.ScopeBuiltin:
			c.emit(node, internal.OpGetBuiltin, symbol.Index)
		case internal.ScopeFree:
			c.emit(node, internal.OpGetFree, symbol.Index)
		}
	case *internal.ArrayLit:
		for _, elem := range node.Elements {
			if err := c.Compile(elem); err != nil {
				return err
			}
		}
		c.emit(node, internal.OpArray, len(node.Elements))
	case *internal.MapLit:
		for _, elt := range node.Elements {
			// key
			if len(elt.Key) > MaxStringLen {
				return c.error(node, ErrStringLimit)
			}
			c.emit(node, internal.OpConstant,
				c.addConstant(&String{Value: elt.Key}))

			// value
			if err := c.Compile(elt.Value); err != nil {
				return err
			}
		}
		c.emit(node, internal.OpMap, len(node.Elements)*2)

	case *internal.SelectorExpr: // selector on RHS side
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
		if err := c.Compile(node.Sel); err != nil {
			return err
		}
		c.emit(node, internal.OpIndex)
	case *internal.IndexExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
		if err := c.Compile(node.Index); err != nil {
			return err
		}
		c.emit(node, internal.OpIndex)
	case *internal.SliceExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
		if node.Low != nil {
			if err := c.Compile(node.Low); err != nil {
				return err
			}
		} else {
			c.emit(node, internal.OpNull)
		}
		if node.High != nil {
			if err := c.Compile(node.High); err != nil {
				return err
			}
		} else {
			c.emit(node, internal.OpNull)
		}
		c.emit(node, internal.OpSliceIndex)
	case *internal.FuncLit:
		c.enterScope()

		for _, p := range node.Type.Params.List {
			s := c.symbolTable.Define(p.Name)

			// function arguments is not assigned directly.
			s.LocalAssigned = true
		}

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		// code optimization
		c.optimizeFunc(node)

		freeSymbols := c.symbolTable.FreeSymbols()
		numLocals := c.symbolTable.MaxSymbols()
		instructions, sourceMap := c.leaveScope()

		for _, s := range freeSymbols {
			switch s.Scope {
			case internal.ScopeLocal:
				if !s.LocalAssigned {
					// Here, the closure is capturing a local variable that's
					// not yet assigned its value. One example is a local
					// recursive function:
					//
					//   func() {
					//     foo := func(x) {
					//       // ..
					//       return foo(x-1)
					//     }
					//   }
					//
					// which translate into
					//
					//   0000 GETL    0
					//   0002 CLOSURE ?     1
					//   0006 DEFL    0
					//
					// . So the local variable (0) is being captured before
					// it's assigned the value.
					//
					// Solution is to transform the code into something like
					// this:
					//
					//   func() {
					//     foo := undefined
					//     foo = func(x) {
					//       // ..
					//       return foo(x-1)
					//     }
					//   }
					//
					// that is equivalent to
					//
					//   0000 NULL
					//   0001 DEFL    0
					//   0003 GETL    0
					//   0005 CLOSURE ?     1
					//   0009 SETL    0
					//
					c.emit(node, internal.OpNull)
					c.emit(node, internal.OpDefineLocal, s.Index)
					s.LocalAssigned = true
				}
				c.emit(node, internal.OpGetLocalPtr, s.Index)
			case internal.ScopeFree:
				c.emit(node, internal.OpGetFreePtr, s.Index)
			}
		}

		compiledFunction := &CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Type.Params.List),
			VarArgs:       node.Type.Params.VarArgs,
			SourceMap:     sourceMap,
		}
		if len(freeSymbols) > 0 {
			c.emit(node, internal.OpClosure,
				c.addConstant(compiledFunction), len(freeSymbols))
		} else {
			c.emit(node, internal.OpConstant, c.addConstant(compiledFunction))
		}
	case *internal.ReturnStmt:
		if c.symbolTable.Parent(true) == nil {
			// outside the function
			return c.errorf(node, "return not allowed outside function")
		}

		if node.Result == nil {
			c.emit(node, internal.OpReturn, 0)
		} else {
			if err := c.Compile(node.Result); err != nil {
				return err
			}
			c.emit(node, internal.OpReturn, 1)
		}
	case *internal.CallExpr:
		if err := c.Compile(node.Func); err != nil {
			return err
		}
		for _, arg := range node.Args {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}
		c.emit(node, internal.OpCall, len(node.Args))
	case *internal.ImportExpr:
		if node.ModuleName == "" {
			return c.errorf(node, "empty module name")
		}

		if mod := c.modules.Get(node.ModuleName); mod != nil {
			v, err := mod.Import(node.ModuleName)
			if err != nil {
				return err
			}

			switch v := v.(type) {
			case []byte: // module written in Tengo
				compiled, err := c.compileModule(node,
					node.ModuleName, node.ModuleName, v)
				if err != nil {
					return err
				}
				c.emit(node, internal.OpConstant, c.addConstant(compiled))
				c.emit(node, internal.OpCall, 0)
			case Object: // builtin module
				c.emit(node, internal.OpConstant, c.addConstant(v))
			default:
				panic(fmt.Errorf("invalid import value type: %T", v))
			}
		} else if c.allowFileImport {
			moduleName := node.ModuleName
			if !strings.HasSuffix(moduleName, ".tengo") {
				moduleName += ".tengo"
			}

			modulePath, err := filepath.Abs(moduleName)
			if err != nil {
				return c.errorf(node, "module file path error: %s",
					err.Error())
			}

			if err := c.checkCyclicImports(node, modulePath); err != nil {
				return err
			}

			moduleSrc, err := ioutil.ReadFile(moduleName)
			if err != nil {
				return c.errorf(node, "module file read error: %s",
					err.Error())
			}

			compiled, err := c.compileModule(node,
				moduleName, modulePath, moduleSrc)
			if err != nil {
				return err
			}
			c.emit(node, internal.OpConstant, c.addConstant(compiled))
			c.emit(node, internal.OpCall, 0)
		} else {
			return c.errorf(node, "module '%s' not found", node.ModuleName)
		}
	case *internal.ExportStmt:
		// export statement must be in top-level scope
		if c.scopeIndex != 0 {
			return c.errorf(node, "export not allowed inside function")
		}

		// export statement is simply ignore when compiling non-module code
		if c.parent == nil {
			break
		}
		if err := c.Compile(node.Result); err != nil {
			return err
		}
		c.emit(node, internal.OpImmutable)
		c.emit(node, internal.OpReturn, 1)
	case *internal.ErrorExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
		c.emit(node, internal.OpError)
	case *internal.ImmutableExpr:
		if err := c.Compile(node.Expr); err != nil {
			return err
		}
		c.emit(node, internal.OpImmutable)
	case *internal.CondExpr:
		if err := c.Compile(node.Cond); err != nil {
			return err
		}

		// first jump placeholder
		jumpPos1 := c.emit(node, internal.OpJumpFalsy, 0)
		if err := c.Compile(node.True); err != nil {
			return err
		}

		// second jump placeholder
		jumpPos2 := c.emit(node, internal.OpJump, 0)

		// update first jump offset
		curPos := len(c.currentInstructions())
		c.changeOperand(jumpPos1, curPos)
		if err := c.Compile(node.False); err != nil {
			return err
		}

		// update second jump offset
		curPos = len(c.currentInstructions())
		c.changeOperand(jumpPos2, curPos)
	}
	return nil
}

// Bytecode returns a compiled bytecode.
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		FileSet: c.file.Set(),
		MainFunction: &CompiledFunction{
			Instructions: append(c.currentInstructions(), internal.OpSuspend),
			SourceMap:    c.currentSourceMap(),
		},
		Constants: c.constants,
	}
}

// EnableFileImport enables or disables module loading from local files.
// Local file modules are disabled by default.
func (c *Compiler) EnableFileImport(enable bool) {
	c.allowFileImport = enable
}

func (c *Compiler) compileAssign(
	node internal.Node,
	lhs, rhs []internal.Expr,
	op token.Token,
) error {
	numLHS, numRHS := len(lhs), len(rhs)
	if numLHS > 1 || numRHS > 1 {
		return c.errorf(node, "tuple assignment not allowed")
	}

	// resolve and compile left-hand side
	ident, selectors := resolveAssignLHS(lhs[0])
	numSel := len(selectors)

	if op == token.Define && numSel > 0 {
		// using selector on new variable does not make sense
		return c.errorf(node, "operator ':=' not allowed with selector")
	}

	symbol, depth, exists := c.symbolTable.Resolve(ident)
	if op == token.Define {
		if depth == 0 && exists {
			return c.errorf(node, "'%s' redeclared in this block", ident)
		}
		symbol = c.symbolTable.Define(ident)
	} else {
		if !exists {
			return c.errorf(node, "unresolved reference '%s'", ident)
		}
	}

	// +=, -=, *=, /=
	if op != token.Assign && op != token.Define {
		if err := c.Compile(lhs[0]); err != nil {
			return err
		}
	}

	// compile RHSs
	for _, expr := range rhs {
		if err := c.Compile(expr); err != nil {
			return err
		}
	}

	switch op {
	case token.AddAssign:
		c.emit(node, internal.OpBinaryOp, int(token.Add))
	case token.SubAssign:
		c.emit(node, internal.OpBinaryOp, int(token.Sub))
	case token.MulAssign:
		c.emit(node, internal.OpBinaryOp, int(token.Mul))
	case token.QuoAssign:
		c.emit(node, internal.OpBinaryOp, int(token.Quo))
	case token.RemAssign:
		c.emit(node, internal.OpBinaryOp, int(token.Rem))
	case token.AndAssign:
		c.emit(node, internal.OpBinaryOp, int(token.And))
	case token.OrAssign:
		c.emit(node, internal.OpBinaryOp, int(token.Or))
	case token.AndNotAssign:
		c.emit(node, internal.OpBinaryOp, int(token.AndNot))
	case token.XorAssign:
		c.emit(node, internal.OpBinaryOp, int(token.Xor))
	case token.ShlAssign:
		c.emit(node, internal.OpBinaryOp, int(token.Shl))
	case token.ShrAssign:
		c.emit(node, internal.OpBinaryOp, int(token.Shr))
	}

	// compile selector expressions (right to left)
	for i := numSel - 1; i >= 0; i-- {
		if err := c.Compile(selectors[i]); err != nil {
			return err
		}
	}

	switch symbol.Scope {
	case internal.ScopeGlobal:
		if numSel > 0 {
			c.emit(node, internal.OpSetSelGlobal, symbol.Index, numSel)
		} else {
			c.emit(node, internal.OpSetGlobal, symbol.Index)
		}
	case internal.ScopeLocal:
		if numSel > 0 {
			c.emit(node, internal.OpSetSelLocal, symbol.Index, numSel)
		} else {
			if op == token.Define && !symbol.LocalAssigned {
				c.emit(node, internal.OpDefineLocal, symbol.Index)
			} else {
				c.emit(node, internal.OpSetLocal, symbol.Index)
			}
		}

		// mark the symbol as local-assigned
		symbol.LocalAssigned = true
	case internal.ScopeFree:
		if numSel > 0 {
			c.emit(node, internal.OpSetSelFree, symbol.Index, numSel)
		} else {
			c.emit(node, internal.OpSetFree, symbol.Index)
		}
	default:
		panic(fmt.Errorf("invalid assignment variable scope: %s",
			symbol.Scope))
	}
	return nil
}

func (c *Compiler) compileLogical(node *internal.BinaryExpr) error {
	// left side term
	if err := c.Compile(node.LHS); err != nil {
		return err
	}

	// jump position
	var jumpPos int
	if node.Token == token.LAnd {
		jumpPos = c.emit(node, internal.OpAndJump, 0)
	} else {
		jumpPos = c.emit(node, internal.OpOrJump, 0)
	}

	// right side term
	if err := c.Compile(node.RHS); err != nil {
		return err
	}

	c.changeOperand(jumpPos, len(c.currentInstructions()))
	return nil
}

func (c *Compiler) compileForStmt(stmt *internal.ForStmt) error {
	c.symbolTable = c.symbolTable.Fork(true)
	defer func() {
		c.symbolTable = c.symbolTable.Parent(false)
	}()

	// init statement
	if stmt.Init != nil {
		if err := c.Compile(stmt.Init); err != nil {
			return err
		}
	}

	// pre-condition position
	preCondPos := len(c.currentInstructions())

	// condition expression
	postCondPos := -1
	if stmt.Cond != nil {
		if err := c.Compile(stmt.Cond); err != nil {
			return err
		}
		// condition jump position
		postCondPos = c.emit(stmt, internal.OpJumpFalsy, 0)
	}

	// enter loop
	loop := c.enterLoop()

	// body statement
	if err := c.Compile(stmt.Body); err != nil {
		c.leaveLoop()
		return err
	}

	c.leaveLoop()

	// post-body position
	postBodyPos := len(c.currentInstructions())

	// post statement
	if stmt.Post != nil {
		if err := c.Compile(stmt.Post); err != nil {
			return err
		}
	}

	// back to condition
	c.emit(stmt, internal.OpJump, preCondPos)

	// post-statement position
	postStmtPos := len(c.currentInstructions())
	if postCondPos >= 0 {
		c.changeOperand(postCondPos, postStmtPos)
	}

	// update all break/continue jump positions
	for _, pos := range loop.Breaks {
		c.changeOperand(pos, postStmtPos)
	}
	for _, pos := range loop.Continues {
		c.changeOperand(pos, postBodyPos)
	}
	return nil
}

func (c *Compiler) compileForInStmt(stmt *internal.ForInStmt) error {
	c.symbolTable = c.symbolTable.Fork(true)
	defer func() {
		c.symbolTable = c.symbolTable.Parent(false)
	}()

	// for-in statement is compiled like following:
	//
	//   for :it := iterator(iterable); :it.next();  {
	//     k, v := :it.get()  // DEFINE operator
	//
	//     ... body ...
	//   }
	//
	// ":it" is a local variable but will be conflict with other user variables
	// because character ":" is not allowed.

	// init
	//   :it = iterator(iterable)
	itSymbol := c.symbolTable.Define(":it")
	if err := c.Compile(stmt.Iterable); err != nil {
		return err
	}
	c.emit(stmt, internal.OpIteratorInit)
	if itSymbol.Scope == internal.ScopeGlobal {
		c.emit(stmt, internal.OpSetGlobal, itSymbol.Index)
	} else {
		c.emit(stmt, internal.OpDefineLocal, itSymbol.Index)
	}

	// pre-condition position
	preCondPos := len(c.currentInstructions())

	// condition
	//  :it.HasMore()
	if itSymbol.Scope == internal.ScopeGlobal {
		c.emit(stmt, internal.OpGetGlobal, itSymbol.Index)
	} else {
		c.emit(stmt, internal.OpGetLocal, itSymbol.Index)
	}
	c.emit(stmt, internal.OpIteratorNext)

	// condition jump position
	postCondPos := c.emit(stmt, internal.OpJumpFalsy, 0)

	// enter loop
	loop := c.enterLoop()

	// assign key variable
	if stmt.Key.Name != "_" {
		keySymbol := c.symbolTable.Define(stmt.Key.Name)
		if itSymbol.Scope == internal.ScopeGlobal {
			c.emit(stmt, internal.OpGetGlobal, itSymbol.Index)
		} else {
			c.emit(stmt, internal.OpGetLocal, itSymbol.Index)
		}
		c.emit(stmt, internal.OpIteratorKey)
		if keySymbol.Scope == internal.ScopeGlobal {
			c.emit(stmt, internal.OpSetGlobal, keySymbol.Index)
		} else {
			c.emit(stmt, internal.OpDefineLocal, keySymbol.Index)
		}
	}

	// assign value variable
	if stmt.Value.Name != "_" {
		valueSymbol := c.symbolTable.Define(stmt.Value.Name)
		if itSymbol.Scope == internal.ScopeGlobal {
			c.emit(stmt, internal.OpGetGlobal, itSymbol.Index)
		} else {
			c.emit(stmt, internal.OpGetLocal, itSymbol.Index)
		}
		c.emit(stmt, internal.OpIteratorValue)
		if valueSymbol.Scope == internal.ScopeGlobal {
			c.emit(stmt, internal.OpSetGlobal, valueSymbol.Index)
		} else {
			c.emit(stmt, internal.OpDefineLocal, valueSymbol.Index)
		}
	}

	// body statement
	if err := c.Compile(stmt.Body); err != nil {
		c.leaveLoop()
		return err
	}

	c.leaveLoop()

	// post-body position
	postBodyPos := len(c.currentInstructions())

	// back to condition
	c.emit(stmt, internal.OpJump, preCondPos)

	// post-statement position
	postStmtPos := len(c.currentInstructions())
	c.changeOperand(postCondPos, postStmtPos)

	// update all break/continue jump positions
	for _, pos := range loop.Breaks {
		c.changeOperand(pos, postStmtPos)
	}
	for _, pos := range loop.Continues {
		c.changeOperand(pos, postBodyPos)
	}
	return nil
}

func (c *Compiler) checkCyclicImports(
	node internal.Node,
	modulePath string,
) error {
	if c.modulePath == modulePath {
		return c.errorf(node, "cyclic module import: %s", modulePath)
	} else if c.parent != nil {
		return c.parent.checkCyclicImports(node, modulePath)
	}
	return nil
}

func (c *Compiler) compileModule(
	node internal.Node,
	moduleName, modulePath string,
	src []byte,
) (*CompiledFunction, error) {
	if err := c.checkCyclicImports(node, modulePath); err != nil {
		return nil, err
	}

	compiledModule, exists := c.loadCompiledModule(modulePath)
	if exists {
		return compiledModule, nil
	}

	modFile := c.file.Set().AddFile(moduleName, -1, len(src))
	p := internal.NewParser(modFile, src, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}

	// inherit builtin functions
	symbolTable := internal.NewSymbolTable()
	for _, sym := range c.symbolTable.BuiltinSymbols() {
		symbolTable.DefineBuiltin(sym.Index, sym.Name)
	}

	// no global scope for the module
	symbolTable = symbolTable.Fork(false)

	// compile module
	moduleCompiler := c.fork(modFile, modulePath, symbolTable)
	if err := moduleCompiler.Compile(file); err != nil {
		return nil, err
	}

	// code optimization
	moduleCompiler.optimizeFunc(node)
	compiledFunc := moduleCompiler.Bytecode().MainFunction
	compiledFunc.NumLocals = symbolTable.MaxSymbols()
	c.storeCompiledModule(modulePath, compiledFunc)
	return compiledFunc, nil
}

func (c *Compiler) loadCompiledModule(
	modulePath string,
) (mod *CompiledFunction, ok bool) {
	if c.parent != nil {
		return c.parent.loadCompiledModule(modulePath)
	}
	mod, ok = c.compiledModules[modulePath]
	return
}

func (c *Compiler) storeCompiledModule(
	modulePath string,
	module *CompiledFunction,
) {
	if c.parent != nil {
		c.parent.storeCompiledModule(modulePath, module)
	}
	c.compiledModules[modulePath] = module
}

func (c *Compiler) enterLoop() *internal.Loop {
	loop := &internal.Loop{}
	c.loops = append(c.loops, loop)
	c.loopIndex++
	if c.trace != nil {
		c.printTrace("LOOPE", c.loopIndex)
	}
	return loop
}

func (c *Compiler) leaveLoop() {
	if c.trace != nil {
		c.printTrace("LOOPL", c.loopIndex)
	}
	c.loops = c.loops[:len(c.loops)-1]
	c.loopIndex--
}

func (c *Compiler) currentLoop() *internal.Loop {
	if c.loopIndex >= 0 {
		return c.loops[c.loopIndex]
	}
	return nil
}

func (c *Compiler) currentInstructions() []byte {
	return c.scopes[c.scopeIndex].Instructions
}

func (c *Compiler) currentSourceMap() map[int]internal.Pos {
	return c.scopes[c.scopeIndex].SourceMap
}

func (c *Compiler) enterScope() {
	scope := internal.CompilationScope{
		SymbolInit: make(map[string]bool),
		SourceMap:  make(map[int]internal.Pos),
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = c.symbolTable.Fork(false)
	if c.trace != nil {
		c.printTrace("SCOPE", c.scopeIndex)
	}
}

func (c *Compiler) leaveScope() (
	instructions []byte,
	sourceMap map[int]internal.Pos,
) {
	instructions = c.currentInstructions()
	sourceMap = c.currentSourceMap()
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Parent(true)
	if c.trace != nil {
		c.printTrace("SCOPL", c.scopeIndex)
	}
	return
}

func (c *Compiler) fork(
	file *internal.SourceFile,
	modulePath string,
	symbolTable *internal.SymbolTable,
) *Compiler {
	child := NewCompiler(file, symbolTable, nil, c.modules, c.trace)
	child.modulePath = modulePath // module file path
	child.parent = c              // parent to set to current compiler
	return child
}

func (c *Compiler) error(node internal.Node, err error) error {
	return &internal.CompilerError{
		FileSet: c.file.Set(),
		Node:    node,
		Err:     err,
	}
}

func (c *Compiler) errorf(
	node internal.Node,
	format string,
	args ...interface{},
) error {
	return &internal.CompilerError{
		FileSet: c.file.Set(),
		Node:    node,
		Err:     fmt.Errorf(format, args...),
	}
}

func (c *Compiler) addConstant(o Object) int {
	if c.parent != nil {
		// module compilers will use their parent's constants array
		return c.parent.addConstant(o)
	}
	c.constants = append(c.constants, o)
	if c.trace != nil {
		c.printTrace(fmt.Sprintf("CONST %04d %s", len(c.constants)-1, o))
	}
	return len(c.constants) - 1
}

func (c *Compiler) addInstruction(b []byte) int {
	posNewIns := len(c.currentInstructions())
	c.scopes[c.scopeIndex].Instructions = append(
		c.currentInstructions(), b...)
	return posNewIns
}

func (c *Compiler) replaceInstruction(pos int, inst []byte) {
	copy(c.currentInstructions()[pos:], inst)
	if c.trace != nil {
		c.printTrace(fmt.Sprintf("REPLC %s",
			internal.FormatInstructions(
				c.scopes[c.scopeIndex].Instructions[pos:], pos)[0]))
	}
}

func (c *Compiler) changeOperand(opPos int, operand ...int) {
	op := c.currentInstructions()[opPos]
	inst := internal.MakeInstruction(op, operand...)
	c.replaceInstruction(opPos, inst)
}

// optimizeFunc performs some code-level optimization for the current function
// instructions. It also removes unreachable (dead code) instructions and adds
// "returns" instruction if needed.
func (c *Compiler) optimizeFunc(node internal.Node) {
	// any instructions between RETURN and the function end
	// or instructions between RETURN and jump target position
	// are considered as unreachable.

	// pass 1. identify all jump destinations
	dsts := make(map[int]bool)
	iterateInstructions(c.scopes[c.scopeIndex].Instructions,
		func(pos int, opcode internal.Opcode, operands []int) bool {
			switch opcode {
			case internal.OpJump, internal.OpJumpFalsy,
				internal.OpAndJump, internal.OpOrJump:
				dsts[operands[0]] = true
			}
			return true
		})

	// pass 2. eliminate dead code
	var newInsts []byte
	posMap := make(map[int]int) // old position to new position
	var dstIdx int
	var deadCode bool
	iterateInstructions(c.scopes[c.scopeIndex].Instructions,
		func(pos int, opcode internal.Opcode, operands []int) bool {
			switch {
			case opcode == internal.OpReturn:
				if deadCode {
					return true
				}
				deadCode = true
			case dsts[pos]:
				dstIdx++
				deadCode = false
			case deadCode:
				return true
			}
			posMap[pos] = len(newInsts)
			newInsts = append(newInsts,
				internal.MakeInstruction(opcode, operands...)...)
			return true
		})

	// pass 3. update jump positions
	var lastOp internal.Opcode
	var appendReturn bool
	endPos := len(c.scopes[c.scopeIndex].Instructions)
	iterateInstructions(newInsts,
		func(pos int, opcode internal.Opcode, operands []int) bool {
			switch opcode {
			case internal.OpJump, internal.OpJumpFalsy, internal.OpAndJump,
				internal.OpOrJump:
				newDst, ok := posMap[operands[0]]
				if ok {
					copy(newInsts[pos:],
						internal.MakeInstruction(opcode, newDst))
				} else if endPos == operands[0] {
					// there's a jump instruction that jumps to the end of
					// function compiler should append "return".
					appendReturn = true
				} else {
					panic(fmt.Errorf("invalid jump position: %d", newDst))
				}
			}
			lastOp = opcode
			return true
		})
	if lastOp != internal.OpReturn {
		appendReturn = true
	}

	// pass 4. update source map
	newSourceMap := make(map[int]internal.Pos)
	for pos, srcPos := range c.scopes[c.scopeIndex].SourceMap {
		newPos, ok := posMap[pos]
		if ok {
			newSourceMap[newPos] = srcPos
		}
	}
	c.scopes[c.scopeIndex].Instructions = newInsts
	c.scopes[c.scopeIndex].SourceMap = newSourceMap

	// append "return"
	if appendReturn {
		c.emit(node, internal.OpReturn, 0)
	}
}

func (c *Compiler) emit(
	node internal.Node,
	opcode internal.Opcode,
	operands ...int,
) int {
	filePos := internal.NoPos
	if node != nil {
		filePos = node.Pos()
	}

	inst := internal.MakeInstruction(opcode, operands...)
	pos := c.addInstruction(inst)
	c.scopes[c.scopeIndex].SourceMap[pos] = filePos
	if c.trace != nil {
		c.printTrace(fmt.Sprintf("EMIT  %s",
			internal.FormatInstructions(
				c.scopes[c.scopeIndex].Instructions[pos:], pos)[0]))
	}
	return pos
}

func (c *Compiler) printTrace(a ...interface{}) {
	const (
		dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
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

func resolveAssignLHS(
	expr internal.Expr,
) (name string, selectors []internal.Expr) {
	switch term := expr.(type) {
	case *internal.SelectorExpr:
		name, selectors = resolveAssignLHS(term.Expr)
		selectors = append(selectors, term.Sel)
		return
	case *internal.IndexExpr:
		name, selectors = resolveAssignLHS(term.Expr)
		selectors = append(selectors, term.Index)
	case *internal.Ident:
		name = term.Name
	}
	return
}

func iterateInstructions(
	b []byte,
	fn func(pos int, opcode internal.Opcode, operands []int) bool,
) {
	for i := 0; i < len(b); i++ {
		numOperands := internal.OpcodeOperands[b[i]]
		operands, read := internal.ReadOperands(numOperands, b[i+1:])
		if !fn(i, b[i], operands) {
			break
		}
		i += read
	}
}

func tracec(c *Compiler, msg string) *Compiler {
	c.printTrace(msg, "{")
	c.indent++
	return c
}

func untracec(c *Compiler) {
	c.indent--
	c.printTrace("}")
}
