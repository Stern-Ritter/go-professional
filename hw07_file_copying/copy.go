package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset int64, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("open file %s: %w", fromPath, ErrUnsupportedFile)
	}
	defer fromFile.Close()

	toFile, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return fmt.Errorf("open file %s: %w", toPath, ErrUnsupportedFile)
	}
	defer toFile.Close()

	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return fmt.Errorf("get stat file %s: %w", fromPath, ErrUnsupportedFile)
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

	progressReader := bar.NewProxyReader(io.LimitReader(fromFile, limit))

	_, err = io.Copy(toFile, progressReader)
	if err != nil {
		return fmt.Errorf("copy file %s to file %s: %w", fromPath, toPath, err)
	}

	return nil
}
