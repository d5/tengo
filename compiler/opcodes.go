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
	OpArray
	OpMap
	OpIndex
	OpSliceIndex
	OpCall
	OpReturn
	OpReturnValue
	OpGetGlobal
	OpSetGlobal
	OpSetSelGlobal
	OpGetLocal
	OpSetLocal
	OpDefineLocal
	OpSetSelLocal
	OpGetFree
	OpSetFree
	OpSetSelFree
	OpGetBuiltin
	OpClosure
	OpIteratorInit
	OpIteratorNext
	OpIteratorKey
	OpIteratorValue
)
