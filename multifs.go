package tengo

import (
	"errors"
	"io/fs"
)

// MultiFS is an fs.FS that tries each member filesystem in order.
// Open returns the first successful result. A non-ErrNotExist error from any
// member stops the search immediately and is returned to the caller.
//
// It is the idiomatic way to layer multiple virtual filesystems when using
// Script.SetImportFS:
//
//	s.SetImportFS(tengo.MultiFS{baseFS, userFS})
//
// Relative imports work correctly across members: the import directory context
// is a plain string prefix that is FS-agnostic, so a module found in any
// member can resolve its own relative imports across the full set.
type MultiFS []fs.FS

// Open implements fs.FS. Each member is tried in slice order; the first
// successful Open wins. If every member returns fs.ErrNotExist the error is
// wrapped in an fs.PathError for the requested name.
func (m MultiFS) Open(name string) (fs.File, error) {
	for _, fsys := range m {
		f, err := fsys.Open(name)
		if err == nil {
			return f, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			// Propagate genuine errors immediately.
			return nil, err
		}
	}
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}
