package main

import (
	"bufio"
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

func readEnvValue(fpath string) (*EnvValue, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if scanner := bufio.NewScanner(file); scanner.Scan() {
		str := scanner.Text()
		str = strings.TrimRight(str, " \t")
		str = strings.ReplaceAll(str, string([]byte{0x00}), "\n")

		return &EnvValue{Value: str, NeedRemove: false}, nil
	}

	return &EnvValue{NeedRemove: true}, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirStat, err := os.Stat(dir)
	if err != nil {
		return nil, ErrIsNotExist
	}
	if !dirStat.IsDir() {
		return nil, ErrIsNotDir
	}

	envs := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrUnsupportableDir
	}

	for _, file := range files {
		fpath := filepath.Join(dir, file.Name())

		fpStat, err := os.Stat(fpath)
		if err != nil {
			return nil, err
		}

		if fpStat.Mode().IsRegular() && !strings.Contains(file.Name(), "=") {
			env, err := readEnvValue(fpath)
			if err != nil {
				return nil, err
			}
			envs[file.Name()] = *env
		}
	}

	return envs, nil
}
