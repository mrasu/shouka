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

type TfGenerator struct {
	file   *file
	config *configs.Config
}

func NewTfGenerator(file *file, config *configs.Config) *TfGenerator {
	return &TfGenerator{
		file:   file,
		config: config,
	}
}

func (tg *TfGenerator) Generate() error {
	if err := tg.ensureDirectoryExistence(); err != nil {
		return err
	}

	if err := tg.writeResourceFile("main"); err != nil {
		return err
	}

	if err := tg.writeResourceFile("code_deploy"); err != nil {
		return err
	}

	if err := tg.writeResourceFile("ecs"); err != nil {
		return err
	}

	if err := tg.writeResourceFile("load_balancer"); err != nil {
		return err
	}

	if err := tg.writeResourceFile("variables"); err != nil {
		return err
	}

	if tg.config.Resources.CloudWatch.RequiresTemplate() {
		if err := tg.writeResourceFile("cloud_watch"); err != nil {
			return err
		}
	}

	if tg.config.Resources.Ecr.RequiresTemplate() {
		if err := tg.writeResourceFile("ecr"); err != nil {
			return err
		}
	}

	if tg.config.Resources.Iam.RequiresTemplate() {
		if err := tg.writeResourceFile("iam"); err != nil {
			return err
		}
	}

	if tg.config.Resources.SecurityGroup.RequiresTemplate() {
		if err := tg.writeResourceFile("security_group"); err != nil {
			return err
		}
	}

	if tg.config.Resources.Subnet.RequiresTemplate() {
		if err := tg.writeResourceFile("subnet"); err != nil {
			return err
		}
	}

	if tg.config.Resources.Vpc.RequiresTemplate() {
		if err := tg.writeResourceFile("vpc"); err != nil {
			return err
		}
	}

	return nil
}

func (tg *TfGenerator) ensureDirectoryExistence() error {
	if err := os.MkdirAll(tg.outputDir(), 0776); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create terraforms directory: %s", tg.outputDir()))
	}

	return nil
}

func (tg *TfGenerator) outputDir() string {
	return path.Join(tg.config.Directory, "terraforms")
}

func (tg *TfGenerator) writeResourceFile(resource string) error {
	tmplFile := resource + ".tf.gotmpl"
	tmpl, err := tg.load(tmplFile)
	if err != nil {
		return err
	}

	writer := bytes.Buffer{}
	if err := tmpl.Execute(&writer, tg.config.Resources); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write template: %s", tmplFile))
	}

	if err := tg.writeFile(resource+".tf", writer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (tg *TfGenerator) load(name string) (*template.Template, error) {
	return tg.file.load("templates/terraforms", name)
}

func (tg *TfGenerator) writeFile(name string, bytes []byte) error {
	return tg.file.writeFile(tg.outputDir(), name, bytes)
}
