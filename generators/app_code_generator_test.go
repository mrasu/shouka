package generators_test

import (
	"testing"

	"github.com/mrasu/shouka/tests"

	"github.com/stretchr/testify/assert"

	"github.com/mrasu/shouka/generators"

	"github.com/mrasu/shouka/configs"
	"github.com/mrasu/shouka/generators/templates"
)

func TestAppCodeGenerator_Generate(t *testing.T) {
	acg, data, df := prepareApp(t)

	assert.NoError(t, acg.Generate(data))

	expectedFiles := listTemplateFilesUnderDirectory(t, "../templates")

	for _, f := range expectedFiles {
		assert.Equal(t, "test-sk-prefix: "+f, df.GetWritten("/tmp/dummy/"+f))
	}
	assert.ElementsMatch(t, joinDir("/tmp/dummy/", expectedFiles), df.WrittenFiles())
}

func prepareApp(t *testing.T) (*generators.AppCodeGenerator, *generators.Data, *tests.DummyFileReadWriter) {
	t.Helper()

	files := listTemplateFilesUnderDirectory(t, "../templates")
	contents := map[string]string{}
	for _, f := range files {
		contents["templates/"+f+".gotmpl"] = "{{.SkPrefix}}: " + f
	}

	df := tests.NewDummyFileReadWriter(contents)
	file := templates.NewFileSystemWithReadWriter(df)
	config := &configs.Config{
		Directory: "/tmp/dummy",
		SkPrefix:  "test-sk-prefix",
	}

	acg := generators.NewAppCodeGenerator(file, config)

	data := generators.NewData(config)

	return acg, data, df
}
