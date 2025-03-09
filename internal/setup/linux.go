package setup

import (
    "fmt"
    "os"

    "github.com/valet-sh/valet-sh-installer/internal/utils"
)

func InstallLinuxDependencies(logFile *os.File) error {
    if err := utils.RunCommand("sudo", []string{"apt-get", "update"}, logFile); err != nil {
        return fmt.Errorf("failed to update apt: %w", err)
    }
    if err := utils.RunCommand("sudo", []string{"apt-get", "install", "-y", "git", "python3", "python3-venv"}, logFile); err != nil {
        return fmt.Errorf("failed to install dependencies: %w", err)
    }

    return nil
}
