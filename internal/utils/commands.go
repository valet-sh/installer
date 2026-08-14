package utils

import (
	"os"
	"os/exec"
)

func RunCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = LogFile
	cmd.Stderr = LogFile
	return cmd.Run()
}

func RunInteractiveCommand(command string, args []string) error {
	return RunInteractiveCommandWithEnv(command, args, nil)
}

func RunInteractiveCommandWithEnv(command string, args []string, extraEnv []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if len(extraEnv) > 0 {
		cmd.Env = append(os.Environ(), extraEnv...)
	}
	return cmd.Run()
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
