package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile               = errors.New("unsupported file")
	ErrOffsetExceedsFileSize         = errors.New("offset exceeds file size")
	ErrOffsetIsNegativ               = errors.New("offset is negative")
	ErrLimitIsNegativ                = errors.New("limit is negative")
	ErrSourceAndDesinationAreTheSame = errors.New("source and destination are the same")
)

type ProgressBar struct {
	Total int64
	Max   int64
}

func (pb *ProgressBar) Write(p []byte) (int, error) {
	n := len(p)
	pb.Total += int64(n)
	pb.PrintProgress()
	return n, nil
}

func (pb ProgressBar) PrintProgress() {
	var percent int64
	if pb.Max == 0 {
		percent = 100
	} else {
		percent = pb.Total / pb.Max * 100
	}
	fmt.Printf("\r%s  %d%%", strings.Repeat("#", int(percent)), percent)
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

	toFileStat, err := os.Stat(toPath)
	if err == nil {
		if os.SameFile(fromFileStat, toFileStat) {
			return ErrSourceAndDesinationAreTheSame
		}
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
	pb := &ProgressBar{Max: realLimit}

	_, err = source.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.CopyN(dest, io.TeeReader(source, pb), realLimit)
	return err
}
