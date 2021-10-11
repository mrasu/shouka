package templates

import "io/fs"

type FileReadWriter interface {
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte, fs.FileMode) error
}
