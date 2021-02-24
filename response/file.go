package response

import "net/http"

type File struct {
	filepath string
}

func (rf *File) Handle(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, rf.filepath)
	return nil
}

func NewFile(filepath string) *File {
	return &File{
		filepath: filepath,
	}
}
