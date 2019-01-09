package parser

import "github.com/d5/tengo/scanner"

type Error struct {
	Pos scanner.FilePos
	Msg string
}

func (e Error) Error() string {
	if e.Pos.Filename != "" || e.Pos.IsValid() {
		return e.Pos.String() + ": " + e.Msg
	}

	return e.Msg
}
