package generators

import (
	"github.com/mrasu/shouka/configs"
	"github.com/mrasu/shouka/generators/templates"
)

type AppCodeGenerator struct {
	file *templates.FileSystem

	outputDir string
}

func NewAppCodeGenerator(file *templates.FileSystem, config *configs.Config) *AppCodeGenerator {
	return &AppCodeGenerator{
		file: file,

		outputDir: config.Directory,
	}
}

func (acg *AppCodeGenerator) Generate(data *Data) error {
	filenames := []string{
		".dockerignore",
		".gitignore",
		"appspec.yml",
		"Dockerfile",
		"README.md",
		"go.mod",
		"main.go",
		"main_test.go",
	}

	for _, f := range filenames {
		if err := acg.writeTemplateFile(f, data); err != nil {
			return err
		}
	}

	return nil
}

func (acg *AppCodeGenerator) writeTemplateFile(file string, data *Data) error {
	writer, err := acg.file.LoadTemplate("templates", file+".gotmpl", data.appCode)
	if err != nil {
		return err
	}

	if err := acg.writeFile(file, writer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (acg *AppCodeGenerator) writeFile(name string, bytes []byte) error {
	return acg.file.WriteFile(acg.outputDir, name, bytes)
}
