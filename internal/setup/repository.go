package setup

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/valet-sh/valet-sh-installer/constants"
    "github.com/valet-sh/valet-sh-installer/internal/git"
)

func SetupRepository(logFile *os.File) error {
    if _, err := os.Stat(constants.VshBasePath); err == nil {
        fmt.Println("Removing existing repository...")
        if err := os.RemoveAll(constants.VshBasePath); err != nil {
            return fmt.Errorf("failed to remove existing repository: %w", err)
        }
    }

    if _, err := os.Stat(filepath.Join(constants.VshBasePath, ".git")); os.IsNotExist(err) {
        if err := git.CloneRepository(constants.VshGithubRepoUrl, constants.VshBasePath); err != nil {
            return fmt.Errorf("failed to clone repository: %w", err)
        }
    }
    return nil
}
