package parser

import (
	"reflect"
	"regexp"
)

func AppendManyNew[T any](v []T, as ...any) []T {
	l := len(v)
	for _, a := range as {
		switch a.(type) {
		case T:
			l += 1
		case []T:
			l += len(a.([]T))
		default:
			panic("type error:" + reflect.TypeOf(a).String())
		}
	}
	n := make([]T, l)
	n = n[:0]
	n = append(n, v...)
	for _, a := range as {
		switch a.(type) {
		case T:
			n = append(n, a.(T))
		case []T:
			n = append(n, a.([]T)...)
		default:
			panic("type error:" + reflect.TypeOf(a).String())
		}
	}
	return n
}

func regexFindAll(s0, frx s) []s {
	return regexp.MustCompile(frx).FindAllString(s0, -1)
}
