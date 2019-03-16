package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/objects"
)

func TestCompiler_Compile(t *testing.T) {
	expect(t, `1 + 2`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expect(t, `1; 2`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpPop),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expect(t, `1 - 2`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSub),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expect(t, `1 * 2`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpMul),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expect(t, `2 / 1`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpDiv),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(2),
				intObject(1))))

	expect(t, `true`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpTrue),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray()))

	expect(t, `false`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpFalse),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray()))

	expect(t, `1 > 2`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpGreaterThan),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expect(t, `1 < 2`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpGreaterThan),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(2),
				intObject(1))))

	expect(t, `1 == 2`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpEqual),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expect(t, `1 != 2`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpNotEqual),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expect(t, `true == false`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpTrue),
				compiler.MakeInstruction(compiler.OpFalse),
				compiler.MakeInstruction(compiler.OpEqual),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray()))

	expect(t, `true != false`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpTrue),
				compiler.MakeInstruction(compiler.OpFalse),
				compiler.MakeInstruction(compiler.OpNotEqual),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray()))

	expect(t, `-1`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpMinus),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1))))

	expect(t, `!true`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpTrue),
				compiler.MakeInstruction(compiler.OpLNot),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray()))

	expect(t, `if true { 10 }; 3333`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpTrue),         // 0000
				compiler.MakeInstruction(compiler.OpJumpFalsy, 8), // 0001
				compiler.MakeInstruction(compiler.OpConstant, 0),  // 0004
				compiler.MakeInstruction(compiler.OpPop),          // 0007
				compiler.MakeInstruction(compiler.OpConstant, 1),  // 0008
				compiler.MakeInstruction(compiler.OpPop)),         // 0011
			objectsArray(
				intObject(10),
				intObject(3333))))

	expect(t, `if (true) { 10 } else { 20 }; 3333;`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpTrue),          // 0000
				compiler.MakeInstruction(compiler.OpJumpFalsy, 11), // 0001
				compiler.MakeInstruction(compiler.OpConstant, 0),   // 0004
				compiler.MakeInstruction(compiler.OpPop),           // 0007
				compiler.MakeInstruction(compiler.OpJump, 15),      // 0008
				compiler.MakeInstruction(compiler.OpConstant, 1),   // 0011
				compiler.MakeInstruction(compiler.OpPop),           // 0014
				compiler.MakeInstruction(compiler.OpConstant, 2),   // 0015
				compiler.MakeInstruction(compiler.OpPop)),          // 0018
			objectsArray(
				intObject(10),
				intObject(20),
				intObject(3333))))

	expect(t, `"kami"`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				stringObject("kami"))))

	expect(t, `"ka" + "mi"`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				stringObject("ka"),
				stringObject("mi"))))

	expect(t, `a := 1; b := 2; a += b`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSetGlobal, 1),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 1),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expect(t, `a := 1; b := 2; a /= b`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSetGlobal, 1),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 1),
				compiler.MakeInstruction(compiler.OpDiv),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expect(t, `[]`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpArray, 0),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray()))

	expect(t, `[1, 2, 3]`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpArray, 3),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3))))

	expect(t, `[1 + 2, 3 - 4, 5 * 6]`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSub),
				compiler.MakeInstruction(compiler.OpConstant, 4),
				compiler.MakeInstruction(compiler.OpConstant, 5),
				compiler.MakeInstruction(compiler.OpMul),
				compiler.MakeInstruction(compiler.OpArray, 3),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(4),
				intObject(5),
				intObject(6))))

	expect(t, `{}`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpMap, 0),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray()))

	expect(t, `{a: 2, b: 4, c: 6}`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpConstant, 4),
				compiler.MakeInstruction(compiler.OpConstant, 5),
				compiler.MakeInstruction(compiler.OpMap, 6),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				stringObject("a"),
				intObject(2),
				stringObject("b"),
				intObject(4),
				stringObject("c"),
				intObject(6))))

	expect(t, `{a: 2 + 3, b: 5 * 6}`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpConstant, 4),
				compiler.MakeInstruction(compiler.OpConstant, 5),
				compiler.MakeInstruction(compiler.OpMul),
				compiler.MakeInstruction(compiler.OpMap, 4),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				stringObject("a"),
				intObject(2),
				intObject(3),
				stringObject("b"),
				intObject(5),
				intObject(6))))

	expect(t, `[1, 2, 3][1 + 1]`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpArray, 3),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpConstant, 4),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpIndex),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(1),
				intObject(1))))

	expect(t, `{a: 2}[2 - 1]`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpMap, 2),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSub),
				compiler.MakeInstruction(compiler.OpIndex),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				stringObject("a"),
				intObject(2),
				intObject(2),
				intObject(1))))

	expect(t, `[1, 2, 3][:]`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpArray, 3),
				compiler.MakeInstruction(compiler.OpNull),
				compiler.MakeInstruction(compiler.OpNull),
				compiler.MakeInstruction(compiler.OpSliceIndex),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3))))

	expect(t, `[1, 2, 3][0 : 2]`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpArray, 3),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpConstant, 4),
				compiler.MakeInstruction(compiler.OpSliceIndex),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(0),
				intObject(2))))

	expect(t, `[1, 2, 3][:2]`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpArray, 3),
				compiler.MakeInstruction(compiler.OpNull),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSliceIndex),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(2))))

	expect(t, `[1, 2, 3][0:]`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpArray, 3),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpNull),
				compiler.MakeInstruction(compiler.OpSliceIndex),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(0))))

	expect(t, `func() { return 5 + 10 }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(5),
				intObject(10),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpConstant, 1),
					compiler.MakeInstruction(compiler.OpAdd),
					compiler.MakeInstruction(compiler.OpReturnValue)))))

	expect(t, `func() { 5 + 10 }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(5),
				intObject(10),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpConstant, 1),
					compiler.MakeInstruction(compiler.OpAdd),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpReturn)))))

	expect(t, `func() { 1; 2 }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpConstant, 1),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpReturn)))))

	expect(t, `func() { 1; return 2 }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpConstant, 1),
					compiler.MakeInstruction(compiler.OpReturnValue)))))

	expect(t, `func() { if(true) { return 1 } else { return 2 } }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpTrue),           // 0000
					compiler.MakeInstruction(compiler.OpJumpFalsy, 11),  // 0001
					compiler.MakeInstruction(compiler.OpConstant, 0),    // 0004
					compiler.MakeInstruction(compiler.OpReturnValue),    // 0007
					compiler.MakeInstruction(compiler.OpJump, 15),       // 0008
					compiler.MakeInstruction(compiler.OpConstant, 1),    // 0011
					compiler.MakeInstruction(compiler.OpReturnValue))))) // 0014

	expect(t, `func() { 1; if(true) { 2 } else { 3 }; 4 }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 4),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(4),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),   // 0000
					compiler.MakeInstruction(compiler.OpPop),           // 0003
					compiler.MakeInstruction(compiler.OpTrue),          // 0004
					compiler.MakeInstruction(compiler.OpJumpFalsy, 15), // 0005
					compiler.MakeInstruction(compiler.OpConstant, 1),   // 0008
					compiler.MakeInstruction(compiler.OpPop),           // 0011
					compiler.MakeInstruction(compiler.OpJump, 19),      // 0012
					compiler.MakeInstruction(compiler.OpConstant, 2),   // 0015
					compiler.MakeInstruction(compiler.OpPop),           // 0018
					compiler.MakeInstruction(compiler.OpConstant, 3),   // 0019
					compiler.MakeInstruction(compiler.OpPop),           // 0022
					compiler.MakeInstruction(compiler.OpReturn)))))     // 0023

	expect(t, `func() { }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				compiledFunction(0, 0, compiler.MakeInstruction(compiler.OpReturn)))))

	expect(t, `func() { 24 }()`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpCall, 0),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(24),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpReturn)))))

	expect(t, `func() { return 24 }()`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpCall, 0),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(24),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpReturnValue)))))

	expect(t, `noArg := func() { 24 }; noArg();`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpCall, 0),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(24),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpReturn)))))

	expect(t, `noArg := func() { return 24 }; noArg();`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpCall, 0),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(24),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpReturnValue)))))

	expect(t, `n := 55; func() { n };`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(55),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpGetGlobal, 0),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpReturn)))))

	expect(t, `func() { n := 55; return n }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(55),
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpDefineLocal, 0),
					compiler.MakeInstruction(compiler.OpGetLocal, 0),
					compiler.MakeInstruction(compiler.OpReturnValue)))))

	expect(t, `func() { a := 55; b := 77; return a + b }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(55),
				intObject(77),
				compiledFunction(2, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpDefineLocal, 0),
					compiler.MakeInstruction(compiler.OpConstant, 1),
					compiler.MakeInstruction(compiler.OpDefineLocal, 1),
					compiler.MakeInstruction(compiler.OpGetLocal, 0),
					compiler.MakeInstruction(compiler.OpGetLocal, 1),
					compiler.MakeInstruction(compiler.OpAdd),
					compiler.MakeInstruction(compiler.OpReturnValue)))))

	expect(t, `f1 := func(a) { return a }; f1(24);`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpCall, 1),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				compiledFunction(1, 1,
					compiler.MakeInstruction(compiler.OpGetLocal, 0),
					compiler.MakeInstruction(compiler.OpReturnValue)),
				intObject(24))))

	expect(t, `f1 := func(a, b, c) { a; b; return c; }; f1(24, 25, 26);`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpCall, 3),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				compiledFunction(3, 3,
					compiler.MakeInstruction(compiler.OpGetLocal, 0),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpGetLocal, 1),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpGetLocal, 2),
					compiler.MakeInstruction(compiler.OpReturnValue)),
				intObject(24),
				intObject(25),
				intObject(26))))

	expect(t, `func() { n := 55; n = 23; return n }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(55),
				intObject(23),
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpDefineLocal, 0),
					compiler.MakeInstruction(compiler.OpConstant, 1),
					compiler.MakeInstruction(compiler.OpSetLocal, 0),
					compiler.MakeInstruction(compiler.OpGetLocal, 0),
					compiler.MakeInstruction(compiler.OpReturnValue)))))
	expect(t, `len([]);`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpGetBuiltin, 4),
				compiler.MakeInstruction(compiler.OpArray, 0),
				compiler.MakeInstruction(compiler.OpCall, 1),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray()))

	expect(t, `func() { return len([]) }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpGetBuiltin, 4),
					compiler.MakeInstruction(compiler.OpArray, 0),
					compiler.MakeInstruction(compiler.OpCall, 1),
					compiler.MakeInstruction(compiler.OpReturnValue)))))

	expect(t, `func(a) { func(b) { return a + b } }`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				compiledFunction(1, 1,
					compiler.MakeInstruction(compiler.OpGetFree, 0),
					compiler.MakeInstruction(compiler.OpGetLocal, 0),
					compiler.MakeInstruction(compiler.OpAdd),
					compiler.MakeInstruction(compiler.OpReturnValue)),
				compiledFunction(1, 1,
					compiler.MakeInstruction(compiler.OpGetLocalPtr, 0),
					compiler.MakeInstruction(compiler.OpClosure, 0, 1),
					compiler.MakeInstruction(compiler.OpPop),
					compiler.MakeInstruction(compiler.OpReturn)))))

	expect(t, `
func(a) {
	return func(b) {
		return func(c) {
			return a + b + c
		}
	}
}`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				compiledFunction(1, 1,
					compiler.MakeInstruction(compiler.OpGetFree, 0),
					compiler.MakeInstruction(compiler.OpGetFree, 1),
					compiler.MakeInstruction(compiler.OpAdd),
					compiler.MakeInstruction(compiler.OpGetLocal, 0),
					compiler.MakeInstruction(compiler.OpAdd),
					compiler.MakeInstruction(compiler.OpReturnValue)),
				compiledFunction(1, 1,
					compiler.MakeInstruction(compiler.OpGetFreePtr, 0),
					compiler.MakeInstruction(compiler.OpGetLocalPtr, 0),
					compiler.MakeInstruction(compiler.OpClosure, 0, 2),
					compiler.MakeInstruction(compiler.OpReturnValue)),
				compiledFunction(1, 1,
					compiler.MakeInstruction(compiler.OpGetLocalPtr, 0),
					compiler.MakeInstruction(compiler.OpClosure, 1, 1),
					compiler.MakeInstruction(compiler.OpReturnValue)))))

	expect(t, `
g := 55;

func() {
	a := 66;

	return func() {
		b := 77;

		return func() {
			c := 88;

			return g + a + b + c;
		}
	}
}`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 6),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(55),
				intObject(66),
				intObject(77),
				intObject(88),
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 3),
					compiler.MakeInstruction(compiler.OpDefineLocal, 0),
					compiler.MakeInstruction(compiler.OpGetGlobal, 0),
					compiler.MakeInstruction(compiler.OpGetFree, 0),
					compiler.MakeInstruction(compiler.OpAdd),
					compiler.MakeInstruction(compiler.OpGetFree, 1),
					compiler.MakeInstruction(compiler.OpAdd),
					compiler.MakeInstruction(compiler.OpGetLocal, 0),
					compiler.MakeInstruction(compiler.OpAdd),
					compiler.MakeInstruction(compiler.OpReturnValue)),
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 2),
					compiler.MakeInstruction(compiler.OpDefineLocal, 0),
					compiler.MakeInstruction(compiler.OpGetFreePtr, 0),
					compiler.MakeInstruction(compiler.OpGetLocalPtr, 0),
					compiler.MakeInstruction(compiler.OpClosure, 4, 2),
					compiler.MakeInstruction(compiler.OpReturnValue)),
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 1),
					compiler.MakeInstruction(compiler.OpDefineLocal, 0),
					compiler.MakeInstruction(compiler.OpGetLocalPtr, 0),
					compiler.MakeInstruction(compiler.OpClosure, 5, 1),
					compiler.MakeInstruction(compiler.OpReturnValue)))))

	expect(t, `for i:=0; i<10; i++ {}`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGreaterThan),
				compiler.MakeInstruction(compiler.OpJumpFalsy, 29),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpJump, 6)),
			objectsArray(
				intObject(0),
				intObject(10),
				intObject(1))))

	expect(t, `m := {}; for k, v in m {}`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpMap, 0),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpIteratorInit),
				compiler.MakeInstruction(compiler.OpSetGlobal, 1),
				compiler.MakeInstruction(compiler.OpGetGlobal, 1),
				compiler.MakeInstruction(compiler.OpIteratorNext),
				compiler.MakeInstruction(compiler.OpJumpFalsy, 37),
				compiler.MakeInstruction(compiler.OpGetGlobal, 1),
				compiler.MakeInstruction(compiler.OpIteratorKey),
				compiler.MakeInstruction(compiler.OpSetGlobal, 2),
				compiler.MakeInstruction(compiler.OpGetGlobal, 1),
				compiler.MakeInstruction(compiler.OpIteratorValue),
				compiler.MakeInstruction(compiler.OpSetGlobal, 3),
				compiler.MakeInstruction(compiler.OpJump, 13)),
			objectsArray()))

	expect(t, `a := 0; a == 0 && a != 1 || a < 1`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpSetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpEqual),
				compiler.MakeInstruction(compiler.OpAndJump, 23),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpNotEqual),
				compiler.MakeInstruction(compiler.OpOrJump, 33),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGreaterThan),
				compiler.MakeInstruction(compiler.OpPop)),
			objectsArray(
				intObject(0),
				intObject(0),
				intObject(1),
				intObject(1))))

	expectError(t, `import("user1")`, "no such file or directory") // unknown module name

	expectError(t, `
r["x"] = {
    @a:1,
    @b:1,
    @c:1,
    @d:1,
    @e:1,
    @f:1,
    @g:1,
    @h:1,
    @i:1,
    @j:1,
    @k:1
}
`, "Parse Error: illegal character U+0040 '@'\n\tat test:3:5 (and 10 more errors)") // too many errors
}

func concat(instructions ...[]byte) []byte {
	var concat []byte
	for _, i := range instructions {
		concat = append(concat, i...)
	}

	return concat
}

func bytecode(instructions []byte, constants []objects.Object) *compiler.Bytecode {
	return &compiler.Bytecode{
		FileSet:      source.NewFileSet(),
		MainFunction: &objects.CompiledFunction{Instructions: instructions},
		Constants:    constants,
	}
}

func expect(t *testing.T, input string, expected *compiler.Bytecode) (ok bool) {
	actual, trace, err := traceCompile(input, nil)

	defer func() {
		if !ok {
			for _, tr := range trace {
				t.Log(tr)
			}
		}
	}()

	if !assert.NoError(t, err) {
		return
	}

	ok = equalBytecode(t, expected, actual)

	return
}

func expectError(t *testing.T, input, expected string) (ok bool) {
	_, trace, err := traceCompile(input, nil)

	defer func() {
		if !ok {
			for _, tr := range trace {
				t.Log(tr)
			}
		}
	}()

	if !assert.Error(t, err) {
		return
	}

	if !assert.True(t, strings.Contains(err.Error(), expected), "expected error string: %s, got: %s", expected, err.Error()) {
		return
	}

	ok = true

	return
}

func equalBytecode(t *testing.T, expected, actual *compiler.Bytecode) bool {
	return assert.Equal(t, expected.MainFunction, actual.MainFunction) &&
		equalConstants(t, expected.Constants, actual.Constants)
}

func equalConstants(t *testing.T, expected, actual []objects.Object) bool {
	if !assert.Equal(t, len(expected), len(actual)) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if !assert.Equal(t, expected[i], actual[i]) {
			return false
		}
	}

	return true
}

type tracer struct {
	Out []string
}

func (o *tracer) Write(p []byte) (n int, err error) {
	o.Out = append(o.Out, string(p))
	return len(p), nil
}

func traceCompile(input string, symbols map[string]objects.Object) (res *compiler.Bytecode, trace []string, err error) {
	fileSet := source.NewFileSet()
	file := fileSet.AddFile("test", -1, len(input))

	p := parser.NewParser(file, []byte(input), nil)

	symTable := compiler.NewSymbolTable()
	for name := range symbols {
		symTable.Define(name)
	}
	for idx, fn := range objects.Builtins {
		symTable.DefineBuiltin(idx, fn.Name)
	}

	tr := &tracer{}
	c := compiler.NewCompiler(file, symTable, nil, nil, tr)
	parsed, err := p.ParseFile()
	if err != nil {
		return
	}

	err = c.Compile(parsed)
	{
		trace = append(trace, fmt.Sprintf("Compiler Trace:\n%s", strings.Join(tr.Out, "")))

		bytecode := c.Bytecode()
		trace = append(trace, fmt.Sprintf("Compiled Constants:\n%s", strings.Join(bytecode.FormatConstants(), "\n")))
		trace = append(trace, fmt.Sprintf("Compiled Instructions:\n%s\n", strings.Join(bytecode.FormatInstructions(), "\n")))
	}
	if err != nil {
		return
	}

	res = c.Bytecode()

	return
}

func objectsArray(o ...objects.Object) []objects.Object {
	return o
}

func intObject(v int64) *objects.Int {
	return &objects.Int{Value: v}
}

func stringObject(v string) *objects.String {
	return &objects.String{Value: v}
}

func compiledFunction(numLocals, numParams int, insts ...[]byte) *objects.CompiledFunction {
	return &objects.CompiledFunction{Instructions: concat(insts...), NumLocals: numLocals, NumParameters: numParams}
}
