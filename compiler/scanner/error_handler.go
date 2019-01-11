package scanner

import "github.com/d5/tengo/compiler/source"

type ErrorHandler func(pos source.FilePos, msg string)
