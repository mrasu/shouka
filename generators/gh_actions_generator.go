package generators

import (
	"path"
	"path/filepath"

	"github.com/mrasu/shouka/generators/templates"

	"github.com/mrasu/shouka/configs"
)

type GhActionsGenerator struct {
	file *templates.FileSystem

	baseDir string
}

func NewGhActionsGenerator(file *templates.FileSystem, config *configs.Config) *GhActionsGenerator {
	return &GhActionsGenerator{
		file: file,

		baseDir: path.Join(config.Directory, ".github"),
	}
}

func (gag *GhActionsGenerator) Generate(data *Data) error {
	filenames := []string{
		"workflows/test.yml",
		"workflows/release.yml",
		"actions/add_tag_to_ecr_image/action.yml",
		"actions/login_aws/action.yml",
		"actions/login_ecr/action.yml",
	}
	for _, f := range filenames {
		if err := ensureDirectoryExistence(path.Join(gag.baseDir, filepath.Dir(f))); err != nil {
			return err
		}

		if err := gag.writeTemplateFile(f, data); err != nil {
			return err
		}
	}

	return nil
}

func (gag *GhActionsGenerator) writeTemplateFile(file string, data *Data) error {
	writer, err := gag.file.LoadTemplate("templates/.github", file+".gotmpl", data.ghActions)
	if err != nil {
		return err
	}

	if err := gag.writeFile(file, writer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (gag *GhActionsGenerator) writeFile(name string, bytes []byte) error {
	return gag.file.WriteFile(gag.baseDir, name, bytes)
}
