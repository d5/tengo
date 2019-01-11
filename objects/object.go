package objects

import "github.com/d5/tengo/compiler/token"

type Object interface {
	TypeName() string
	String() string
	BinaryOp(op token.Token, rhs Object) (Object, error)
	IsFalsy() bool
	Equals(Object) bool
	Copy() Object
}
