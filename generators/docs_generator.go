package generators

import (
	"fmt"
	"path"

	"github.com/mrasu/shouka/configs"
	"github.com/mrasu/shouka/generators/templates"
)

type Combination int

const (
	FargateGithubActions Combination = iota + 1
)

var fargateGithubActionsTldr = `You can set up CI/CD with ECS(Fargate) and GitHub Actions.

1. Run terraform
2. Add words at ` + "`%%SK-CHANGE-REQUIRED%%`" + `
3. Push to GitHub
`

type DocsGenerator struct {
	file *templates.FileSystem

	baseDir string
}

type docsTmplData struct {
	Tldr string
}

func NewDocsGenerator(file *templates.FileSystem, config *configs.Config) *DocsGenerator {
	return &DocsGenerator{
		file: file,

		baseDir: path.Join(config.Directory, "docs"),
	}
}

func (dg *DocsGenerator) Generate(data *Data) (string, error) {
	if err := ensureDirectoryExistence(path.Join(dg.baseDir)); err != nil {
		return "", err
	}

	tmplData := dg.createTmplData(data)

	filenames := []string{
		"README.md",
	}
	for _, f := range filenames {
		if err := dg.writeTemplateFile(f, tmplData); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s\nYou can read more at docs/README.md", tmplData.Tldr), nil
}

func (dg *DocsGenerator) createTmplData(data *Data) *docsTmplData {
	switch data.docs.tldr {
	case FargateGithubActions:
		return &docsTmplData{Tldr: fargateGithubActionsTldr}
	default:
		panic("Unknown enum")
	}
}

func (dg *DocsGenerator) writeTemplateFile(file string, tmplData *docsTmplData) error {
	writer, err := dg.file.LoadTemplate("templates/docs", file+".gotmpl", tmplData)
	if err != nil {
		return err
	}

	if err := dg.writeFile(file, writer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (dg *DocsGenerator) writeFile(name string, bytes []byte) error {
	return dg.file.WriteFile(dg.baseDir, name, bytes)
}
