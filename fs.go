package flow

import (
	"net/http"
	"os"
)

type onlyFilesFS struct {
	fs http.FileSystem
}

type restrictedFile struct {
	http.File
}

// Dir returns a http.Filesystem that can be used by http.FileServer().
//
//It is used internally in router.Static().
// if listDirectory == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, listDirectory bool) http.FileSystem {
	fs := http.Dir(root)
	if listDirectory {
		return fs
	}

	return &onlyFilesFS{fs}
}

// Open conforms to http.Filesystem.
func (fs onlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	return restrictedFile{f}, nil
}

// Readdir overrides the http.File default implementation and
// disables directory listing
func (f restrictedFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}