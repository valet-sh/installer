package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "net/http"
    "strings"

    "github.com/spf13/cobra"

    "github.com/valet-sh/valet-sh-installer/constants"
    "github.com/valet-sh/valet-sh-installer/internal/git"
    "github.com/valet-sh/valet-sh-installer/internal/runtime"
    "github.com/valet-sh/valet-sh-installer/internal/utils"
)

var updateCmd = &cobra.Command{
    Use:  "update",
    Short: "Update valet-sh to the latest version",
    Long: `Update valet-sh to the latest version`,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return runUpdate()
    },
}

func init() {
}

func runUpdate() error {
    repoPath := constants.ValetBasePath

    if err := checkIfRepoExists(repoPath); err != nil {
        return err
    }

    nextChannelEnabled := false
    enableNextFilePath := filepath.Join(constants.ValetEtcPath, constants.NextBranchFile)

    if _, err := os.Stat(enableNextFilePath); err == nil {
        nextChannelEnabled = true
    }

    if nextChannelEnabled {
        fmt.Println("Using next channel (development) for update")
        return updateNextBranch(repoPath)
    } else {
        fmt.Println("Using stable channel for update")
        return updateStableBranch(repoPath)
    }
}

func checkIfRepoExists(repoPath string) error {
    _, err := os.Stat(repoPath)
    if os.IsNotExist(err) {
        return fmt.Errorf("valet-sh not found in %s", repoPath)
    }
    return nil
}

func updateNextBranch(repoPath string) error {
    fmt.Println("Updating valet-sh to the latest version on the next branch")

    if err := git.CheckoutBranch(repoPath, "next"); err != nil {
        return fmt.Errorf("Failed to checkout next branch: %w", err)
    }

    if err := git.PullLatest(repoPath); err != nil {
        return fmt.Errorf("Failed to pull latest changes: %w", err)
    }

    return runtimeUpdate()
}

func updateStableBranch(repoPath string) error {
    fmt.Println("Updating valet-sh to the latest version on the stable branch")

    if err := git.FetchTags(repoPath); err != nil {
        return fmt.Errorf("Failed to fetch tags: %w", err)
    }

    currentRelease, err := git.GetCurrentReleaseTag(repoPath)
    if err != nil {
        return fmt.Errorf("Failed to get current release: %w", err)
    }

    tags, err := git.GetAllTags(repoPath)
    if err != nil {
        return fmt.Errorf("Failed to get all tags: %w", err)
    }

    if len(tags) == 0 {
        return fmt.Errorf("No valid releases found")
    }

    majorVersion := constants.ValetMajorVersion
    semverRegex := fmt.Sprintf(`^(%s)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(\-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`, majorVersion)
    validVersions := git.FilterTagsSemver(tags, semverRegex)

    if len(validVersions) == 0 {
        return fmt.Errorf("No valid releases found")
    }

    latestVersion := validVersions[0]
    latestTag := latestVersion
    if !strings.HasPrefix(tags[0], "v") && strings.HasPrefix(latestVersion, "v") {
        latestTag = latestVersion[1:]
    }

    fmt.Printf("Latest version available: %s | current: %s\n", latestTag, currentRelease)

    compareResult := utils.CompareVersions(currentRelease, latestTag)

    if compareResult < 0 {
        fmt.Printf("Updating valet-sh from version %s to %s\n", currentRelease, latestTag)
        if err := git.CheckoutBranch(repoPath, latestTag); err != nil {
            return fmt.Errorf("Failed to checkout version %s: %w", latestTag, err)
        }
        fmt.Printf("valet-sh successfully updated to version %s\n", latestTag)
    } else {
        fmt.Printf("Already on latest version %s\n", currentRelease)
    }

    return runtimeUpdate()
}

func runtimeUpdate() error {
    status, err := runtime.CheckRuntime()
    if err != nil {
        return fmt.Errorf("Failed to check runtime: %w", err)
    }

    if status.NeedsUpdate || status.PackageChanged {
        fmt.Printf("Updating valet-sh to version %s\n", status.CurrentVersion)
        fmt.Println("Updating runtime")

        url := fmt.Sprintf("https://github.com/valet-sh/runtime/releases/download/%s/%s.tar.gz", status.TargetVersion, status.CurrentPackage)

        fmt.Printf("Check if runtime release '%s' exists\n", url)
        resp, err := http.Head(url)
        if err != nil {
            return fmt.Errorf("failed to check runtime release: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("runtime release not found: %s", url)
        }

        tmpDir := constants.ValetVenvTmpPath
        fmt.Println("Cleaning up temporary directory")
        os.RemoveAll(tmpDir)

        venvDir := constants.ValetVenvPath
        _, err = os.Stat(venvDir)
        venvExists := !os.IsNotExist(err)

        if venvExists {
            fmt.Println("Backing up current venv")
            err = os.Rename(venvDir, tmpDir)
            if err != nil {
                return fmt.Errorf("failed to backup current venv: %w", err)
            }
        }

        fmt.Printf("Downloading and unpacking new runtime '%s' ", status.CurrentVersion)
        respDownload, err := http.Get(url)
        if err != nil {
            return fmt.Errorf("failed to download runtime release: %w", err)
        }
        defer respDownload.Body.Close()

        if respDownload.StatusCode != http.StatusOK {
            return fmt.Errorf("Bad status code:  %s", respDownload.Status)
        }

        if err != nil {
            return fmt.Errorf("failed to create venv directory: %w", err)
        }

        err = utils.Untar(constants.ValetPath, respDownload.Body)
        if err != nil {
            if venvExists {
                os.RemoveAll(venvDir)
                os.Rename(tmpDir, venvDir)
            }
            return fmt.Errorf("failed to extract runtime: %w", err)
        }

        venvVersion := status.CurrentPackage + "-" + status.CurrentVersion
        err = os.WriteFile(filepath.Join(constants.ValetVenvPath, constants.VersionFileName),
            []byte(venvVersion), 0644)
        if err != nil {
            return fmt.Errorf("failed to update version file: %w", err)
        }

        fmt.Printf("\n - Runtime '%s' updated successfully\n", status.CurrentVersion)

        if venvExists {
            os.RemoveAll(tmpDir)
        }

    } else {
        fmt.Println("valet-sh runtime is already up to date")
    }

    return nil
}
