package utils

import (
    "fmt"
    "os"
    "os/exec"
)

func RequestSudoAccess() error {
    cmd := exec.Command("sudo", "true")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    if err := cmd.Run(); err != nil {
        fmt.Println("Error with sudo:", err)
        return fmt.Errorf("failed to get sudo access: %w", err)
    }
    return nil
}
