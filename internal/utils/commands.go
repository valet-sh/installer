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

func PathExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
