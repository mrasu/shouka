package tests

import (
	"fmt"
	"io/fs"

	"github.com/pkg/errors"
)

type DummyFileReadWriter struct {
	contents        map[string]string
	writtenContents map[string]string
}

func NewDummyFileReadWriter(contents map[string]string) *DummyFileReadWriter {
	return &DummyFileReadWriter{
		contents:        contents,
		writtenContents: map[string]string{},
	}
}

func (df *DummyFileReadWriter) ReadFile(filepath string) ([]byte, error) {
	if c, ok := df.contents[filepath]; ok {
		return []byte(c), nil
	}
	return nil, errors.New(fmt.Sprintf("Not registered path for DummyFileReadWriter.ReadFile. filepath: %s", filepath))
}

func (df *DummyFileReadWriter) WriteFile(filepath string, contents []byte, perm fs.FileMode) error {
	df.writtenContents[filepath] = string(contents)
	return nil
}

func (df *DummyFileReadWriter) GetWritten(filepath string) string {
	return df.writtenContents[filepath]
}

func (df *DummyFileReadWriter) WrittenFiles() (files []string) {
	for f := range df.writtenContents {
		files = append(files, f)
	}
	return
}

func (df *DummyFileReadWriter) WrittenContent(filepath string) string {
	return df.writtenContents[filepath]
}
