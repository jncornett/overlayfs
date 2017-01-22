package overlayfs_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/jncornett/overlayfs"
)

func TestEmptyOverlay(t *testing.T) {
	tests := []string{"", "foo", "foo/bar.txt"}
	fs := overlayfs.NewOverlayFs()
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			f, err := fs.Open(test)
			if f != nil {
				t.Error("expected f to be nil")
			}
			if os.ErrNotExist != err {
				t.Errorf("expected err to be %v, got %v", os.ErrNotExist, err)
			}
		})
	}
}

type mockFile struct {
	overlay string
	name    string
}

func (f *mockFile) Close() error                       { return nil }
func (f *mockFile) Read([]byte) (int, error)           { return 0, nil }
func (f *mockFile) Seek(int64, int) (int64, error)     { return 0, nil }
func (f *mockFile) Readdir(int) ([]os.FileInfo, error) { return nil, nil }
func (f *mockFile) Stat() (os.FileInfo, error)         { return nil, nil }

var _ http.File = &mockFile{}

func mockFs(label string, files ...string) overlayfs.FileSystemFunc {
	return func(name string) (http.File, error) {
		for _, f := range files {
			if name == f {
				return &mockFile{overlay: label, name: name}, nil
			}
		}
		return nil, os.ErrNotExist
	}
}

func TestOverlay(t *testing.T) {
	tests := []struct {
		name     string
		overlays []http.FileSystem
		filename string
		found    bool
		source   string
	}{
		{
			name: "notpresent",
			overlays: []http.FileSystem{
				mockFs("a", "one.txt", "two.txt"),
				mockFs("b", "three.txt", "four.txt"),
			},
			filename: "five.txt",
			found:    false,
		},
		{
			name: "firstlayer",
			overlays: []http.FileSystem{
				mockFs("a", "one.txt", "two.txt"),
				mockFs("b", "two.txt", "three.txt"),
			},
			filename: "two.txt",
			found:    true,
			source:   "a",
		},
		{
			name: "secondlayer",
			overlays: []http.FileSystem{
				mockFs("a", "one.txt", "two.txt"),
				mockFs("b", "two.txt", "three.txt"),
			},
			filename: "three.txt",
			found:    true,
			source:   "b",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs := overlayfs.NewOverlayFs(test.overlays...)
			f, err := fs.Open(test.filename)
			if test.found {
				if err != nil {
					t.Fatal("expected err to be nil")
				}
				if f == nil {
					t.Fatal("expected f to not be nil")
				}
				mf, ok := f.(*mockFile)
				if !ok {
					t.Fatal("expected f to be of type *mockFile")
				}
				if test.filename != mf.name {
					t.Errorf("expected filename to be %v, got %v", test.filename, mf.name)
				}
				if test.source != mf.overlay {
					t.Errorf("expected source overlay to be %v, got %v", test.source, mf.overlay)
				}
			} else {
				if os.ErrNotExist != err {
					t.Errorf("expected err to be %v, got %v", os.ErrNotExist, err)
				}
			}
		})
	}
}
