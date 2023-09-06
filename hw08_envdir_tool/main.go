package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrTooFewParameters = errors.New("too few parameters")
	ErrIsNotDir         = errors.New("first param is not dir")
	ErrUnsupportableDir = errors.New("unsupportable dir")
	ErrIsNotExist       = errors.New("is not exist")
)

const (
	SuccessExitCode   = 0
	UnsuccessExitCode = 1
)

func errCode(err error) int {
	fmt.Println(err)
	return UnsuccessExitCode
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(ErrTooFewParameters)
		os.Exit(UnsuccessExitCode)
	}
	dir := os.Args[1]

	envs, err := ReadDir(dir)
	if err != nil {
		os.Exit(errCode(err))
	}
	os.Exit(RunCmd(os.Args[2:], envs))
}
