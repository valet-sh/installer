package runtime

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "runtime"
    "bufio"

    "github.com/valet-sh/valet-sh-installer/constants"
)

func CheckRuntime() error {
    fmt.Println("Checking runtime")

    runtimePath := filepath.Join(constants.ValetBasePath, constants.RuntimeFileName)
    versionPath := filepath.Join(constants.ValetVenvPath, constants.VersionFileName)

    for _, path := range []string{runtimePath, versionPath} {
        fmt.Printf("Check if '%s' exists\n", path)
        if _, err := os.Stat(path); err != nil {
            if os.IsNotExist(err) {
                return fmt.Errorf("file %s does not exist", path)
            }
            return fmt.Errorf("error checking file %s: %w", path, err)
        }
    }

    fmt.Println()

    installedRuntimeVersion, err := os.ReadFile(runtimePath)
    if err != nil {
        return fmt.Errorf("failed to read installed runtime version: %w", err)
    }

    targetRuntimeVersion, err := os.ReadFile(versionPath)
    if err != nil {
        return fmt.Errorf("failed to read target runtime version: %w", err)
    }


    osName, osCodename, err := CheckOSVersion()
    if err != nil {
        return err
    }

    arch := getArchitecture()

    var installed_runtime_specific_version string
    var target_runtime_specific_version string

    if osName == "ubuntu" {
        target_runtime_specific_version = (strings.ToLower(osName) + "_" + strings.ToLower(osCodename) + "-" + string(arch))
    } else {
        target_runtime_specific_version = (strings.ToLower(osName) + "-" + string(arch))
    }

    fmt.Printf("valet-sh: Runtime version: %s\n", installedRuntimeVersion)
    fmt.Printf("valet-sh-venv: Runtime Version: %s\n", targetRuntimeVersion)


    fmt.Printf("OS: %s\n", osName)
    fmt.Printf("Codename: %s\n", osCodename)
    fmt.Printf("Architecture: %s\n", arch)
    fmt.Printf("Package Name: %s\n", installed_runtime_specific_version)
    fmt.Printf("Target Package Name: %s\n", target_runtime_specific_version)

    if string(installedRuntimeVersion) != string(targetRuntimeVersion) {
        fmt.Println("\nVenv runtime version is different from the installed runtime version\n")
        fmt.Printf("- new runtime version: %s is required\n", installedRuntimeVersion)
        return UpdateRuntime(string(installedRuntimeVersion))
    }

    fmt.Println("Runtime version is up to date")

    return nil
}

func getArchitecture() string {
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

func UpdateRuntime(targetRuntimeVersion string) error {
    fmt.Printf("Updating runtime to version %s\n", targetRuntimeVersion)

    return nil
}
