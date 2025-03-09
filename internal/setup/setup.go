package setup

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/valet-sh/valet-sh-installer/constants"
    "github.com/valet-sh/valet-sh-installer/internal/utils"
    "github.com/valet-sh/valet-sh-installer/internal/git"
)

func PrepareLogFile() (*os.File, error) {
    setupLogFile, err := os.Create(constants.VshInstallLog)
    if err != nil {
        return nil, fmt.Errorf("failed to create setup log file: %w", err)
    }
    return setupLogFile, nil
}

func InstallVshDependencies(logFile *os.File) error {
    venvPip := filepath.Join(constants.VshVenvPath, "bin", "pip3")

    if err := utils.RunCommand(venvPip, []string{"uninstall", "-y", "valet-sh"}, logFile); err != nil {
        return fmt.Errorf("failed to uninstall existing valet-sh: %w", err)
    }

    if err := utils.RunCommand(venvPip, []string{"install", "--upgrade", "pip", "setuptools==60.8.2", "wheel==0.37.1"}, logFile); err != nil {
        return fmt.Errorf("failed to install base packages: %w", err)
    }

    fmt.Println("Installing ansible and dependencies...")
    if err := utils.RunCommand(venvPip, []string{"install", "-r", filepath.Join(constants.VshBasePath, "requirements.txt")}, logFile); err != nil {
        return fmt.Errorf("failed to install valet-sh requirements: %w", err)
    }

    requirementsYml := filepath.Join(constants.VshBasePath, "requirements.yml")
    if _, err := os.Stat(requirementsYml); err == nil {
        fmt.Println("Installing ansible collections...")
        collectionsPath := filepath.Join(constants.VshBasePath, "collections")
        if err := utils.RunCommand("ansible-galaxy", []string{
            "collection", "install",
            "-r", requirementsYml,
            "-p", collectionsPath,
        }, logFile); err != nil {
            return fmt.Errorf("failed to install ansible collections: %w", err)
        }
    }

    return nil
}

func PrepareVshDirectory(vshUser, vshGroup string, logFile *os.File) error {
    if _, err := os.Stat(constants.VshPath); os.IsNotExist(err) {
        if err := utils.RunCommand("sudo", []string{"mkdir", "-p", constants.VshPath}, logFile); err != nil {
            return fmt.Errorf("failed to create install directory: %w", err)
        }
        if err := utils.RunCommand("sudo", []string{"chown", fmt.Sprintf("%s:%s", vshUser, vshGroup), constants.VshPath}, logFile); err != nil {
            return fmt.Errorf("failed to set permissions on install directory: %w", err)
        }
    }
    return nil
}


func RemoveVshAnsibleFactsFile() error {
    if _, err := os.Stat(constants.VshAnsibleFactsFile); err == nil {
        _ = os.Remove(constants.VshAnsibleFactsFile)
    }

    return nil
}

func RemoveVshRepository() error {
    if _, err := os.Stat(constants.VshBasePath); err == nil {
        fmt.Println("Removing existing repository...")
        if err := os.RemoveAll(constants.VshBasePath); err != nil {
            return fmt.Errorf("failed to remove existing repository: %w", err)
        }
    }

    return nil
}

func RemoveVshVenv() error {
    if _, err := os.Stat(constants.VshVenvPath); err == nil {
        fmt.Println("Removing existing virtual environment...")
        if err := os.RemoveAll(constants.VshVenvPath); err != nil {
            return fmt.Errorf("failed to remove existing virtual environment: %w", err)
        }
    }

    return nil
}

func SetupVshRepository(setupLogFile *os.File) error {
    if _, err := os.Stat(filepath.Join(constants.VshBasePath, ".git")); os.IsNotExist(err) {
        if err := git.CloneRepository(constants.VshGithubRepoUrl, constants.VshBasePath); err != nil {
            return fmt.Errorf("failed to clone repository: %w", err)
        }
    } else {
        if err := utils.RunCommand("git", []string{"-C", constants.VshBasePath, "pull"}, setupLogFile); err != nil {
            return fmt.Errorf("failed to update repository: %w", err)
        }
    }

    return nil
}
