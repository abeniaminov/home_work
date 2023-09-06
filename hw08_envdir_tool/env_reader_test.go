package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected *EnvValue
	}{
		{
			name:     "fail multiline test",
			content:  []byte("foo\nbar\nerl"),
			expected: &EnvValue{Value: "foo", NeedRemove: false},
		},
		{
			name:     "fail 0x00 test",
			content:  []byte("foo" + string([]byte{0x00}) + "erl"),
			expected: &EnvValue{Value: "foo\nerl", NeedRemove: false},
		},
		{
			name:     "fail tabs test",
			content:  []byte("foo\terl\t\t"),
			expected: &EnvValue{Value: "foo\terl", NeedRemove: false},
		},
		{
			name:     "fail spaces test",
			content:  []byte("foo erl  "),
			expected: &EnvValue{Value: "foo erl", NeedRemove: false},
		},
		{
			name:     "fail combination test",
			content:  []byte("foo erl \t" + string([]byte{0x00}) + "boo  \t\t"),
			expected: &EnvValue{Value: "foo erl \t\nboo", NeedRemove: false},
		},
		{
			name:     "fail need remove test",
			content:  nil,
			expected: &EnvValue{Value: "", NeedRemove: true},
		},
	}
	t.Run("path not exists", func(t *testing.T) {
		_, err := ReadDir("./testdata/noexists")
		require.ErrorIs(t, err, ErrIsNotExist, "actual err - %v", err)
	})
	t.Run("path is not dir", func(t *testing.T) {
		_, err := ReadDir("./testdata/env/BAR")
		require.ErrorIs(t, err, ErrIsNotDir, "actual err - %v", err)
	})
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v, err := prepareReadEnv(test.content)
			require.NoError(t, err, "actual err - %v", err)
			require.Equal(t, test.expected, v)
		})
	}
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
