package scanner

import (
	"sort"
	"sync"
)

type FileSet struct {
	mutex sync.RWMutex // protects the file set
	base  int          // base offset for the next file
	files []*File      // list of files in the order added to the set
	last  *File        // cache of last file looked up
}

func NewFileSet() *FileSet {
	return &FileSet{
		base: 1, // 0 == NoPos
	}
}

func (s *FileSet) Base() int {
	s.mutex.RLock()
	b := s.base
	s.mutex.RUnlock()

	return b
}

func (s *FileSet) AddFile(filename string, base, size int) *File {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if base < 0 {
		base = s.base
	}
	if base < s.base || size < 0 {
		panic("illegal base or size")
	}

	f := &File{
		set:   s,
		name:  filename,
		base:  base,
		size:  size,
		lines: []int{0},
	}

	base += size + 1 // +1 because EOF also has a position
	if base < 0 {
		panic("offset overflow (> 2G of source code in file set)")
	}

	// add the file to the file set
	s.base = base
	s.files = append(s.files, f)
	s.last = f

	return f
}

// File returns the file that contains the position p.
// If no such file is found (for instance for p == NoPos),
// the result is nil.
//
func (s *FileSet) File(p Pos) (f *File) {
	if p != NoPos {
		f = s.file(p)
	}

	return
}

// PositionFor converts a Pos p in the fileset into a FilePos value.
func (s *FileSet) Position(p Pos) (pos FilePos) {
	if p != NoPos {
		if f := s.file(p); f != nil {
			return f.position(p)
		}
	}

	return
}

func (s *FileSet) file(p Pos) *File {
	s.mutex.RLock()

	// common case: p is in last file
	if f := s.last; f != nil && f.base <= int(p) && int(p) <= f.base+f.size {
		s.mutex.RUnlock()
		return f
	}

	// p is not in last file - search all files
	if i := searchFiles(s.files, int(p)); i >= 0 {
		f := s.files[i]

		// f.base <= int(p) by definition of searchFiles
		if int(p) <= f.base+f.size {
			s.mutex.RUnlock()
			s.mutex.Lock()
			s.last = f // race is ok - s.last is only a cache
			s.mutex.Unlock()
			return f
		}
	}

	s.mutex.RUnlock()
	return nil
}

func searchFiles(a []*File, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i].base > x }) - 1
}
