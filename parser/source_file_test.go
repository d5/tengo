package parser_test

import (
	"fmt"
	"testing"

	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
	"github.com/d5/tengo/v2/token"
)

func TestSourceFilePosString(t *testing.T) {
	for _, tc := range []struct {
		pos      parser.SourceFilePos
		expected string
	}{
		{parser.SourceFilePos{Filename: "f", Line: 2, Column: 3}, "f:2:3"},
		{parser.SourceFilePos{Filename: "f", Line: 2}, "f:2"},
		{parser.SourceFilePos{Line: 2, Column: 3}, "2:3"},
		{parser.SourceFilePos{Line: 2}, "2"},
		{parser.SourceFilePos{Filename: "f", Column: 3}, "f"},
		{parser.SourceFilePos{Column: 3}, "-"},
		{parser.SourceFilePos{}, "-"},
	} {
		require.Equal(t, tc.expected, tc.pos.String())
	}
}

func TestSourceFileSetFile(t *testing.T) {
	fileSet := parser.NewFileSet()
	a := fileSet.AddFile("a", -1, 10)
	b := fileSet.AddFile("b", 20, 5)
	require.Equal(t, int64(15), fileSet.Size())
	require.Equal(t, 26, fileSet.Base)

	for _, tc := range []struct {
		pos      parser.Pos
		expected *parser.SourceFile
	}{
		{parser.NoPos, nil},
		{1, a},
		{11, a},
		{12, nil},
		{19, nil},
		{20, b},
		{25, b},
		{26, nil},
		{5, a},
	} {
		require.True(t, fileSet.File(tc.pos) == tc.expected, "pos %d", tc.pos)
	}
}

func TestSourceFileSetPosition(t *testing.T) {
	fileSet := parser.NewFileSet()
	a := fileSet.AddFile("a", -1, 10)
	a.AddLine(4)
	fileSet.AddFile("b", 20, 5)

	for _, tc := range []struct {
		pos      parser.Pos
		expected string
		offset   int
	}{
		{parser.NoPos, "-", 0},
		{1, "a:1:1", 0},
		{4, "a:1:4", 3},
		{5, "a:2:1", 4},
		{11, "a:2:7", 10},
		{15, "-", 0},
		{22, "b:1:3", 2},
	} {
		pos := fileSet.Position(tc.pos)
		require.Equal(t, tc.expected, pos.String(), "pos %d", tc.pos)
		require.Equal(t, tc.offset, pos.Offset, "pos %d", tc.pos)
	}
}

func TestSourceFile(t *testing.T) {
	fileSet := parser.NewFileSet()
	fileSet.AddFile("first", -1, 3)
	src := "a\nbc\n\nd"
	file := fileSet.AddFile("second", -1, len(src))
	require.True(t, file.Set() == fileSet)
	require.Equal(t, 1, file.LineCount())

	s := parser.NewScanner(file, []byte(src), nil, 0)
	for {
		if tok, _, _ := s.Scan(); tok == token.EOF {
			break
		}
	}
	require.Equal(t, 4, file.LineCount())
	require.Equal(t, []int{0, 2, 5, 6}, file.Lines)

	require.Equal(t, 0, file.Offset(file.LineStart(1)))
	require.Equal(t, 2, file.Offset(file.LineStart(2)))
	require.Equal(t, 6, file.Offset(file.LineStart(4)))
	require.Equal(t, len(src), file.Offset(file.FileSetPos(len(src))))
	require.Equal(t, "second:3:1", file.Position(file.LineStart(3)).String())
	require.False(t, file.Position(parser.NoPos).IsValid())

	file.AddLine(6)
	file.AddLine(3)
	file.AddLine(len(src))
	require.Equal(t, 4, file.LineCount())
}

func TestSourceFilePanics(t *testing.T) {
	fileSet := parser.NewFileSet()
	file := fileSet.AddFile("a", -1, 10)
	end := file.FileSetPos(10)

	expectPanic(t, "illegal base or size", func() {
		fileSet.AddFile("b", file.Base, 1)
	})
	expectPanic(t, "illegal base or size", func() {
		fileSet.AddFile("b", -1, -1)
	})
	expectPanic(t, "illegal line number (line numbering starts at 1)",
		func() { file.LineStart(0) })
	expectPanic(t, "illegal line number", func() { file.LineStart(2) })
	expectPanic(t, "illegal file offset", func() { file.FileSetPos(11) })
	expectPanic(t, "illegal SourcePos value", func() { file.Offset(end + 1) })
	expectPanic(t, "illegal SourcePos value",
		func() { file.Offset(parser.NoPos) })
	expectPanic(t, "illegal SourcePos value",
		func() { file.Position(end + 1) })
}

func expectPanic(t *testing.T, expected string, fn func()) {
	defer func() {
		require.Equal(t, expected, fmt.Sprint(recover()))
	}()
	fn()
}
