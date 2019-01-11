package parser

import (
	"fmt"
	"sort"

	"github.com/d5/tengo/source"
)

type ErrorList []*Error

func (p *ErrorList) Add(pos source.FilePos, msg string) {
	*p = append(*p, &Error{pos, msg})
}

func (p *ErrorList) Reset() {
	*p = (*p)[0:0]
}

func (p ErrorList) Len() int {
	return len(p)
}

func (p ErrorList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ErrorList) Less(i, j int) bool {
	e := &p[i].Pos
	f := &p[j].Pos

	if e.Filename != f.Filename {
		return e.Filename < f.Filename
	}

	if e.Line != f.Line {
		return e.Line < f.Line
	}

	if e.Column != f.Column {
		return e.Column < f.Column
	}

	return p[i].Msg < p[j].Msg
}

func (p ErrorList) Sort() {
	sort.Sort(p)
}

func (p *ErrorList) RemoveMultiples() {
	sort.Sort(p)

	var last source.FilePos // initial last.Line is != any legal error line

	i := 0
	for _, e := range *p {
		if e.Pos.Filename != last.Filename || e.Pos.Line != last.Line {
			last = e.Pos
			(*p)[i] = e
			i++
		}
	}

	*p = (*p)[0:i]
}

func (p ErrorList) Error() string {
	switch len(p) {
	case 0:
		return "no errors"
	case 1:
		return p[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", p[0], len(p)-1)
}

func (p ErrorList) Err() error {
	if len(p) == 0 {
		return nil
	}

	return p
}
