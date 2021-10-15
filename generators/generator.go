package generators

import (
	"io/fs"

	"github.com/mrasu/shouka/generators/templates"

	"github.com/mrasu/shouka/configs"
)

type Generator struct {
	file   *templates.FileSystem
	config *configs.Config
}

func NewGenerator(fs fs.ReadFileFS, config *configs.Config) *Generator {
	return &Generator{
		file:   templates.NewFileSystem(fs),
		config: config,
	}
}

func (g *Generator) Generate() (string, error) {
	if err := ensureDirectoryExistence(g.config.Directory); err != nil {
		return "", err
	}

	data := NewData(g.config)

	if err := NewAppCodeGenerator(g.file, g.config).Generate(data); err != nil {
		return "", err
	}

	if err := NewTfGenerator(g.file, g.config).Generate(data); err != nil {
		return "", err
	}

	if err := NewGhActionsGenerator(g.file, g.config).Generate(data); err != nil {
		return "", err
	}

	msg, err := NewDocsGenerator(g.file, g.config).Generate(data)
	if err != nil {
		return "", err
	}

	return msg, nil
}
