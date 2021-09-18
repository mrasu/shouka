package generators

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/pkg/errors"
)

type file struct {
	fs *embed.FS
}

func newFile(fs *embed.FS) *file {
	return &file{fs: fs}
}

func (f *file) loadTemplate(tmplDir, filename string, data interface{}) (*bytes.Buffer, error) {
	tmpl, err := f.load(tmplDir, filename)
	if err != nil {
		return nil, err
	}

	writer := bytes.Buffer{}
	if err := tmpl.Execute(&writer, data); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to write template: %s", filename))
	}

	return &writer, nil
}

func (f *file) load(tmplDir, name string) (*template.Template, error) {
	filepath := tmplDir + "/" + name
	text, err := f.fs.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read file: %s", filepath))
	}

	tmpl, err := template.New(name).Parse(string(text))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse: %s", name))
	}

	return tmpl, nil
}

func (f *file) writeFile(outputDir, name string, bytes []byte) error {
	filepath := path.Join(outputDir, name)
	if err := ioutil.WriteFile(filepath, bytes, 0666); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write file: %s", filepath))
	}

	return nil
}
