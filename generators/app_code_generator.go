package generators

import (
	"bytes"
	"fmt"
	"github.com/mrasu/shouka/configs"
	"github.com/pkg/errors"
	"os"
	"path"
	"text/template"
)

type AppCodeGenerator struct {
	file   *file
	config *configs.Config
}

func NewAppCodeGenerator(file *file, config *configs.Config) *AppCodeGenerator {
	return &AppCodeGenerator{
		file:   file,
		config: config,
	}
}

func (acg *AppCodeGenerator) Generate() error {
	if err := acg.ensureDirectoryExistence(); err != nil {
		return err
	}
	if err := acg.writeTemplateFile("Dockerfile"); err != nil {
		return err
	}
	if err := acg.writeTemplateFile("main.go"); err != nil {
		return err
	}

	return nil
}

func (acg *AppCodeGenerator) ensureDirectoryExistence() error {
	if err := os.MkdirAll(acg.outputDir(), 0776); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create app_codes directory: %s", acg.outputDir()))
	}

	return nil
}

func (acg *AppCodeGenerator) outputDir() string {
	return path.Join(acg.config.Directory, "app_codes")
}

func (acg *AppCodeGenerator) writeTemplateFile(file string) error {
	tmplFile := file + ".gotmpl"
	tmpl, err := acg.load(tmplFile)
	if err != nil {
		return err
	}

	writer := bytes.Buffer{}
	if err := tmpl.Execute(&writer, nil); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write template: %s", tmplFile))
	}

	if err := acg.writeFile(file, writer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (acg *AppCodeGenerator) load(name string) (*template.Template, error) {
	return acg.file.load("templates/app_codes", name)
}

func (acg *AppCodeGenerator) writeFile(name string, bytes []byte) error {
	return acg.file.writeFile(acg.outputDir(), name, bytes)
}
