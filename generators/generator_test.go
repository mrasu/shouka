package generators_test

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mrasu/shouka/configs"
	"github.com/mrasu/shouka/generators"
	"github.com/stretchr/testify/assert"
)

func TestCodeGenerator_Generate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "generators")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tg := prepareGenerator(t, tmpDir)

	_, err = tg.Generate()
	assert.NoError(t, err)

	expectedFiles := listTemplateFiles(t, "")

	actualFiles := listFiles(t, tmpDir)
	assert.ElementsMatch(t, expectedFiles, actualFiles)
}

func prepareGenerator(t *testing.T, dir string) *generators.Generator {
	t.Helper()

	config := &configs.Config{
		Directory: dir,
	}

	tg := generators.NewGenerator(&dummyFs{}, config)

	return tg
}

type dummyFs struct{}

func (df *dummyFs) Open(name string) (fs.File, error) {
	return os.Open(path.Join("../", name))
}

func (df *dummyFs) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(path.Join("../", name))
}
