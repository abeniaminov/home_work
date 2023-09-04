package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("fail exec multiline test", func(t *testing.T) {
		content := []byte("foo\nbar\nerl")
		command := []string{"ls", "-a"}
		v := prepareExecutor(command, content)
		expected := 0
		require.Equal(t, expected, v)
		env, ok := os.LookupEnv("FOO")
		require.Equal(t, ok, true)
		require.Equal(t, env, "foo")
	})
	t.Run("fail exec multiline param test", func(t *testing.T) {
		content := []byte("-l\nbar\nerl")
		command := []string{"ls", "${FOO}"}
		v := prepareExecutor(command, content)
		expected := 0
		require.Equal(t, expected, v)
		env, ok := os.LookupEnv("FOO")
		require.Equal(t, ok, true)
		require.Equal(t, env, "-l")
	})
	t.Run("fail exec 0x00 test", func(t *testing.T) {
		content := []byte("foo" + string([]byte{0x00}) + "erl")
		command := []string{"ls", "-l"}
		v := prepareExecutor(command, content)
		expected := 0
		require.Equal(t, expected, v)
		env, ok := os.LookupEnv("FOO")
		require.Equal(t, ok, true)
		require.Equal(t, env, "foo\nerl")
	})
	t.Run("fail exec tabs test", func(t *testing.T) {
		content := []byte("foo\terl\t\t")
		command := []string{"ls", "-l"}
		v := prepareExecutor(command, content)
		expected := 0
		require.Equal(t, expected, v)
		env, ok := os.LookupEnv("FOO")
		require.Equal(t, ok, true)
		require.Equal(t, env, "foo\terl")
	})
	t.Run("fail exec spaces test", func(t *testing.T) {
		content := []byte("foo erl  ")
		command := []string{"ls", "-l"}
		v := prepareExecutor(command, content)
		expected := 0
		require.Equal(t, expected, v)
		env, ok := os.LookupEnv("FOO")
		require.Equal(t, ok, true)
		require.Equal(t, env, "foo erl")
	})
	t.Run("fail exec combination test", func(t *testing.T) {
		content := []byte("foo erl \t" + string([]byte{0x00}) + "boo  \t\t")
		command := []string{"ls", "-l"}
		v := prepareExecutor(command, content)
		expected := 0
		require.Equal(t, expected, v)
		env, ok := os.LookupEnv("FOO")
		require.Equal(t, ok, true)
		require.Equal(t, env, "foo erl \t\nboo")
	})
	t.Run("fail exec need remove test", func(t *testing.T) {
		command := []string{"ls", "-l"}
		v := prepareExecutor(command, nil)
		expected := 0
		require.Equal(t, expected, v)
		_, ok := os.LookupEnv("FOO")
		require.Equal(t, ok, false)
	})
}

func prepareExecutor(args []string, fileContent []byte) (returnCode int) {
	tmpdir, err := os.MkdirTemp("/tmp", "testdir")
	if err != nil {
		return 1
	}
	defer os.RemoveAll(tmpdir)

	f, err := os.Create(filepath.Join(tmpdir, "FOO"))
	if err != nil {
		return 1
	}
	if _, err := f.Write(fileContent); err != nil {
		return 1
	}

	if err = f.Close(); err != nil {
		return 1
	}
	envs, err := ReadDir(tmpdir)
	if err != nil {
		return 1
	}

	return RunCmd(args, envs)
}