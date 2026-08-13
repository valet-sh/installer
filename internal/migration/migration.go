package migration

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/valet-sh/valet-sh-installer/constants"
	"github.com/valet-sh/valet-sh-installer/internal/git"
	"github.com/valet-sh/valet-sh-installer/internal/prechecks"
	"github.com/valet-sh/valet-sh-installer/internal/utils"
)

const (
	apiTimeout = 3 * time.Second
)

func Cli() error {
	if _, err := os.Stat(constants.VshOldCliPath); err == nil {
		if err := os.Remove(constants.VshOldCliPath); err != nil {
			return fmt.Errorf("failed to remove old valet.sh CLI: %w", err)
		}
	}

	latestCliVersion, err := git.FetchLatestCliTag(apiTimeout)
	if err != nil {
		return fmt.Errorf("failed to fetch latest CLI version: %w", err)
	}

	if err := getLatestCliBinary(latestCliVersion); err != nil {
		return fmt.Errorf("failed to download latest binary: %w", err)
	}

	return nil
}

func getLatestCliBinary(version string) error {
	binaryName := getBinaryName()
	if binaryName == "" {
		return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	binaryURL := fmt.Sprintf(
		"https://github.com/valet-sh/go-cli/releases/download/%s/%s",
		version,
		binaryName,
	)
	client := &http.Client{Timeout: apiTimeout}
	resp, err := client.Get(binaryURL)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download binary: HTTP %d", resp.StatusCode)
	}
	if err := utils.RunCommand("sudo", []string{"mkdir", "-p", constants.VshCliPath}); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	currentUser, err := prechecks.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	if err := utils.RunCommand("sudo", []string{"chown", currentUser, constants.VshCliPath}); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	tempFile, err := os.CreateTemp(constants.VshCliPath, "valet-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write binary: %w", err)
	}
	tempFile.Close()

	if err := os.Rename(tempFile.Name(), constants.VshCliBinaryPath); err != nil {
		return fmt.Errorf("failed to move binary: %w", err)
	}

	if err := os.Chmod(constants.VshCliBinaryPath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}
	if err := UpdateSymlink(constants.VshCliBinaryPath); err != nil {
		return fmt.Errorf("failed to update symlink: %w", err)
	}
	return nil
}

func getBinaryName() string {
	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return "valet-darwin-amd64"
		case "arm64":
			return "valet-darwin-arm64"
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "valet-linux-amd64"
		case "arm64":
			return "valet-linux-arm64"
		}
	}
	return ""
}

func UpdateSymlink(newBinaryPath string) error {
	symlinkPath := "/usr/local/bin/valet.sh"

	if err := utils.RunCommand("sudo", []string{"ln", "-sf", newBinaryPath, symlinkPath}); err != nil {
		return fmt.Errorf("failed to update symlink: %w", err)
	}

	return nil
}
