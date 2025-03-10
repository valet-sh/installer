package setup

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/valet-sh/valet-sh-installer/constants"
    "github.com/valet-sh/valet-sh-installer/internal/utils"
)

func CreateSymlinks(vshUser string, logFile *os.File) error {
    localBinPath := "/usr/local/bin"
    if _, err := os.Stat(localBinPath); os.IsNotExist(err) {
        if err := utils.RunCommand("sudo", []string{"mkdir", "-p", localBinPath}, logFile); err != nil {
            return fmt.Errorf("failed to create local bin directory: %w", err)
        }
        if err := utils.RunCommand("sudo", []string{"chown", vshUser, localBinPath}, logFile); err != nil {
            return fmt.Errorf("failed to set permissions on /usr/local/bin directory: %w", err)
        }
    }

    vshBinPath := filepath.Join(constants.VshVenvPath, "bin", "valet.sh")
    if err := utils.RunCommand("sudo", []string{"ln", "-sf", vshBinPath, "/usr/local/bin/valet.sh"}, logFile); err != nil {
        return fmt.Errorf("failed to create symlink: %w", err)
    }

    return nil
}
