package compiler_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
)

func TestSymbolTable(t *testing.T) {
	/*
		GLOBAL
		[0] a
		[1] b

			LOCAL 1
			[0] d

				LOCAL 2
				[0] e
				[1] f

				LOCAL 2 BLOCK 1
				[2] g
				[3] h

				LOCAL 2 BLOCK 2
				[2] i
				[3] j
				[4] k

			LOCAL 1 BLOCK 1
			[1] l
			[2] m
			[3] n
			[4] o
			[5] p

				LOCAL 3
				[0] q
				[1] r
	*/

	global := symbolTable()
	assert.Equal(t, globalSymbol("a", 0), global.Define("a"))
	assert.Equal(t, globalSymbol("b", 1), global.Define("b"))

	local1 := global.Fork(false)
	assert.Equal(t, localSymbol("d", 0), local1.Define("d"))

	local1Block1 := local1.Fork(true)
	assert.Equal(t, localSymbol("l", 1), local1Block1.Define("l"))
	assert.Equal(t, localSymbol("m", 2), local1Block1.Define("m"))
	assert.Equal(t, localSymbol("n", 3), local1Block1.Define("n"))
	assert.Equal(t, localSymbol("o", 4), local1Block1.Define("o"))
	assert.Equal(t, localSymbol("p", 5), local1Block1.Define("p"))

	local2 := local1.Fork(false)
	assert.Equal(t, localSymbol("e", 0), local2.Define("e"))
	assert.Equal(t, localSymbol("f", 1), local2.Define("f"))

	local2Block1 := local2.Fork(true)
	assert.Equal(t, localSymbol("g", 2), local2Block1.Define("g"))
	assert.Equal(t, localSymbol("h", 3), local2Block1.Define("h"))

	local2Block2 := local2.Fork(true)
	assert.Equal(t, localSymbol("i", 2), local2Block2.Define("i"))
	assert.Equal(t, localSymbol("j", 3), local2Block2.Define("j"))
	assert.Equal(t, localSymbol("k", 4), local2Block2.Define("k"))

	local3 := local1Block1.Fork(false)
	assert.Equal(t, localSymbol("q", 0), local3.Define("q"))
	assert.Equal(t, localSymbol("r", 1), local3.Define("r"))

	assert.Equal(t, 2, global.MaxSymbols())
	assert.Equal(t, 6, local1.MaxSymbols())
	assert.Equal(t, 6, local1Block1.MaxSymbols())
	assert.Equal(t, 5, local2.MaxSymbols())
	assert.Equal(t, 4, local2Block1.MaxSymbols())
	assert.Equal(t, 5, local2Block2.MaxSymbols())
	assert.Equal(t, 2, local3.MaxSymbols())

	resolveExpect(t, global, "a", globalSymbol("a", 0), 0)
	resolveExpect(t, local1, "d", localSymbol("d", 0), 0)
	resolveExpect(t, local1, "a", globalSymbol("a", 0), 1)
	resolveExpect(t, local3, "a", globalSymbol("a", 0), 3)
	resolveExpect(t, local3, "d", freeSymbol("d", 0), 2)
	resolveExpect(t, local3, "r", localSymbol("r", 1), 0)
	resolveExpect(t, local2Block2, "k", localSymbol("k", 4), 0)
	resolveExpect(t, local2Block2, "e", localSymbol("e", 0), 1)
	resolveExpect(t, local2Block2, "b", globalSymbol("b", 1), 3)
}

func symbol(name string, scope compiler.SymbolScope, index int) *compiler.Symbol {
	return &compiler.Symbol{
		Name:  name,
		Scope: scope,
		Index: index,
	}
}

func globalSymbol(name string, index int) *compiler.Symbol {
	return symbol(name, compiler.ScopeGlobal, index)
}

func localSymbol(name string, index int) *compiler.Symbol {
	return symbol(name, compiler.ScopeLocal, index)
}

func freeSymbol(name string, index int) *compiler.Symbol {
	return symbol(name, compiler.ScopeFree, index)
}

func symbolTable() *compiler.SymbolTable {
	return compiler.NewSymbolTable()
}

func resolveExpect(t *testing.T, symbolTable *compiler.SymbolTable, name string, expectedSymbol *compiler.Symbol, expectedDepth int) {
	actualSymbol, actualDepth, ok := symbolTable.Resolve(name)
	assert.True(t, ok)
	assert.Equal(t, expectedSymbol, actualSymbol)
	assert.Equal(t, expectedDepth, actualDepth)
}
