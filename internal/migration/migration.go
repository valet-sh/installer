package migration

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/valet-sh/valet-sh-installer/constants"
	"github.com/valet-sh/valet-sh-installer/internal/git"
	"github.com/valet-sh/valet-sh-installer/internal/prechecks"
	"github.com/valet-sh/valet-sh-installer/internal/utils"
)

const (
	apiTimeout      = 3 * time.Second
	downloadTimeout = 2 * time.Minute
)

func Cli() error {
	latestCliVersion, err := git.FetchLatestCliTag(apiTimeout)
	if err != nil {
		return fmt.Errorf("failed to fetch latest CLI version: %w", err)
	}

	if err := getLatestCliBinary(latestCliVersion); err != nil {
		return fmt.Errorf("failed to download latest binary: %w", err)
	}

	if _, err := os.Stat(constants.VshOldCliPath); err == nil {
		if err := os.Remove(constants.VshOldCliPath); err != nil {
			return fmt.Errorf("failed to remove old valet.sh CLI: %w", err)
		}
	}

	return nil
}

func RunInstall(migrationConfirmed bool) error {
	utils.Println("Running install on the valet.sh Go-CLI")

	var env []string
	if migrationConfirmed {
		env = []string{"VALET_MIGRATE=1"}
	}

	if err := utils.RunInteractiveCommandWithEnv(constants.VshCliBinaryPath, []string{"install"}, env); err != nil {
		return fmt.Errorf("failed to run install: %w", err)
	}

	return nil
}

func getLatestCliBinary(version string) error {
	binaryName := getBinaryName()
	if binaryName == "" {
		return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	releaseURL := fmt.Sprintf("https://github.com/valet-sh/go-cli/releases/download/%s", version)

	client := &http.Client{Timeout: downloadTimeout}
	resp, err := client.Get(fmt.Sprintf("%s/%s", releaseURL, binaryName))
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

	hasher := sha256.New()
	if _, err := io.Copy(io.MultiWriter(tempFile, hasher), resp.Body); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write binary: %w", err)
	}
	tempFile.Close()

	expectedChecksum, err := fetchExpectedChecksum(releaseURL, binaryName)
	if err != nil {
		return fmt.Errorf("failed to fetch checksum: %w", err)
	}
	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch for downloaded CLI binary: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	if err := os.Chmod(tempFile.Name(), 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	if err := os.Rename(tempFile.Name(), constants.VshCliBinaryPath); err != nil {
		return fmt.Errorf("failed to move binary: %w", err)
	}

	if err := UpdateSymlink(constants.VshCliBinaryPath); err != nil {
		return fmt.Errorf("failed to update symlink: %w", err)
	}
	return nil
}

func fetchExpectedChecksum(releaseURL string, binaryName string) (string, error) {
	client := &http.Client{Timeout: downloadTimeout}
	resp, err := client.Get(releaseURL + "/checksums.txt")
	if err != nil {
		return "", fmt.Errorf("failed to download checksums: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download checksums: HTTP %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[1] == binaryName {
			return fields[0], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("checksum for %s not found in checksums.txt", binaryName)
}

func getBinaryName() string {
	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return "valet-darwin-arm64"
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "valet-linux-amd64"
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
