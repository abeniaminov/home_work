package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("path not exists", func(t *testing.T) {
		_, err := ReadDir("./testdata/noexists")
		require.ErrorIs(t, err, ErrIsNotExist, "actual err - %v", err)
	})
	t.Run("path is not dir", func(t *testing.T) {
		_, err := ReadDir("./testdata/env/BAR")
		require.ErrorIs(t, err, ErrIsNotDir, "actual err - %v", err)
	})
	t.Run("fail multiline test", func(t *testing.T) {
		content := []byte("foo\nbar\nerl")
		v, err := prepareReadEnv(content)
		if err != nil {
			t.Fatal(err)
		}
		expected := &EnvValue{Value: "foo", NeedRemove: false}
		require.Equal(t, expected, v)
	})
	t.Run("fail 0x00 test", func(t *testing.T) {
		content := []byte("foo" + string([]byte{0x00}) + "erl")
		v, err := prepareReadEnv(content)
		if err != nil {
			t.Fatal(err)
		}
		expected := &EnvValue{Value: "foo\nerl", NeedRemove: false}
		require.Equal(t, expected, v)
	})
	t.Run("fail tabs test", func(t *testing.T) {
		content := []byte("foo\terl\t\t")
		v, err := prepareReadEnv(content)
		if err != nil {
			t.Fatal(err)
		}
		expected := &EnvValue{Value: "foo\terl", NeedRemove: false}
		require.Equal(t, expected, v)
	})
	t.Run("fail spaces test", func(t *testing.T) {
		content := []byte("foo erl  ")
		v, err := prepareReadEnv(content)
		if err != nil {
			t.Fatal(err)
		}
		expected := &EnvValue{Value: "foo erl", NeedRemove: false}
		require.Equal(t, expected, v)
	})
	t.Run("fail combination test", func(t *testing.T) {
		content := []byte("foo erl \t" + string([]byte{0x00}) + "boo  \t\t")
		v, err := prepareReadEnv(content)
		if err != nil {
			t.Fatal(err)
		}
		expected := &EnvValue{Value: "foo erl \t\nboo", NeedRemove: false}
		require.Equal(t, expected, v)
	})
	t.Run("fail need remove test", func(t *testing.T) {
		v, err := prepareReadEnv(nil)
		if err != nil {
			t.Fatal(err)
		}
		expected := &EnvValue{Value: "", NeedRemove: true}
		require.Equal(t, expected, v)
	})
}

func prepareReadEnv(fileContent []byte) (*EnvValue, error) {
	tmpfile, err := os.CreateTemp("/tmp", "testfile.")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())

	if fileContent != nil {
		if _, err := tmpfile.Write(fileContent); err != nil {
			return nil, err
		}
	}

	if err := tmpfile.Close(); err != nil {
		return nil, err
	}

	return readEnvValue(tmpfile.Name())
}
