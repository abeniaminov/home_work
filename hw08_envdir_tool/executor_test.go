package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name     string
		command  []string
		content  []byte
		exitCode int
		expected string
	}{
		{
			name:     "fail exec multiline test",
			command:  []string{"ls", "-a"},
			content:  []byte("foo\nbar\nerl"),
			exitCode: 0,
			expected: "foo",
		},
		{
			name:     "fail exec multiline param test",
			command:  []string{"ls", "${FOO}"},
			content:  []byte("-l\nbar\nerl"),
			exitCode: 0,
			expected: "-l",
		},
		{
			name:     "fail exec 0x00 test",
			command:  []string{"ls", "-a"},
			content:  []byte("foo" + string([]byte{0x00}) + "erl"),
			exitCode: 0,
			expected: "foo\nerl",
		},
		{
			name:     "fail exec tabs test",
			command:  []string{"ls", "-a"},
			content:  []byte("foo\terl\t\t"),
			exitCode: 0,
			expected: "foo\terl",
		},
		{
			name:     "fail exec spaces test",
			command:  []string{"ls", "-a"},
			content:  []byte("foo erl  "),
			exitCode: 0,
			expected: "foo erl",
		},
		{
			name:     "fail exec combination test",
			command:  []string{"ls", "-a"},
			content:  []byte("foo erl \t" + string([]byte{0x00}) + "boo  \t\t"),
			exitCode: 0,
			expected: "foo erl \t\nboo",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := prepareExecutor(test.command, test.content)
			require.Equal(t, test.exitCode, v)
			env, ok := os.LookupEnv("FOO")
			require.True(t, ok)
			require.Equal(t, env, test.expected)
		})
	}
	t.Run("fail exec need remove test", func(t *testing.T) {
		command := []string{"ls", "-l"}
		v := prepareExecutor(command, nil)
		exitCode := 0
		require.Equal(t, exitCode, v)
		_, ok := os.LookupEnv("FOO")
		require.False(t, ok)
	})
	t.Run("fail exec test", func(t *testing.T) {
		command := []string{"ls", "mura"}
		v := prepareExecutor(command, nil)
		require.GreaterOrEqual(t, v, 1)
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
