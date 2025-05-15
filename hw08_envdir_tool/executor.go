package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdArgs []string, env Environment) (returnCode int) {
	err := prepareEnvironment(env)
	if err != nil {
		return -1
	}

	command := cmdArgs[0]
	args := cmdArgs[1:]
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return -1
	}

	return cmd.ProcessState.ExitCode()
}

func prepareEnvironment(env Environment) error {
	for envVar, envValue := range env {
		if envValue.NeedRemove {
			err := os.Unsetenv(envVar)
			if err != nil {
				return err
			}
		} else {
			err := os.Setenv(envVar, envValue.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
