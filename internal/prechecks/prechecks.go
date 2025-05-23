package prechecks

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/valet-sh/valet-sh-installer/internal/utils"

	"github.com/valet-sh/valet-sh-installer/constants"
)

func CheckForValet() error {
	_, err := os.Stat(constants.VshPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("valet-sh does not exists, please run `valet-sh-installer install` first")
	}

	return nil
}

func CheckForEtcDirectory() error {
	if _, err := os.Stat(constants.VshEtcPath); os.IsNotExist(err) {
		err := os.MkdirAll(constants.VshEtcPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create etc directory: %w", err)
		}
	}
	return nil
}

func CheckForValetReleaseChannelFile() error {
	ReleaseChannelFilePath := filepath.Join(constants.VshEtcPath, constants.ReleaseChannelFile)

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

func GetCurrentUser() (string, error) {
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		return "", fmt.Errorf("failed to get current user")
	}
	return currentUser, nil
}

func CheckNotRoot() error {
	currentUser, err := user.Current()
	if err != nil {
		utils.Println("Error determining current user:", err)
		os.Exit(1)
	}

	if currentUser.Uid == "0" {
		color.Red.Println("This application should not be run with sudo or as root. Please run as a regular user.")

		os.Exit(1)
	}

	return nil
}

func RemoveOldCollectionDir() error {
	collectionDir := filepath.Join(constants.VshBasePath, "collections")
	if _, err := os.Stat(collectionDir); os.IsNotExist(err) {
		return nil
	}

	err := os.RemoveAll(collectionDir)
	if err != nil {
		return fmt.Errorf("failed to remove old collection directory: %w", err)
	}

	return nil
}
