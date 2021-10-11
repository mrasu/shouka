package templates

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"text/template"

	"github.com/pkg/errors"
)

type FileSystem struct {
	filer FileReadWriter
}

func NewFileSystem(fs fs.ReadFileFS) *FileSystem {
	return NewFileSystemWithReadWriter(
		&templateFileReadWriter{fs: fs},
	)
}

func NewFileSystemWithReadWriter(filer FileReadWriter) *FileSystem {
	return &FileSystem{
		filer: filer,
	}
}

func (fs *FileSystem) LoadTemplate(tmplDir, filename string, data interface{}) (*bytes.Buffer, error) {
	tmpl, err := fs.load(tmplDir, filename)
	if err != nil {
		return nil, err
	}

	writer := bytes.Buffer{}
	if err := tmpl.Execute(&writer, data); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to write template: %s", filename))
	}

	return &writer, nil
}

func (fs *FileSystem) load(tmplDir, name string) (*template.Template, error) {
	filepath := tmplDir + "/" + name
	text, err := fs.filer.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read FileSystem: %s", filepath))
	}

	tmpl, err := template.New(name).Parse(string(text))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse: %s", name))
	}

	return tmpl, nil
}

func (fs *FileSystem) WriteFile(outputDir, name string, bytes []byte) error {
	filepath := path.Join(outputDir, name)
	if err := fs.filer.WriteFile(filepath, bytes, 0666); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write FileSystem: %s", filepath))
	}

	return nil
}
