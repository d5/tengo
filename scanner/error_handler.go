package scanner

import "github.com/d5/tengo/source"

type ErrorHandler func(pos source.FilePos, msg string)
