package prechecks

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/shirou/gopsutil/v4/host"

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

type CheckResult struct {
	Passed  bool
	Message string
	Details string
}

func Requirements3xCheck() error {
	checks := []CheckResult{
		CheckOSVersionRequirements(),
		CheckSystemArchitecture(),
	}

	var failures []string
	for _, result := range checks {
		if !result.Passed {
			failures = append(failures, fmt.Sprintf("%s\n  %s", result.Message, result.Details))
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf("%s", strings.Join(failures, "\n"))
	}

	return nil
}

func CheckOSVersionRequirements() CheckResult {
	osInfo, err := host.Info()
	if err != nil {
		return CheckResult{
			Passed:  false,
			Message: "Failed to retrieve OS information",
			Details: err.Error(),
		}
	}

	currentVersion, err := version.NewVersion(osInfo.PlatformVersion)
	if err != nil {
		return CheckResult{
			Passed:  false,
			Message: "Failed to parse current OS version",
			Details: fmt.Sprintf("Version: %s, Error: %v", osInfo.PlatformVersion, err),
		}
	}

	minVersion, minVersionStr := getMinOSVersion(runtime.GOOS)
	if minVersion == nil {
		return CheckResult{
			Passed:  false,
			Message: "Unsupported operating system",
			Details: fmt.Sprintf("OS: %s. Only macOS and Linux are supported", runtime.GOOS),
		}
	}

	if !currentVersion.GreaterThanOrEqual(minVersion) {
		return CheckResult{
			Passed:  false,
			Message: "Unsupported OS version",
			Details: fmt.Sprintf("Current: %s, Required: %s or higher", currentVersion, minVersionStr),
		}
	}

	return CheckResult{Passed: true}
}

func CheckSystemArchitecture() CheckResult {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	supportedArch, err := getSupportedArchitecture(osName)
	if err != nil {
		return CheckResult{
			Passed:  false,
			Message: "Unsupported operating system",
			Details: err.Error(),
		}
	}

	if arch != supportedArch {
		return CheckResult{
			Passed:  false,
			Message: "Unsupported architecture",
			Details: fmt.Sprintf("Current: %s, Required: %s on %s", arch, supportedArch, osName),
		}
	}

	return CheckResult{Passed: true}
}

func getMinOSVersion(osName string) (*version.Version, string) {
	var minVersionStr string

	switch osName {
	case "darwin":
		minVersionStr = constants.Vsh3xMinMacOSVersion
	case "linux":
		minVersionStr = constants.Vsh3xMinLinuxVersion
	default:
		return nil, ""
	}

	minVersion, err := version.NewVersion(minVersionStr)
	if err != nil {
		return nil, minVersionStr
	}

	return minVersion, minVersionStr
}

func getSupportedArchitecture(osName string) (string, error) {
	switch osName {
	case "darwin":
		return "arm64", nil
	case "linux":
		return "amd64", nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", osName)
	}
}

var migrationText = `
┌────────────────────────────────────────────────────────────────────┐
│ valet.sh — Migration Required (v2.x Detected)                      │
└────────────────────────────────────────────────────────────────────┘

  WARNING: Existing services will be uninstalled and NO database
  dumps will be generated automatically.

  Please backup all critical database and application data.

  Post-Migration Environment:
  • Base Services Installed: Nginx, Dnsmasq, Mailpit, Container Runtime
    (Podman/Apple Container)
  • Optional Services Removed: PHP, MySQL, MariaDB, RabbitMQ, etc.
    (Must be re-installed on demand)

  For migration support and manual backup steps, visit:
  https://valet.sh/3.x/how-to-articles/migrating-from-2.x-to-3.x


`

var ErrMigrationCanceled = errors.New("migration canceled")

func CheckMigration() (bool, error) {
	_, err := os.Stat(constants.VshServiceFile)
	serviceExists := err == nil
	_, err = os.Stat(constants.VshBundlesFile)
	bundleExists := err == nil

	if serviceExists || (serviceExists && bundleExists) {
		fmt.Printf("\033[0;33m\033[1m%s\033[0;0m\n", migrationText)

		utils.Printf("Type 'migrate' to proceed, or press [Enter] to cancel/skip: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("failed to read input: %w", err)
		}

		trimmedInput := strings.TrimSpace(input)

		if trimmedInput == "" || trimmedInput != "migrate" {
			fmt.Println("\n Migration canceled. No changes have been made.")
			return false, ErrMigrationCanceled
		}

		return true, nil
	}

	return false, nil
}
