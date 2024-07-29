package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	testDir := "./testdata"
	outputDir := filepath.Join(testDir, "output")
	expectedDir := filepath.Join(testDir, "expected")

	err := os.MkdirAll(outputDir, 0o755)
	require.NoError(t, err, "unexpected error creating test data output directory")
	defer os.RemoveAll(outputDir)

	tests := []struct {
		name          string
		fromPath      string
		toPath        string
		expectedPath  string
		offset        int64
		limit         int64
		expectedError error
	}{
		{
			name:          "file copy full",
			fromPath:      filepath.Join(testDir, "source.txt"),
			toPath:        filepath.Join(outputDir, "full_copy.txt"),
			expectedPath:  filepath.Join(expectedDir, "full_copy.txt"),
			offset:        0,
			limit:         0,
			expectedError: nil,
		},
		{
			name:          "file copy with offset",
			fromPath:      filepath.Join(testDir, "source.txt"),
			toPath:        filepath.Join(outputDir, "offset_copy.txt"),
			expectedPath:  filepath.Join(expectedDir, "offset_copy.txt"),
			offset:        12,
			limit:         0,
			expectedError: nil,
		},
		{
			name:          "file copy with limit",
			fromPath:      filepath.Join(testDir, "source.txt"),
			toPath:        filepath.Join(outputDir, "limit_copy.txt"),
			expectedPath:  filepath.Join(expectedDir, "limit_copy.txt"),
			offset:        0,
			limit:         10,
			expectedError: nil,
		},
		{
			name:          "file not found",
			fromPath:      filepath.Join(testDir, "not_found.txt"),
			toPath:        filepath.Join(outputDir, "error_copy.txt"),
			expectedPath:  "",
			offset:        0,
			limit:         0,
			expectedError: ErrUnsupportedFile,
		},
		{
			name:          "file copy with offset exceeds file size",
			fromPath:      filepath.Join(testDir, "source.txt"),
			toPath:        filepath.Join(outputDir, "offset_exceeds_copy.txt"),
			expectedPath:  "",
			offset:        1000,
			limit:         0,
			expectedError: ErrOffsetExceedsFileSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.fromPath, tt.toPath, tt.offset, tt.limit)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError, "should return error: %v, but got: %v", tt.expectedError, err)
			} else {
				assert.NoError(t, err, "shouldn`t return error, but got: %v", err)

				expectedFileContent, err := os.ReadFile(tt.expectedPath)
				require.NoError(t, err, "unexpected error reading expected file: %v", tt.expectedPath)
				gotFileContent, err := os.ReadFile(tt.toPath)
				require.NoError(t, err, "unexpected error reading got file: %v", tt.toPath)

				assert.Equal(t, expectedFileContent, gotFileContent,
					"content copied file should be: '%s',got: '%s'", string(expectedFileContent), string(gotFileContent))
			}
		})
	}
}
