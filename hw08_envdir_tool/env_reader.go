package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dirName string) (Environment, error) {
	dir, err := filepath.Abs(dirName)
	if err != nil {
		return nil, fmt.Errorf("get absolute path of directory %s: %w", dirName, err)
	}
	dirInfo, err := os.Stat(dir)
	if os.IsNotExist(err) || !dirInfo.IsDir() {
		return nil, fmt.Errorf("directory %s does not exist", dirName)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("get directory %s files: %w", dirName, err)
	}

	env := make(Environment)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := strings.TrimSpace(file.Name())
		err := checkFilename(fileName)
		if err != nil {
			return nil, fmt.Errorf("invalid file %s: %w", fileName, err)
		}

		filePath := filepath.Join(dir, fileName)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, fmt.Errorf("get file %s info: %w", filePath, err)
		}

		envValue := EnvValue{}

		fileSize := fileInfo.Size()
		if fileSize == 0 {
			envValue.NeedRemove = true
		} else {
			firstLine, err := readFirstLine(filePath)
			if err != nil {
				return nil, fmt.Errorf("read first line from file %s: %w", filePath, err)
			}
			envValue.Value = firstLine
		}

		env[fileName] = envValue
	}

	return env, nil
}

func checkFilename(filename string) error {
	if strings.Contains(filename, "=") {
		return fmt.Errorf("filename contains '='")
	}
	return nil
}

func readFirstLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read file %s: %w", filePath, err)
	}
	normalizedLine := normalizeLine(line)
	return normalizedLine, nil
}

func normalizeLine(line string) string {
	lineBytes := []byte(line)
	lineBytes = bytes.ReplaceAll(lineBytes, []byte{0x00}, []byte("\n"))
	lineBytes = bytes.TrimRight(lineBytes, " \n\r\t")
	return string(lineBytes)
}
