package compiler

// Opcode represents a single byte operation code.
type Opcode byte

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
	OpArray                          // Array literal
	OpMap                            // Map literal
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
	OpGetFree                        // Get free variables
	OpSetFree                        // Set free variables
	OpSetSelFree                     // Set free variables using selectors
	OpGetBuiltin                     // Get builtin function
	OpClosure                        // Push closure
	OpIteratorInit                   // Iterator init
	OpIteratorNext                   // Iterator next
	OpIteratorKey                    // Iterator key
	OpIteratorValue                  // Iterator value
)
