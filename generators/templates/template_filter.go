package templates

import (
	"io/fs"
	"io/ioutil"
)

type templateFileReadWriter struct {
	fs fs.ReadFileFS
}

func (tf *templateFileReadWriter) ReadFile(filepath string) ([]byte, error) {
	return tf.fs.ReadFile(filepath)
}

func (tf *templateFileReadWriter) WriteFile(filepath string, bytes []byte, perm fs.FileMode) error {
	return ioutil.WriteFile(filepath, bytes, perm)
}
