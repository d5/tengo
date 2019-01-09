package compiler

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpRem
	OpBAnd
	OpBOr
	OpBXor
	OpBShiftLeft
	OpBShiftRight
	OpBAndNot
	OpBComplement
	OpPop
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterThanEqual
	OpMinus
	OpLNot
	OpJumpFalsy
	OpAndJump
	OpOrJump
	OpJump
	OpNull
	OpGetGlobal
	OpSetGlobal
	OpSetSelGlobal
	OpArray
	OpMap
	OpIndex
	OpSliceIndex
	OpCall
	OpReturn
	OpReturnValue
	OpGetLocal
	OpSetLocal
	OpSetSelLocal
	OpGetBuiltin
	OpClosure
	OpGetFree
	OpSetFree
	OpSetSelFree
	OpIteratorInit
	OpIteratorNext
	OpIteratorKey
	OpIteratorValue
)
