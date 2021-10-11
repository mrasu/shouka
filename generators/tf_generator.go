package generators

import (
	"path"

	"github.com/mrasu/shouka/generators/templates"

	"github.com/mrasu/shouka/configs"
)

type TfGenerator struct {
	file   *templates.FileSystem
	config *configs.Config

	outputDir string
}

func NewTfGenerator(file *templates.FileSystem, config *configs.Config) *TfGenerator {
	return &TfGenerator{
		file:   file,
		config: config,

		outputDir: path.Join(config.Directory, "terraforms"),
	}
}

func (tg *TfGenerator) Generate(data *Data) error {
	if err := ensureDirectoryExistence(tg.outputDir); err != nil {
		return err
	}

	if err := tg.writeTemplateFile(".gitignore", data); err != nil {
		return err
	}

	if err := tg.writeResourceFile("main", data); err != nil {
		return err
	}

	if err := tg.writeResourceFile("code_deploy", data); err != nil {
		return err
	}

	if err := tg.writeResourceFile("ecs", data); err != nil {
		return err
	}

	if err := tg.writeResourceFile("load_balancer", data); err != nil {
		return err
	}

	if err := tg.writeResourceFile("variables", data); err != nil {
		return err
	}

	if tg.config.Resources.CloudWatch.RequiresTemplate() {
		if err := tg.writeResourceFile("cloud_watch", data); err != nil {
			return err
		}
	}

	if tg.config.Resources.Ecr.RequiresTemplate() {
		if err := tg.writeResourceFile("ecr", data); err != nil {
			return err
		}
	}

	if tg.config.Resources.Iam.RequiresTemplate() {
		if err := tg.writeResourceFile("iam", data); err != nil {
			return err
		}
	}

	if tg.config.Resources.SecurityGroup.RequiresTemplate() {
		if err := tg.writeResourceFile("security_group", data); err != nil {
			return err
		}
	}

	if tg.config.Resources.Subnet.RequiresTemplate() {
		if err := tg.writeResourceFile("subnet", data); err != nil {
			return err
		}
	}

	if tg.config.Resources.Vpc.RequiresTemplate() {
		if err := tg.writeResourceFile("vpc", data); err != nil {
			return err
		}
	}

	return nil
}

func (tg *TfGenerator) writeTemplateFile(filename string, data *Data) error {
	writer, err := tg.file.LoadTemplate("templates/terraforms", filename+".gotmpl", data.Resources)
	if err != nil {
		return err
	}

	if err := tg.writeFile(filename, writer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (tg *TfGenerator) writeResourceFile(resource string, data *Data) error {
	return tg.writeTemplateFile(resource+".tf", data)
}

func (tg *TfGenerator) writeFile(name string, bytes []byte) error {
	return tg.file.WriteFile(tg.outputDir, name, bytes)
}
