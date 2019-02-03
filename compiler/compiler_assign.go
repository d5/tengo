package compiler

import (
	"errors"
	"fmt"

	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
)

func (c *Compiler) compileAssign(lhs, rhs []ast.Expr, op token.Token) error {
	numLHS, numRHS := len(lhs), len(rhs)
	if numLHS < numRHS {
		// # of LHS must be >= # of RHS
		return fmt.Errorf("assigntment count error: %d < %d", numLHS, numRHS)
	}
	if numLHS > 1 {
		// TODO: until we fully implement the tuple assignment
		return fmt.Errorf("tuple assignment not implemented")
	}
	//if numLHS > 1 && op != token.Assign && op != token.Define {
	//	return fmt.Errorf("invalid operator for tuple assignment: %s", op.String())
	//}

	// resolve and compile left-hand side
	ident, selectors, err := resolveAssignLHS(lhs[0])
	if err != nil {
		return err
	}

	numSel := len(selectors)

	if op == token.Define && numSel > 0 {
		// using selector on new variable does not make sense
		return errors.New("cannot use selector with ':='")
	}

	symbol, depth, exists := c.symbolTable.Resolve(ident)
	if op == token.Define {
		if depth == 0 && exists {
			return fmt.Errorf("'%s' redeclared in this block", ident)
		}

		symbol = c.symbolTable.Define(ident)
	} else {
		if !exists {
			return fmt.Errorf("unresolved reference '%s'", ident)
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
		c.emit(OpAdd)
	case token.SubAssign:
		c.emit(OpSub)
	case token.MulAssign:
		c.emit(OpMul)
	case token.QuoAssign:
		c.emit(OpDiv)
	case token.RemAssign:
		c.emit(OpRem)
	case token.AndAssign:
		c.emit(OpBAnd)
	case token.OrAssign:
		c.emit(OpBOr)
	case token.AndNotAssign:
		c.emit(OpBAndNot)
	case token.XorAssign:
		c.emit(OpBXor)
	case token.ShlAssign:
		c.emit(OpBShiftLeft)
	case token.ShrAssign:
		c.emit(OpBShiftRight)
	}

	// compile selector expressions (right to left)
	for i := numSel - 1; i >= 0; i-- {
		if err := c.Compile(selectors[i]); err != nil {
			return err
		}
	}

	switch symbol.Scope {
	case ScopeGlobal:
		if numSel > 0 {
			c.emit(OpSetSelGlobal, symbol.Index, numSel)
		} else {
			c.emit(OpSetGlobal, symbol.Index)
		}
	case ScopeLocal:
		if numSel > 0 {
			c.emit(OpSetSelLocal, symbol.Index, numSel)
		} else {
			if op == token.Define && !symbol.LocalAssigned {
				c.emit(OpDefineLocal, symbol.Index)
			} else {
				c.emit(OpSetLocal, symbol.Index)
			}
		}

		// mark the symbol as local-assigned
		symbol.LocalAssigned = true
	case ScopeFree:
		if numSel > 0 {
			c.emit(OpSetSelFree, symbol.Index, numSel)
		} else {
			c.emit(OpSetFree, symbol.Index)
		}
	default:
		return fmt.Errorf("invalid assignment variable scope: %s", symbol.Scope)
	}

	return nil
}

func resolveAssignLHS(expr ast.Expr) (name string, selectors []ast.Expr, err error) {
	switch term := expr.(type) {
	case *ast.SelectorExpr:
		name, selectors, err = resolveAssignLHS(term.Expr)
		if err != nil {
			return
		}

		selectors = append(selectors, term.Sel)

		return

	case *ast.IndexExpr:
		name, selectors, err = resolveAssignLHS(term.Expr)
		if err != nil {
			return
		}

		selectors = append(selectors, term.Index)

	case *ast.Ident:
		return term.Name, nil, nil
	}

	return
}
