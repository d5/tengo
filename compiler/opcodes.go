package compiler

// Opcode represents a single byte operation code.
type Opcode = byte

// List of opcodes
const (
	OpConstant         Opcode = iota // Load constant
	OpAdd                            // Add
	OpSub                            // Sub
	OpMul                            // Multiply
	OpDiv                            // Divide
	OpRem                            // Remainder
	OpBAnd                           // bitwise AND
	OpBOr                            // bitwise OR
	OpBXor                           // bitwise XOR
	OpBShiftLeft                     // bitwise shift left
	OpBShiftRight                    // bitwise shift right
	OpBAndNot                        // bitwise AND NOT
	OpBComplement                    // bitwise complement
	OpPop                            // Pop
	OpTrue                           // Push true
	OpFalse                          // Push false
	OpEqual                          // Equal ==
	OpNotEqual                       // Not equal !=
	OpGreaterThan                    // Greater than >=
	OpGreaterThanEqual               // Greater than or equal to >=
	OpMinus                          // Minus -
	OpLNot                           // Logical not !
	OpJumpFalsy                      // Jump if falsy
	OpAndJump                        // Logical AND jump
	OpOrJump                         // Logical OR jump
	OpJump                           // Jump
	OpNull                           // Push null
	OpArray                          // Array object
	OpMap                            // Map object
	OpError                          // Error object
	OpImmutable                      // Immutable object
	OpIndex                          // Index operation
	OpSliceIndex                     // Slice operation
	OpCall                           // Call function
	OpReturn                         // Return
	OpReturnValue                    // Return value
	OpGetGlobal                      // Get global variable
	OpSetGlobal                      // Set global variable
	OpSetSelGlobal                   // Set global variable using selectors
	OpGetLocal                       // Get local variable
	OpSetLocal                       // Set local variable
	OpDefineLocal                    // Define local variable
	OpSetSelLocal                    // Set local variable using selectors
	OpGetFreePtr                     // Get free variable pointer object
	OpGetFree                        // Get free variables
	OpSetFree                        // Set free variables
	OpGetLocalPtr                    // Get local variable as a pointer
	OpSetLocalPtr                    // Set local variable as a pointer
	OpSetSelFree                     // Set free variables using selectors
	OpGetBuiltin                     // Get builtin function
	OpGetBuiltinModule               // Get builtin module
	OpClosure                        // Push closure
	OpIteratorInit                   // Iterator init
	OpIteratorNext                   // Iterator next
	OpIteratorKey                    // Iterator key
	OpIteratorValue                  // Iterator value
)

// OpcodeNames is opcode names.
var OpcodeNames = [...]string{
	OpConstant:         "CONST",
	OpPop:              "POP",
	OpTrue:             "TRUE",
	OpFalse:            "FALSE",
	OpAdd:              "ADD",
	OpSub:              "SUB",
	OpMul:              "MUL",
	OpDiv:              "DIV",
	OpRem:              "REM",
	OpBAnd:             "AND",
	OpBOr:              "OR",
	OpBXor:             "XOR",
	OpBAndNot:          "ANDN",
	OpBShiftLeft:       "SHL",
	OpBShiftRight:      "SHR",
	OpBComplement:      "NEG",
	OpEqual:            "EQL",
	OpNotEqual:         "NEQ",
	OpGreaterThan:      "GTR",
	OpGreaterThanEqual: "GEQ",
	OpMinus:            "NEG",
	OpLNot:             "NOT",
	OpJumpFalsy:        "JMPF",
	OpAndJump:          "ANDJMP",
	OpOrJump:           "ORJMP",
	OpJump:             "JMP",
	OpNull:             "NULL",
	OpGetGlobal:        "GETG",
	OpSetGlobal:        "SETG",
	OpSetSelGlobal:     "SETSG",
	OpArray:            "ARR",
	OpMap:              "MAP",
	OpError:            "ERROR",
	OpImmutable:        "IMMUT",
	OpIndex:            "INDEX",
	OpSliceIndex:       "SLICE",
	OpCall:             "CALL",
	OpReturn:           "RET",
	OpReturnValue:      "RETVAL",
	OpGetLocal:         "GETL",
	OpSetLocal:         "SETL",
	OpDefineLocal:      "DEFL",
	OpSetSelLocal:      "SETSL",
	OpGetBuiltin:       "BUILTIN",
	OpGetBuiltinModule: "BLTMOD",
	OpClosure:          "CLOSURE",
	OpGetFreePtr:       "GETFP",
	OpGetFree:          "GETF",
	OpSetFree:          "SETF",
	OpGetLocalPtr:      "GETLP",
	OpSetLocalPtr:      "SETLP",
	OpSetSelFree:       "SETSF",
	OpIteratorInit:     "ITER",
	OpIteratorNext:     "ITNXT",
	OpIteratorKey:      "ITKEY",
	OpIteratorValue:    "ITVAL",
}

// OpcodeOperands is the number of operands.
var OpcodeOperands = [...][]int{
	OpConstant:         {2},
	OpPop:              {},
	OpTrue:             {},
	OpFalse:            {},
	OpAdd:              {},
	OpSub:              {},
	OpMul:              {},
	OpDiv:              {},
	OpRem:              {},
	OpBAnd:             {},
	OpBOr:              {},
	OpBXor:             {},
	OpBAndNot:          {},
	OpBShiftLeft:       {},
	OpBShiftRight:      {},
	OpBComplement:      {},
	OpEqual:            {},
	OpNotEqual:         {},
	OpGreaterThan:      {},
	OpGreaterThanEqual: {},
	OpMinus:            {},
	OpLNot:             {},
	OpJumpFalsy:        {2},
	OpAndJump:          {2},
	OpOrJump:           {2},
	OpJump:             {2},
	OpNull:             {},
	OpGetGlobal:        {2},
	OpSetGlobal:        {2},
	OpSetSelGlobal:     {2, 1},
	OpArray:            {2},
	OpMap:              {2},
	OpError:            {},
	OpImmutable:        {},
	OpIndex:            {},
	OpSliceIndex:       {},
	OpCall:             {1},
	OpReturn:           {},
	OpReturnValue:      {},
	OpGetLocal:         {1},
	OpSetLocal:         {1},
	OpDefineLocal:      {1},
	OpSetSelLocal:      {1, 1},
	OpGetBuiltin:       {1},
	OpGetBuiltinModule: {},
	OpClosure:          {2, 1},
	OpGetFreePtr:       {1},
	OpGetFree:          {1},
	OpSetFree:          {1},
	OpGetLocalPtr:      {1},
	OpSetLocalPtr:      {1},
	OpSetSelFree:       {1, 1},
	OpIteratorInit:     {},
	OpIteratorNext:     {},
	OpIteratorKey:      {},
	OpIteratorValue:    {},
}

// ReadOperands reads operands from the bytecode.
func ReadOperands(numOperands []int, ins []byte) (operands []int, offset int) {
	for _, width := range numOperands {
		switch width {
		case 1:
			operands = append(operands, int(ins[offset]))
		case 2:
			operands = append(operands, int(ins[offset+1])|int(ins[offset])<<8)
		}

		offset += width
	}

	return
}
