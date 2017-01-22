package overlayfs

import (
	"net/http"
	"os"
)

// FileSystemFunc is an adapter to allow the use of ordinary functions as
// HTTP file systems.
type FileSystemFunc func(name string) (http.File, error)

// Open calls f(name)
func (f FileSystemFunc) Open(name string) (http.File, error) { return f(name) }

// NewOverlayFs returns a FileSystemFunc that iterates through filesystems,
// repeatedly calling Open(name) on them until one of the calls returns a
// nil error value. If the function exhausts the filesystem list without
// successfully finding a file, it returns os.ErrNotExist as the error value.
func NewOverlayFs(filesystems ...http.FileSystem) FileSystemFunc {
	return func(name string) (f http.File, err error) {
		for _, fs := range filesystems {
			f, err = fs.Open(name)
			if err == nil {
				return // found!
			}
		}
		err = os.ErrNotExist
		return
	}
}
