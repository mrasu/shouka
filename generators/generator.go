package generators

import (
	"embed"

	"github.com/mrasu/shouka/configs"
)

type Generator struct {
	file   *file
	config *configs.Config
}

func NewGenerator(fs *embed.FS, config *configs.Config) *Generator {
	return &Generator{
		file:   newFile(fs),
		config: config,
	}
}

func (g *Generator) Generate() error {
	if g.config.CreatesAppCode {
		if err := NewAppCodeGenerator(g.file, g.config).Generate(); err != nil {
			return err
		}
	}

	if err := NewTfGenerator(g.file, g.config).Generate(); err != nil {
		return err
	}

	return nil
}
