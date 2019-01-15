package compiler

// Definition represents an Opcode name and
// the number of operands.
type Definition struct {
	Name     string
	Operands []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:         {Name: "CONST", Operands: []int{2}},
	OpPop:              {Name: "POP", Operands: []int{}},
	OpTrue:             {Name: "TRUE", Operands: []int{}},
	OpFalse:            {Name: "FALSE", Operands: []int{}},
	OpAdd:              {Name: "ADD", Operands: []int{}},
	OpSub:              {Name: "SUB", Operands: []int{}},
	OpMul:              {Name: "MUL", Operands: []int{}},
	OpDiv:              {Name: "DIV", Operands: []int{}},
	OpRem:              {Name: "REM", Operands: []int{}},
	OpBAnd:             {Name: "AND", Operands: []int{}},
	OpBOr:              {Name: "OR", Operands: []int{}},
	OpBXor:             {Name: "XOR", Operands: []int{}},
	OpBAndNot:          {Name: "ANDN", Operands: []int{}},
	OpBShiftLeft:       {Name: "SHL", Operands: []int{}},
	OpBShiftRight:      {Name: "SHR", Operands: []int{}},
	OpBComplement:      {Name: "NEG", Operands: []int{}},
	OpEqual:            {Name: "EQL", Operands: []int{}},
	OpNotEqual:         {Name: "NEQ", Operands: []int{}},
	OpGreaterThan:      {Name: "GTR", Operands: []int{}},
	OpGreaterThanEqual: {Name: "GEQ", Operands: []int{}},
	OpMinus:            {Name: "NEG", Operands: []int{}},
	OpLNot:             {Name: "NOT", Operands: []int{}},
	OpJumpFalsy:        {Name: "JMPF", Operands: []int{2}},
	OpAndJump:          {Name: "ANDJMP", Operands: []int{2}},
	OpOrJump:           {Name: "ORJMP", Operands: []int{2}},
	OpJump:             {Name: "JMP", Operands: []int{2}},
	OpNull:             {Name: "NULL", Operands: []int{}},
	OpGetGlobal:        {Name: "GETG", Operands: []int{2}},
	OpSetGlobal:        {Name: "SETG", Operands: []int{2}},
	OpSetSelGlobal:     {Name: "SETSG", Operands: []int{2, 1}},
	OpArray:            {Name: "ARR", Operands: []int{2}},
	OpMap:              {Name: "MAP", Operands: []int{2}},
	OpIndex:            {Name: "INDEX", Operands: []int{}},
	OpSliceIndex:       {Name: "SLICE", Operands: []int{}},
	OpCall:             {Name: "CALL", Operands: []int{1}},
	OpReturn:           {Name: "RET", Operands: []int{}},
	OpReturnValue:      {Name: "RETVAL", Operands: []int{1}},
	OpGetLocal:         {Name: "GETL", Operands: []int{1}},
	OpSetLocal:         {Name: "SETL", Operands: []int{1}},
	OpDefineLocal:      {Name: "DEFL", Operands: []int{1}},
	OpSetSelLocal:      {Name: "SETSL", Operands: []int{1, 1}},
	OpGetBuiltin:       {Name: "BUILTIN", Operands: []int{1}},
	OpClosure:          {Name: "CLOSURE", Operands: []int{2, 1}},
	OpGetFree:          {Name: "GETF", Operands: []int{1}},
	OpSetFree:          {Name: "SETF", Operands: []int{1}},
	OpSetSelFree:       {Name: "SETSF", Operands: []int{1, 1}},
	OpIteratorInit:     {Name: "ITER", Operands: []int{}},
	OpIteratorNext:     {Name: "ITNXT", Operands: []int{}},
	OpIteratorKey:      {Name: "ITKEY", Operands: []int{}},
	OpIteratorValue:    {Name: "ITVAL", Operands: []int{}},
}

// Lookup returns a Definition of a given opcode.
func Lookup(opcode Opcode) (def *Definition, ok bool) {
	def, ok = definitions[opcode]

	return
}
