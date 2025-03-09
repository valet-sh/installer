package setup

import (
    "fmt"
    "os"
    "strings"
    goruntime "runtime"

    "github.com/valet-sh/valet-sh-installer/constants"
    "github.com/valet-sh/valet-sh-installer/internal/utils"
    "github.com/valet-sh/valet-sh-installer/internal/runtime"
)

func SetupPythonEnvironment(logFile *os.File) error {
    pythonBin := "python3"
    if goruntime.GOOS == "darwin" {
        arch := runtime.GetArchitecture()
        if strings.HasPrefix(arch, "arm") {
            pythonBin = "/opt/homebrew/opt/python@3.10/bin/python3.10"
        } else {
            pythonBin = "/usr/local/opt/python@3.10/bin/python3.10"
        }
    }

    if _, err := os.Stat(constants.VshVenvPath); err == nil {
        fmt.Println("Removing existing virtual environment...")
        if err := os.RemoveAll(constants.VshVenvPath); err != nil {
            return fmt.Errorf("failed to remove existing virtual environment: %w", err)
        }
    }

    if _, err := os.Stat(constants.VshVenvPath); os.IsNotExist(err) {
        fmt.Printf("Creating virtual environment at %s\n", constants.VshVenvPath)
        venvArgs := []string{"-m", "venv"}
        if goruntime.GOOS == "linux" {
            venvArgs = append(venvArgs, "--system-site-packages")
        }
        venvArgs = append(venvArgs, constants.VshVenvPath)
        if err := utils.RunCommand(pythonBin, venvArgs, logFile); err != nil {
            return fmt.Errorf("failed to create virtual environment: %w", err)
        }
    }
    return nil
}
