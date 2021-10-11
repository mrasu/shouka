package generators_test

import (
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func listFiles(t *testing.T, dir string) []string {
	t.Helper()

	var files []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			base := strings.Replace(path, dir+"/", "", 1)
			files = append(files, base)
		}
		return nil
	})
	require.NoError(t, err)

	return files
}

func listTemplateFiles(t *testing.T, dir string) []string {
	t.Helper()

	files := listFiles(t, path.Join("../templates", dir))
	var tFiles []string
	for _, f := range files {
		base := strings.Replace(f, ".gotmpl", "", 1)
		tFiles = append(tFiles, path.Join(dir, base))
	}

	return tFiles
}

func listTemplateFilesUnderDirectory(t *testing.T, dir string) []string {
	var files []string
	dirFiles, err := ioutil.ReadDir(dir)
	require.NoError(t, err)
	for _, f := range dirFiles {
		if !f.IsDir() {
			base := strings.Replace(f.Name(), ".gotmpl", "", 1)
			files = append(files, base)
		}
	}

	return files
}

func joinDir(dir string, files []string) (res []string) {
	for _, f := range files {
		res = append(res, path.Join(dir, f))
	}

	return
}
