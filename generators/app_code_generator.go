package generators

import (
	"github.com/mrasu/shouka/configs"
)

type AppCodeGenerator struct {
	file   *file
	config *configs.Config

	outputDir string
}

func NewAppCodeGenerator(file *file, config *configs.Config) *AppCodeGenerator {
	return &AppCodeGenerator{
		file:   file,
		config: config,

		outputDir: config.Directory,
	}
}

func (acg *AppCodeGenerator) Generate(data *data) error {
	files := []string{
		".dockerignore",
		".gitignore",
		"appspec.yml",
		"Dockerfile",
		"go.mod",
		"main.go",
		"main_test.go",
	}

	for _, f := range files {
		if err := acg.writeTemplateFile(f, data); err != nil {
			return err
		}
	}

	return nil
}

func (acg *AppCodeGenerator) writeTemplateFile(file string, data *data) error {
	writer, err := acg.file.loadTemplate("templates", file+".gotmpl", data.appCode)
	if err != nil {
		return err
	}

	if err := acg.writeFile(file, writer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (acg *AppCodeGenerator) writeFile(name string, bytes []byte) error {
	return acg.file.writeFile(acg.outputDir, name, bytes)
}
