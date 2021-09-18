package generators

import (
	"embed"

	"github.com/mrasu/shouka/configs"
)

type Generator struct {
	file   *file
	config *configs.Config
	data   *data
}

func NewGenerator(fs *embed.FS, config *configs.Config) *Generator {
	return &Generator{
		file:   newFile(fs),
		config: config,
	}
}

func (g *Generator) Generate() error {
	if err := ensureDirectoryExistence(g.config.Directory); err != nil {
		return err
	}

	data := newData(g.config)

	if err := NewAppCodeGenerator(g.file, g.config).Generate(data); err != nil {
		return err
	}

	if err := NewTfGenerator(g.file, g.config).Generate(data); err != nil {
		return err
	}

	if err := NewGhActionsGenerator(g.file, g.config).Generate(data); err != nil {
		return err
	}

	return nil
}
