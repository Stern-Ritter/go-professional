package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFromPathEqualToPath   = errors.New("source and destination paths cannot be equal")
)

func Copy(fromPath string, toPath string, offset int64, limit int64) error {
	err := checkEqualPaths(fromPath, toPath)
	if err != nil {
		return err
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("open file %s: %w", fromPath, ErrUnsupportedFile)
	}
	defer fromFile.Close()

	toFile, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open file %s: %w", toPath, ErrUnsupportedFile)
	}
	defer toFile.Close()

	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return fmt.Errorf("get info about file %s: %w", fromPath, ErrUnsupportedFile)
	}

	fromFileSize := fromFileInfo.Size()
	if offset > fromFileSize {
		return fmt.Errorf("file %s: %w", fromPath, ErrOffsetExceedsFileSize)
	}

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek file %s: %w", fromPath, err)
	}

	if limit == 0 || limit > fromFileSize-offset {
		limit = fromFileSize - offset
	}

	bar := pb.Full.Start64(limit)
	defer bar.Finish()

	progressReader := bar.NewProxyReader(fromFile)

	_, err = io.CopyN(toFile, progressReader, limit)
	if err != nil {
		return fmt.Errorf("copy file %s to file %s: %w", toPath, toPath, err)
	}

	return nil
}

func checkEqualPaths(fromPath string, toPath string) error {
	absFromPath, err := filepath.Abs(filepath.Clean(fromPath))
	if err != nil {
		return fmt.Errorf("get absolute path %s: %w", fromPath, err)
	}

	absToPath, err := filepath.Abs(filepath.Clean(toPath))
	if err != nil {
		return fmt.Errorf("get absolute path %s: %w", toPath, err)
	}

	if absFromPath == absToPath {
		return fmt.Errorf("from path %s, to path %s: %w", fromPath, toPath, ErrFromPathEqualToPath)
	}

	return nil
}
