package assert

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

// NoError asserts err is not an error.
func NoError(t *testing.T, err error, msg ...interface{}) bool {
	if err == nil {
		return true
	}

	return failExpectedActual(t, "no error", err, msg...)
}

// Error asserts err is an error.
func Error(t *testing.T, err error, msg ...interface{}) bool {
	if err != nil {
		return true
	}

	return failExpectedActual(t, "error", err, msg...)
}

// Nil asserts v is nil.
func Nil(t *testing.T, v interface{}, msg ...interface{}) bool {
	if v == nil {
		return true
	}

	return failExpectedActual(t, "nil", v, msg...)
}

// True asserts v is true.
func True(t *testing.T, v bool, msg ...interface{}) bool {
	if v {
		return true
	}

	return failExpectedActual(t, "true", v, msg...)
}

// False asserts vis false.
func False(t *testing.T, v bool, msg ...interface{}) bool {
	if !v {
		return true
	}

	return failExpectedActual(t, "false", v, msg...)
}

// NotNil asserts v is not nil.
func NotNil(t *testing.T, v interface{}, msg ...interface{}) bool {
	if v != nil {
		return true
	}

	return failExpectedActual(t, "not nil", v, msg...)
}

// IsType asserts expected and actual are of the same type.
func IsType(t *testing.T, expected, actual interface{}, msg ...interface{}) bool {
	if reflect.TypeOf(expected) == reflect.TypeOf(actual) {
		return true
	}

	return failExpectedActual(t, reflect.TypeOf(expected), reflect.TypeOf(actual), msg...)
}

// Equal asserts expected and actual are equal.
func Equal(t *testing.T, expected, actual interface{}, msg ...interface{}) bool {
	if expected == nil {
		return Nil(t, actual, "expected nil, but got not nil")
	}
	if !NotNil(t, actual, "expected not nil, but got nil") {
		return false
	}
	if !IsType(t, expected, actual, msg...) {
		return false
	}

	switch expected := expected.(type) {
	case int:
		if expected != actual.(int) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case int64:
		if expected != actual.(int64) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case float64:
		if expected != actual.(float64) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case string:
		if expected != actual.(string) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case []byte:
		if bytes.Compare(expected, actual.([]byte)) != 0 {
			return failExpectedActual(t, string(expected), string(actual.([]byte)), msg...)
		}
	case []int:
		if !equalIntSlice(expected, actual.([]int)) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case bool:
		if expected != actual.(bool) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case rune:
		if expected != actual.(rune) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case compiler.Symbol:
		if !equalSymbol(expected, actual.(compiler.Symbol)) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case source.Pos:
		if expected != actual.(source.Pos) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case token.Token:
		if expected != actual.(token.Token) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case []objects.Object:
		return equalObjectSlice(t, expected, actual.([]objects.Object), msg...)
	case *objects.Int:
		return Equal(t, expected.Value, actual.(*objects.Int).Value, msg...)
	case *objects.Float:
		return Equal(t, expected.Value, actual.(*objects.Float).Value, msg...)
	case *objects.String:
		return Equal(t, expected.Value, actual.(*objects.String).Value, msg...)
	case *objects.Char:
		return Equal(t, expected.Value, actual.(*objects.Char).Value, msg...)
	case *objects.Bool:
		if expected != actual {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case *objects.ReturnValue:
		return Equal(t, expected.Value, actual.(objects.ReturnValue).Value, msg...)
	case *objects.Array:
		return equalObjectSlice(t, expected.Value, actual.(*objects.Array).Value, msg...)
	case *objects.ImmutableArray:
		return equalObjectSlice(t, expected.Value, actual.(*objects.ImmutableArray).Value, msg...)
	case *objects.Bytes:
		if bytes.Compare(expected.Value, actual.(*objects.Bytes).Value) != 0 {
			return failExpectedActual(t, string(expected.Value), string(actual.(*objects.Bytes).Value), msg...)
		}
	case *objects.Map:
		return equalObjectMap(t, expected.Value, actual.(*objects.Map).Value, msg...)
	case *objects.ImmutableMap:
		return equalObjectMap(t, expected.Value, actual.(*objects.ImmutableMap).Value, msg...)
	case *objects.CompiledFunction:
		return equalCompiledFunction(t, expected, actual.(*objects.CompiledFunction), msg...)
	case *objects.Closure:
		return equalClosure(t, expected, actual.(*objects.Closure), msg...)
	case *objects.Undefined:
		if expected != actual {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case *objects.Error:
		return Equal(t, expected.Value, actual.(*objects.Error).Value, msg...)
	case objects.Object:
		if !expected.Equals(actual.(objects.Object)) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	case error:
		if expected != actual.(error) {
			return failExpectedActual(t, expected, actual, msg...)
		}
	default:
		panic(fmt.Errorf("type not implemented: %T", expected))
	}

	return true
}

// Fail marks the function as having failed but continues execution.
func Fail(t *testing.T, msg ...interface{}) bool {
	t.Logf("\nError trace:\n\t%s\n%s", strings.Join(errorTrace(), "\n\t"), message(msg...))

	t.Fail()

	return false
}

func failExpectedActual(t *testing.T, expected, actual interface{}, msg ...interface{}) bool {
	var addMsg string
	if len(msg) > 0 {
		addMsg = "\nMessage:  " + message(msg...)
	}

	t.Logf("\nError trace:\n\t%s\nExpected: %v\nActual:   %v%s",
		strings.Join(errorTrace(), "\n\t"),
		expected, actual,
		addMsg)

	t.Fail()

	return false
}

func message(formatArgs ...interface{}) string {
	var format string
	var args []interface{}
	if len(formatArgs) > 0 {
		format = formatArgs[0].(string)
	}
	if len(formatArgs) > 1 {
		args = formatArgs[1:]
	}

	return fmt.Sprintf(format, args...)
}

func equalIntSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func equalSymbol(a, b compiler.Symbol) bool {
	return a.Name == b.Name &&
		a.Index == b.Index &&
		a.Scope == b.Scope
}

func equalObjectSlice(t *testing.T, expected, actual []objects.Object, msg ...interface{}) bool {
	// TODO: this test does not differentiate nil vs empty slice

	if !Equal(t, len(expected), len(actual), msg...) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if !Equal(t, expected[i], actual[i], msg...) {
			return false
		}
	}

	return true
}

func equalObjectMap(t *testing.T, expected, actual map[string]objects.Object, msg ...interface{}) bool {
	if !Equal(t, len(expected), len(actual), msg...) {
		return false
	}

	for key, expectedVal := range expected {
		actualVal := actual[key]

		if !Equal(t, expectedVal, actualVal, msg...) {
			return false
		}
	}

	return true
}

func equalCompiledFunction(t *testing.T, expected, actual objects.Object, msg ...interface{}) bool {
	expectedT := expected.(*objects.CompiledFunction)
	actualT := actual.(*objects.CompiledFunction)

	return Equal(t, expectedT.Instructions, actualT.Instructions, msg...)
}

func equalClosure(t *testing.T, expected, actual objects.Object, msg ...interface{}) bool {
	expectedT := expected.(*objects.Closure)
	actualT := actual.(*objects.Closure)

	if !Equal(t, expectedT.Fn, actualT.Fn, msg...) {
		return false
	}

	if !Equal(t, len(expectedT.Free), len(actualT.Free), msg...) {
		return false
	}

	for i := 0; i < len(expectedT.Free); i++ {
		if !Equal(t, *expectedT.Free[i], *actualT.Free[i], msg...) {
			return false
		}
	}

	return true
}
