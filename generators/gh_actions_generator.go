package generators

import (
	"path"
	"path/filepath"

	"github.com/mrasu/shouka/configs"
)

type GhActionsGenerator struct {
	file *file

	baseDir string
}

func NewGhActionsGenerator(file *file, config *configs.Config) *GhActionsGenerator {
	return &GhActionsGenerator{
		file: file,

		baseDir: path.Join(config.Directory, ".github"),
	}
}

func (gag *GhActionsGenerator) Generate(data *data) error {
	files := []string{
		"workflows/test.yml",
		"workflows/release.yml",
		"actions/add_tag_to_ecr_image/action.yml",
		"actions/login_aws/action.yml",
		"actions/login_ecr/action.yml",
	}
	for _, f := range files {
		if err := ensureDirectoryExistence(path.Join(gag.baseDir, filepath.Dir(f))); err != nil {
			return err
		}

		if err := gag.writeTemplateFile(f, data); err != nil {
			return err
		}
	}

	return nil
}

func (gag *GhActionsGenerator) writeTemplateFile(file string, data *data) error {
	writer, err := gag.file.loadTemplate("templates/.github", file+".gotmpl", data.ghActions)
	if err != nil {
		return err
	}

	if err := gag.writeFile(file, writer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (gag *GhActionsGenerator) writeFile(name string, bytes []byte) error {
	return gag.file.writeFile(gag.baseDir, name, bytes)
}
