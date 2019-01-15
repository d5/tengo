package source

import (
	"sync"
)

// File represents a source file.
type File struct {
	set   *FileSet
	name  string // file name as provided to AddFile
	base  int    // Pos value range for this file is [base...base+size]
	size  int    // file size as provided to AddFile
	mutex sync.Mutex
	lines []int // lines contains the offset of the first character for each line (the first entry is always 0)
}

// Name returns the file name.
func (f *File) Name() string {
	return f.name
}

// Base returns the base position of the file.
func (f *File) Base() int {
	return f.base
}

// Size returns the size of the file.
func (f *File) Size() int {
	return f.size
}

// LineCount returns the current number of lines.
func (f *File) LineCount() int {
	f.mutex.Lock()
	n := len(f.lines)
	f.mutex.Unlock()

	return n
}

// AddLine adds a new line.
func (f *File) AddLine(offset int) {
	f.mutex.Lock()
	if i := len(f.lines); (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
	f.mutex.Unlock()
}

// LineStart returns the position of the first character in the line.
func (f *File) LineStart(line int) Pos {
	if line < 1 {
		panic("illegal line number (line numbering starts at 1)")
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	if line > len(f.lines) {
		panic("illegal line number")
	}

	return Pos(f.base + f.lines[line-1])
}

// FileSetPos returns the position in the file set.
func (f *File) FileSetPos(offset int) Pos {
	if offset > f.size {
		panic("illegal file offset")
	}

	return Pos(f.base + offset)
}

// Offset translates the file set position into the file offset.
func (f *File) Offset(p Pos) int {
	if int(p) < f.base || int(p) > f.base+f.size {
		panic("illegal Pos value")
	}

	return int(p) - f.base
}

// Position translates the file set position into the file position.
func (f *File) Position(p Pos) (pos FilePos) {
	if p != NoPos {
		if int(p) < f.base || int(p) > f.base+f.size {
			panic("illegal Pos value")
		}

		pos = f.position(p)
	}

	return
}

func (f *File) position(p Pos) (pos FilePos) {
	offset := int(p) - f.base
	pos.Offset = offset
	pos.Filename, pos.Line, pos.Column = f.unpack(offset)

	return
}

func (f *File) unpack(offset int) (filename string, line, column int) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	filename = f.name
	if i := searchInts(f.lines, offset); i >= 0 {
		line, column = i+1, offset-f.lines[i]+1
	}

	return
}

func searchInts(a []int, x int) int {
	// This function body is a manually inlined version of:
	//   return sort.Search(len(a), func(i int) bool { return a[i] > x }) - 1
	i, j := 0, len(a)
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h
		// i ≤ h < j
		if a[h] <= x {
			i = h + 1
		} else {
			j = h
		}
	}

	return i - 1
}
