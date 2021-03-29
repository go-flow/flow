package response

import "net/http"

type File struct {
	filepath string
}

func NewFile(filepath string) *File {
	return &File{
		filepath: filepath,
	}
}

func (File) Status() int {
	return http.StatusOK
}

func (rf *File) Handle(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, rf.filepath)
	return nil
}
