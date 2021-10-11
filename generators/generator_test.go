package generators_test

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
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

	assert.NoError(t, tg.Generate())

	expectedFiles := excludeDirs(listTemplateFiles(t, ""), []string{"docs"})

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

func excludeDirs(files, exDirs []string) []string {
	var res []string
	for _, f := range files {
		match := false
		for _, d := range exDirs {
			if strings.HasPrefix(f, d+"/") {
				match = true
				break
			}
		}
		if !match {
			res = append(res, f)
		}
	}

	return res
}

type dummyFs struct{}

func (df *dummyFs) Open(name string) (fs.File, error) {
	return os.Open(path.Join("../", name))
}

func (df *dummyFs) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(path.Join("../", name))
}
