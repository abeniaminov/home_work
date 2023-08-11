package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetIsNegativ       = errors.New("offset is negative")
	ErrLimitIsNegativ        = errors.New("limit is negative")
)

type ProgressBar struct {
	Total int
	Max   int
}

func (pb *ProgressBar) Write(p []byte) (int, error) {
	n := len(p)
	pb.Total += n
	pb.PrintProgress()
	return n, nil
}

func (pb ProgressBar) PrintProgress() {
	persent := pb.Total / pb.Max * 100
	fmt.Printf("\r%s  %d%%", strings.Repeat("#", persent), persent)
	fmt.Printf("\rCopied... %d bytes ", pb.Total)
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	min := func(x, y int64) int64 {
		if x > y {
			return y
		}
		return x
	}

	if offset < 0 {
		return ErrOffsetIsNegativ
	}
	if limit < 0 {
		return ErrLimitIsNegativ
	}

	fromFileStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !fromFileStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	sourceFileSize := fromFileStat.Size()
	if offset > sourceFileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = sourceFileSize
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer dest.Close()

	realLimit := min(limit, sourceFileSize-offset)
	pb := &ProgressBar{Max: int(realLimit)}

	source.Seek(offset, 0)
	if _, err := io.CopyN(dest, io.TeeReader(source, pb), realLimit); err != nil {
		return err
	}

	return nil
}
