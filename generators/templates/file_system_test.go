package templates_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/mrasu/shouka/generators/templates"
	"github.com/mrasu/shouka/tests"
)

func TestFileSystem_LoadTemplate(t *testing.T) {
	contents := map[string]string{
		"dir/dummy.gotmpl": "{{.Foo}}: {{.Bar}}",
	}

	df := tests.NewDummyFileReadWriter(contents)
	fs := templates.NewFileSystemWithReadWriter(df)

	data := struct {
		Foo string
		Bar string
	}{
		Foo: "foo",
		Bar: "bar",
	}

	bytes, err := fs.LoadTemplate("dir", "dummy.gotmpl", data)
	require.NoError(t, err)

	assert.Equal(t, "foo: bar", bytes.String())
}

func TestFileSystem_WriteFile(t *testing.T) {
	df := tests.NewDummyFileReadWriter(map[string]string{})
	fs := templates.NewFileSystemWithReadWriter(df)

	require.NoError(t, fs.WriteFile("foo", "bar.txt", []byte("hello world")))

	assert.Equal(t, "hello world", df.WrittenContent("foo/bar.txt"))
}
