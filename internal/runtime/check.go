package runtime

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/valet-sh/valet-sh-updater/constants"
)

func CheckRuntime() error {
    fmt.Println("Checking runtime")

    runtimePath := filepath.Join(constants.ValetBasePath, constants.RuntimeFileName)
    versionPath := filepath.Join(constants.ValetBasePath, constants.VersionFileName)

    for _, path := range []string{runtimePath, versionPath} {
        fmt.Printf("Check if '%s' exists\n", path)
        if _, err := os.Stat(path); err != nil {
            if os.IsNotExist(err) {
                return fmt.Errorf("file %s does not exist", path)
            }
            return fmt.Errorf("error checking file %s: %w", path, err)
        }
    }

    runtimeVersion, err := os.ReadFile(runtimePath)
    if err != nil {
        return fmt.Errorf("failed to read runtime version: %w", err)
    }

    version, err := os.ReadFile(versionPath)
    if err != nil {
        return fmt.Errorf("failed to read version: %w", err)
    }

    fmt.Printf("Runtime version: %s\n", runtimeVersion)
    fmt.Printf("Version: %s\n", version)

    return nil
}
