package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("file from is unsupported", func(t *testing.T) {
		err := Copy("/", "/tmp/mm.txt", 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})
	t.Run("offset is negative", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/mm.txt", -1, 0)
		require.Truef(t, errors.Is(err, ErrOffsetIsNegativ), "actual err - %v", err)
	})
	t.Run("limit is negative", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/mm.txt", 0, -1)
		require.Truef(t, errors.Is(err, ErrLimitIsNegativ), "actual err - %v", err)
	})
	t.Run("offset is large", func(t *testing.T) {
		fromFileStat, _ := os.Stat("testdata/input.txt")
		sourceFileSize := fromFileStat.Size()
		err := Copy("testdata/input.txt", "/tmp/mm.txt", sourceFileSize+1, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})
	t.Run("offset eq file size", func(t *testing.T) {
		fromFileStat, _ := os.Stat("testdata/input.txt")
		sourceFileSize := fromFileStat.Size()
		Copy("testdata/input.txt", "/tmp/mm.txt", sourceFileSize, 0)
		toFileStat, _ := os.Stat("/tmp/mm.txt")
		destFileSize := toFileStat.Size()
		require.Equal(t, 0, int(destFileSize), "wrong dest file size")
		os.Remove(toFileStat.Name())
	})
	t.Run("copy real limit", func(t *testing.T) {
		fromFileStat, _ := os.Stat("testdata/input.txt")
		sourceFileSize := fromFileStat.Size()
		Copy("testdata/input.txt", "/tmp/mm.txt", sourceFileSize-2, 0)
		toFileStat, _ := os.Stat("/tmp/mm.txt")
		destFileSize := toFileStat.Size()
		require.Equal(t, 2, int(destFileSize), "wrong dest file size")
		os.Remove(toFileStat.Name())
	})
	t.Run("limit is larger than file size", func(t *testing.T) {
		fromFileStat, _ := os.Stat("testdata/input.txt")
		sourceFileSize := fromFileStat.Size()
		Copy("testdata/input.txt", "/tmp/mm.txt", 0, sourceFileSize+10)
		toFileStat, _ := os.Stat("/tmp/mm.txt")
		destFileSize := toFileStat.Size()
		require.Equal(t, int(sourceFileSize), int(destFileSize), "wrong dest file size")
		os.Remove(toFileStat.Name())
	})
}
