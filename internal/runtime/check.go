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

type RuntimeStatus struct {
    NeedsUpdate     bool
    CurrentVersion  string
    TargetVersion   string
    PackageName     string
}

func CheckRuntime() (*RuntimeStatus, error) {
    fmt.Println("Checking runtime")

    runtimePath := filepath.Join(constants.ValetBasePath, constants.RuntimeFileName)
    versionPath := filepath.Join(constants.ValetVenvPath, constants.VersionFileName)

    for _, path := range []string{runtimePath, versionPath} {
        fmt.Printf("Check if '%s' exists\n", path)
        if _, err := os.Stat(path); err != nil {
            if os.IsNotExist(err) {
                return nil, fmt.Errorf("file %s does not exist", path)
            }
            return nil, fmt.Errorf("error checking file %s: %w", path, err)
        }
    }

    fmt.Println()

    installedRuntimeVersion, err := os.ReadFile(runtimePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read installed runtime version: %w", err)
    }

    targetRuntimeVersion, err := os.ReadFile(versionPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read target runtime version: %w", err)
    }

    osName, osCodename, err := CheckOSVersion()
    if err != nil {
        return nil, err
    }

    arch := getArchitecture()
    packageName := BuildPackageName(osName, osCodename, arch)

    status := &RuntimeStatus{
        CurrentVersion: strings.TrimSpace(string(installedRuntimeVersion)),
        TargetVersion:  strings.TrimSpace(string(targetRuntimeVersion)),
        PackageName:    packageName,
        NeedsUpdate:    strings.TrimSpace(string(installedRuntimeVersion)) != strings.TrimSpace(string(targetRuntimeVersion)),
    }

    fmt.Printf("valet-sh: Runtime version: %s\n", status.CurrentVersion)
    fmt.Printf("valet-sh-venv: Runtime Version: %s\n", status.TargetVersion)
    fmt.Printf("OS: %s\n", osName)
    fmt.Printf("Codename: %s\n", osCodename)
    fmt.Printf("Architecture: %s\n", arch)
    fmt.Printf("Package Name: %s\n", status.PackageName)

    if status.NeedsUpdate {
        fmt.Println("\nVenv runtime version is different from the installed runtime version\n")
        fmt.Printf("- new runtime version: %s is required\n", status.CurrentVersion)
    } else {
        fmt.Println("Runtime version is up to date")
    }

    return status, nil
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

func BuildPackageName(osName, osCodename, arch string) string {
    if osName == "ubuntu" {
        return strings.ToLower(osName) + "_" + strings.ToLower(osCodename) + "-" + string(arch)
    }
    return strings.ToLower(osName) + "-" + string(arch)
}
