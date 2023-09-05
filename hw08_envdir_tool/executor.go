package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for name, e := range env {
		if e.NeedRemove {
			if err := os.Unsetenv(name); err != nil {
				return errCode(err)
			}
		} else {
			if err := os.Setenv(name, e.Value); err != nil {
				return errCode(err)
			}
		}
	}
	args := cmdExpEnv(cmd[1:])
	eCmd := exec.Command(cmd[0], args...) // #nosec G204
	eCmd.Stdout = os.Stdout
	eCmd.Stdin = os.Stdin
	eCmd.Stderr = os.Stderr

	if err := eCmd.Run(); err != nil {
		var exerr *exec.ExitError
		if errors.As(err, &exerr) {
			return exerr.ExitCode()
		}
		return errCode(err)
	}

	return SuccessExitCode
}

func cmdExpEnv(cmd []string) []string {
	for i, c := range cmd {
		cmd[i] = os.ExpandEnv(c)
	}
	return cmd
}
