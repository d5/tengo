package tengo

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrStackOverflow is a stack overflow error.
	ErrStackOverflow = errors.New("stack overflow")

	// ErrObjectAllocLimit is an objects allocation limit error.
	ErrObjectAllocLimit = errors.New("object allocation limit exceeded")

	// ErrIndexOutOfBounds is an error where a given index is out of the
	// bounds.
	ErrIndexOutOfBounds = errors.New("index out of bounds")

	// ErrInvalidIndexType represents an invalid index type.
	ErrInvalidIndexType = errors.New("invalid index type")

	// ErrInvalidIndexValueType represents an invalid index value type.
	ErrInvalidIndexValueType = errors.New("invalid index value type")

	// ErrInvalidIndexOnError represents an invalid index on error.
	ErrInvalidIndexOnError = errors.New("invalid index on error")

	// ErrInvalidOperator represents an error for invalid operator usage.
	ErrInvalidOperator = errors.New("invalid operator")

	// ErrBytesLimit represents an error where the size of bytes value exceeds
	// the limit.
	ErrBytesLimit = errors.New("exceeding bytes size limit")

	// ErrStringLimit represents an error where the size of string value
	// exceeds the limit.
	ErrStringLimit = errors.New("exceeding string size limit")

	// ErrNotIndexable is an error where an Object is not indexable.
	ErrNotIndexable = errors.New("not indexable")

	// ErrNotIndexAssignable is an error where an Object is not index
	// assignable.
	ErrNotIndexAssignable = errors.New("not index-assignable")

	// ErrNotImplemented is an error where an Object has not implemented a
	// required method.
	ErrNotImplemented = errors.New("not implemented")

	// ErrInvalidRangeStep is an error where the step parameter is less than or equal to 0 when using builtin range function.
	ErrInvalidRangeStep = errors.New("range step must be greater than 0")
)

// ErrInvalidReturnValueCount represents an invalid return value count error.
type ErrInvalidReturnValueCount struct {
	Expected int
	Actual   int
}

func (e ErrInvalidReturnValueCount) Error() string {
	return fmt.Sprintf("invalid return value count, expected %d, actual %d",
		e.Expected, e.Actual)
}

// ErrInvalidArgumentCount represents an invalid argument count error.
type ErrInvalidArgumentCount struct {
	Min    int
	Max    int
	Actual int
}

func (e ErrInvalidArgumentCount) Error() string {
	if e.Max < 0 {
		return fmt.Sprintf("invalid variadic argument count, expected at least %d, actual %d",
			e.Min, e.Actual)
	}
	if e.Min == e.Max {
		return fmt.Sprintf("invalid argument count, expected %d, actual %d",
			e.Min, e.Actual)
	}
	return fmt.Sprintf("invalid argument count, expected between %d and %d, actual %d",
		e.Min, e.Max, e.Actual)
}

// ErrInvalidArgumentType represents an invalid argument type error.
type ErrInvalidArgumentType struct {
	Index    int
	Expected string
	Actual   string
}

// TNs is a shorthand alias for typenames
type TNs = []string

// AnyTN is a typename when any type can be accepted
const AnyTN = "any"

func (e ErrInvalidArgumentType) Error() string {
	return fmt.Sprintf("invalid type for argument at index %d: expected %s, actual %s",
		e.Index, e.Expected, e.Actual)
}

// CheckArgs wraps a callable variadic function with flexible argument type checking logic
// when an argument might have several allowable types.
func CheckArgs(fn CallableFunc, min, max int, expected ...TNs) CallableFunc {
	if min < 0 || (max >= 0 && min > max) {
		panic("invalid min max arg count values")
	}
	if max < 0 && len(expected) < min {
		// variadic func
		panic("invalid expected len, must be at least min when max < 0")
	}
	if max >= 0 && (len(expected) < min || len(expected) > max) {
		panic("invalid expected len, must be between min and max")
	}
	return func(actual ...Object) (Object, error) {
		if (max >= 0 && len(actual) > max) ||
			(len(actual) < min) {
			return nil, ErrInvalidArgumentCount{
				Min:    min,
				Max:    max,
				Actual: len(actual),
			}
		}
		for i := range actual {
			var expectedTypes []string
			if i > len(expected)-1 {
				// variadic args
				expectedTypes = expected[len(expected)-1]
			} else {
				// required args
				expectedTypes = expected[i]
			}
			if len(expectedTypes) == 0 || (len(expectedTypes) == 1 && expectedTypes[0] == AnyTN) {
				// any type allowed
				continue
			}
			for j := range expectedTypes {
				if expectedTypes[j] == actual[i].TypeName() {
					break
				} else if j+1 == len(expectedTypes) {
					return nil, ErrInvalidArgumentType{
						Index:    i,
						Expected: strings.Join(expectedTypes, "/"),
						Actual:   actual[i].TypeName(),
					}
				}
			}
		}
		return fn(actual...)
	}
}

// CheckOptArgs wraps a callable function with strict argument type checking for optional arguments
func CheckOptArgs(fn CallableFunc, min, max int, expected ...string) CallableFunc {
	flexExpected := make([][]string, 0, len(expected))
	for _, tn := range expected {
		flexExpected = append(flexExpected, []string{tn})
	}
	return CheckArgs(fn, min, max, flexExpected...)
}

// CheckAnyArgs wraps a callable function with argument count checking,
// this is only useful to functions that use tengo.ToSomeType(i int, args ...Object) (type, error)
// functions as they do more complex type conversions with their own type checking internally.
func CheckAnyArgs(fn CallableFunc, minMax ...int) CallableFunc {
	min := 0
	if len(minMax) > 0 {
		min = minMax[0]
	}
	max := min
	if len(minMax) > 1 {
		max = minMax[1]
	}
	var expected [][]string
	if max < 0 {
		expected = make([][]string, min+1)
	} else {
		expected = make([][]string, max)
	}
	return CheckArgs(fn, min, max, expected...)
}

// CheckStrictArgs wraps a callable function with strict argument type checking logic
func CheckStrictArgs(fn CallableFunc, expected ...string) CallableFunc {
	return CheckOptArgs(fn, len(expected), len(expected), expected...)
}
