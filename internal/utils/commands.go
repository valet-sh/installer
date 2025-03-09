package utils

import (
    "os"
    "os/exec"
)

func RunCommand(command string, args []string, logFile *os.File) error {
    cmd := exec.Command(command, args...)
    cmd.Stdout = logFile
    cmd.Stderr = logFile
    return cmd.Run()
}

func PathExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
