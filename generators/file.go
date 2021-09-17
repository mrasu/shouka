package generators

import (
	"embed"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
	"text/template"
)

type file struct {
	fs *embed.FS
}

func newFile(fs *embed.FS) *file {
	return &file{fs: fs}
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
