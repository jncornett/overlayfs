# overlayfs
A stupidly simple HTTP overlay file system for use with golang's "net/http" package.

[![GoDoc](https://godoc.org/github.com/jncornett/overlayfs?status.svg)](https://godoc.org/github.com/jncornett/overlayfs)

## usage
The following example serves an overlay file system based on files in `/tmp` first, and if the file is not found it will look in `/etc`.
```go
import (
	"net/http"
	
	"github.com/jncornett/overlayfs"
)

func main() {
	http.Handle("/", overlayfs.NewOverlayFs(
		http.FileServer(http.Dir("/tmp")),
		http.FileServer(http.Dir("/etc")),
	))
	http.ListenAndServe(":8080", nil)
}
```
