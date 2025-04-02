package runtime

import (
	"bufio"
	"fmt"
	"github.com/valet-sh/valet-sh-installer/internal/utils"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/valet-sh/valet-sh-installer/constants"
)

type RuntimeStatus struct {
	NeedsUpdate    bool
	PackageChanged bool
	CurrentVersion string
	TargetVersion  string
	CurrentPackage string
	PackageName    string
}

func CheckRuntime() (*RuntimeStatus, error) {
	utils.Println("Checking runtime...")

	runtimeVersionPaths := map[string]string{
		"runtime": filepath.Join(constants.VshBasePath, constants.RuntimeFileName),
		"version": filepath.Join(constants.VshVenvPath, constants.VersionFileName),
	}
	var installedRuntimeVersion, targetRuntimeVersion []byte

	for fileType, path := range runtimeVersionPaths {
		utils.Printf("Check if '%s' exists\n", path)

		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				if fileType == "runtime" {

					return nil, fmt.Errorf("runtime file %s does not exist", path)
				} else if fileType == "version" {
					utils.Println("Version file missing, using version from valet-sh")
					targetRuntimeVersion = installedRuntimeVersion
					continue
				}
			}
			return nil, fmt.Errorf("error checking file %s: %w", path, err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", fileType, err)
		}
		utils.Printf("Read content from %s: %s\n", path, string(content))

		if fileType == "runtime" {
			installedRuntimeVersion = content
		} else if fileType == "version" {
			targetRuntimeVersion = content
		}
	}

	if len(installedRuntimeVersion) == 0 {
		return nil, fmt.Errorf("failed to determine installed runtime version")
	}
	if len(targetRuntimeVersion) == 0 {
		return nil, fmt.Errorf("failed to determine target runtime version")
	}

	osName, osCodename, err := CheckOSVersion()
	if err != nil {
		return nil, err
	}

	arch := GetArchitecture()
	packageName := BuildPackageName(osName, osCodename, arch)

	currentVersion := strings.TrimSpace(string(installedRuntimeVersion))

	targetParts := strings.Split(strings.TrimSpace(string(targetRuntimeVersion)), "-")
	targetVersion := targetParts[len(targetParts)-1]
	targetPackage := strings.Join(targetParts[:len(targetParts)-1], "-")

	status := &RuntimeStatus{
		CurrentVersion: currentVersion,
		TargetVersion:  targetVersion,
		CurrentPackage: packageName,
		PackageName:    targetPackage,
		NeedsUpdate:    currentVersion != targetVersion,
		PackageChanged: packageName != targetPackage,
	}

	utils.Printf("valet-sh: Runtime version: %s\n", status.CurrentVersion)
	utils.Printf("valet-sh-venv: Runtime Version: %s\n", status.TargetVersion)
	utils.Printf("OS: %s\n", osName)
	utils.Printf("Codename: %s\n", osCodename)
	utils.Printf("Architecture: %s\n", arch)
	utils.Printf("Current Package: %s\n", status.CurrentPackage)
	utils.Printf("Target Package: %s\n", status.PackageName)

	if status.NeedsUpdate && status.PackageChanged {
		utils.Println("\nBoth runtime version and package need to be updated")
		utils.Printf("Version update from %s to %s and package change from %s to %s required\n",
			status.TargetVersion, status.CurrentVersion, status.PackageName, status.CurrentPackage)
	} else if status.NeedsUpdate {
		utils.Println("\nVenv runtime version is different from the installed runtime version")
		utils.Printf("New runtime version %s is required\n", status.CurrentVersion)
	} else if status.PackageChanged {
		utils.Println("\nRuntime package needs to be updated")
		utils.Printf("Package change from %s to %s is required\n", status.PackageName, status.CurrentPackage)
	} else {
		utils.Println("Runtime version and package are up to date\n")
	}

	return status, nil
}

func GetArchitecture() string {
	arch := runtime.GOARCH

	archMapping := map[string]string{
		"amd64": "x86_64",
		"386":   "i386",
		"arm64": "aarch64",
	}

	if mappedArch, exists := archMapping[arch]; exists {
		return mappedArch
	}

	return arch
}

var supportedDistros = map[string]bool{
	"ubuntu":     true,
	"linux mint": true,
}

func isSupported(distro string) bool {
	_, ok := supportedDistros[distro]
	return ok
}

func CheckOSVersion() (string, string, error) {
	switch runtime.GOOS {
	case "linux":
		osInfo, err := parseOSRelease()
		if err != nil {
			return "", "", fmt.Errorf("failed to check OS info: %w", err)
		}

		if isSupported(strings.ToLower(osInfo.name)) {
			return strings.ToLower(osInfo.name), osInfo.version, nil
		}

		return "", "", fmt.Errorf("unsupported Linux distribution: %s", osInfo.name)
	case "darwin":
		return "macos", "", nil
	default:
		return "", "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

type osInfo struct {
	name    string
	version string
}

func parseOSRelease() (osInfo, error) {
	const osReleasePath = "/etc/os-release"

	content, err := os.ReadFile(osReleasePath)
	if err != nil {
		return osInfo{}, fmt.Errorf("failed to read %s: %w", osReleasePath, err)
	}

	var info osInfo
	var ubuntuCodename string

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "NAME="):
			info.name = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"'")
		case strings.HasPrefix(line, "VERSION_CODENAME="):
			info.version = strings.Trim(strings.TrimPrefix(line, "VERSION_CODENAME="), "\"'")
		case strings.HasPrefix(line, "UBUNTU_CODENAME="):
			ubuntuCodename = strings.Trim(strings.TrimPrefix(line, "UBUNTU_CODENAME="), "\"'")
		}
	}

	if strings.ToLower(info.name) == "" {
		return osInfo{}, fmt.Errorf("OS name not found in %s", osReleasePath)
	}

	if strings.ToLower(info.name) == "linux mint" && ubuntuCodename != "" {
		info.version = ubuntuCodename
		info.name = "ubuntu"
	}

	return info, nil
}

func BuildPackageName(osName, osCodename, arch string) string {
	if osName == "ubuntu" {
		return strings.ToLower(osName) + "_" + strings.ToLower(osCodename) + "-" + string(arch)
	}
	return strings.ToLower(osName) + "-" + string(arch)
}
