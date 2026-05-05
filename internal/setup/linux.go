package setup

import (
    "fmt"

    "github.com/valet-sh/valet-sh-installer/internal/utils"
)

func InstallLinuxDependencies() error {
    if err := utils.RunCommand("sudo", []string{"apt-get", "update"}); err != nil {
        return fmt.Errorf("failed to update apt: %w", err)
    }
    if err := utils.RunCommand("sudo", []string{"apt-get", "install", "-y", "git", "python3", "python3-venv"}); err != nil {
        return fmt.Errorf("failed to install dependencies: %w", err)
    }

    return nil
}
