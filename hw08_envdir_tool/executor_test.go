package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	defer func() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	stdoutFile, err := os.CreateTemp("", "stdout")
	require.NoError(t, err, "unexpected error creating temp stdout file")
	defer os.Remove(stdoutFile.Name())

	stderrFile, err := os.CreateTemp("", "stderr")
	require.NoError(t, err, "unexpected error creating temp stderr file")
	defer os.Remove(stderrFile.Name())

	os.Stdout = stdoutFile
	os.Stderr = stderrFile

	cmd := []string{"sh", "-c", "echo test; echo TEST_ENV=$TEST_ENV >&2"}
	env := Environment{
		"TEST_ENV": {Value: "test env value", NeedRemove: false},
	}

	returnCode := RunCmd(cmd, env)

	stdoutFile.Seek(0, io.SeekStart)
	stdoutContent, err := io.ReadAll(stdoutFile)
	require.NoError(t, err, "unexpected error reading temp stdout file")

	stderrFile.Seek(0, io.SeekStart)
	stderrContent, err := io.ReadAll(stderrFile)
	require.NoError(t, err, "unexpected error reading temp stderr file")

	expectedStdout := "test\n"
	expectedStderr := "TEST_ENV=test env value\n"

	assert.Equal(t, 0, returnCode, "expected return code to be 0")
	assert.Equal(t, expectedStdout, string(stdoutContent), "stdout should contain: %s but got: %s",
		expectedStdout, string(stdoutContent))
	assert.Equal(t, expectedStderr, string(stderrContent), "stderr should contain: %s but got: %s",
		expectedStderr, string(stderrContent))
}
