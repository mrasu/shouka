package generators_test

import (
	"testing"

	"github.com/mrasu/shouka/tests"

	"github.com/mrasu/shouka/configs"
	"github.com/mrasu/shouka/generators"
	"github.com/mrasu/shouka/generators/templates"
	"github.com/stretchr/testify/assert"
)

func TestGhActionsGenerator_Generate(t *testing.T) {
	gag, data, df := prepareGh(t)

	assert.NoError(t, gag.Generate(data))

	expectedFiles := listTemplateFiles(t, ".github")

	for _, f := range expectedFiles {
		assert.Equal(t, "dummy-region: "+f, df.GetWritten("/tmp/dummy/"+f))
	}
	assert.ElementsMatch(t, joinDir("/tmp/dummy/", expectedFiles), df.WrittenFiles())
}

func prepareGh(t *testing.T) (*generators.GhActionsGenerator, *generators.Data, *tests.DummyFileReadWriter) {
	t.Helper()

	files := listTemplateFiles(t, ".github")
	contents := map[string]string{}
	for _, f := range files {
		contents["templates/"+f+".gotmpl"] = "{{.AwsRegion}}: " + f
	}

	df := tests.NewDummyFileReadWriter(contents)
	file := templates.NewFileSystemWithReadWriter(df)
	config := &configs.Config{
		Directory: "/tmp/dummy",
		Resources: configs.ResourceConfig{Region: "dummy-region"},
	}

	acg := generators.NewGhActionsGenerator(file, config)

	data := generators.NewData(config)

	return acg, data, df
}
