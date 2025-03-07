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

func CheckForValetReleaseChannelFile() error {
    ReleaseChannelFilePath := filepath.Join(constants.ValetEtcPath, constants.ReleaseChannelFile)

    _, err := os.Stat(ReleaseChannelFilePath)
    if os.IsNotExist(err) {
        _, err := os.Create(ReleaseChannelFilePath)
        if err != nil {
            return fmt.Errorf("failed to create release channel file: %w", err)
        }
        releaseChannelStableVersion := constants.ValetStableVersion
        err = os.WriteFile(ReleaseChannelFilePath, []byte(releaseChannelStableVersion), 0644)
        if err != nil {
            return fmt.Errorf("failed to write release channel file: %w", err)
        }
    }

    return nil
}
