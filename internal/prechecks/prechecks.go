package prechecks

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/valet-sh/valet-sh-installer/constants"
)

func CheckForValet() error {
    _, err := os.Stat(constants.ValetPath)
    if os.IsNotExist(err) {
        return fmt.Errorf("Valet-sh does not exists, please run `valet-sh-installer install first`")
    }

    return nil
}

func CheckForEtcDirectory() error {
    if _, err := os.Stat(constants.ValetEtcPath); os.IsNotExist(err) {
        err := os.MkdirAll(constants.ValetEtcPath, 0755)
        if err != nil {
            return fmt.Errorf("failed to create etc directory: %w", err)
        }
    }
    return nil
}

func CheckForValetMajorReleaseFile() error {
    majorReleaseFilePath := filepath.Join(constants.ValetEtcPath, constants.MajorReleaseFile)

    _, err := os.Stat(majorReleaseFilePath)
    if os.IsNotExist(err) {
        _, err := os.Create(majorReleaseFilePath)
        if err != nil {
            return fmt.Errorf("failed to create major release file: %w", err)
        }
        majorVersion := constants.ValetMajorVersion
        err = os.WriteFile(majorReleaseFilePath, []byte(majorVersion), 0644)
        if err != nil {
            return fmt.Errorf("failed to write major release file: %w", err)
        }
    }

    return nil
}
