package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	dir := t.TempDir()

	file1 := filepath.Join(dir, "ENV1")
	file2 := filepath.Join(dir, "ENV2")

	writeFile(t, file1, "value1")
	writeFile(t, file2, "value2")

	env, err := ReadDir(dir)
	require.NoError(t, err, "unexpected error reading directory")
	assert.Len(t, env, 2, "expected 2 environment variables")
	assert.Equal(t, EnvValue{Value: "value1"}, env["ENV1"], "unexpected value for ENV1")
	assert.Equal(t, EnvValue{Value: "value2"}, env["ENV2"], "unexpected value for ENV2")
}

func TestReadDir_EmptyFile(t *testing.T) {
	dir := t.TempDir()

	filePath := filepath.Join(dir, "ENV3")
	file, err := os.Create(filePath)
	require.NoError(t, err, "unexpected error creating file")
	err = file.Close()
	require.NoError(t, err, "unexpected error closing file")

	env, err := ReadDir(dir)
	require.NoError(t, err, "unexpected error reading directory")
	assert.Len(t, env, 1, "expected 1 environment variable")
	assert.Equal(t, EnvValue{NeedRemove: true}, env["ENV3"], "unexpected value for ENV3")
}

func TestReadFirstLine(t *testing.T) {
	file, err := os.CreateTemp("", "test")
	require.NoError(t, err, "unexpected error creating temp file")
	defer os.Remove(file.Name())

	_, err = file.WriteString("line1\nline2\n")
	require.NoError(t, err, "unexpected error writing to file")

	file.Close()

	line, err := readFirstLine(file.Name())
	require.NoError(t, err, "unexpected error reading first line")
	assert.Equal(t, "line1", line, "unexpected line read from file")
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	file, err := os.Create(path)
	require.NoError(t, err, "unexpected error creating file")
	defer file.Close()

	_, err = file.WriteString(content + "\n")
	require.NoError(t, err, "unexpected error writing to file")
}
